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

	// creating rendered objects
	canvas := makeCanvas(DOM)
	ball := makeBall(canvas)

	// binding rendered object to actual DOM
	DOM.Body.Call("appendChild", canvas)

	// attempt to receive from empty channel
	// allow daemon style process
	<-emptyChannel
}

func setup(DOM browser.DOM) {
	// javascript rendering func
	var renderer js.Func

	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		update(DOM)
		return nil
	})

	// using browser animation frame
	DOM.Window.Call("requestAnimationFrame", renderer)

	// let's handle that mouse pointer down
	var mouseEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "pointerdown", mouseEventHandler)
}

func update(DOM browser.DOM) {
	fmt.Printf("updating frame")
}

func clickEvent(DOM browser.DOM, event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go DOM.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}

func makeCanvas(DOM browser.DOM) js.Value {
	fmt.Printf("%v", DOM)

	canvas := DOM.Document.Call("createElement", "canvas")

	canvas.Set("height", DOM.Size.Height)
	canvas.Set("width", DOM.Size.Width)

	return canvas
}

func makeBall(canvas js.Value) js.Value {
	ball := canvas.Call("getContext", "2d")
	ball.Set("fillStyle", "red")
	return ball
}