package wolfenstein

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
)

type GameState struct {
	level     []int
	mapSize   int
	blockSize int

	player Player
}

type Player struct {
	position Point
	delta    Point
}

type Point struct {
	x     float64
	y     float64
	angle float64
}

func NewGameState(width, height int) (*GameState, error) {
	var gs GameState

	// silly level
	gs.level = []int{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 1, 0, 0, 0, 0, 1,
		1, 0, 1, 0, 0, 0, 0, 1,
		1, 0, 1, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
	}

	gs.mapSize = 8
	gs.blockSize = 64

	gs.player = Player{
		position: Point{
			float64(gs.mapSize * gs.blockSize / 2),
			float64(gs.mapSize * gs.blockSize / 2),
			0.0,
		},
		delta: Point{0, 0, 0.0},
	}

	gs.updateDelta()

	return &gs, nil
}

func (gs *GameState) GetMapSize() int {
	return gs.mapSize
}

func (gs *GameState) GetLevel() []int {
	return gs.level
}

func (gs *GameState) GetPlayer() Player {
	return gs.player
}

func (gs *GameState) GetPlayerPosition() (x, y, deltaX, deltaY float64) {
	return gs.player.position.x, gs.player.position.y, gs.player.delta.x, gs.player.delta.y
}

func (gs *GameState) GetBlockSize() int {
	return gs.blockSize
}

func (gs *GameState) GetPlayerAngle() float64 {
	return gs.player.position.angle
}

func (gs *GameState) MoveUp() {
	gs.player.position.x += gs.player.delta.x
	gs.player.position.y += gs.player.delta.y
}

func (gs *GameState) MoveDown() {
	gs.player.position.x -= gs.player.delta.x
	gs.player.position.y -= gs.player.delta.y
}

func (gs *GameState) MoveLeft() {
	gs.player.position.angle -= 0.1

	if gs.player.position.angle < 0 {
		gs.player.position.angle += 2 * math.Pi
	}

	gs.updateDelta()
}

func (gs *GameState) MoveRight() {
	gs.player.position.angle += 0.1

	if gs.player.position.angle > 2*math.Pi {
		gs.player.position.angle -= 2 * math.Pi
	}

	gs.updateDelta()
}

func (gs *GameState) updateDelta() {
	gs.player.delta.x = math.Cos(gs.player.position.angle) * 5
	gs.player.delta.y = math.Sin(gs.player.position.angle) * 5
}

func (gs *GameState) RenderRay(gc *draw2dimg.GraphicContext) {
	// player position as origin
	posX := gs.player.position.x
	posY := gs.player.position.y
	posAngle := gs.player.position.angle

	blockSize := gs.GetBlockSize()

	var rayX, rayY float64

	for x := 0; x < 1; x++ {
		//which cell of the map we're in
		mapX := int(math.Trunc(posX / float64(blockSize)))*blockSize

		// right
		if posAngle < math.Pi/2 || posAngle > 3*math.Pi/2 {
			rayX = posX + (float64(blockSize) - (posX - float64(mapX)))
			rayY = posY + (float64(blockSize)-(posX-float64(mapX)))*math.Tan(posAngle)
		}

		// left
		if posAngle > math.Pi/2 && posAngle < 3*math.Pi/2 {
			rayX = posX - (float64(blockSize) - (posX - float64(mapX)))
			rayY = posY + (float64(blockSize)-(posX-float64(mapX)))*-math.Tan(posAngle)
		}

		gc.SetFillColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
		gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
		draw2dkit.Circle(gc, rayX, rayY, 2)
		gc.FillStroke()
	}
}
