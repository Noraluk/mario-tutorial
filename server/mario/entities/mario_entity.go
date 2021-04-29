package marioEntity

import (
	"server/common"
	"server/config"
)

type Mario struct {
	X        int             `json:"x"`
	Y        int             `json:"y"`
	Width    float64         `json:"width"`
	Height   float64         `json:"height"`
	Position common.Position `json:"position"`
	Velocity Velocity        `json:"velocity"`
	Action   string          `json:"action"`
	Corner   Corner          `json:"corner"`
}

type Velocity struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Corner struct {
	TopLeft               common.Position `json:"topLeft"`
	TopRight              common.Position `json:"topRight"`
	BottomLeft            common.Position `json:"bottomLeft"`
	BottomRight           common.Position `json:"bottomRight"`
	CenterLeft            common.Position `json:"centerLeft"`
	CenterRight           common.Position `json:"centerFront"`
	IsTopLeftCollide      bool            `json:"isTopLeftCollide"`
	IsTopRightCollide     bool            `json:"isTopRightCollide"`
	IsBottomLeftCollide   bool            `json:"isBottomLeftCollide"`
	IsBottomRighttCollide bool            `json:"isBottomRightCollide"`
	IsCenterLeftCollide   bool            `json:"isCenterLeftCollide"`
	IsCenterRightCollide  bool            `json:"isCenterRightCollide"`
}

func (e *Mario) SetCorner(config config.Config) {
	nextPositionX := int(e.Position.X)
	nextPositionY := int(e.Position.Y)

	e.Corner.BottomLeft = common.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY) + e.Height,
	}
	e.Corner.IsBottomLeftCollide = config.GetCollider(e.Corner.BottomLeft)

	e.Corner.BottomRight = common.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY) + e.Height,
	}
	e.Corner.IsBottomRighttCollide = config.GetCollider(e.Corner.BottomRight)

	e.Corner.TopLeft = common.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY),
	}
	e.Corner.IsTopLeftCollide = config.GetCollider(e.Corner.TopLeft)

	e.Corner.TopRight = common.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY),
	}
	e.Corner.IsTopRightCollide = config.GetCollider(e.Corner.TopRight)

	e.Corner.CenterRight = common.Position{
		X: float64(nextPositionX) + e.Width,
		Y: float64(nextPositionY) + e.Height/2,
	}
	e.Corner.IsCenterRightCollide = config.GetCollider(e.Corner.CenterRight)

	e.Corner.CenterLeft = common.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY) + e.Height/2,
	}
	e.Corner.IsCenterLeftCollide = config.GetCollider(e.Corner.CenterLeft)
}
