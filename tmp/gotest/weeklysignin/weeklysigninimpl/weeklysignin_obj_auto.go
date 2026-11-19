// DO NOT MODIFY !!!
// DO NOT MODIFY !!!
// DO NOT MODIFY !!!
// Important things should be repeated three times.
// This is a auto-generated file. Generate by auto_code/auto_code_data.go
package weeklysigninimpl

import (
	"gitlab.babeltime.com/serverfw/logicfw"
	
	"gitlab.babeltime.com/c_server/majorserver"
	"sync"
	"gitlab.babeltime.com/serverfw/nads"

	
	"fmt"
	
)

var gTWeeklySignInObjDummy TWeeklySignInObj
var gTWeeklySignInObjLock logicfw.TLockId

var poolTWeeklySignInObj = sync.Pool{
	New: func() interface{} {
		return new(TWeeklySignInObj)
	},
}

func NewTWeeklySignInObj() *TWeeklySignInObj {
	obj := poolTWeeklySignInObj.Get().(*TWeeklySignInObj)
	*obj = gTWeeklySignInObjDummy
	return obj
}



func (this *TWeeklySignInObj) CheckUpdate() error {
	
	if this.dbStruct != nil && this.dbStruct.IsDataChanged() {
		if this.dbStruct.HasLock() {
			lockId := this.context.GetLockIdByData(this.dbStruct)
			if !this.context.HasLocked(lockId) {
				err := fmt.Errorf("not lock:%v when Get data:%v", lockId, this.dbStruct)
				this.context.GetLogger().CtxFatal(this.context.GetNativeCtx(), err.Error())
				return err
			}
		}
	}
	
	
	return nil
}

func (this *TWeeklySignInObj) Update() error {
	this.context.GetLogger().CtxTrace(this.context.GetNativeCtx(), "Update TWeeklySignInObj[%d]", this.GetId())
	
	if this.dbStruct != nil && this.dbStruct.IsDataChanged() {
		this.context.GetLogger().CtxTrace(this.context.GetNativeCtx(), "Update dataChanged = %v", this.dbStruct.IsDataChanged())
		err := this.context.StoreUpdate(this.dbStruct)
		if err != nil {
			return err
		}
		this.dbStruct.ResetDataChanged()
		this.dbStruct.SetUpdateToStore()
	}else {
		this.context.GetLogger().CtxTrace(this.context.GetNativeCtx(), "Update not dataChanged")
	}
	
	

	this.safeRunHook(func() {
		this.doAfterUpdate()
	})

	return nil
}

func (this *TWeeklySignInObj) GetAllTableOp() map[string]nads.ITableOperator {
	opMap := make(map[string]nads.ITableOperator)

	
	if this.dbStruct != nil {
		name := this.dbStruct.GetTableName()
		opMap[name] = this.dbStruct
	}
	
	

	return opMap
}

func (this *TWeeklySignInObj) GetMaxDataToken() uint32{
	
	maxToken := uint32(0)
	
	
	
	if this.dbStruct != nil && this.dbStruct.IsDataChanged() {
		if maxToken < this.dbStruct.GetCurDataToken() {
			maxToken = this.dbStruct.GetCurDataToken()
		}
	}
	
	
	return maxToken
}

func (this *TWeeklySignInObj) Release() {
	if this.isRelease() {
		this.context.GetLogger().CtxWarning(this.context.GetNativeCtx(), "Duplicate Release TWeeklySignInObj")
		return
	}
	this.setRelease()
	
	this.context.GetLogger().CtxTrace(this.context.GetNativeCtx(), "Release TWeeklySignInObj[%d]", this.GetId())
	

	
	if this.dbStruct != nil {
		this.dbStruct.Put()
	}
	

	

	

	this.safeRunHook(func() {
		this.doWhenRelease()
	})
	
	datakey := fmt.Sprintf("%s#%v", "weeklysignin.IWeeklySignInData", this.GetId())
	this.context.RemoveData(datakey)
	*this = gTWeeklySignInObjDummy
	poolTWeeklySignInObj.Put(this)
}


func (this *TWeeklySignInObj) SetContext(context *majorserver.TContext) {
	this.context = context
}

func (this *TWeeklySignInObj) GetCmcSize() int {
	
	return this.dbStruct.GetCmcSize()
}

// 填充可预料的最复杂数据，用于估算开销
func (this *TWeeklySignInObj) FillMaxData() error {
	 //return this.fillMaxData()
	return nil
}

// 防止提供给业务层的hook发生panic，如doAfterUpdate
func (this *TWeeklySignInObj) safeRunHook(hook func()) {
	defer func() {
		if r := recover(); r != nil {
			this.context.GetLogger().Fatal("TWeeklySignInObj hook:%v reover with:%v", hook, r)
		}
	}()
	hook()
}

func (this *TWeeklySignInObj) Put() {
	if this.isRelease() {
		this.context.GetLogger().Warning("Duplicate Release TWeeklySignInObj")
		return
	}
	this.setRelease()

	
	if this.dbStruct != nil {
		this.dbStruct.Put()
	}
	

	

	

	poolTWeeklySignInObj.Put(this)
}

func (this *TWeeklySignInObj) MarkReadOnly() {
	
	if this.dbStruct != nil {
		this.dbStruct.MarkReadOnly()
	}
	

	

	
	this.context.GetLogger().Trace("TWeeklySignInObj[%d] mark permission to read only", this.GetId())
	
}

func (this *TWeeklySignInObj) MarkAccessDenied() {
	
	if this.dbStruct != nil {
		this.dbStruct.MarkAccessDenied()
	}
	

	
}
