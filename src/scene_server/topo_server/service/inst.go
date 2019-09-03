/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/operation"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateInst create a new inst
func (s *Service) CreateInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	objID := pathParams("bk_obj_id")
	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("failed to search the inst, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	if data.Exists("BatchInfo") {
		/*
		   BatchInfo data format:
		    {
		      "BatchInfo": {
		        "4": { // excel line number
		          "bk_inst_id": 1,
		          "bk_inst_key": "a22",
		          "bk_inst_name": "a11",
		          "bk_version": "121",
		          "import_from": "1"
		        },
		      "input_type": "excel"
		    }
		*/
		batchInfo := new(operation.InstBatchInfo)
		if err := data.MarshalJSONInto(batchInfo); err != nil {
			blog.Errorf("create instance failed, import object[%s] instance batch, but got invalid BatchInfo:[%v], err: %+v, rid: %s", objID, batchInfo, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		setInst, err := s.Core.InstOperation().CreateInstBatch(params, obj, batchInfo)
		if nil != err {
			blog.Errorf("failed to create new object %s, %s, rid: %s", objID, err.Error(), params.ReqID)
			return nil, err
		}

		// auth register new created
		if len(setInst.SuccessCreated) != 0 {
			if err := s.AuthManager.RegisterInstancesByID(params.Context, params.Header, objID, setInst.SuccessCreated...); err != nil {
				blog.Errorf("create instance success, but register instances to iam failed, instances: %+v, err: %+v, rid: %s", setInst.SuccessCreated, err, params.ReqID)
				return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
			}
		}

		// auth update registered instances
		if len(setInst.SuccessUpdated) == 0 {
			if err := s.AuthManager.UpdateRegisteredInstanceByID(params.Context, params.Header, objID, setInst.SuccessUpdated...); err != nil {
				blog.Errorf("update registered instances to iam failed, err: %+v, rid: %s", err, params.ReqID)
				return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
			}
		}

		return setInst, nil
	}

	setInst, err := s.Core.InstOperation().CreateInst(params, obj, data)
	if nil != err {
		blog.Errorf("failed to create a new %s, %s, rid: %s", objID, err.Error(), params.ReqID)
		return nil, err
	}

	instanceID, err := setInst.GetInstID()
	if err != nil {
		blog.Errorf("create instance failed, unexpected error, create instance success, but get id failed, instance: %+v, err: %+v, rid: %s", setInst, err, params.ReqID)
		return nil, err
	}

	// auth: register instances to iam
	if err := s.AuthManager.RegisterInstancesByID(params.Context, params.Header, objID, instanceID); err != nil {
		blog.Errorf("create instance success, but register instance to iam failed, instance: %d, err: %+v, rid: %s", instanceID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return setInst.ToMapStr(), nil
}

func (s *Service) DeleteInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, pathParams("bk_obj_id"))
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	deleteCondition := &operation.OpCondition{}
	if err := data.MarshalJSONInto(deleteCondition); nil != err {
		return nil, err
	}

	// auth: deregister resources
	if err := s.AuthManager.DeregisterInstanceByRawID(params.Context, params.Header, obj.GetObjectID(), deleteCondition.Delete.InstID...); err != nil {
		blog.Errorf("batch delete instance failed, deregister instance failed, instID: %d, err: %s, rid: %s", deleteCondition.Delete.InstID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	return nil, s.Core.InstOperation().DeleteInstByInstID(params, obj, deleteCondition.Delete.InstID, true)
}

// DeleteInst delete the inst
func (s *Service) DeleteInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	if "batch" == pathParams("inst_id") {
		return s.DeleteInsts(params, pathParams, queryParams, data)
	}

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "inst id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, pathParams("bk_obj_id"))
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	// auth: deregister resources
	if err := s.AuthManager.DeregisterInstanceByRawID(params.Context, params.Header, obj.GetObjectID(), instID); err != nil {
		blog.Errorf("delete instance failed, deregister instance failed, instID: %d, err: %s, rid: %s", instID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.InstOperation().DeleteInstByInstID(params, obj, []int64{instID}, true)
	return nil, err
}

func (s *Service) UpdateInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")

	updateCondition := &operation.OpCondition{}
	if err := data.MarshalJSONInto(updateCondition); nil != err {
		blog.Errorf("[api-inst] failed to parse the input data(%v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}

	// check inst_id field to be not empty, is dangerous for empty inst_id field, which will update or delete all instance
	for idx, item := range updateCondition.Update {
		if item.InstID == 0 {
			blog.Errorf("update instance failed, %d's update item's field `inst_id` emtpy, rid: %s", idx, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommParamsInvalid)
		}
	}
	for idx, instID := range updateCondition.Delete.InstID {
		if instID == 0 {
			blog.Errorf("update instance failed, %d's delete item's field `inst_id` emtpy, rid: %s", idx, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommParamsInvalid)
		}
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	instanceIDs := make([]int64, 0)
	for _, item := range updateCondition.Update {
		instanceIDs = append(instanceIDs, item.InstID)
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(item.InstID)
		err = s.Core.InstOperation().UpdateInst(params, item.InstInfo, obj, cond, item.InstID)
		if nil != err {
			blog.Errorf("[api-inst] failed to update the object(%s) inst (%d),the data (%#v), error info is %s, rid: %s", obj.Object().ObjectID, item.InstID, data, err.Error(), params.ReqID)
			return nil, err
		}
	}

	// auth: deregister resources
	if err := s.AuthManager.UpdateRegisteredInstanceByID(params.Context, params.Header, objID, instanceIDs...); err != nil {
		blog.Errorf("update inst success, but update register to iam failed, instanceIDs: %+v, err: %+v, rid: %s", instanceIDs, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return nil, nil
}

// UpdateInst update the inst
func (s *Service) UpdateInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	if "batch" == pathParams("inst_id") {
		return s.UpdateInsts(params, pathParams, queryParams, data)
	}

	objID := pathParams("bk_obj_id")
	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "inst id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	// this is a special logic for mainline object instance.
	// for auth reason, the front's request add metadata for mainline model's instance update.
	// but actually, it's should not add metadata field in the request.
	// so, we need remove it from the data if it's a mainline model instance.
	yes, err := s.Core.AssociationOperation().IsMainlineObject(params, objID)
	if err != nil {
		return nil, err
	}
	if yes {
		data.Remove("metadata")
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	err = s.Core.InstOperation().UpdateInst(params, data, obj, cond, instID)
	if nil != err {
		blog.Errorf("[api-inst] failed to update the object(%s) inst (%s),the data (%#v), error info is %s, rid: %s", obj.Object().ObjectID, pathParams("inst_id"), data, err.Error(), params.ReqID)
		return nil, err
	}

	// auth: deregister resources
	if err := s.AuthManager.UpdateRegisteredInstanceByID(params.Context, params.Header, objID, instID); err != nil {
		blog.Error("update inst failed, authorization failed, instID: %d, err: %+v, rid: %s", instID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return nil, err
}

// SearchInst search the inst
func (s *Service) SearchInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	//	if nil != params.MetaData {
	//		data.Set(metadata.BKMetadata, *params.MetaData)
	//	}
	// construct the query inst condition
	queryCond := &paraparse.SearchParams{
		Condition: mapstr.New(),
	}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.Core.InstOperation().FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstAndAssociationDetail search the inst with association details
func (s *Service) SearchInstAndAssociationDetail(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")
	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	// construct the query inst condition
	queryCond := &paraparse.SearchParams{
		Condition: mapstr.New(),
	}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.Core.InstOperation().FindInst(params, obj, query, true)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstByObject search the inst of the object
func (s *Service) SearchInstByObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	objID := pathParams("bk_obj_id")
	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	queryCond := &paraparse.SearchParams{
		Condition: mapstr.New(),
	}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start
	cnt, instItems, err := s.Core.InstOperation().FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstByAssociation search inst by the association inst
func (s *Service) SearchInstByAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	objID := pathParams("bk_obj_id")

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	cnt, instItems, err := s.Core.InstOperation().FindInstByAssociationInst(params, obj, data)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstByInstID search the inst by inst ID
func (s *Service) SearchInstByInstID(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoInstSelectFailed, err.Error())
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	queryCond := &metadata.QueryInput{}
	queryCond.Condition = cond.ToMapStr()

	cnt, instItems, err := s.Core.InstOperation().FindInst(params, obj, queryCond, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}

// SearchInstChildTopo search the child inst topo for a inst
func (s *Service) SearchInstChildTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_object_id")

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, err
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	query := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.Core.InstOperation().FindInstChildTopo(params, obj, instID, query)
	return instItems, err

}

// SearchInstTopo search the inst topo
func (s *Service) SearchInstTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	objID := pathParams("bk_obj_id")
	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("search inst topo failed, path parameter inst_id invalid, object: %s inst_id: %s, err: %+v, rid: %s", objID, pathParams("inst_id"), err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", pathParams("bk_obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	query := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.Core.InstOperation().FindInstTopo(params, obj, instID, query)

	return instItems, err
}
