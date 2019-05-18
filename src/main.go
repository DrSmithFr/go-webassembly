package main

import (
	"fmt"
	"github.com/DrSmithFr/go-webassembly/src/api"
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"syscall/js"
)

func main() {
	// loading DOM to memory
	DOM := browser.LoadDOM()

	// setting up everything
	bindEvents(DOM)
	bindAnimationFrame(DOM)

	// Exposing our API
	wrapper := api.GetJavascriptObject(DOM)

	// give javascript control over the Go WASM API
	callbackJavascript(DOM, wrapper)

	// allow daemon style process
	emptyChanToKeepAppRunning := make(chan bool)
	<-emptyChanToKeepAppRunning
}

func bindEvents(DOM browser.DOM)  {
	// let's handle that mouse pointer down
	var mouseEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "pointerdown", mouseEventHandler)
}

func bindAnimationFrame(DOM browser.DOM)  {
	// javascript rendering func
	var renderer js.Func

	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		update(DOM)
		DOM.Window.Call("requestAnimationFrame", renderer)
		return nil
	})

	// using browser animation frame
	DOM.Window.Call("requestAnimationFrame", renderer)
}

func callbackJavascript(DOM browser.DOM, ApiWrapper js.Value)  {
	callback := DOM.Document.Get("onWasmLoad")

	if callback.Type() == js.TypeFunction {
		callback.Invoke(ApiWrapper)
	} else {
		panic("document.onWasmLoad() is undefined in current DOM")
	}
}

func update(DOM browser.DOM) {
	fmt.Printf("updating frame")
}

func clickEvent(DOM browser.DOM, event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go DOM.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}
