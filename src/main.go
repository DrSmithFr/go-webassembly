package main

import (
	"fmt"
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"syscall/js"
)

var cvs *browser.Canvas2d

var width float64
var height float64

var gs = gameState{laserSize: 35, directionX: 13.7, directionY: -13.7, laserX: 40, laserY: 40}
type gameState struct{ laserX, laserY, directionX, directionY, laserSize float64 }

func main() {
	// loading DOM to memory
	DOM := browser.LoadDOM()

	// setting up everything
	bindEvents(DOM)

	cvs, _ = browser.NewCanvas2d(false)
	cvs.Create(
		js.Global().Get("innerWidth").Int(),
		js.Global().Get("innerHeight").Int(),
	)

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(120, Render)

	// allow daemon style process
	emptyChanToKeepAppRunning := make(chan bool)
	<-emptyChanToKeepAppRunning
}

func bindEvents(DOM browser.DOM) {
	// let's handle that mouse pointer down
	var resizeEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resizeEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "resize", resizeEventHandler)

	// let's handle that mouse pointer down
	var mouseEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "pointerdown", mouseEventHandler)
}

func resizeEvent(DOM browser.DOM, event js.Value) {
	windowsWidth := js.Global().Get("innerWidth").Int()
	windowsHeight := js.Global().Get("innerHeight").Int()

	cvs.SetSize(windowsWidth, windowsHeight)

	width = float64(windowsWidth)
	height = float64(windowsHeight)

	go DOM.Log(fmt.Sprintf("resizeEvent x:%d y:%d", windowsWidth, windowsHeight))
}

func clickEvent(DOM browser.DOM, event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go DOM.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}

func Render(gc *draw2dimg.GraphicContext) bool {

	if gs.laserX+gs.directionX > width-gs.laserSize || gs.laserX+gs.directionX < gs.laserSize {
		gs.directionX = -gs.directionX
	}
	if gs.laserY+gs.directionY > height-gs.laserSize || gs.laserY+gs.directionY < gs.laserSize {
		gs.directionY = -gs.directionY
	}

	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	// move red laser
	gs.laserX += gs.directionX
	gs.laserY += gs.directionY

	// draws red 🔴 laser
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0xff, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0xff, 0xff})

	gc.BeginPath()
	// gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
	draw2dkit.Circle(gc, gs.laserX, gs.laserY, gs.laserSize)
	gc.FillStroke()
	gc.Close()

	return true
}
