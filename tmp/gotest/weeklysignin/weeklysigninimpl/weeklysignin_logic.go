package weeklysigninimpl

import (
	"gitlab.babeltime.com/c_server/majorserver"
	"gitlab.babeltime.com/c_server/majorserver/module/public"
	"gitlab.babeltime.com/c_server/majorserver/module/weeklysignin"
	"gitlab.babeltime.com/c_server/proto/go/pb"
	"time"
)

func getWeeklySignInfo(aContext *majorserver.TContext, aUid uint64, aRsp *pb.WeeklySignInfoRsp) {
	obj := getWeeklySignInObj(aContext, aUid)

	//检查请求时间，以及是否重置db
	checkTime(aContext, obj)

	//返回数据库信息
	signInfo := pb.NewWeeklySignInfo()
	signInfo.SetUid(obj.dbStruct.GetUid())
	signInfo.SetFiveConsRewardFlag(obj.dbStruct.GetGetFiveConsReward())
	signInfo.SetSignData(obj.dbStruct.GetSignData().ToArray())
	aRsp.SetInfo(signInfo)
}

// 用户签到
func doSignIn(aContext *majorserver.TContext, aUid uint64, aRsp *pb.WeeklySignInRsp) {
	logger := aContext.GetLogger()
	obj := getWeeklySignInObj(aContext, aUid)

	checkTime(aContext, obj)

	//周几
	reqTime := aContext.Unix()
	dayIdx := getWeekDayIdx(reqTime)

	statusCode := weeklysignin.SUCCESS
	dbTime := obj.dbStruct.GetSignData().At(dayIdx - 1)
	if dbTime == weeklysignin.INVALID_TIME {
		//今日未签到，进行签到
		err := obj.dbStruct.GetSignData().Set(dayIdx-1, reqTime)
		if err != nil {
			logger.ThrowFake("doSignIn SetSignTime err:%S", err)
		}
		getReward(aContext, aUid, dayIdx, false)
	} else {
		//今日已签到
		logger.Debug("repeat sign in.today has signed in:%d", dbTime)
		statusCode = weeklysignin.REPEAT_SIGNIN
	}

	aRsp.SetStatusCode(int32(statusCode))
}

// 补签
func makeUpSignIn(aContext *majorserver.TContext, aUid uint64, aReqIdx uint32, aRsp *pb.WeeklySignInRsp) {
	logger := aContext.GetLogger()
	obj := getWeeklySignInObj(aContext, aUid)

	checkTime(aContext, obj)

	nowTime := aContext.Unix()
	nowIdx := getWeekDayIdx(nowTime)

	statusCode := weeklysignin.SUCCESS

	//补签日期>=今日
	if int(aReqIdx) >= int(nowIdx) {
		statusCode = weeklysignin.INVALID_MAKEUP
		logger.ThrowFake("req makeup dayId is invalid. ReqDay:%d,Today:%d", aReqIdx, nowIdx)
	}

	dbTime := obj.dbStruct.GetSignData().At(int(aReqIdx) - 1)
	if dbTime == weeklysignin.INVALID_TIME {
		//请求补签的日期没有签到，补签
		err := obj.dbStruct.GetSignData().Set(int(aReqIdx)-1, nowTime)
		if err != nil {
			logger.ThrowFake("doSignIn SetSignTime err:%S", err)
		}
		getReward(aContext, aUid, int(aReqIdx), false)
	} else {
		//补签重复，无效
		statusCode = weeklysignin.REPEAT_MAKEUP
		logger.Debug("req makeup day has already sign in. Sign time:%d, Req Time:%d", dbTime, nowTime)
	}
	aRsp.SetStatusCode(int32(statusCode))
}

// 领取连签奖励
func getFiveConsReward(aContext *majorserver.TContext, aUid uint64, aRsp *pb.WeeklySignInRsp) {
	logger := aContext.GetLogger()
	obj := getWeeklySignInObj(aContext, aUid)

	checkTime(aContext, obj)

	statusCode := weeklysignin.SUCCESS
	//没有在服务器判断是否达到连签5天的条件
	if int(obj.dbStruct.GetGetFiveConsReward()) == weeklysignin.HAS_RECEIVED_FIVECONS {
		logger.Debug("FiveCons Reward has received.")
		statusCode = weeklysignin.REPEAT_FIVECONS
	} else {
		//成功领取
		getReward(aContext, aUid, 1, true)
		obj.dbStruct.SetGetFiveConsReward(int32(weeklysignin.HAS_RECEIVED_FIVECONS))
	}

	aRsp.SetStatusCode(int32(statusCode))
}

// 检查请求时间，以及是否重置db
func checkTime(aContext *majorserver.TContext, aObj *TWeeklySignInObj) {
	logger := aContext.GetLogger()
	dbWeekStartTime := aObj.dbStruct.GetWeekStartTime()
	dbWeekEndTime := dbWeekStartTime + weeklysignin.WEEK_SECONDS
	reqTime := aContext.Unix()

	//进入新周：请求时间>记录的结束时间
	if reqTime > dbWeekEndTime {
		// 重置数据
		// 只根据请求时间就直接刷新吗
		aObj.RefreshWeek(aContext)
		logger.Debug("req time is new week. db refresh")
	}

	//请求时间太旧
	if reqTime < dbWeekStartTime {
		logger.ThrowFake("req time is invalid. req time:%d, dbtime:%d", reqTime, dbWeekStartTime)
	}
}

//func modifySignData(aContext *majorserver.TContext, aObj *TWeeklySignInObj, aIdx int, aNowTime int64) {
//	logger := aContext.GetLogger()
//
//	arrSignData := aObj.dbStruct.GetSignData().ToArray()
//	arrSignData[aIdx] = aNowTime
//
//	aObj.dbStruct.GetSignData().Set()
//}

func getWeeklySignInObj(aContext *majorserver.TContext, aUid uint64) *TWeeklySignInObj {
	logger := aContext.GetLogger()
	obj, ok := aContext.FindWeeklySignInData(aUid).(*TWeeklySignInObj)
	if !ok {
		logger.CtxThrowFake(aContext.GetNativeCtx(), "FindWeeklySignInData error uid:%v", aUid)
	}
	return obj
}

// 传入时间戳，返回周几 （1到7 即 周一到周日）
func getWeekDayIdx(aUnixTime int64) int {
	dayIdx := time.Unix(aUnixTime, 0).Weekday()
	if dayIdx == 0 {
		dayIdx = 7
	}
	return int(dayIdx)
}

// 玩家获得奖励
func getReward(aContext *majorserver.TContext, aUid uint64, aDayId int, aFiveConsFlag bool) {
	logger := aContext.GetLogger()
	conf := aContext.GetConfWeeklySignInByIndex(aDayId - 1)
	rewardId := conf.GetDailyReward()
	if aFiveConsFlag == true {
		rewardId = conf.GetFiveConsecutiveReward()
	}
	logger.Debug("Weeklysignin getReward. reward id:%d", rewardId)
	public.DoRewardWithId(aContext, aUid, rewardId, 1)
}
