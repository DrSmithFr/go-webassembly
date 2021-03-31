package wolfenstein

import "math"

type GameState struct {
	level []int

	mapX  int
	mapY  int

	blockSize int

	playerX float64
	playerY float64

	playerDeltaX float64
	playerDeltaY float64

	playerAngle float64
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

	gs.mapX = 8
	gs.mapY = 8

	// use all screen
	gs.blockSize = int(math.Min(float64(width), float64(height))) / 8

	gs.playerX = float64(width / 2)
	gs.playerY = float64(height / 2)

	gs.playerAngle = 0
	gs.updateDelta()

	return &gs, nil
}

func (gs *GameState) GetMapSize() (x int,y int) {
	return gs.mapX, gs.mapY
}

func (gs *GameState) GetLevel() []int {
	return gs.level
}

func (gs *GameState) GetPlayerPosition()(x, y, deltaX, deltaY float64)  {
	return gs.playerX, gs.playerY, gs.playerDeltaX, gs.playerDeltaY
}

func (gs *GameState) GetBlockSize() int  {
	return gs.blockSize
}

func (gs *GameState) GetPlayerAngle() float64  {
	return gs.playerAngle
}

func (gs *GameState) MoveUp() {
	gs.playerX += gs.playerDeltaX
	gs.playerY += gs.playerDeltaY
}

func (gs *GameState) MoveDown() {
	gs.playerX -= gs.playerDeltaX
	gs.playerY -= gs.playerDeltaY
}

func (gs *GameState) MoveLeft() {
	gs.playerAngle -= 0.1

	if gs.playerAngle < 0 {
		gs.playerAngle += 2 * math.Pi
	}

	gs.updateDelta()
}

func (gs *GameState) MoveRight() {
	gs.playerAngle += 0.1

	if gs.playerAngle > 2 * math.Pi {
		gs.playerAngle -= 2 * math.Pi
	}

	gs.updateDelta()
}

func (gs *GameState)updateDelta()  {
	gs.playerDeltaX = math.Cos(gs.playerAngle) * 5
	gs.playerDeltaY = math.Sin(gs.playerAngle) * 5
}