package marioEntity

import (
	bgEntity "server/background/entities"
	"server/config"
)

type Mario struct {
	X        int               `json:"x"`
	Y        int               `json:"y"`
	Width    float64           `json:"width"`
	Height   float64           `json:"height"`
	Position bgEntity.Position `json:"position"`
	Velocity Velocity          `json:"velocity"`
	Action   string            `json:"action"`
	Corner   Corner            `json:"corner"`
}

type Velocity struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Corner struct {
	TopLeft         bgEntity.Position `json:"topLeft"`
	TopRight        bgEntity.Position `json:"topRight"`
	BottomLeft      bgEntity.Position `json:"bottomLeft"`
	BottomRight     bgEntity.Position `json:"bottomRight"`
	CenterLeft      bgEntity.Position `json:"centerLeft"`
	CenterRight     bgEntity.Position `json:"centerFront"`
	TopLeftTile     string            `json:"topLeftTile"`
	TopRightTile    string            `json:"topRightTile"`
	BottomLeftTile  string            `json:"bottomLeftTile"`
	BottomRightTile string            `json:"bottomRightTile"`
	CenterLeftTile  string            `json:"centerLeftTile"`
	CenterRightTile string            `json:"centerRightTile"`
}

func (e *Mario) SetCorner(config config.Config) {
	nextPositionX := int(e.Position.X)
	nextPositionY := int(e.Position.Y)

	e.Corner.BottomLeft = bgEntity.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY) + e.Height,
	}
	e.Corner.BottomLeftTile = config.GetCollider(e.Corner.BottomLeft)

	e.Corner.BottomRight = bgEntity.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY) + e.Height,
	}
	e.Corner.BottomRightTile = config.GetCollider(e.Corner.BottomRight)

	e.Corner.TopLeft = bgEntity.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY),
	}
	e.Corner.TopLeftTile = config.GetCollider(e.Corner.TopLeft)

	e.Corner.TopRight = bgEntity.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY),
	}
	e.Corner.TopRightTile = config.GetCollider(e.Corner.TopRight)

	e.Corner.CenterRight = bgEntity.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY) + e.Height/2,
	}
	e.Corner.CenterRightTile = config.GetCollider(e.Corner.CenterRight)

	e.Corner.CenterLeft = bgEntity.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY) + e.Height/2,
	}
	e.Corner.CenterLeftTile = config.GetCollider(e.Corner.CenterLeft)
}
