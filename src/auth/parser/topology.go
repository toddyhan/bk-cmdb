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

package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (ps *parseStream) topology() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.business().
		mainline().
		associationType().
		objectAssociation().
		objectInstanceAssociation().
		objectInstance().
		object().
		ObjectClassification().
		objectAttributeGroup().
		objectAttribute().
		ObjectModule().
		ObjectSet().
		objectUnique().
		audit().
		instanceAudit().
		privilege().
		fullTextSearch()

	return ps
}

var (
	createBusinessRegexp             = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/?$`)
	updateBusinessRegexp             = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/[0-9]+/?$`)
	deleteBusinessRegexp             = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/[0-9]+/?$`)
	findBusinessRegexp               = regexp.MustCompile(`^/api/v3/biz/search/[^\s/]+/?$`)
	findResourcePoolBusinessRegexp   = regexp.MustCompile(`^/api/v3/biz/default/[^\s/]+/search/?$`)
	createResourcePoolBusinessRegexp = regexp.MustCompile(`^/api/v3/biz/default/[^\s/]+/?$`)
	updateBusinessStatusRegexp       = regexp.MustCompile(`^/api/v3/biz/status/[^\s/]+/[^\s/]+/[0-9]+/?$`)
)

const findReducedBusinessList = `/api/v3/biz/with_reduced`

func (ps *parseStream) business() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find reduced business list for the user with any business resources
	if ps.hitPattern(findReducedBusinessList, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create business, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(createBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(createResourcePoolBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// update business, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(updateBusinessRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid update business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Update,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// update business enable status, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(updateBusinessRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("invalid update business enable status request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business enable status, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Update,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// delete business, this is not a normalize api.
	// TODO: update this api format
	if ps.hitRegexp(updateBusinessRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid delete business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Delete,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// find business, this is not a normalize api.
	// TODO: update this api format
	if ps.hitRegexp(findBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.FindMany,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	if ps.hitRegexp(findResourcePoolBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.FindMany,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	if ps.hitRegexp(updateBusinessStatusRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Archive,
					InstanceID: bizID,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	return ps
}

const (
	createMainlineObjectPattern = "/api/v3/topo/model/mainline"
)

var (
	deleteMainlineObjectRegexp                      = regexp.MustCompile(`^/api/v3/topo/model/mainline/owners/[^\s/]+/objectids/[^\s/]+/?$`)
	findMainlineObjectTopoRegexp                    = regexp.MustCompile(`^/api/v3/topo/model/[^\s/]+/?$`)
	findMainlineInstanceTopoRegexp                  = regexp.MustCompile(`^/api/v3/topo/inst/[^\s/]+/[0-9]+/?$`)
	findMainineSubInstanceTopoRegexp                = regexp.MustCompile(`^/api/v3/topo/inst/child/[^\s/]+/[^\s/]+/[0-9]+/[0-9]+/?$`)
	findMainlineIdleFaultModuleRegexp               = regexp.MustCompile(`^/api/v3/topo/internal/[^\s/]+/[0-9]+/?$`)
	findMainlineIdleFaultModuleWithStatisticsRegexp = regexp.MustCompile(`^/api/v3/topo/internal/[^\s/]+/[0-9]+/with_statistics/?$`)
)

func (ps *parseStream) mainline() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create mainline object operation.
	if ps.hitPattern(createMainlineObjectPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete mainline object operation
	if ps.hitRegexp(deleteMainlineObjectRegexp, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Delete,
				},
			},
		}

		return ps
	}

	// get mainline object operation
	if ps.hitRegexp(findMainlineObjectTopoRegexp, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.MainlineModelTopology,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find mainline instance topology operation.
	if ps.hitRegexp(findMainlineInstanceTopoRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline instance topology, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstanceTopology,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find mainline object instance's sub-instance topology operation.
	if ps.hitRegexp(findMainineSubInstanceTopoRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("find mainline object's sub instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline object's sub instance topology, but got invalid business id %s", ps.RequestCtx.Elements[7])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstanceTopology,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find mainline internal idle and fault module operation.
	if ps.hitRegexp(findMainlineIdleFaultModuleRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline idle and fault module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline idle and fault module, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find mainline internal idle and fault module with statistics operation.
	if ps.hitRegexp(findMainlineIdleFaultModuleWithStatisticsRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find mainline idle and fault module with statistics, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline idle and fault module with statistics, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.Find,
				},
			},
		}

		return ps
	}
	return ps
}

const (
	findManyAssociationKindPattern = "/api/v3/topo/association/type/action/search"
	createAssociationKindPattern   = "/api/v3/topo/association/type/action/create"
)

var (
	updateAssociationKindRegexp = regexp.MustCompile(`^/api/v3/topo/association/type/[0-9]+/action/update$`)
	deleteAssociationKindRegexp = regexp.MustCompile(`^/api/v3/topo/association/type/[0-9]+/action/delete$`)
)

func (ps *parseStream) associationType() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find association kind operation
	if ps.hitPattern(findManyAssociationKindPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create association kind operation
	if ps.hitPattern(createAssociationKindPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update association kind operation
	if ps.hitRegexp(updateAssociationKindRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Update,
					InstanceID: kindID,
				},
			},
		}

		return ps
	}

	// delete association kind operation
	if ps.hitRegexp(deleteAssociationKindRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Delete,
					InstanceID: kindID,
				},
			},
		}

		return ps
	}

	return ps
}

const (
	findObjectAssociationPattern                    = "/api/v3/object/association/action/search"
	createObjectAssociationPattern                  = "/api/v3/object/association/action/create"
	findObjectAssociationWithAssociationKindPattern = "/api/v3/topo/association/type/action/search/batch"
)

var (
	updateObjectAssociationRegexp = regexp.MustCompile(`^/api/v3/object/association/[0-9]+/action/update$`)
	deleteObjectAssociationRegexp = regexp.MustCompile(`^/api/v3/object/association/[0-9]+/action/delete$`)
)

func (ps *parseStream) objectAssociation() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// search object association operation
	if ps.RequestCtx.URI == findObjectAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create object association operation
	if ps.RequestCtx.URI == createObjectAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update object association operation
	if ps.hitRegexp(updateObjectAssociationRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("update object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelAssociation,
					Action:     meta.Update,
					InstanceID: assoID,
				},
			},
		}
		return ps
	}

	// delete object association operation
	if ps.hitRegexp(deleteObjectAssociationRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("delete object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelAssociation,
					Action:     meta.Delete,
					InstanceID: assoID,
				},
			},
		}
		return ps
	}

	// find object association with a association kind list.
	if ps.RequestCtx.URI == findObjectAssociationWithAssociationKindPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectInstanceAssociationPattern   = "/api/v3/inst/association/action/search"
	createObjectInstanceAssociationPattern = "/api/v3/inst/association/action/create"
)

var (
	deleteObjectInstanceAssociationRegexp = regexp.MustCompile("/api/v3/inst/association/[0-9]+/action/delete")
)

func (ps *parseStream) objectInstanceAssociation() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find object instance's association operation.
	if ps.RequestCtx.URI == findObjectInstanceAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create object's instance association operation.
	if ps.RequestCtx.URI == createObjectInstanceAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceAssociation,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete object's instance association operation.
	if ps.hitRegexp(deleteObjectInstanceAssociationRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("delete object instance association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelInstanceAssociation,
					Action:     meta.Delete,
					InstanceID: assoID,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectBatchRegexp = `/api/v3/object/search/batch`
)

var (
	createObjectInstanceRegexp          = regexp.MustCompile(`^/api/v3/inst/[^\s/]+/[^\s/]+/?$`)
	findObjectInstanceRegexp            = regexp.MustCompile(`^/api/v3/inst/association/search/owner/[^\s/]+/object/[^\s/]+/?$`)
	updateObjectInstanceRegexp          = regexp.MustCompile(`^/api/v3/inst/[^\s/]+/[^\s/]+/[0-9]+/?$`)
	updateObjectInstanceBatchRegexp     = regexp.MustCompile(`^/api/v3/inst/[^\s/]+/[^\s/]+/batch$`)
	deleteObjectInstanceBatchRegexp     = regexp.MustCompile(`^/api/v3/inst/[^\s/]+/[^\s/]+/batch$`)
	deleteObjectInstanceRegexp          = regexp.MustCompile(`^/api/v3/inst/[^\s/]+/[^\s/]+/[0-9]+/?$`)
	findObjectInstanceSubTopologyRegexp = regexp.MustCompile(`^/api/v3/inst/association/topo/search/owner/[^\s/]+/object/[^\s/]+/inst/[0-9]+/?$`)
	findObjectInstanceTopologyRegexp    = regexp.MustCompile(`^/api/v3/inst/association/topo/search/owner/[^\s/]+/object/[^\s/]+/inst/[0-9]+/?$`)
	findBusinessInstanceTopologyRegexp  = regexp.MustCompile(`^/api/v3/topo/inst/[^\s/]+/[0-9]+/?$`)
	findObjectInstancesRegexp           = regexp.MustCompile(`^/api/v3/inst/search/owner/[^\s/]+/object/[^\s/]+/?$`)
	findObjectInstancesDetailRegexp     = regexp.MustCompile(`^/api/v3/inst/search/owner/[^\s/]+/object/[^\s/]+/detail/?$`)
)

func (ps *parseStream) objectInstance() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object instance operation.
	if ps.hitRegexp(createObjectInstanceRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.Create,
				},
				BusinessID: bizID,
			},
		}
		return ps
	}

	// find object instance operation.
	if ps.hitRegexp(findObjectInstanceRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("search object instance, but got invalid url")
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.Find,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[8],
					},
				},
			},
		}
		return ps
	}

	// update object instance operation.
	if ps.hitRegexp(updateObjectInstanceRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update object instance, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object instance, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelInstance,
					Action:     meta.Update,
					InstanceID: instID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[4],
					},
				},
			},
		}
		return ps
	}

	// update object instance batch operation.
	if ps.hitRegexp(updateObjectInstanceBatchRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update object instance batch, but got invalid url")
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.UpdateMany,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[4],
					},
				},
			},
		}
		return ps
	}

	// delete object instance batch operation.
	if ps.hitRegexp(deleteObjectInstanceBatchRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object instance batch, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.DeleteMany,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[4],
					},
				},
			},
		}
		return ps
	}

	// delete object instance operation.
	if ps.hitRegexp(deleteObjectInstanceRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object instance, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelInstance,
					Action:     meta.Delete,
					InstanceID: instID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[4],
					},
				},
			},
		}
		return ps
	}

	// find object instance topology operation
	if ps.hitRegexp(findObjectInstanceSubTopologyRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 12 {
			ps.err = errors.New("find object instance topology, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[11], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find object instance topology, but got invalid instance id %s", ps.RequestCtx.Elements[11])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelInstanceTopology,
					Action:     meta.Find,
					InstanceID: instID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[9],
					},
				},
			},
		}
		return ps
	}

	// find object instance fully topology operation.
	if ps.hitRegexp(findObjectInstanceTopologyRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 12 {
			ps.err = errors.New("find object instance topology, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[11], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find object instance, but get instance id %s", ps.RequestCtx.Elements[11])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelInstanceTopology,
					Action:     meta.Find,
					InstanceID: instID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[9],
					},
				},
			},
		}

		return ps
	}

	// find business instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find business instance topology, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.MainlineInstanceTopology,
					Action:     meta.Find,
					InstanceID: bizID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: string(meta.Business),
					},
				},
			},
		}
		return ps
	}

	// find object's instance list operation
	if ps.hitRegexp(findObjectInstancesRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find object's instance  list, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[7],
					},
				},
			},
		}
		return ps
	}

	// find object/s instance list details operation.
	if ps.hitRegexp(findObjectInstancesDetailRegexp, http.MethodPost) {
		// TODO: parse these query condition
		models, err := ps.getModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[7]})
		if err != nil {
			ps.err = err
			return ps
		}

		for _, model := range models {
			bizID, err := metadata.BizIDFromMetadata(model.Metadata)
			if err != nil {
				ps.err = err
				return ps
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: bizID,
					Basic: meta.Basic{
						Type:   meta.ModelInstance,
						Action: meta.FindMany,
					},
					Layers: []meta.Item{
						{
							Type: meta.Model,
							Name: ps.RequestCtx.Elements[7],
						},
					},
				},
			)
		}

		return ps
	}

	if ps.hitPattern(findObjectBatchRegexp, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil && err != metadata.LabelKeyNotExistError {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectPattern       = "/api/v3/object"
	findObjectsPattern        = "/api/v3/objects"
	findObjectTopologyPattern = "/api/v3/objects/topo"
	createObjectBatchPattern  = "/api/v3/object/batch"
	objectStatistics          = "/api/v3/object/statistics"
)

var (
	deleteObjectRegexp                = regexp.MustCompile(`^/api/v3/object/[0-9]+/?$`)
	updateObjectRegexp                = regexp.MustCompile(`^/api/v3/object/[0-9]+/?$`)
	findObjectTopologyGraphicRegexp   = regexp.MustCompile(`^/api/v3/objects/topographics/scope_type/[^\s/]+/scope_id/[^\s/]+/action/search$`)
	updateObjectTopologyGraphicRegexp = regexp.MustCompile(`^/api/v3/objects/topographics/scope_type/[^\s/]+/scope_id/[^\s/]+/action/[a-z]+/?$`)
)

func (ps *parseStream) object() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create common object operation.
	if ps.hitPattern(createObjectPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// create common object batch operation.
	if ps.hitPattern(createObjectBatchPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("import object, but parse biz id failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.UpdateMany,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(objectStatistics, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// delete object operation
	if ps.hitRegexp(deleteObjectRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 4 {
			ps.err = errors.New("delete object, but got invalid url")
			return ps
		}

		objID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but got invalid object's id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Delete,
					InstanceID: objID,
				},
			},
		}
		return ps
	}

	// update object operation.
	if ps.hitRegexp(updateObjectRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 4 {
			ps.err = errors.New("update object, but got invalid url")
			return ps
		}

		objID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object, but got invalid object's id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Update,
					InstanceID: objID,
				},
			},
		}
		return ps
	}

	// get object operation.
	if ps.hitPattern(findObjectsPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find object's topology operation.
	if ps.hitPattern(findObjectTopologyPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find object's topology graphic operation.
	if ps.hitRegexp(findObjectTopologyGraphicRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// update object's topology graphic operation.
	if ps.hitRegexp(updateObjectTopologyGraphicRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectClassificationPattern   = "/api/v3/object/classification"
	findObjectClassificationListPattern = "/api/v3/object/classifications"
)

var (
	deleteObjectClassificationRegexp         = regexp.MustCompile("^/api/v3/object/classification/[0-9]+/?$")
	updateObjectClassificationRegexp         = regexp.MustCompile("^/api/v3/object/classification/[0-9]+/?$")
	findObjectsBelongsToClassificationRegexp = regexp.MustCompile(`^/api/v3/object/classification/[^\s/]+/objects$`)
)

func (ps *parseStream) ObjectClassification() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's classification operation.
	if ps.hitPattern(createObjectClassificationPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete object's classification operation.
	if ps.hitRegexp(deleteObjectClassificationRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, but got invalid object's id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Delete,
					InstanceID: classID,
				},
			},
		}
		return ps
	}

	// update object's classification operation.
	if ps.hitRegexp(updateObjectClassificationRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object classification, but got invalid object's  classification id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Update,
					InstanceID: classID,
				},
			},
		}
		return ps
	}

	// find object's classification list operation.
	if ps.hitPattern(findObjectClassificationListPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find all the objects belongs to a classification
	if ps.hitRegexp(findObjectsBelongsToClassificationRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributeGroupPattern         = "/api/v3/objectatt/group/new"
	updateObjectAttributeGroupPattern         = "/api/v3/objectatt/group/update"
	updateObjectAttributeGroupPropertyPattern = "/api/v3/objectatt/group/property"
)

var (
	findObjectAttributeGroupRegexp     = regexp.MustCompile(`^/api/v3/objectatt/group/property/owner/[^\s/]+/object/[^\s/]+/?$`)
	deleteObjectAttributeGroupRegexp   = regexp.MustCompile(`^/api/v3/objectatt/group/groupid/[0-9]+/?$`)
	removeAttributeAwayFromGroupRegexp = regexp.MustCompile(`^/api/v3/objectatt/group/owner/[^\s/]+/object/[^\s/]+/propertyids/[^\s/]+/groupids/[^\s/]+/?$`)
)

func (ps *parseStream) objectAttributeGroup() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's attribute group operation.
	if ps.hitPattern(createObjectAttributeGroupPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// find object's attribute group operation.
	if ps.hitRegexp(findObjectAttributeGroupRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("find object's attribute group, but got invalid uri")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Find,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[8],
					},
				},
			},
		}
		return ps
	}

	// update object's attribute group operation.
	if ps.hitPattern(updateObjectAttributeGroupPattern, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(updateObjectAttributeGroupPropertyPattern, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	// delete object's attribute group operation.
	if ps.hitRegexp(deleteObjectAttributeGroupRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object's attribute group, but got invalid url")
			return ps
		}

		groupID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object's attribute group, but got invalid group's id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelAttributeGroup,
					Action:     meta.Delete,
					InstanceID: groupID,
				},
			},
		}
		return ps
	}

	// remove a object's attribute away from a group.
	if ps.hitRegexp(removeAttributeAwayFromGroupRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 12 {
			ps.err = errors.New("remove a object attribute away from a group, but got invalid uri")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Delete,
					Name:   ps.RequestCtx.Elements[11],
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributePattern = "/api/v3/object/attr"
	findObjectAttributePattern   = "/api/v3/object/attr/search"
)

var (
	deleteObjectAttributeRegexp = regexp.MustCompile(`^/api/v3/object/attr/[0-9]+/?$`)
	updateObjectAttributeRegexp = regexp.MustCompile(`^/api/v3/object/attr/[0-9]+/?$`)
)

func (ps *parseStream) objectAttribute() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's attribute operation.
	if ps.hitPattern(createObjectAttributePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete object's attribute operation.
	if ps.hitRegexp(deleteObjectAttributeRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Delete,
					InstanceID: attrID,
				},
			},
		}
		return ps
	}

	// update object attribute operation
	if ps.hitRegexp(updateObjectAttributeRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Update,
					InstanceID: attrID,
				},
			},
		}
		return ps
	}

	// get object's attribute operation.
	if ps.hitPattern(findObjectAttributePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	createModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/?$`)
	deleteModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)
	updateModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)
	findModuleRegexp   = regexp.MustCompile(`^/api/v3/module/search/[^\s/]+/[0-9]+/[0-9]+/?$`)
)

func (ps *parseStream) ObjectModule() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create module
	if ps.hitRegexp(createModuleRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("create module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.Create,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// delete module operation.
	if ps.hitRegexp(deleteModuleRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		moduleID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid module id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelModule,
					Action:     meta.Delete,
					InstanceID: moduleID,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// update module operation.
	if ps.hitRegexp(updateModuleRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		moduleID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid module id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelModule,
					Action:     meta.Update,
					InstanceID: moduleID,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// find module operation.
	if ps.hitRegexp(findModuleRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module, but got invalid set id %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelSet,
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	return ps
}

var (
	createSetRegexp     = regexp.MustCompile(`^/api/v3/set/[0-9]+/?$`)
	deleteSetRegexp     = regexp.MustCompile(`^/api/v3/set/[0-9]+/[0-9]+/?$`)
	deleteManySetRegexp = regexp.MustCompile(`^/api/v3/set/[0-9]+/batch$`)
	updateSetRegexp     = regexp.MustCompile(`^/api/v3/set/[0-9]+/[0-9]+/?$`)
	findSetRegexp       = regexp.MustCompile(`^/api/v3/set/search/[^\s/]+/[0-9]+/?$`)
)

func (ps *parseStream) ObjectSet() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create set
	if ps.hitRegexp(createSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 4 {
			ps.err = errors.New("create set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete set operation.
	if ps.hitRegexp(deleteSetRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelSet,
					Action:     meta.Delete,
					InstanceID: setID,
				},
			},
		}
		return ps
	}

	// delete many set operation.
	if ps.hitRegexp(deleteManySetRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete set list, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set list, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.DeleteMany,
				},
			},
		}
		return ps
	}

	// update set operation.
	if ps.hitRegexp(updateSetRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update set, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelSet,
					Action:     meta.Update,
					InstanceID: setID,
				},
			},
		}
		return ps
	}

	// find set operation.
	if ps.hitRegexp(findSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find set, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	createObjectUniqueRegexp = regexp.MustCompile(`^/api/v3/object/[^\s/]+/unique/action/create$`)
	updateObjectUniqueRegexp = regexp.MustCompile(`^/api/v3/object/[^\s/]+/unique/[0-9]+/action/update$`)
	deleteObjectUniqueRegexp = regexp.MustCompile(`^/api/v3/object/[^\s/]+/unique/[0-9]+/action/delete$`)
	findObjectUniqueRegexp   = regexp.MustCompile(`^/api/v3/object/[^\s/]+/unique/action/search$`)
)

func (ps *parseStream) objectUnique() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// add object unique operation.
	if ps.hitRegexp(createObjectUniqueRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.Create,
					Name:   ps.RequestCtx.Elements[3],
				},
			},
		}
		return ps
	}

	// update object unique operation.
	if ps.hitRegexp(updateObjectUniqueRegexp, http.MethodPut) {
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Update,
					InstanceID: uniqueID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[3],
					},
				},
			},
		}
		return ps
	}

	// delete object unique operation.
	if ps.hitRegexp(deleteObjectUniqueRegexp, http.MethodDelete) {
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Delete,
					InstanceID: uniqueID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[3],
					},
				},
			},
		}
		return ps
	}

	// find object unique operation.
	if ps.hitRegexp(findObjectUniqueRegexp, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[5],
					},
				},
			},
		}
		return ps
	}

	return ps
}

var (
	searchAuditlog               = `/api/v3/audit/search`
	searchInstanceAuditlogRegexp = regexp.MustCompile(`^/api/v3/object/[^\s/]+/audit/search/?$`)
)

func (ps *parseStream) audit() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// add object unique operation.
	if ps.hitPattern(searchAuditlog, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type: meta.AuditLog,
					// audit authorization in topo scene layer
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

func (ps *parseStream) instanceAudit() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// add object unique operation.
	if ps.hitRegexp(searchInstanceAuditlogRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type: meta.AuditLog,
					// instance audit authorization by instance
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	findPrivilege = regexp.MustCompile(`^/api/v3/topo/privilege/.*$`)
	postPrivilege = regexp.MustCompile(`^/api/v3/topo/privilege/.*$`)
)

func (ps *parseStream) privilege() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findPrivilege, http.MethodGet) || ps.hitRegexp(findPrivilege, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	fullTextSearchPattern = "/api/v3/find/full_text"
)

func (ps *parseStream) fullTextSearch() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findPrivilege, http.MethodGet) || ps.hitPattern(fullTextSearchPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}
