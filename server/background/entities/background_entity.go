package bgEntity

import "server/common"

type Range struct {
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
}

func NewRange(x1, x2, y1, y2 int) Range {
	return Range{X1: x1, X2: x2, Y1: y1, Y2: y2}
}

type Background struct {
	Tile      string          `json:"tile"`
	IsCollide bool            `json:"isCollide"`
	Position  common.Position `json:"position"`
	Ranges    []Range         `json:"ranges"`
	Animation Animation       `json:"animation"`
}

func NewBackground(tile string, isCollide bool, position common.Position, ranges []Range, animation Animation) Background {
	return Background{
		Tile:      tile,
		IsCollide: isCollide,
		Position:  position,
		Ranges:    ranges,
		Animation: animation,
	}
}

type Animation struct {
	Frames   []common.Position `json:"frames"`
	Current  int               `json:"current"`
	Duration int               `json:"duration"`
}

type Level struct {
	Backgrounds []Background `json:"backgrounds"`
}

type TileCollider struct {
	Name string          `json:"name"`
	Tile common.Position `json:"tile"`
}
