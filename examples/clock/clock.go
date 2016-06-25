package main

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Clock struct {
	vecty.Core
	time        string
	unmountChan chan bool
}

func (c *Clock) Mount() {
	go func() {
		for {
			select {
			case <-c.unmountChan:
				return
			case <-time.After(1 * time.Second):
				c.time = time.Now().UTC().String()
				vecty.Rerender(c)
			}
		}
	}()
}

func (c *Clock) Unmount() {
	c.unmountChan <- true
}

func (c *Clock) Render() *vecty.HTML {
	return elem.Div(
		vecty.Text(fmt.Sprintf("Time: %s", c.time)),
	)
}

func NewClock() *Clock {
	return &Clock{
		time:        time.Now().UTC().String(),
		unmountChan: make(chan bool, 1),
	}
}

type Page struct {
	vecty.Core
}

func (p *Page) Render() *vecty.HTML {
	return elem.Body(NewClock())
}

func main() {
	vecty.RenderBody(&Page{})
}
