package gojahost

import (
	"time"

	"github.com/dop251/goja"
	"github.com/gost-dom/browser/clock"
	"github.com/gost-dom/browser/dom"
)

type eventLoopWrapper struct {
	baseInstanceWrapper[*clock.Clock]
}

func newEventLoopWrapper(instance *GojaContext) eventLoopWrapper {
	return eventLoopWrapper{newBaseInstanceWrapper[*clock.Clock](instance)}
}

func (w eventLoopWrapper) initializeWindows(prototype *goja.Object, _ *goja.Runtime) {
	prototype.DefineDataProperty(
		"setInterval",
		w.ctx.vm.ToValue(w.setInterval),
		goja.FLAG_FALSE,
		goja.FLAG_TRUE,
		goja.FLAG_TRUE,
	)
	// prototype.Set("setInterval", w.ctx.vm.ToValue(w.setTimeout))
	prototype.DefineDataProperty(
		"setTimeout",
		w.ctx.vm.ToValue(w.setTimeout),
		goja.FLAG_FALSE,
		goja.FLAG_TRUE,
		goja.FLAG_TRUE,
	)
	prototype.DefineDataProperty(
		"clearTimeout",
		w.ctx.vm.ToValue(w.clearTimeout),
		goja.FLAG_FALSE,
		goja.FLAG_TRUE,
		goja.FLAG_TRUE,
	)
	prototype.DefineDataProperty(
		"clearInterval",
		w.ctx.vm.ToValue(w.clearInterval),
		goja.FLAG_FALSE,
		goja.FLAG_TRUE,
		goja.FLAG_TRUE,
	)
}

func (l eventLoopWrapper) setTimeout(c goja.FunctionCall) goja.Value {
	f, ok := goja.AssertFunction(c.Argument(0))
	if !ok {
		panic(l.ctx.vm.NewTypeError("setTimeout: Argument must be a function"))
	}
	delay := c.Argument(1).ToInteger()
	handle := l.ctx.clock.AddSafeTask(
		clock.Relative(time.Millisecond*time.Duration(delay)),
		func() {
			_, err := f(l.ctx.vm.GlobalObject())
			if err != nil {
				l.ctx.window.DispatchEvent(dom.NewErrorEvent(err))
			}
		},
	)
	return l.vm().ToValue(handle)
}

func (l eventLoopWrapper) clearTimeout(c goja.FunctionCall) goja.Value {
	id := c.Argument(0).ToInteger()
	l.ctx.clock.Cancel(clock.TaskHandle(id))
	return nil
}

func (l eventLoopWrapper) setInterval(c goja.FunctionCall) goja.Value {
	f, ok := goja.AssertFunction(c.Argument(0))
	if !ok {
		panic(l.ctx.vm.NewTypeError("setTimeout: Argument must be a function"))
	}
	delay := c.Argument(1).ToInteger()
	handle := l.ctx.clock.AddRepeat(
		time.Millisecond*time.Duration(delay),
		func() {
			_, err := f(l.ctx.vm.GlobalObject())
			if err != nil {
				l.ctx.window.DispatchEvent(dom.NewErrorEvent(err))
			}
		},
	)
	return l.vm().ToValue(handle)
}

func (l eventLoopWrapper) clearInterval(c goja.FunctionCall) goja.Value {
	id := c.Argument(0).ToInteger()
	l.ctx.clock.Cancel(clock.TaskHandle(id))
	return nil
}
