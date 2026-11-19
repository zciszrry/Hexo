package weeklysigninimpl

import (
	"gitlab.babeltime.com/c_server/majorserver/module/handler"
	"gitlab.babeltime.com/c_server/majorserver/module/weeklysignin"
)

type TInWeeklySignInModule struct {
	handler.TContextHandler
}

func (this *TInWeeklySignInModule) FindWeeklySignInData(aUid uint64) weeklysignin.IWeeklySignInData {
	return findWeeklySignInData(this.GetContext(), aUid)
}
