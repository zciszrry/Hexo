package weeklysigninimpl

import (
	"gitlab.babeltime.com/c_server/majorserver/module/handler"
	"gitlab.babeltime.com/c_server/proto/go/pb"
)

type TRpcWeeklySignInModule struct {
	handler.TContextHandler
}

// 获取签到信息
func (this *TRpcWeeklySignInModule) GetWeeklySignInfo(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInfoRsp) {
	context := this.GetContext()
	logger := context.GetLogger()
	uid := context.GetSession().GetUid()
	err := context.LockUserData(uid)
	if err != nil {
		logger.ThrowFake("LockData uid:%d err:%s", uid, err)
	}

	getWeeklySignInfo(context, uid, aRsp)
}

// 签到
func (this *TRpcWeeklySignInModule) DoSignIn(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInRsp) {
	context := this.GetContext()
	logger := context.GetLogger()
	uid := context.GetSession().GetUid()
	err := context.LockUserData(uid)
	if err != nil {
		logger.ThrowFake("LockData uid:%d err:%s", uid, err)
	}

	doSignIn(context, uid, aRsp)
}

// 补签
func (this *TRpcWeeklySignInModule) MakeUpSignIn(aReq *pb.MakeUpSignInReq, aRsp *pb.WeeklySignInRsp) {
	context := this.GetContext()
	logger := context.GetLogger()
	uid := context.GetSession().GetUid()
	err := context.LockUserData(uid)
	if err != nil {
		logger.ThrowFake("LockData uid:%d err:%s", uid, err)
	}

	idx := aReq.GetDayId()
	makeUpSignIn(context, uid, idx, aRsp)
}

// 领取连签奖励
func (this *TRpcWeeklySignInModule) GetFiveConsReward(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInRsp) {
	context := this.GetContext()
	logger := context.GetLogger()
	uid := context.GetSession().GetUid()
	err := context.LockUserData(uid)
	if err != nil {
		logger.ThrowFake("LockData uid:%d err:%s", uid, err)
	}

	getFiveConsReward(context, uid, aRsp)
}
