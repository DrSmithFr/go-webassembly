package main

import (
	"fmt"
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/DrSmithFr/go-webassembly/src/wolfenstein"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
	"syscall/js"
)

var DOM *browser.DOM
var cvs *browser.Canvas2d
var gs *wolfenstein.GameState

type move struct {
	up    bool
	down  bool
	left  bool
	right bool
}

var keyboard = move{false, false, false, false}

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
	var keydownEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keydownEvent(DOM, args[0])
		return nil
	})

	var keyupEventHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyupEvent(DOM, args[0])
		return nil
	})

	DOM.Document.Call("addEventListener", "keydown", keydownEventHandler)
	DOM.Document.Call("addEventListener", "keyup", keyupEventHandler)

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

func keydownEvent(DOM browser.DOM, event js.Value) {
	code := event.Get("code").String()

	switch code {
	case "ArrowUp", "KeyW":
		keyboard.up = true
	case "ArrowDown", "KeyS":
		keyboard.down = true
	case "ArrowRight", "KeyD":
		keyboard.right = true
	case "ArrowLeft", "KeyA":
		keyboard.left = true
	}

	//go DOM.Log(fmt.Sprintf("key down:%s", code))
}

func keyupEvent(DOM browser.DOM, event js.Value) {
	code := event.Get("code").String()

	switch code {
	case "ArrowUp", "KeyW":
		keyboard.up = false
	case "ArrowDown", "KeyS":
		keyboard.down = false
	case "ArrowRight", "KeyD":
		keyboard.right = false
	case "ArrowLeft", "KeyA":
		keyboard.left = false
	}

	//go DOM.Log(fmt.Sprintf("key up:%s", code))
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
	renderRayCasting(gc)
	handleMove()

	return true
}

func renderRayCasting(gc *draw2dimg.GraphicContext) {
	gc.BeginPath()

	var rayX, rayY, rayAngle float64
	var rayTargetX, rayTargetY float64
	var mapX, mapY, mapIndex int
	var dof int

	level := gs.GetLevel()
	blocSize := gs.GetBlockSize()
	mapSizeX, mapSizeY := gs.GetMapSize()
	playerX, playerY, _, _ := gs.GetPlayerPosition()

	rayAngle = gs.GetPlayerAngle()

	for rayN := 0; rayN < 1; rayN++ {
		// check Horizontal
		dof = 0
		aTan := -1 / math.Tan(rayAngle)

		if rayAngle > math.Pi {
			// looking up
			rayY = math.Trunc(playerY/float64(blocSize))*float64(blocSize) - 1
			rayX = (playerY-rayY)*aTan + playerX

			rayTargetY = - float64(blocSize)
			rayTargetX = - rayTargetY * aTan

			gc.FillStroke()
		} else if rayAngle < math.Pi {
			// looking down (ok)
			rayY = math.Trunc(playerY/float64(blocSize))*float64(blocSize) + float64(blocSize)
			rayX = (playerY-rayY)*aTan + playerX

			rayTargetY = float64(blocSize)
			rayTargetX = - rayTargetY * aTan
		}

		if rayAngle == 0 || rayAngle == math.Pi {
			rayX = playerX
			rayY = playerY
			dof = 8
		}

		for ; dof < 8; {
			mapX = int(math.Trunc(rayX / float64(blocSize)))
			mapY = int(math.Trunc((rayY) / float64(blocSize)))

			mapIndex = mapY*mapSizeX + mapX

			// hit wall
			if mapIndex > 0 && mapIndex < mapSizeX*mapSizeY && level[mapIndex] == 1 {
				dof = 8
			} else {
				rayX += rayTargetX
				rayY += rayTargetY
				dof++
			}
		}

		gc.SetFillColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
		gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0xff, 0xff})

		gc.BeginPath()
		gc.MoveTo(playerX, playerY)
		gc.LineTo(rayX, rayY)
		gc.Close()
		gc.FillStroke()

		// check Vertical
		dof = 0
		nTan := -math.Tan(rayAngle)
		P2 := math.Pi / 2
		P3 := 3*P2

		if rayAngle > P2 && rayAngle < P3 {
			// looking left
			rayX = math.Trunc(playerX/float64(blocSize))*float64(blocSize) - 1
			rayY = (playerX-rayX)*nTan + playerY

			rayTargetX = - float64(blocSize)
			rayTargetY = - rayTargetX * nTan

			gc.FillStroke()
		} else if rayAngle < P2 || rayAngle > P3 {
			// looking right
			rayX = math.Trunc(playerX/float64(blocSize))*float64(blocSize) + float64(blocSize)
			rayY = (playerX-rayX)*nTan + playerY

			rayTargetX = float64(blocSize)
			rayTargetY = - rayTargetX * nTan
		}

		if rayAngle == 0 || rayAngle == math.Pi {
			rayX = playerX
			rayY = playerY
			dof = 8
		}

		for ; dof < 8; {
			mapX = int(math.Trunc(rayX / float64(blocSize)))
			mapY = int(math.Trunc((rayY) / float64(blocSize)))

			mapIndex = mapY*mapSizeX + mapX

			// hit wall
			if mapIndex > 0 && mapIndex < mapSizeX*mapSizeY && level[mapIndex] == 1 {
				dof = 8
			} else {
				rayX += rayTargetX
				rayY += rayTargetY
				dof++
			}
		}

		gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
		gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

		gc.BeginPath()
		gc.MoveTo(playerX, playerY)
		gc.LineTo(rayX, rayY)
		gc.Close()
		gc.FillStroke()
	}
}

func handleMove() {
	if keyboard.up {
		gs.MoveUp()
	} else if keyboard.down {
		gs.MoveDown()
	}

	if keyboard.right {
		gs.MoveRight()
	} else if keyboard.left {
		gs.MoveLeft()
	}
}

func renderLevel(gc *draw2dimg.GraphicContext) {
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.BeginPath()

	level := gs.GetLevel()
	blockSize := gs.GetBlockSize()
	mapX, mapY := gs.GetMapSize()

	for y := 0; y < mapY; y++ {
		for x := 0; x < mapX; x++ {
			if level[x+y*mapY] == 0 {
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

	// draw player on screen
	playerX, playerY, playerDeltaX, playerDeltaY := gs.GetPlayerPosition()
	draw2dkit.Circle(gc, playerX, playerY, 5)
	gc.FillStroke()

	// draw player direction
	gc.BeginPath()
	gc.MoveTo(playerX, playerY)
	gc.LineTo(playerX+playerDeltaX*5, playerY+playerDeltaY*5)
	gc.Close()
	gc.FillStroke()
}
