package bgEntity

import "server/common"

type Ranges struct {
	X1       int         `json:"x1"`
	X2       int         `json:"x2"`
	Y1       int         `json:"y1"`
	Y2       int         `json:"y2"`
	TileSize common.Size `json:"tileSize"`
}

type Background struct {
	Tile      string          `json:"tile"`
	IsCollide bool            `json:"isCollide"`
	Position  common.Position `json:"position"`
	Ranges    []Ranges        `json:"ranges"`
}

type Level struct {
	Backgrounds []Background `json:"backgrounds"`
}

type TileCollider struct {
	Name string          `json:"name"`
	Tile common.Position `json:"tile"`
}
