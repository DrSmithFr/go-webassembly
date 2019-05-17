package main

import (
	"./browser"
	"fmt"
	"syscall/js"
)

func main() {
	emptyChannel := make(chan bool)

	// loading DOM to memory
	DOM := browser.LoadDOM()

	// setting up everything
	setup(DOM)

	// attempt to receive from empty channel
	// allow daemon style process
	<-emptyChannel
}

func setup(DOM browser.DOM) {
	// let's handle that mouse pointer down
	var mouseEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "pointerdown", mouseEventHandler)
}

func clickEvent(DOM browser.DOM, event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go DOM.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}