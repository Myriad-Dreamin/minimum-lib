package hook

import (
	"github.com/Myriad-Dreamin/minimum-lib/controller"
)

type Hook struct {
	funcs []HookFunc
}

type HookFunc func(c controller.MContext) bool

func (hook *Hook) Use(hookFunc HookFunc) {
	hook.funcs = append(hook.funcs, hookFunc)
}

func (hook *Hook) Consume(c controller.MContext) bool {
	for _, f := range hook.funcs {
		if !f(c) {
			return false
		}
	}
	return true
}
