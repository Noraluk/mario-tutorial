package marioEntity

import (
	"server/common"
	"server/config"
)

type Mario struct {
	Action   Action          `json:"action"`
	Position common.Position `json:"position"`
	Velocity Velocity        `json:"velocity"`
	Movement string          `json:"movement"`
	Corner   Corner          `json:"corner"`
}

func NewMario(action Action) *Mario {
	return &Mario{Action: action, Position: common.Position{X: 0, Y: 0}, Velocity: Velocity{X: 0, Y: 0.1}}
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
		Y: float64(nextPositionY) + float64(e.Action.Size.Height),
	}
	e.Corner.IsBottomLeftCollide = config.GetCollider(e.Corner.BottomLeft)

	e.Corner.BottomRight = common.Position{
		X: float64(nextPositionX + e.Action.Size.Width),
		Y: float64(nextPositionY + e.Action.Size.Height),
	}
	e.Corner.IsBottomRighttCollide = config.GetCollider(e.Corner.BottomRight)

	e.Corner.TopLeft = common.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY),
	}
	e.Corner.IsTopLeftCollide = config.GetCollider(e.Corner.TopLeft)

	e.Corner.TopRight = common.Position{
		X: float64(nextPositionX + e.Action.Size.Width),
		Y: float64(nextPositionY),
	}
	e.Corner.IsTopRightCollide = config.GetCollider(e.Corner.TopRight)

	e.Corner.CenterRight = common.Position{
		X: float64(nextPositionX + e.Action.Size.Width),
		Y: float64(nextPositionY + e.Action.Size.Height/2),
	}
	e.Corner.IsCenterRightCollide = config.GetCollider(e.Corner.CenterRight)

	e.Corner.CenterLeft = common.Position{
		X: float64(nextPositionX),
		Y: float64(nextPositionY + e.Action.Size.Height/2),
	}
	e.Corner.IsCenterLeftCollide = config.GetCollider(e.Corner.CenterLeft)
}

type Action struct {
	Name  string          `json:"name"`
	Image common.Position `json:"image"`
	Size  common.Size     `json:"size"`
}

func NewActions() []Action {
	defaultSize := common.Size{Width: 16, Height: 16}

	actions := []Action{}
	actions = append(actions, Action{Name: "stand", Image: common.Position{X: 276, Y: 44}, Size: defaultSize})
	actions = append(actions, Action{Name: "run-1", Image: common.Position{X: 290, Y: 44}, Size: defaultSize})
	actions = append(actions, Action{Name: "run-2", Image: common.Position{X: 304, Y: 43}, Size: defaultSize})
	actions = append(actions, Action{Name: "run-3", Image: common.Position{X: 321, Y: 44}, Size: defaultSize})

	return actions
}
