package vecty

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

// EventListener is markup that specifies a callback function to be invoked when
// the named DOM event is fired.
type EventListener struct {
	Name               string
	Listener           func(*Event)
	callPreventDefault bool
	wrapper            func(jsEvent *js.Object)
}

// PreventDefault prevents the default behavior of the event from occuring.
//
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/preventDefault.
func (l *EventListener) PreventDefault() *EventListener {
	l.callPreventDefault = true
	return l
}

func (l *EventListener) Apply(h *HTML) {
	h.EventListeners = append(h.EventListeners, l)
}

// Event represents a DOM event.
type Event struct {
	Target *js.Object
}

type MarkupOrComponentOrHTML interface{}

func apply(m MarkupOrComponentOrHTML, h *HTML) {
	if m == nil {
		return
	}
	switch m := m.(type) {
	case Markup:
		m.Apply(h)
	case Component:
		h.Children = append(h.Children, m)
	case *HTML:
		h.Children = append(h.Children, m)
	default:
		panic(fmt.Sprintf("vecty: invalid type %T does not match MarkupOrComponent interface", m))
	}
}

type Markup interface {
	Apply(h *HTML)
}

type markupFunc func(h *HTML)

func (m markupFunc) Apply(h *HTML) { m(h) }

func Style(key, value string) Markup {
	return markupFunc(func(h *HTML) {
		if h.Styles == nil {
			h.Styles = make(map[string]string)
		}
		h.Styles[key] = value
	})
}

func Property(key string, value interface{}) Markup {
	return markupFunc(func(h *HTML) {
		if h.Properties == nil {
			h.Properties = make(map[string]interface{})
		}
		h.Properties[key] = value
	})
}

func Data(key, value string) Markup {
	return markupFunc(func(h *HTML) {
		h.Dataset[key] = value
	})
}

// ClassMap is markup that specifies classes to be applied to an element if
// their boolean value are true.
type ClassMap map[string]bool

// Apply implements the Markup interface.
func (m ClassMap) Apply(h *HTML) {
	var classes []string
	for name, active := range m {
		if active {
			classes = append(classes, name)
		}
	}
	Property("className", strings.Join(classes, " ")).Apply(h)
}

type List []MarkupOrComponentOrHTML

func (l List) Apply(h *HTML) {
	for _, m := range l {
		apply(m, h)
	}
}

// If returns nil if cond is false, otherwise it returns the given markup.
func If(cond bool, markup ...MarkupOrComponentOrHTML) MarkupOrComponentOrHTML {
	if cond {
		return List(markup)
	}
	return nil
}
