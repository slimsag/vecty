package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/examples/issue60/store"
)

type TestView struct {
	vecty.Core
}

func (p *TestView) Render() *vecty.HTML {
	if store.Type() == "p" {
		return elem.Paragraph(vecty.Text(fmt.Sprintf("P: Count = %d", store.Count())))
	} else {
		return elem.Div(vecty.Text(fmt.Sprintf("DIV: Count = %d", store.Count())))
	}
}
