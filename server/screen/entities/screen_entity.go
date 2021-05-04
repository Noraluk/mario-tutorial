package screenEntity

import (
	bgEntity "server/background/entities"
	"server/common"
	marioEntity "server/mario/entities"
)

type Camera struct {
	Size     common.Size     `json:"size"`
	Position common.Position `json:"position"`
}

func NewCamera() *Camera {
	return &Camera{Position: common.Position{X: 0, Y: 0}, Size: common.Size{Width: 256, Height: 256}}
}

type Screen struct {
	Backgrounds []bgEntity.Background `json:"backgrounds"`
	Camera      *Camera               `json:"camera"`
	Mario       *marioEntity.Mario    `json:"mario"`
}

func NewScreen(bgs []bgEntity.Background, camera *Camera, mario *marioEntity.Mario) *Screen {
	return &Screen{
		Backgrounds: bgs,
		Camera:      camera,
		Mario:       mario,
	}
}

func (sc *Screen) SetBackgrounds(bgs []bgEntity.Background) {
	sc.Backgrounds = bgs
}

func (sc *Screen) SetCamera(camera *Camera) {
	sc.Camera = camera
}

func (sc *Screen) SetMario(mario *marioEntity.Mario) {
	sc.Mario = mario
}
