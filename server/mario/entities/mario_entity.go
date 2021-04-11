package marioEntity

import bgEntity "server/background/entities"

type Mario struct {
	X        int               `json:"x"`
	Y        int               `json:"y"`
	Width    int               `json:"width"`
	Height   int               `json:"height"`
	Position bgEntity.Position `json:"position"`
	Velocity Velocity          `json:"velocity"`
}

type Velocity struct {
	X int `json:"x"`
	Y int `JSON:"y"`
}
