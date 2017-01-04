package store

import (
	"github.com/gopherjs/vecty/examples/issue60/actions"
	"github.com/gopherjs/vecty/examples/issue60/dispatcher"
	"github.com/gopherjs/vecty/storeutil"
)

var (
	t         string
	c         int
	Listeners = storeutil.NewListenerRegistry()
)

func init() {
	dispatcher.Register(onAction)
}

func Count() int {
	return c
}

func Type() string {
	return t
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.Increment:
		c++

	case *actions.Change:
		t = a.Type

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}
