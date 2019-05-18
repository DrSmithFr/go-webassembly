package mapper

import (
	"image"
	"syscall/js"
)

func ImageToImageData(img *image.RGBA) js.Value  {
	// ClampedArray need by ImageData constructor
	array := js.TypedArrayOf(img.Pix)
	data := js.Global().Get("Uint8ClampedArray").New(array)

	// creating the javascript image data
	size := img.Bounds().Size()
	imageData := js.Global().Get("ImageData").New(data, size.X, size.Y)

	return imageData
}
