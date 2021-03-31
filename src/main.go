package main

import (
	"fmt"
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/DrSmithFr/go-webassembly/src/wolfenstein"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"syscall/js"
)

var DOM *browser.DOM
var cvs *browser.Canvas2d
var gs *wolfenstein.GameState

var width float64
var height float64

func main() {
	// loading DOM to memory
	DOM = browser.LoadDOM()

	// setting up everything
	bindEvents(*DOM)

	// create canvas
	cvs, _ = browser.NewCanvas2d(false)
	cvs.Create(
		js.Global().Get("innerWidth").Int(),
		js.Global().Get("innerHeight").Int(),
	)

	// create gameState
	gs, _ = wolfenstein.NewGameState(cvs.Width(), cvs.Height())

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	// starting rendering
	cvs.Start(120, Render)

	// allow daemon style process
	emptyChanToKeepAppRunning := make(chan bool)
	<-emptyChanToKeepAppRunning
}

func bindEvents(DOM browser.DOM) {
	// let's handle windows resize
	var resizeEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resizeEvent(DOM, args[0])
		return nil
	})

	DOM.Window.Call("addEventListener", "resize", resizeEventHandler)

	// let's handle key down
	var keyboardEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyboardEvent(DOM, args[0])
		return nil
	})

	DOM.Document.Call("addEventListener", "keydown", keyboardEventHandler)

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

func keyboardEvent(DOM browser.DOM, event js.Value) {
	code := event.Get("code").String()

	switch code {
	case "ArrowUp", "KeyW":
		gs.MoveUp()
	case "ArrowDown", "KeyS":
		gs.MoveDown()
	case "ArrowRight", "KeyD":
		gs.MoveRight()
	case "ArrowLeft", "KeyA":
		gs.MoveLeft()
	}

	go DOM.Log(fmt.Sprintf("key press:%s", code))
}

func clickEvent(DOM browser.DOM, event js.Value) {
	mouseX := event.Get("clientX").Int()
	mouseY := event.Get("clientY").Int()

	go DOM.Log(fmt.Sprintf("mouseEvent x:%d y:%d", mouseX, mouseY))
}

func Render(gc *draw2dimg.GraphicContext) bool {
	// render default color
	gc.SetFillColor(color.RGBA{0x18, 0x18, 0x18, 0xff})
	gc.Clear()

	renderLevel(gc)
	renderPlayer(gc)

	return true
}

func renderLevel(gc *draw2dimg.GraphicContext) {
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.BeginPath()

	level := gs.GetLevel()
	blockSize := gs.GetBlockSize()

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if level[x*8+y] == 0 {
				// avoid useless rendering
				continue
			}

			draw2dkit.Rectangle(
				gc,
				float64(x*blockSize+1),
				float64(y*blockSize+1),
				float64(x*blockSize+blockSize-1),
				float64(y*blockSize+blockSize-1),
			)
			gc.FillStroke()
		}
	}
}

func renderPlayer(gc *draw2dimg.GraphicContext) {
	// draw player on screen
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.BeginPath()

	playerX, playerY, playerDeltaX, playerDeltaY := gs.GetPlayerPosition()
	draw2dkit.Circle(gc, playerX, playerY, 5)
	gc.FillStroke()

	// player direction
	// draw player on screen
	gc.BeginPath()
	gc.MoveTo(playerX, playerY)
	gc.LineTo(playerX + playerDeltaX * 5, playerY + playerDeltaY * 5)
	gc.Close()
	gc.FillStroke()
}
