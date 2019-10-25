/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * plat represent cloud plat here
 */

func (am *AuthManager) CollectAllPlats(ctx context.Context, header http.Header) ([]PlatSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{}),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDPlat, &cond)
	if err != nil {
		blog.V(3).Infof("get all plats, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get all plats, err: %+v", err)
	}
	plats := make([]PlatSimplify, 0)
	for _, cls := range result.Data.Info {
		plat := PlatSimplify{}
		_, err = plat.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get all plat failed, err: %+v", err)
		}
		plats = append(plats, plat)
	}
	return plats, nil
}

func (am *AuthManager) collectPlatByIDs(ctx context.Context, header http.Header, platIDs ...int64) ([]PlatSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	platIDs = util.IntArrayUnique(platIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKSubAreaField).In(platIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDPlat, &cond)
	if err != nil {
		blog.V(3).Infof("get plats by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	plats := make([]PlatSimplify, 0)
	for _, cls := range result.Data.Info {
		plat := PlatSimplify{}
		_, err = plat.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get plat by id failed, err: %+v", err)
		}
		plats = append(plats, plat)
	}
	return plats, nil
}

func (am *AuthManager) MakeResourcesByPlat(header http.Header, action meta.Action, plats ...PlatSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, plat := range plats {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Plat,
				Name:       plat.BKCloudNameField,
				InstanceID: plat.BKCloudIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) AuthorizeByPlat(ctx context.Context, header http.Header, action meta.Action, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByPlat(header, action, plats...)

	return am.authorize(ctx, header, 0, resources...)
}

func (am *AuthManager) AuthorizeByPlatIDs(ctx context.Context, header http.Header, action meta.Action, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plat by id failed, err: %+d", err)
	}
	return am.AuthorizeByPlat(ctx, header, action, plats...)
}

func (am *AuthManager) UpdateRegisteredPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredPlatByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredPlat(ctx, header, plats...)
}

func (am *AuthManager) UpdateRegisteredPlatByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredPlat(ctx, header, plats...)
}

func (am *AuthManager) DeregisterPlatByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister plats failed, get plats by id failed, err: %+v", err)
	}
	return am.DeregisterPlat(ctx, header, plats...)
}

func (am *AuthManager) RegisterPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterPlatByID(ctx context.Context, header http.Header, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(platIDs) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	return am.RegisterPlat(ctx, header, plats...)
}

func (am *AuthManager) DeregisterPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterPlatByID(ctx context.Context, header http.Header, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(platIDs) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	return am.DeregisterPlat(ctx, header, plats...)
}
