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

package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

func backup(ctx context.Context, db dal.RDB, opt *option) error {
	dir := filepath.Dir(opt.position)
	now := time.Now().Format("2006_01_02_15_04_05")
	file := filepath.Join(dir, "backup_bk_biz_"+now+".json")
	exportOpt := *opt
	exportOpt.position = file
	exportOpt.mini = false
	exportOpt.scope = "all"
	err := export(ctx, db, &exportOpt)
	if nil != err {
		return err
	}
	fmt.Println("%s business has been backup to \033[35m"+file+"\033[0m", opt.bizName)
	return nil
}

func importBKBiz(ctx context.Context, db dal.RDB, opt *option) error {
	file, err := os.OpenFile(opt.position, os.O_RDONLY, os.ModePerm)
	if nil != err {
		return err
	}
	defer file.Close()

	tar := new(Topo)
	json.NewDecoder(file).Decode(tar)
	if nil != err {
		return err
	}

	cur, err := getBKTopo(ctx, db, opt)
	if err != nil {
		return fmt.Errorf("get src topo faile %s", err.Error())
	}

	if !opt.dryrun {
		err = backup(ctx, db, opt)
		if err != nil {
			return fmt.Errorf("backup faile %s", err)
		}
	}

	if tar.BizTopo != nil {
		// topo check
		if !compareSlice(tar.Mainline, cur.Mainline) {
			return fmt.Errorf("different topo mainline found, your expecting import topo is [%s], but the existing topo is [%s]",
				strings.Join(tar.Mainline, "->"), strings.Join(cur.Mainline, "->"))
		}

		// walk blueking biz and get difference
		ipt := newImporter(ctx, db, opt)
		if err := ipt.walk(true, tar.BizTopo); err != nil {
			blog.Errorf("walk biz topo failed, err: %+v", err)
		}

		// walk to create new node
		err := tar.BizTopo.walk(func(node *Node) error {
			if node.mark == actionCreate {
				fmt.Printf("--- \033[34m%s %s %+v\033[0m\n", node.mark, node.ObjID, node.Data)
				if !opt.dryrun {
					err := db.Table(common.GetInstTableName(node.ObjID)).Insert(ctx, node.Data)
					if nil != err {
						return fmt.Errorf("insert to %s, data:%+v, error: %s", node.ObjID, node.Data, err.Error())
					}
				}
			}
			if node.mark == actionUpdate {
				fmt.Printf("--- \033[36m%s %s %+v\033[0m\n", node.mark, node.ObjID, node.Data)
				if !opt.dryrun {
					instID, err := node.getInstID()
					if nil != err {
						return err
					}
					updateCondition := map[string]interface{}{
						common.GetInstIDField(node.ObjID): instID,
					}

					err = db.Table(common.GetInstTableName(node.ObjID)).Update(ctx, updateCondition, node.Data)
					if nil != err {
						return fmt.Errorf("update to %s by %+v data:%+v, error: %s", node.ObjID, updateCondition, node.Data, err.Error())
					}
				}
			}
			return nil
		})

		if err != nil {
			blog.Errorf("tar biz topo walk failed, err: %+v", err)
		}

		// walk to delete unuse node
		for objID, sdeletes := range ipt.sdelete {
			for _, sdelete := range sdeletes {
				// fmt.Printf("\n--- \033[36mdelete parent node %s %+v\033[0m\n", objID, sdelete)

				instID, err := getInt64(sdelete[common.GetInstIDField(objID)])
				if nil != err {
					return err
				}

				err = cur.BizTopo.walk(func(node *Node) error {
					nodeID, err := node.getInstID()
					if nil != err {
						return err
					}
					if node.ObjID == objID && nodeID == instID {
						childErr := node.walk(func(child *Node) error {
							childID, err := child.getInstID()
							if nil != err {
								return err
							}
							if child.ObjID == common.BKInnerObjIDModule {
								// if should delete module then check whether it has host
								moduleHostCondition := map[string]interface{}{
									common.BKModuleIDField: childID,
								}
								count, err := db.Table(common.BKTableNameModuleHostConfig).Find(moduleHostCondition).Count(ctx)
								if nil != err {
									return fmt.Errorf("get host count error: %s", err.Error())
								}
								if count > 0 {
									return fmt.Errorf("there are %d hosts binded to module %v, please unbind them first and try again ", count, node.Data[common.BKModuleNameField])
								}
							}

							deleteCondition := map[string]interface{}{
								common.GetInstIDField(child.ObjID): childID,
							}
							switch child.ObjID {
							case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
							default:
								deleteCondition[common.BKObjIDField] = child.ObjID
							}
							fmt.Printf("--- \033[31mdelete %s %+v by %+v\033[0m\n", child.ObjID, child.Data, deleteCondition)
							if !opt.dryrun {

								err = db.Table(common.GetInstTableName(child.ObjID)).Delete(ctx, deleteCondition)
								if nil != err {
									return fmt.Errorf("delete %s by %+v, error: %s", child.ObjID, deleteCondition, err.Error())
								}
							}
							return nil
						})
						if childErr != nil {
							return childErr
						}
					}
					return nil
				})
				if err != nil && err.Error() != "break" {
					return err
				}

			}
		}
	}

	bizID, err := cur.BizTopo.getInstID()
	if nil != err {
		return err
	}

	return importProcess(ctx, db, opt, cur.ProcTopos, tar.ProcTopos, bizID)
}

func importProcess(ctx context.Context, db dal.RDB, opt *option, cur, tar *ProcessTopo, bizID uint64) (err error) {
	if tar == nil {
		return nil
	}

	curProcs := map[string]*Process{}
	for _, topo := range cur.Processes {
		curProcs[topo.Data[common.BKProcessNameField].(string)] = topo
	}
	tarProcs := map[string]*Process{}
	for _, topo := range tar.Processes {
		tarProcs[topo.Data[common.BKProcessNameField].(string)] = topo

		topo.Data[common.BKAppIDField] = bizID
		topo.Data[common.BKOwnerIDField] = opt.OwnerID
		curTopo := curProcs[topo.Data[common.BKProcessNameField].(string)]
		if curTopo != nil {
			procID, err := getInt64(curTopo.Data[common.BKProcessIDField])
			if nil != err {
				return fmt.Errorf("cur process has no bk_process_id field, %s", err.Error())
			}

			topo.Data[common.BKProcessIDField] = procID
			cond := getModifyCondition(topo.Data, []string{common.BKProcessIDField})
			// fmt.Printf("--- \033[36mupdate process by %+v, data: %+v\033[0m\n", condition, topo.Data)
			if !opt.dryrun {
				err = db.Table(common.BKTableNameBaseProcess).Update(ctx, cond, topo.Data)
				if nil != err {
					return fmt.Errorf("insert process data: %+v, error: %s", topo.Data, err.Error())
				}
			}

			// add missing module
			for _, moduleName := range topo.Modules {
				if inSlice(moduleName, curTopo.Modules) {
					continue
				}
				procMod := ProModule{}
				procMod.ModuleName = moduleName
				procMod.BizID = bizID
				procMod.ProcessID = procID
				procMod.OwnerID = opt.OwnerID
				fmt.Printf("--- \033[34minsert process module data: %+v\033[0m\n", procMod)
				if !opt.dryrun {
					err = db.Table(common.BKTableNameProcModule).Insert(ctx, &procMod)
					if nil != err {
						return fmt.Errorf("insert process module data: %+v, error: %s", topo.Data, err.Error())
					}
				}
			}
			// delete unused module map
			for _, curModule := range curTopo.Modules {
				if inSlice(curModule, topo.Modules) {
					continue
				}
				delCondition := map[string]interface{}{
					common.BKModuleNameField: curModule,
					common.BKProcessIDField:  procID,
				}
				fmt.Printf("--- \033[31mdelete process module by %+v\033[0m\n", delCondition)
				if !opt.dryrun {
					err = db.Table(common.BKTableNameProcModule).Delete(ctx, delCondition)
					if nil != err {
						return fmt.Errorf("delete process module by %+v, error: %v", delCondition, err)
					}
				}
			}

		} else {
			nid, err := db.NextSequence(ctx, common.BKTableNameBaseProcess)
			if nil != err {
				return fmt.Errorf("GetIncID for prcess faile, error: %s ", err.Error())
			}
			topo.Data[common.BKProcessIDField] = nid
			fmt.Printf("--- \033[34minsert process data: %+v\033[0m\n", topo.Data)
			if !opt.dryrun {
				err = db.Table(common.BKTableNameBaseProcess).Insert(ctx, topo.Data)
				if nil != err {
					return fmt.Errorf("insert process data: %+v, error: %s", topo.Data, err.Error())
				}
			}
			for _, moduleName := range topo.Modules {
				procMod := ProModule{}
				procMod.ModuleName = moduleName
				procMod.BizID = bizID
				procMod.ProcessID = nid
				procMod.OwnerID = opt.OwnerID
				fmt.Printf("--- \033[34minsert process module data: %+v\033[0m\n", topo.Data)
				if !opt.dryrun {
					err = db.Table(common.BKTableNameProcModule).Insert(ctx, &procMod)
					if nil != err {
						return fmt.Errorf("insert process module data: %+v, error: %s", topo.Data, err.Error())
					}
				}
			}
		}
	}

	// remove unused process
	for key, proc := range curProcs {
		if tarProcs[key] == nil {
			delCondition := map[string]interface{}{
				common.BKProcessIDField: proc.Data[common.BKProcessIDField],
			}
			fmt.Printf("--- \033[31mdelete process by %+v\033[0m\n", delCondition)
			if !opt.dryrun {
				err = db.Table(common.BKTableNameBaseProcess).Delete(ctx, delCondition)
				if nil != err {
					return fmt.Errorf("delete process by %+v, error: %s", delCondition, err.Error())
				}
			}
			fmt.Printf("--- \033[31mdelete process module by %+v\033[0m\n", delCondition)
			if !opt.dryrun {
				err = db.Table(common.BKTableNameProcModule).Delete(ctx, delCondition)
				if nil != err {
					return fmt.Errorf("delete process module by %+v, error: %s", delCondition, err.Error())
				}
			}
		}
	}

	return nil
}

func getModifyCondition(data map[string]interface{}, keys []string) map[string]interface{} {
	cond := map[string]interface{}{}
	for _, key := range keys {
		cond[key] = data[key]
	}
	return cond
}

var ignoreKeys = map[string]bool{
	"_id":                   true,
	"create_time":           true,
	common.BKInstParentStr:  true,
	"default":               true,
	common.BKAppIDField:     true,
	common.BKSetIDField:     true,
	common.BKProcessIDField: true,
	common.BKInstIDField:    true,
}

// compareSlice returns whether slice a,b exactly equal
func compareSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}

type importer struct {
	screate  map[string][]*Node
	supdate  map[string][]*Node
	sdelete  map[string][]map[string]interface{}
	bizID    uint64
	setID    uint64
	parentID uint64
	ownerID  string

	ctx context.Context
	db  dal.RDB
	opt *option
}

func newImporter(ctx context.Context, db dal.RDB, opt *option) *importer {
	return &importer{
		screate:  map[string][]*Node{},
		supdate:  map[string][]*Node{},
		sdelete:  map[string][]map[string]interface{}{},
		bizID:    0,
		setID:    0,
		parentID: 0,
		ownerID:  "",

		ctx: ctx,
		db:  db,
		opt: opt,
	}
}

func (ipt *importer) walk(includeRoot bool, node *Node) error {
	if node.mark != "" {
		return nil
	}
	if includeRoot {
		switch node.ObjID {
		case common.BKInnerObjIDApp:
			cond := getModifyCondition(node.Data, []string{common.BKAppNameField})
			app := map[string]interface{}{}
			err := ipt.db.Table(common.GetInstTableName(node.ObjID)).Find(cond).One(ipt.ctx, &app)
			if nil != err {
				return fmt.Errorf("get blueking business by %+v error: %s", cond, err.Error())
			}
			bizID, err := getInt64(app[common.BKAppIDField])
			if nil != err {
				return fmt.Errorf("get blueking bizID faile, data: %+v, error: %s", app, err.Error())
			}
			ownerID, ok := app[common.BKOwnerIDField].(string)
			if !ok {
				return fmt.Errorf("get blueking nk_suppplier_account faile, data: %+v, error: %v", app, err)
			}
			node.Data[common.BKAppIDField] = bizID
			node.Data[common.BKOwnerIDField] = ownerID
			ipt.bizID = bizID
			ipt.parentID = bizID
			ipt.ownerID = ownerID
			if !containsMap(app, node.Data) {
				node.mark = actionUpdate
			}
		case common.BKInnerObjIDSet:
			node.Data[common.BKOwnerIDField] = ipt.ownerID
			node.Data[common.BKAppIDField] = ipt.bizID
			node.Data[common.BKInstParentStr] = ipt.parentID
			node.Data[common.BKDefaultField] = 0
			cond := getModifyCondition(node.Data, []string{common.BKSetNameField, common.BKInstParentStr})
			set := map[string]interface{}{}
			err := ipt.db.Table(common.GetInstTableName(node.ObjID)).Find(cond).One(ipt.ctx, &set)
			if nil != err && !ipt.db.IsNotFoundError(err) {
				return fmt.Errorf("get set by %+v error: %s", cond, err.Error())
			}
			if ipt.db.IsNotFoundError(err) {
				node.mark = actionCreate
				nid, err := ipt.db.NextSequence(ipt.ctx, common.GetInstTableName(node.ObjID))
				if nil != err {
					return fmt.Errorf("GetIncID error: %s", err.Error())
				}
				node.Data[common.BKSetIDField] = nid
				ipt.parentID = nid
				ipt.setID = nid
			} else {
				if !containsMap(set, node.Data) {
					node.mark = actionUpdate
				}
				setID, err := getInt64(set[common.BKSetIDField])
				if nil != err {
					return fmt.Errorf("get setID faile, data: %+v, error: %s", set, err.Error())
				}
				node.Data[common.BKSetIDField] = setID
				ipt.parentID = setID
				ipt.setID = setID
			}
		case common.BKInnerObjIDModule:
			node.Data[common.BKOwnerIDField] = ipt.ownerID
			node.Data[common.BKAppIDField] = ipt.bizID
			node.Data[common.BKSetIDField] = ipt.setID
			node.Data[common.BKInstParentStr] = ipt.parentID
			node.Data[common.BKDefaultField] = 0
			cond := getModifyCondition(node.Data, []string{common.BKModuleNameField, common.BKInstParentStr})
			module := map[string]interface{}{}
			err := ipt.db.Table(common.GetInstTableName(node.ObjID)).Find(cond).One(ipt.ctx, &module)
			if nil != err && !ipt.db.IsNotFoundError(err) {
				return fmt.Errorf("get module by %+v error: %s", cond, err.Error())
			}
			if ipt.db.IsNotFoundError(err) {
				node.mark = actionCreate
				nid, err := ipt.db.NextSequence(ipt.ctx, common.GetInstTableName(node.ObjID))
				if nil != err {
					return fmt.Errorf("GetIncID error: %s", err.Error())
				}
				node.Data[common.BKModuleIDField] = nid
			} else {
				if !containsMap(module, node.Data) {
					node.mark = actionUpdate
				}
				moduleID, err := getInt64(module[common.BKModuleIDField])
				if nil != err {
					return fmt.Errorf("get moduleID faile, data: %+v, error: %s", module, err.Error())
				}
				node.Data[common.BKModuleIDField] = moduleID
			}
		default:
			node.Data[common.BKOwnerIDField] = ipt.ownerID
			node.Data[common.BKInstParentStr] = ipt.parentID
			cond := getModifyCondition(node.Data, []string{node.getInstNameField(), common.BKInstParentStr})
			cond[common.BKObjIDField] = node.ObjID
			inst := map[string]interface{}{}
			err := ipt.db.Table(common.GetInstTableName(node.ObjID)).Find(cond).One(ipt.ctx, &inst)
			if nil != err && !ipt.db.IsNotFoundError(err) {
				return fmt.Errorf("get inst by %+v error: %s", cond, err.Error())
			}
			if ipt.db.IsNotFoundError(err) {
				node.mark = actionCreate
				nid, err := ipt.db.NextSequence(ipt.ctx, common.GetInstTableName(node.ObjID))
				if nil != err {
					return fmt.Errorf("GetIncID error: %s", err.Error())
				}
				node.Data[common.GetInstIDField(node.ObjID)] = nid
				ipt.parentID = nid
			} else {
				if !containsMap(inst, node.Data) {
					node.mark = actionUpdate
				}
				instID, err := getInt64(inst[common.GetInstIDField(node.ObjID)])
				if nil != err {
					return fmt.Errorf("get instID faile, data: %+v, error: %s", inst, err.Error())
				}
				node.Data[common.GetInstIDField(node.ObjID)] = instID
				ipt.parentID = instID
			}
		}

		// fetch data that should delete
		if node.ObjID != common.BKInnerObjIDModule {
			childTableName := common.GetInstTableName(node.getChildObjID())
			instID, err := node.getInstID()
			if nil != err {
				return fmt.Errorf("get instID faile, data: %+v, error: %v", node, err)
			}

			childFilter := map[string]interface{}{
				common.BKInstParentStr: instID,
				node.getChildInstNameField(): map[string]interface{}{
					common.BKDBNIN: node.getChildInstNames(),
				},
			}
			switch node.getChildObjID() {
			case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
				childFilter[common.BKDefaultField] = map[string]interface{}{
					common.BKDBLTE: 0,
				}
			default:
				childFilter[common.BKObjIDField] = node.getChildObjID()
			}
			shouldDelete := make([]map[string]interface{}, 0)
			err = ipt.db.Table(childTableName).Find(childFilter).All(ipt.ctx, &shouldDelete)
			if nil != err {
				return fmt.Errorf("get child of %+v error: %s", childFilter, err.Error())
			}
			if len(shouldDelete) > 0 {
				// fmt.Printf("found %d should delete %s by %+v\n, parent %+v \n", len(shouldDelete), node.getChildObjID(), childCondition, node.Data)
				ipt.sdelete[node.getChildObjID()] = append(ipt.sdelete[node.getChildObjID()], shouldDelete...)
			}
		}

		if node.mark == actionCreate {
			ipt.screate[node.getChildObjID()] = append(ipt.screate[node.getChildObjID()], node)
			parentID := ipt.parentID
			bizID := ipt.bizID
			setID := ipt.setID
			err := ipt.walk(false, node)
			if nil != err {
				return nil
			}
			ipt.parentID = parentID
			ipt.bizID = bizID
			ipt.setID = setID
		}
		if node.mark == actionUpdate {
			ipt.supdate[node.getChildObjID()] = append(ipt.supdate[node.getChildObjID()], node)
		}
	}

	parentID := ipt.parentID
	bizID := ipt.bizID
	setID := ipt.setID
	for _, child := range node.Children {
		if err := ipt.walk(true, child); nil != err {
			return err
		}
		ipt.parentID = parentID
		ipt.bizID = bizID
		ipt.setID = setID
	}

	return nil
}

// getModelAttributes returns the model attributes
func getModelAttributes(ctx context.Context, db dal.RDB, opt *option, objIDs []string) (modelAttributes map[string][]metadata.Attribute, modelKeys map[string][]string, err error) {
	cond := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			"$in": objIDs,
		},
	}

	attributes := make([]metadata.Attribute, 0)
	err = db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &attributes)
	if nil != err {
		return nil, nil, fmt.Errorf("faile to getModelAttributes for %v, error: %s", objIDs, err.Error())
	}

	modelAttributes = map[string][]metadata.Attribute{}
	modelKeys = map[string][]string{}
	for _, att := range attributes {
		if att.IsRequired {
			modelKeys[att.ObjectID] = append(modelKeys[att.ObjectID], att.PropertyID)
		}
		modelAttributes[att.ObjectID] = append(modelAttributes[att.ObjectID], att)
	}

	for index := range modelKeys {
		sort.Strings(modelKeys[index])
	}

	return modelAttributes, modelKeys, nil
}

func inSlice(sub string, slice []string) bool {
	for _, s := range slice {
		if s == sub {
			return true
		}
	}
	return false
}

// compare tar to src, returns whether src contains sunb
func containsMap(src, sub map[string]interface{}) bool {
	for key := range sub {
		if !equalIgnoreLength(src[key], sub[key]) {
			return false
		}
	}
	return true
}

func equalIgnoreLength(src, tar interface{}) bool {
	if src == tar {
		return true
	}

	if srcInt, err := getInt64(src); err == nil {
		if tarInt, err := getInt64(tar); err == nil {
			if srcInt == tarInt {
				return true
			}
		}
	}
	return false
}
