package api

import (
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/DrSmithFr/go-webassembly/src/mapper"
	"github.com/DrSmithFr/go-webassembly/src/rendering"
	"syscall/js"
)

func GetJavascriptObject(DOM browser.DOM) js.Value  {
	wrapper := js.Global().Get("Object").New()

	generateImage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		img := rendering.MakeImage(int(DOM.Size.Width), int(DOM.Size.Height))
		return mapper.ImageToImageData(img)
	})

	// @param {int} width
	// @param {int} height
	// @return {ImageData}
	wrapper.Set("generate", generateImage)

	return js.ValueOf(wrapper)
}
