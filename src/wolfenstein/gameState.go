package wolfenstein

import "math"

type GameState struct {
	level []int
	blockSize int

	playerX int
	playerY int
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

	// use all screen
	gs.blockSize = int(math.Min(float64(width), float64(height))) / 8

	gs.playerX = width / 2
	gs.playerY = height / 2

	return &gs, nil
}

func (gs *GameState) GetLevel() []int {
	return gs.level
}

func (gs *GameState) GetPlayerPosition()(width int, height int)  {
	return gs.playerX, gs.playerY
}

func (gs *GameState) GetBlockSize() int  {
	return gs.blockSize
}