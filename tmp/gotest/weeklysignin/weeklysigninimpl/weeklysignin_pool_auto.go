//DO NOT MODIFY GENERATED CODE !!!
//DO NOT MODIFY GENERATED CODE !!!
//DO NOT MODIFY GENERATED CODE !!!
//This is a auto-generated file. Generate by auto_code/auto_code_rpc.go
package weeklysigninimpl

import(
	"sync"
	
	"gitlab.babeltime.com/serverfw/logicfw"
)

var g_dummy_TRpcWeeklySignInModule TRpcWeeklySignInModule

func init() {
	pool := &sync.Pool {
		New:func() interface{} {
			return new(TRpcWeeklySignInModule)
		},
	}
	
	fn := func(d interface{}) {
		o, ok := d.(*TRpcWeeklySignInModule)
		if ok{
			*o = g_dummy_TRpcWeeklySignInModule
			pool.Put(d)
		}
	}
	logicfw.RegisterPool(pool, fn,  "weeklysignin")
}
