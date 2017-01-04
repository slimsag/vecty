package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/examples/issue60/actions"
	"github.com/gopherjs/vecty/examples/issue60/dispatcher"
)

type PageView struct {
	vecty.Core
}

func (p *PageView) Render() *vecty.HTML {
	return elem.Body(
		&TestView{},
		elem.Button(
			event.Click(func(ev *vecty.Event) {
				dispatcher.Dispatch(&actions.Increment{})
			}),
			vecty.Text("Increment"),
		),
		elem.Button(
			event.Click(func(ev *vecty.Event) {
				dispatcher.Dispatch(&actions.Change{"div"})
			}),
			vecty.Text("Div"),
		),
		elem.Button(
			event.Click(func(ev *vecty.Event) {
				dispatcher.Dispatch(&actions.Change{"p"})
			}),
			vecty.Text("Paragraph"),
		),
	)
}
