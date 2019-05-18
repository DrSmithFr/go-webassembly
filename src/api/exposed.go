package api

import (
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/DrSmithFr/go-webassembly/src/mapper"
	"github.com/DrSmithFr/go-webassembly/src/rendering"
	"syscall/js"
)

func GetJavascriptObject(DOM browser.DOM) js.Value  {
	wrapper := js.Global().Get("Object").New()

	// @param {int} width
	// @param {int} height
	// @return {ImageData}
	wrapper.Set("goImage", generateImage())

	return js.ValueOf(wrapper)
}

func generateImage() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0]
		height := args[1]

		img := rendering.MakeImage(width.Int(), height.Int())

		return mapper.ImageToImageData(img)
	})
}
