package main

import (
	"fmt"
	"github.com/DrSmithFr/go-webassembly/src/browser"
	"github.com/DrSmithFr/go-webassembly/src/wolfenstein"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
	"runtime"
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

	DOM.Log(fmt.Sprintf("number of thread: %d", runtime.NumCPU()))

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

	if false {
		renderRayCasting(gc)
	}else {
		gs.RenderRay(gc)
	}

	handleMove()

	return true
}

func renderRayCasting(gc *draw2dimg.GraphicContext) {
	gc.BeginPath()

	var rayX, rayY, rayAngle float64
	var rayTargetX, rayTargetY float64
	var mapX, mapY, mapIndex int
	var dof int
	var distT float64

	level := gs.GetLevel()
	blocSize := gs.GetBlockSize()
	mapSize := gs.GetMapSize()
	playerX, playerY, _, _ := gs.GetPlayerPosition()

	oneRadian := 0.0174533
	rayAngle = gs.GetPlayerAngle()

	if rayAngle < 0 {
		rayAngle += 2 * math.Pi
	} else if rayAngle > 2*math.Pi {
		rayAngle -= 2 * math.Pi
	}

	for rayN := 0; rayN < 1; rayN++ {
		// check Horizontal
		dof = 0
		aTan := -1 / math.Tan(rayAngle)

		distH := 1000000.0
		hx := playerX
		hy := playerY

		if rayAngle > math.Pi {
			// looking up
			rayY = math.Trunc(playerY/float64(blocSize))*float64(blocSize) - 1
			rayX = (playerY-rayY)*aTan + playerX

			rayTargetY = - float64(blocSize)
			rayTargetX = - rayTargetY * aTan
		} else if rayAngle < math.Pi {
			// looking down
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

			mapIndex = mapY*mapSize + mapX

			// hit wall
			if mapIndex > 0 && mapIndex < mapSize*mapSize && level[mapIndex] == 1 {
				dof = 8
				hx = rayX
				hy = rayY
				distH = dist(playerX, playerY, hx, hy, rayAngle)
			} else {
				rayX += rayTargetX
				rayY += rayTargetY
				dof++
			}
		}

		// check Vertical
		dof = 0
		nTan := -math.Tan(rayAngle)
		P2 := math.Pi / 2
		P3 := 3 * P2

		distV := 1000000.0
		vx := playerX
		vy := playerY

		if rayAngle > P2 && rayAngle < P3 {
			// looking left
			rayX = math.Trunc(playerX/float64(blocSize))*float64(blocSize) - 1
			rayY = (playerX-rayX)*nTan + playerY

			rayTargetX = - float64(blocSize)
			rayTargetY = - rayTargetX * nTan
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

			mapIndex = mapY*mapSize + mapX

			// hit wall
			if mapIndex > 0 && mapIndex < mapSize*mapSize && level[mapIndex] == 1 {
				vx = rayX
				vy = rayY
				distV = dist(playerX, playerY, vx, vy, rayAngle)
				dof = 8
			} else {
				rayX += rayTargetX
				rayY += rayTargetY
				dof++
			}
		}

		// vertical wall
		if distV < distH {
			rayX = vx
			rayY = vy
			distT = distV

			gc.SetFillColor(color.RGBA{0xE5, 0x00, 0x00, 0xff})
			gc.SetStrokeColor(color.RGBA{0xE5, 0x00, 0x00, 0xff})
		}

		// horizontal wall
		if distH < distV {
			rayX = hx
			rayY = hy
			distT = distH

			gc.SetFillColor(color.RGBA{0xb2, 0x00, 0x00, 0xff})
			gc.SetStrokeColor(color.RGBA{0xb2, 0x00, 0x00, 0xff})
		}

		// render raycast
		gc.BeginPath()
		gc.MoveTo(playerX, playerY)
		gc.LineTo(rayX, rayY)
		gc.Close()
		gc.FillStroke()

		// render 3D walls
		ca := gs.GetPlayerAngle() - rayAngle
		if ca < 0 {
			ca += 2 * math.Pi
		}
		if ca > 2*math.Pi {
			ca -= 2 * math.Pi
		}
		distT = distT * math.Cos(ca) // fix fisheye

		lineH := float64(mapSize*320) / distT
		if lineH > 320 {
			lineH = 320
		}

		lineOffset := (320 / 2) - lineH/2

		draw2dkit.Rectangle(
			gc,
			float64(rayN*8+530),
			lineOffset,
			float64(rayN*8+530)+8,
			lineH+lineOffset,
		)
		gc.FillStroke()

		// updating angle for next ray
		rayAngle += oneRadian

		if rayAngle < 0 {
			rayAngle += 2 * math.Pi
		} else if rayAngle > 2*math.Pi {
			rayAngle -= 2 * math.Pi
		}
	}
}

func dist(aX, aY, bX, bY, angle float64) float64 {
	return math.Sqrt((bX-aX)*(bX-aX) + (bY-aY)*(bY-aY))
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
	mapSize := gs.GetMapSize()

	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			if level[x+y*mapSize] == 0 {
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
