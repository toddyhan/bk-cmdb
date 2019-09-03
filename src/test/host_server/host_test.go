package host_server_test

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host test", func() {
	var bizId, setId, moduleId, idleModuleId, faultModuleId int64
	var hostId, hostId1, hostId2, hostId3 int64

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'cc_biz'", func() {
			test.ClearDatabase()

			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cc_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("cc_biz"))
			bizId = int64(rsp.Data["bk_biz_id"].(float64))
		})

		It("create set", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        bizId,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizId,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_set_name"].(string)).To(Equal("test"))
			Expect(int64(rsp.Data["bk_parent_id"].(float64))).To(Equal(bizId))
			Expect(int64(rsp.Data["bk_biz_id"].(float64))).To(Equal(bizId))
			setId = int64(rsp.Data["bk_set_id"].(float64))
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_module_name"].(string)).To(Equal("cc_module"))
			Expect(int64(rsp.Data["bk_set_id"].(float64))).To(Equal(setId))
			Expect(int64(rsp.Data["bk_parent_id"].(float64))).To(Equal(setId))
			moduleId = int64(rsp.Data["bk_module_id"].(float64))
		})
	})

	Describe("add host test", func() {
		It("add host using api", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "1.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search host created using api", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"1.0.0.1"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.1"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_api_asset_1"))
			hostId1 = int64(data["bk_host_id"].(float64))
		})

		It("add host using excel", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"5": map[string]interface{}{
						"bk_host_innerip": "1.0.0.2",
						"bk_asset_id":     "addhost_excel_asset_1",
						"bk_host_name":    "1.0.0.2",
					},
				},
				"input_type": "excel",
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search host created using excel", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"1.0.0.2"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_excel_asset_1"))
			Expect(data["bk_host_name"].(string)).To(Equal("1.0.0.2"))
			hostId2 = int64(data["bk_host_id"].(float64))
		})

		It("search host using multiple ips", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data: []string{
						"1.0.0.1",
						"1.0.0.2",
					},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		// This api is marked as deprecated
		// It("add host using agent", func() {
		// 	input := map[string]interface{}{
		// 		"host_info": map[string]interface{}{
		// 			"bk_host_innerip": "1.0.0.3",
		// 			"bk_asset_id":     "addhost_agent_asset_1",
		// 			"bk_cloud_id":     0,
		// 		},
		// 	}
		// 	rsp, err := hostServerClient.AddHostFromAgent(context.Background(), header, input)
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(rsp.Result).To(Equal(true))
		// })

		// It("search host created using agent", func() {
		// 	input := &params.HostCommonSearch{
		// 		AppID: int(bizId),
		// 		Ip: params.IPInfo{
		// 			Data:  []string{"1.0.0.3"},
		// 			Exact: 1,
		// 			Flag:  "bk_host_innerip|bk_host_outerip",
		// 		},
		// 		Page: params.PageInfo{
		// 			Sort: "bk_host_id",
		// 		},
		// 	}
		// 	rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(rsp.Result).To(Equal(true))
		// 	Expect(rsp.Data.Count).To(Equal(1))
		// 	data := rsp.Data.Info[0]["host"].(map[string]interface{})
		// 	Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.3"))
		// 	Expect(data["bk_asset_id"].(string)).To(Equal("addhost_agent_asset_1"))
		// })

		It("add host to resource", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "1.0.0.4",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.4"))
			hostId = int64(data["bk_host_id"].(float64))
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			for _, data := range rsp.Data {
				if data.PropertyID == "bk_host_innerip" {
					Expect(data.PropertyValue).To(Equal("1.0.0.4"))
					break
				}
			}
		})
	})

	Describe("transfer host test", func() {
		It("search biz host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.1"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
		})

		It("transfer resourcehost to idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.1"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
			Expect(rsp.Data.Info[2]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.4"))
		})

		It("transfer host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.4"))
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host same module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			host := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(host["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
			Expect(int64(host["bk_host_id"].(float64))).To(Equal(hostId2))
			module := rsp.Data.Info[0]["module"].([]interface{})[0].(map[string]interface{})
			Expect(module["bk_module_name"].(string)).To(Equal("cc_module"))
			Expect(int64(module["bk_module_id"].(float64))).To(Equal(moduleId))
		})

		It("clone host", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   1,
				OrgIP:   "1.0.0.1",
				DstIP:   "1.0.0.5",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search cloned host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Ip: params.IPInfo{
					Data:  []string{"1.0.0.5"},
					Exact: 0,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.5"))
		})

		It("get instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", strconv.FormatInt(bizId, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(2))
			idleModuleId = rsp.Data.Module[0].ModuleID
			faultModuleId = rsp.Data.Module[1].ModuleID
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer host to fault module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					faultModuleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.2"))
		})

		It("search transfered module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer fault host to idle module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					idleModuleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.4"))
		})

		It("transfer host to idle module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect("1.0.0.2").To(SatisfyAny(Equal(data["bk_host_innerip"].(string)), Equal(data1["bk_host_innerip"].(string))))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("transfer host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
		})

		It("search resource host change start limit", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort:  "bk_host_id",
					Start: 2,
					Limit: 2,
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(len(rsp.Data.Info)).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.5"))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})

		It("sync host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "1.0.0.6",
						"bk_asset_id":     "host_sync_asset_1",
						"bk_cloud_id":     0,
					},
				},
				"bk_biz_id": bizId,
				"bk_module_id": []int64{
					idleModuleId,
				},
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.6"))
			hostId3 = int64(data["bk_host_id"].(float64))
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId3,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("1.0.0.4"))
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data1["bk_host_innerip"].(string)).To(Equal("1.0.0.6"))
		})

		It("move all module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId,
				SetID:         setId,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("update idle host", func() {
			input := map[string]interface{}{
				"bk_host_id":   fmt.Sprintf("%v,%v", hostId, hostId3),
				"bk_host_name": "update_host_name",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_name"].(string)).To(Equal("update_host_name"))
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data1["bk_host_name"].(string)).To(Equal("update_host_name"))
		})

		It("delete resource host", func() {
			input := map[string]interface{}{
				"bk_host_id":          fmt.Sprintf("%v,%v", hostId1, hostId2),
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})
	})
})
