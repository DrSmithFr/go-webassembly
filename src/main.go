package main

import (
	"./browser"
	"fmt"
	"syscall/js"
)

// creating dom as global variable for simplicity
var dom browser.DOM

func main() {
	emptyChannel := make(chan bool)

	// loading DOM to memory
	dom = browser.LoadDOM()

	// setting up everything
	setup()

	// attempt to receive from empty channel
	// allow daemon style process
	<-emptyChannel
}

func setup() {
	// javascript rendering func
	var renderer js.Func

	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		update()
		return nil
	})

	// using browser animation frame
	dom.Window.Call("requestAnimationFrame", renderer)

	// let's handle that mouse pointer down
	var mouseEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickEvent(args[0])
		return nil
	})

	dom.Window.Call("addEventListener", "pointerdown", mouseEventHandler)
}

func update() {
	fmt.Printf("updating frame")
}

func clickEvent(event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go dom.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}
