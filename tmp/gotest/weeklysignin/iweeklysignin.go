package weeklysignin

import (
	"gitlab.babeltime.com/c_server/proto/go/pb"
)

// inmodule
type IInWeeklySignInModule interface{}

// indata
type IInWeeklySIgnInData interface {
	//finddata
	FindWeeklySignInData(aUid uint64) IWeeklySignInData
}

type IRpcWeeklySignInModule interface {
	// 获取签到信息
	GetWeeklySignInfo(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInfoRsp)
	//签到
	DoSignIn(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInRsp)
	// 补签
	MakeUpSignIn(aReq *pb.MakeUpSignInReq, aRsp *pb.WeeklySignInRsp)
	// 领取连签奖励
	GetFiveConsReward(aReq *pb.TEmptyReq, aRsp *pb.WeeklySignInRsp)
}

type IWeeklySignInData interface {
	GetId() uint64
	GetKey() string

	GmGetWeekStartTime() int64
}
