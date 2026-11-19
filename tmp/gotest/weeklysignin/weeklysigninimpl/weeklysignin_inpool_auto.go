//DO NOT MODIFY GENERATED CODE !!!
//DO NOT MODIFY GENERATED CODE !!!
//DO NOT MODIFY GENERATED CODE !!!
//This is a auto-generated file. Generate by auto_code/auto_code_data.go
package weeklysigninimpl

import(
	"sync"
	
	"gitlab.babeltime.com/serverfw/logicfw"
)

var g_dummy_TInWeeklySignInModule TInWeeklySignInModule
var g_pool_TInWeeklySignInModule sync.Pool

func init() {
	g_pool_TInWeeklySignInModule = sync.Pool {
		New:func() interface{} {
			return new(TInWeeklySignInModule)
		},
	}

	logicfw.RegisterPool(&g_pool_TInWeeklySignInModule, nil, "weeklysignin_in")
}

func (this *TInWeeklySignInModule) PutInPool() {
	*this = g_dummy_TInWeeklySignInModule
	g_pool_TInWeeklySignInModule.Put(this)
}
