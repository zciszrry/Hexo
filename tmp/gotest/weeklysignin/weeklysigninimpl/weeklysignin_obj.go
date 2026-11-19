package weeklysigninimpl

import (
	"database/sql"
	"fmt"
	"gitlab.babeltime.com/c_server/majorserver/module/weeklysignin"
	"time"

	"gitlab.babeltime.com/c_server/common/db/tables"
	"gitlab.babeltime.com/c_server/majorserver"
)

// dataobj
type TWeeklySignInObj struct {
	TDataObjInner
	dbStruct *tables.TWeeklySignIn
	context  *majorserver.TContext
}

// finddataobj
func findWeeklySignInData(aContext *majorserver.TContext, aUid uint64) weeklysignin.IWeeklySignInData {
	logger := aContext.GetLogger()
	dbData := tables.NewTWeeklySignIn()
	dbData.SetUid(aUid)

	err := aContext.StoreGet(dbData)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.ThrowInter("findWeeklySignInData DbGet err:%s", err.Error())
		}
		initWeeklySignInData(aContext, dbData, aUid)
		dbData.ResetDataChanged()
	}
	obj := NewTWeeklySignInObj()
	obj.dbStruct = dbData
	obj.context = aContext
	return obj
}

func initWeeklySignInData(aContext *majorserver.TContext, aData *tables.TWeeklySignIn, aUid uint64) {
	logger := aContext.GetLogger()
	aData.SetUid(aUid)
	for i := 0; i < 7; i++ {
		aData.GetSignData().Append(weeklysignin.INVALID_TIME)
	}
	aData.SetWeekStartTime(GetWeekStartTime(aContext.Unix()))
	err := aContext.StoreInsert(aData)
	if err != nil {
		err := aContext.StoreGet(aData)
		if err != nil {
			logger.ThrowInter("initWeeklySignInData context.StoreInsert err:%s", err.Error())
		}
	}
}

func (this *TWeeklySignInObj) GetId() uint64        { return this.dbStruct.GetUid() }
func (this *TWeeklySignInObj) GetKey() string       { return fmt.Sprintf("%v", this.dbStruct.GetUid()) }
func (this *TWeeklySignInObj) doWhenRelease()       {}
func (this *TWeeklySignInObj) doAfterUpdate() error { return nil }

// 刷新数据表
func (this *TWeeklySignInObj) RefreshWeek(aContext *majorserver.TContext) {
	logger := aContext.GetLogger()
	dataChanged := this.dbStruct.IsDataChanged()

	startTime := GetWeekStartTime(aContext.Unix())
	this.dbStruct.SetWeekStartTime(startTime)
	this.dbStruct.SetGetFiveConsReward(0)
	//this.dbStruct.GetSignData().Clear()
	for i := 0; i < 7; i++ {
		this.dbStruct.GetSignData().Set(i, weeklysignin.INVALID_TIME)
	}

	if this.dbStruct.GetSignData().Len() != 7 {
		logger.ThrowFake("the db signData have err. arr length:%d", this.dbStruct.GetSignData().Len())
	}

	this.dbStruct.ResetDataChanged()
	if !dataChanged {
		this.dbStruct.ResetDataChanged()
	}
}

// 获得时间戳所在周的开始时间，周一0点0分0秒
func GetWeekStartTime(aTime int64) int64 {
	t := time.Unix(aTime, 0)

	day := t.Weekday()
	if day == 0 {
		day = 7
	}
	startTime := time.Date(t.Year(), t.Month(), t.Day()-(int(day)-1), 0, 0, 0, 0, t.Location())
	return startTime.Unix()
}

func (this *TWeeklySignInObj) GmGetWeekStartTime() int64 {
	return this.dbStruct.GetWeekStartTime()
}
