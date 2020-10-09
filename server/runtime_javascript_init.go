package server

import (
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"strings"
)

const INIT_MODULE_FN_NAME = "InitModule"

type RuntimeJavascriptMatchCallbacks map[string]*jsMatchHandlers

type jsMatchHandlers struct {
	initFn        goja.Callable
	joinAttemptFn goja.Callable
	joinFn        goja.Callable
	leaveFn       goja.Callable
	loopFn        goja.Callable
	terminateFn   goja.Callable
}

type RuntimeJavascriptCallbacks struct {
	Rpc    map[string]goja.Callable
	Before map[string]goja.Callable
	After  map[string]goja.Callable
	Matchmaker       goja.Callable
	TournamentEnd    goja.Callable
	TournamentReset  goja.Callable
	LeaderboardReset goja.Callable
}

type RuntimeJavascriptInitModule struct {
	Logger             *zap.Logger
	Callbacks          *RuntimeJavascriptCallbacks
	MatchCallbacks     *RuntimeJavascriptMatchCallbacks
	announceCallbackFn func(RuntimeExecutionMode, string)
}

func NewRuntimeJavascriptInitModule(logger *zap.Logger, announceCallbackFn func(RuntimeExecutionMode, string)) *RuntimeJavascriptInitModule {
	callbacks := &RuntimeJavascriptCallbacks{
		Rpc:    make(map[string]goja.Callable),
		Before: make(map[string]goja.Callable),
		After:  make(map[string]goja.Callable),
	}

	matchCallbacks := &RuntimeJavascriptMatchCallbacks{}

	return &RuntimeJavascriptInitModule{
		Logger:             logger,
		announceCallbackFn: announceCallbackFn,
		Callbacks:          callbacks,
		MatchCallbacks:     matchCallbacks,
	}
}

func (im *RuntimeJavascriptInitModule) mappings(r *goja.Runtime) map[string]func(goja.FunctionCall) goja.Value {
	return map[string]func(goja.FunctionCall) goja.Value{
		"registerRpc":               im.registerRpc(r),
		"registerReqBefore":         im.registerReqBefore(r),
		"registerReqAfter":          im.registerReqAfter(r),
		"registerRTBefore":          im.registerRTBefore(r),
		"registerRTAfter":           im.registerRTAfter(r),
		"registerMatchmakerMatched": im.registerMatchmakerMatched(r),
		"registerTournamentEnd":     im.registerTournamentEnd(r),
		"registerTournamentReset":   im.registerTournamentReset(r),
		"registerLeaderboardReset":  im.registerLeaderboardReset(r),
		"registerMatch":             im.registerMatch(r),
	}
}

func (im *RuntimeJavascriptInitModule) Constructor(r *goja.Runtime) func(goja.ConstructorCall) *goja.Object {
	return func(call goja.ConstructorCall) *goja.Object {
		for key, fn := range im.mappings(r) {
			call.This.Set(key, fn)
		}

		return nil
	}
}

func (im *RuntimeJavascriptInitModule) registerRpc(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fName := f.Argument(0)
		if goja.IsNull(fName) || goja.IsUndefined(fName) {
			panic(r.NewTypeError("expects a non empty string"))
		}
		key, ok := fName.Export().(string)
		if !ok {
			panic(r.NewTypeError("expects a non empty string"))
		}
		if key == "" {
			panic(r.NewTypeError("expects a non empty string"))
		}

		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		lKey := strings.ToLower(key)
		im.registerCallbackFn(RuntimeExecutionModeRPC, lKey, fn)
		im.announceCallbackFn(RuntimeExecutionModeRPC, lKey)

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerReqBefore(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fName := f.Argument(0)
		if goja.IsNull(fName) || goja.IsUndefined(fName) {
			panic(r.NewTypeError("expects a non empty string"))
		}
		key, ok := fName.Export().(string)
		if !ok {
			panic(r.NewTypeError("expects a non empty string"))
		}
		if key == "" {
			panic(r.NewTypeError("expects a non empty string"))
		}

		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		lKey := strings.ToLower(API_PREFIX + key)
		im.registerCallbackFn(RuntimeExecutionModeBefore, lKey, fn)
		im.announceCallbackFn(RuntimeExecutionModeBefore, lKey)

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerReqAfter(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fName := f.Argument(0)
		if goja.IsNull(fName) || goja.IsUndefined(fName) {
			panic(r.NewTypeError("expects a non empty string"))
		}
		key, ok := fName.Export().(string)
		if !ok {
			panic(r.NewTypeError("expects a non empty string"))
		}
		if key == "" {
			panic(r.NewTypeError("expects a non empty string"))
		}

		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		lKey := strings.ToLower(API_PREFIX + key)
		im.registerCallbackFn(RuntimeExecutionModeAfter, lKey, fn)
		im.announceCallbackFn(RuntimeExecutionModeAfter, lKey)

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerRTBefore(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fName := f.Argument(0)
		if goja.IsNull(fName) || goja.IsUndefined(fName) {
			panic(r.NewTypeError("expects a non empty string"))
		}
		key, ok := fName.Export().(string)
		if !ok {
			panic(r.NewTypeError("expects a non empty string"))
		}
		if key == "" {
			panic(r.NewTypeError("expects a non empty string"))
		}

		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		lKey := strings.ToLower(RTAPI_PREFIX + key)
		im.registerCallbackFn(RuntimeExecutionModeBefore, lKey, fn)
		im.announceCallbackFn(RuntimeExecutionModeBefore, lKey)

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerRTAfter(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fName := f.Argument(0)
		if goja.IsNull(fName) || goja.IsUndefined(fName) {
			panic(r.NewTypeError("expects a non empty string"))
		}
		key, ok := fName.Export().(string)
		if !ok {
			panic(r.NewTypeError("expects a non empty string"))
		}
		if key == "" {
			panic(r.NewTypeError("expects a non empty string"))
		}

		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		lKey := strings.ToLower(RTAPI_PREFIX + key)
		im.registerCallbackFn(RuntimeExecutionModeAfter, lKey, fn)
		im.announceCallbackFn(RuntimeExecutionModeAfter, lKey)

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerMatchmakerMatched(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		im.registerCallbackFn(RuntimeExecutionModeMatchmaker, "", fn)
		im.announceCallbackFn(RuntimeExecutionModeMatchmaker, "")

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerTournamentEnd(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		im.registerCallbackFn(RuntimeExecutionModeTournamentEnd, "", fn)
		im.announceCallbackFn(RuntimeExecutionModeTournamentEnd, "")

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerTournamentReset(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		im.registerCallbackFn(RuntimeExecutionModeTournamentReset, "", fn)
		im.announceCallbackFn(RuntimeExecutionModeTournamentReset, "")

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerLeaderboardReset(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		fn, ok := goja.AssertFunction(f.Argument(1))
		if !ok {
			panic(r.NewTypeError("expects a function"))
		}

		im.registerCallbackFn(RuntimeExecutionModeLeaderboardReset, "", fn)
		im.announceCallbackFn(RuntimeExecutionModeLeaderboardReset, "")

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerMatch(r *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(f goja.FunctionCall) goja.Value {
		name := getJsString(r, f.Argument(0))

		funcObj := f.Argument(1)
		if goja.IsNull(funcObj) || goja.IsUndefined(funcObj) {
			panic(r.NewTypeError("expects an object"))
		}

		funcMap, ok := funcObj.Export().(map[string]interface{})
		if !ok {
			panic(r.NewTypeError("expects an object"))
		}

		functions := &jsMatchHandlers{}

		fnValue, ok := funcMap["match_init"]
		if !ok {
			panic(r.NewTypeError("match_init not found"))
		}
		fn, ok := goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_init value not a valid function"))
		}
		functions.initFn = fn

		fnValue, ok = funcMap["match_join_attempt"]
		if !ok {
			panic(r.NewTypeError("match_join_attempt not found"))
		}
		fn, ok = goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_join_attempt value not a valid function"))
		}
		functions.joinAttemptFn = fn

		fnValue, ok = funcMap["match_join"]
		if !ok {
			panic(r.NewTypeError("match_join not found"))
		}
		fn, ok = goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_join value not a valid function"))
		}
		functions.joinFn = fn

		fnValue, ok = funcMap["match_leave"]
		if !ok {
			panic(r.NewTypeError("match_leave not found"))
		}
		fn, ok = goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_leave value not a valid function"))
		}
		functions.leaveFn = fn

		fnValue, ok = funcMap["match_loop"]
		if !ok {
			panic(r.NewTypeError("match_loop not found"))
		}
		fn, ok = goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_loop value not a valid function"))
		}
		functions.loopFn = fn

		fnValue, ok = funcMap["match_terminate"]
		if !ok {
			panic(r.NewTypeError("match_terminate not found"))
		}
		fn, ok = goja.AssertFunction(r.ToValue(fnValue))
		if !ok {
			panic(r.NewTypeError("match_terminate value not a valid function"))
		}
		functions.terminateFn = fn

		(*im.MatchCallbacks)[name] = functions

		return goja.Undefined()
	}
}

func (im *RuntimeJavascriptInitModule) registerCallbackFn(mode RuntimeExecutionMode, key string, fn goja.Callable) {
	switch mode {
	case RuntimeExecutionModeRPC:
		im.Callbacks.Rpc[key] = fn
	case RuntimeExecutionModeBefore:
		im.Callbacks.Before[key] = fn
	case RuntimeExecutionModeAfter:
		im.Callbacks.After[key] = fn
	case RuntimeExecutionModeMatchmaker:
		im.Callbacks.Matchmaker = fn
	case RuntimeExecutionModeTournamentEnd:
		im.Callbacks.TournamentEnd = fn
	case RuntimeExecutionModeTournamentReset:
		im.Callbacks.TournamentReset = fn
	case RuntimeExecutionModeLeaderboardReset:
		im.Callbacks.LeaderboardReset = fn
	}
}


