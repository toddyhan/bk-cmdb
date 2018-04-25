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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	sourceAPI "configcenter/src/source_controller/api/object"
	"encoding/json"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

// GetHttpResult get http result
func GetHttpResult(req *restful.Request, url, method string, params interface{}) (bool, string, interface{}) {
	var strParams []byte
	switch params.(type) {
	case string:
		strParams = []byte(params.(string))
	default:
		strParams, _ = json.Marshal(params)

	}
	blog.Info("get request url:%s", url)
	blog.Info("get request info  params:%v", string(strParams))
	reply, err := httpcli.ReqHttp(req, url, method, []byte(strParams))
	blog.Info("get request result:%v", string(reply))
	if err != nil {
		blog.Error("http do error, params:%s, error:%s", strParams, err.Error())
		return false, err.Error(), nil
	}

	addReply, err := simplejson.NewJson([]byte(reply))
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), nil
	}
	isSuccess, err := addReply.Get("result").Bool()
	if nil != err || !isSuccess {
		errMsg, _ := addReply.Get("message").String()
		blog.Error("http do error, url:%s, params:%s, error:%s", url, strParams, errMsg)
		return false, errMsg, addReply.Get("data").Interface()
	}
	return true, "", addReply.Get("data").Interface()
}

//GetObjectFields get object fields
func GetObjectFields(forward *sourceAPI.ForwardParam, ownerID, objID, ObjAddr string) map[string]map[string]interface{} {
	data := make(map[string]interface{})
	data[common.BKOwnerIDField] = ownerID
	data[common.BKObjIDField] = objID
	info, _ := json.Marshal(data)
	client := sourceAPI.NewClient(ObjAddr)
	result, _ := client.SearchMetaObjectAtt(forward, []byte(info))
	fields := make(map[string]map[string]interface{})
	for _, j := range result {
		propertyID := j.PropertyID
		fieldType := j.PropertyType
		switch fieldType {
		case common.FiledTypeSingleChar:
			fields[propertyID] = common.KvMap{"default": "", "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeLongChar:
			fields[propertyID] = common.KvMap{"default": "", "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired} //""
		case common.FiledTypeInt:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired} //0
		case common.FiledTypeEnum:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeDate:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeTime:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeUser:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		default:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
			continue
		}

	}
	return fields
}
