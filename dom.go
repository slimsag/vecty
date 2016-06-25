package vecty

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

// Core implements the Context method of the Component interface, and is the
// core/central struct which all Component implementations should embed.
type Core struct {
	prevRender *HTML
	didMount   bool
}

// Context implements the Component interface.
func (c *Core) Context() *Core { return c }

func doMount(c ComponentOrHTML) {
	if _, isHTML := c.(*HTML); isHTML {
		return
	}
	ctx := c.(Component).Context()
	if ctx.didMount {
		return
	}
	ctx.didMount = true

	if m, ok := c.(Mounter); ok {
		m.Mount()
	}
}

// Component represents a single visual component within an application. To
// define a new component simply implement the Render method and embed the Core
// struct:
//
// 	type MyComponent struct {
// 		vecty.Core
// 		... additional component fields (state or properties) ...
// 	}
//
// 	func (c *MyComponent) Render() *vecty.HTML {
// 		... rendering ...
// 	}
//
type Component interface {
	// Render is responsible for building HTML which represents the component.
	Render() *HTML

	// Context returns the components context, which is used internally by
	// Vecty in order to store the previous component render for diffing.
	Context() *Core
}

// ComponentOrHTML represents one of:
//
//  Component
//  *HTML
//
// If the underlying value is not one of these types, the code handling the
// value is expected to panic.
type ComponentOrHTML interface{}

// Mounter is an optional interface that Component's can implement in order to
// receive component mount events.
type Mounter interface {
	// Mount is called directly after the component is mounted. That is,
	// directly after the element has been added to the DOM.
	Mount()
}

// Unmounter is an optional interface that Component's can implement in order
// to receive component unmount events.
type Unmounter interface {
	// Unmount is called directly before the component is unmounted. That is,
	// prior to the element being removed from the DOM.
	Unmount()
}

// Restorer is an optional interface that Component's can implement in order to
// restore state during component reconciliation and also to short-circuit
// the reconciliation of a Component's body.
type Restorer interface {
	// Restore is called when the component should restore itself against a
	// previous instance of a component. The previous component may be nil or
	// of a different type than this Restorer itself, thus a type assertion
	// should be used.
	//
	// If skip = true is returned, restoration of this component's body is
	// skipped. That is, the component is not rerendered. If the component can
	// prove when Restore is called that the HTML rendered by Component.Render
	// would not change, true should be returned.
	Restore(prev Component) (skip bool)
}

// HTML represents some form of HTML: an element with a specific tag, or some
// literal text (a TextNode).
type HTML struct {
	Tag, Text       string
	Styles, Dataset map[string]string
	Properties      map[string]interface{}
	EventListeners  []*EventListener
	Children        []ComponentOrHTML
	Node            *js.Object
}

func (h *HTML) restoreHTML(prev *HTML) {
	h.Node = prev.Node

	// Text modifications.
	if h.Text != prev.Text {
		h.Node.Set("nodeValue", h.Text)
	}

	// Properties
	for name, value := range h.Properties {
		oldValue := prev.Properties[name]
		if value != oldValue || name == "value" || name == "checked" {
			h.Node.Set(name, value)
		}
	}
	for name := range prev.Properties {
		if _, ok := h.Properties[name]; !ok {
			h.Node.Set(name, nil)
		}
	}

	// Styles
	style := h.Node.Get("style")
	for name, value := range h.Styles {
		oldValue := prev.Styles[name]
		if value != oldValue {
			style.Call("setProperty", name, value)
		}
	}
	for name := range prev.Styles {
		if _, ok := h.Styles[name]; !ok {
			style.Call("removeProperty", name)
		}
	}

	for _, l := range prev.EventListeners {
		h.Node.Call("removeEventListener", l.Name, l.wrapper)
	}
	for _, l := range h.EventListeners {
		h.Node.Call("addEventListener", l.Name, l.wrapper)
	}

	// TODO better list element reuse
	for i, nextChild := range h.Children {
		nextChildRender := doRender(nextChild)
		if i >= len(prev.Children) {
			if doRestore(nil, nextChild, nil, nextChildRender) {
				continue
			}
			h.Node.Call("appendChild", nextChildRender.Node)
			doMount(nextChild)
			continue
		}
		prevChild := prev.Children[i]
		prevChildRender, ok := prevChild.(*HTML)
		if !ok {
			// ??? what if prev is HTML
			prevChildRender = prevChild.(Component).Context().prevRender
		}
		if doRestore(prevChild, nextChild, prevChildRender, nextChildRender) {
			continue
		}
		replaceNode(nextChildRender.Node, prevChildRender.Node)
		doMount(nextChild)
	}
	for i := len(h.Children); i < len(prev.Children); i++ {
		prevChild := prev.Children[i]
		prevChildRender, ok := prevChild.(*HTML)
		if !ok {
			// ??? what if prev is HTML
			prevChildRender = prevChild.(Component).Context().prevRender
		}
		removeNode(prevChildRender.Node)
		if unmounter, ok := prevChild.(Unmounter); ok {
			unmounter.Unmount()
		}
	}
}

// Restore implements the Restorer interface.
func (h *HTML) Restore(old ComponentOrHTML) {
	for _, l := range h.EventListeners {
		l.wrapper = func(jsEvent *js.Object) {
			if l.callPreventDefault {
				jsEvent.Call("preventDefault")
			}
			l.Listener(&Event{Target: jsEvent.Get("target")})
		}
	}

	if prev, ok := old.(*HTML); ok && prev != nil {
		h.restoreHTML(prev)
		return
	}

	if h.Tag != "" && h.Text != "" {
		panic("vecty: only one of HTML.Tag or HTML.Text may be set")
	}
	if h.Tag != "" {
		h.Node = js.Global.Get("document").Call("createElement", h.Tag)
	} else if h.Text != "" {
		h.Node = js.Global.Get("document").Call("createTextNode", h.Text)
	}
	for name, value := range h.Properties {
		h.Node.Set(name, value)
	}
	dataset := h.Node.Get("dataset")
	for name, value := range h.Dataset {
		dataset.Set(name, value)
	}
	style := h.Node.Get("style")
	for name, value := range h.Styles {
		style.Call("setProperty", name, value)
	}
	for _, l := range h.EventListeners {
		h.Node.Call("addEventListener", l.Name, l.wrapper)
	}
	for _, nextChild := range h.Children {
		nextChildRender, isHTML := nextChild.(*HTML)
		if !isHTML {
			nextChildComp := nextChild.(Component)
			nextChildRender = nextChildComp.Render()
			nextChildComp.Context().prevRender = nextChildRender
		}

		if doRestore(nil, nextChild, nil, nextChildRender) {
			continue
		}
		h.Node.Call("appendChild", nextChildRender.Node)
		doMount(nextChild)
	}
}

// Tag returns an HTML element with the given tag name. Generally, this
// function is not used directly but rather the elem subpackage (which is type
// safe) is used instead.
func Tag(tag string, m ...MarkupOrComponentOrHTML) *HTML {
	h := &HTML{
		Tag: tag,
	}
	for _, m := range m {
		apply(m, h)
	}
	return h
}

// Text returns a TextNode with the given literal text. Because the returned
// HTML represents a TextNode, the text does not have to be escaped (arbitrary
// user input fed into this function will always be safely rendered).
func Text(text string, m ...MarkupOrComponentOrHTML) *HTML {
	h := &HTML{
		Text: text,
	}
	for _, m := range m {
		apply(m, h)
	}
	return h
}

// Rerender causes the body of the given component (i.e. the HTML returned by
// the Component's Render method) to be re-rendered and subsequently restored.
func Rerender(c Component) {
	prevRender := c.Context().prevRender
	nextRender := doRender(c)
	var prevComponent Component = nil // TODO
	if doRestore(prevComponent, c, prevRender, nextRender) {
		return
	}
	if prevRender != nil {
		replaceNode(nextRender.Node, prevRender.Node)
		doMount(c)
	}
}

func doRender(c ComponentOrHTML) *HTML {
	if h, isHTML := c.(*HTML); isHTML {
		return h
	}
	comp := c.(Component)
	r := comp.Render()
	comp.Context().prevRender = r
	return r
}

func doRestore(prev, next ComponentOrHTML, prevRender, nextRender *HTML) (skip bool) {
	if r, ok := next.(Restorer); ok {
		var p Component
		if prev != nil {
			p = prev.(Component)
		}
		if r.Restore(p) {
			return true
		}
	}
	nextRender.Restore(prevRender)
	return false
}

// RenderBody renders the given component as the document body. The given
// Component's Render method must return a "body" element.
func RenderBody(body Component) {
	nextRender := doRender(body)
	if nextRender.Tag != "body" {
		panic(fmt.Sprintf("vecty: RenderBody expected Component.Render to return a body tag, found %q", nextRender.Tag))
	}
	doRestore(nil, body, nil, nextRender)
	// TODO: doRestore skip == true here probably implies a user code bug
	doc := js.Global.Get("document")
	if doc.Get("readyState").String() == "loading" {
		doc.Call("addEventListener", "DOMContentLoaded", func() { // avoid duplicate body
			doc.Set("body", nextRender.Node)
		})
		return
	}
	doc.Set("body", nextRender.Node)
	doMount(body)
}

// SetTitle sets the title of the document.
func SetTitle(title string) {
	js.Global.Get("document").Set("title", title)
}

// AddStylesheet adds an external stylesheet to the document.
func AddStylesheet(url string) {
	link := js.Global.Get("document").Call("createElement", "link")
	link.Set("rel", "stylesheet")
	link.Set("href", url)
	js.Global.Get("document").Get("head").Call("appendChild", link)
}
