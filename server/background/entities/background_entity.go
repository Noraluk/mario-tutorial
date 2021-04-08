package bg_entity

type Ranges struct {
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Background struct {
	Tile     string   `json:"tile"`
	Position Position `json:"position"`
	Ranges   []Ranges `json:"ranges"`
}

type Level struct {
	Backgrounds []Background `json:"backgrounds"`
}

type TileCollider struct {
	Level string   `json:"name"`
	Tile  Position `json:"tile"`
}
