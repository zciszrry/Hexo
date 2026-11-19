package weeklysignin

// 一周的总秒数：6*24*3600+23*3600+59*60+59
const WEEK_SECONDS int64 = 604799

// 无效时间戳，表示未签到
const INVALID_TIME int64 = -1

// 状态码，用于Rsp
const (
	SUCCESS         int = 1 //成功签到or补签or领取连签奖励
	REPEAT_SIGNIN   int = 2 //重复地签到
	REPEAT_MAKEUP   int = 3 //重复地补签
	INVALID_MAKEUP  int = 4 //无法补签
	REPEAT_FIVECONS int = 5 //重复地领取连签奖励
)

// db.getFiveConsReward						comment:"是否领取过连签奖励;0未领取；1已经领取过"
const (
	NO_RECEIVE_FIVECONS   int = 0
	HAS_RECEIVED_FIVECONS int = 1
)
