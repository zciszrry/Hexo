package test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.babeltime.com/c_server/majorserver"
	"gitlab.babeltime.com/c_server/proto/go/pb"
	"testing"
)

func TestWeeklySignInSuite(t *testing.T) {
	suite.Run(t, new(WeeklySignInSuite))
}

type WeeklySignInSuite struct {
	suite.Suite
	context *majorserver.TContext
	userCli *TClient
}

func (this *WeeklySignInSuite) SetupSuite() {
	this.context = majorserver.MakeROContext()
}

func (this *WeeklySignInSuite) TearDownSuite() {
	this.context.ReleaseContext()
}

func (this *WeeklySignInSuite) SetupTest() {
	userCli := NewLoginClient()
	this.NotNil(userCli)

	userCli.uid = 790010000000198
	this.userCli = userCli

	fmt.Println("In SetUp. Cli Id : ", this.userCli.uid)
}

func GetWeeklySignInfo(userCli *TClient) *pb.WeeklySignInfoRsp {
	req := pb.NewTEmptyReq()
	rsp := pb.NewWeeklySignInfoRsp()
	err := userCli.SendReq("weeklysignin.GetWeeklySignInfo", req, rsp)
	if err != nil {
		fmt.Println("SendReq err : ", err)
		//return nil
	}

	fmt.Println("WeelySignInfo")
	fmt.Println(rsp)
	fmt.Println()

	return rsp
}

func DoSignIn(userCli *TClient) *pb.WeeklySignInRsp {
	req := pb.NewTEmptyReq()
	rsp := pb.NewWeeklySignInRsp()
	err := userCli.SendReq("weeklysignin.DoSignIn", req, rsp)
	if err != nil {
		fmt.Println("SendReq err : ", err)
		//return nil
	}

	fmt.Println("DoSignIn")
	fmt.Println("状态码：", rsp.StatusCode)
	fmt.Println()

	return rsp
}

func MakeUpSignIn(userCli *TClient) *pb.WeeklySignInRsp {
	req := pb.NewMakeUpSignInReq()
	req.SetDayId(uint32(1))
	rsp := pb.NewWeeklySignInRsp()
	err := userCli.SendReq("weeklysignin.DoSignIn", req, rsp)
	if err != nil {
		fmt.Println("SendReq err : ", err)
		//return nil
	}

	fmt.Println("MakeUpSignIn")
	fmt.Println("状态码：", rsp.StatusCode)
	fmt.Println()

	return rsp
}

func GetFiveConsReward(userCli *TClient) *pb.WeeklySignInRsp {
	req := pb.NewTEmptyReq()
	rsp := pb.NewWeeklySignInRsp()
	err := userCli.SendReq("weeklysignin.DoSignIn", req, rsp)
	if err != nil {
		fmt.Println("SendReq err : ", err)
		//return nil
	}

	return rsp
}

func (this *WeeklySignInSuite) TestBuilding() {

}
