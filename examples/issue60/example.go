package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/examples/issue60/components"
	"github.com/gopherjs/vecty/examples/issue60/store"
)

func main() {
	p := &components.PageView{}
	store.Listeners.Add(p, func() {
		vecty.Rerender(p)
	})
	vecty.RenderBody(p)
}
