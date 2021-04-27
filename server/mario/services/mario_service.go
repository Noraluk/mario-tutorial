package marioService

import (
	bgService "server/background/services"
	"server/common"
	"server/config"
	"server/constants"
	marioEntity "server/mario/entities"
)

type Mario interface {
	CanFall(mario *marioEntity.Mario) bool
	IsCeiling(mario *marioEntity.Mario) bool
	MoveRight(mario *marioEntity.Mario)
	MoveLeft(mario *marioEntity.Mario)
}

type mario struct {
	backgroundService bgService.Background
	config            config.Config
}

func New(config config.Config) Mario {
	return &mario{
		backgroundService: bgService.New(config),
		config:            config,
	}
}

func (s *mario) CanFall(mario *marioEntity.Mario) bool {
	currentPixelGroundCount := 0
	for x := int(mario.Corner.BottomLeft.X); x <= int(mario.Corner.BottomRight.X); x++ {
		tile := s.config.GetCollider(common.Position{X: float64(x), Y: float64(int(mario.Corner.BottomLeft.Y))})
		if tile == constants.SKY {
			currentPixelGroundCount++
		}
	}
	if currentPixelGroundCount < 15 {
		mario.Velocity.Y = 0
		return false
	}

	nextY := int(mario.Corner.BottomLeft.Y + mario.Velocity.Y)
	nextPixelGroundCount := 0
	for x := int(mario.Corner.BottomLeft.X); x <= int(mario.Corner.BottomRight.X); x++ {
		tile := s.config.GetCollider(common.Position{X: float64(x), Y: float64(nextY)})
		if tile == constants.SKY {
			nextPixelGroundCount++
		}
	}
	if nextPixelGroundCount < 15 {
		remainder := nextY % constants.TILE_SILE
		distance := float64(nextY-remainder) - (mario.Position.Y + mario.Height)
		mario.Velocity.Y = distance
		return true
	}

	return true
}

func (s *mario) IsCeiling(mario *marioEntity.Mario) bool {
	pixelGroundCount := 0
	for x := int(mario.Corner.TopLeft.X); x <= int(mario.Corner.TopRight.X); x++ {
		tile := s.config.GetCollider(common.Position{X: float64(x), Y: mario.Corner.TopLeft.Y})
		if tile == constants.GROUND {
			pixelGroundCount++
		}
	}

	if pixelGroundCount > 2 {
		mario.Velocity.Y = 0
		return true
	}
	return false
}

func (s *mario) MoveRight(mario *marioEntity.Mario) {
	mario.Velocity.X = 5

	currentPixelGroundCount := 0
	for y := int(mario.Corner.TopRight.Y); y <= int(mario.Corner.BottomRight.Y); y++ {
		tile := s.config.GetCollider(common.Position{X: mario.Corner.TopRight.X, Y: float64(y)})
		if tile == constants.GROUND {
			currentPixelGroundCount++
		}
	}
	if currentPixelGroundCount > 2 {
		mario.Velocity.X = 0
	}

	nextX := int(mario.Corner.TopRight.X + mario.Velocity.X)
	nextPixelGroundCount := 0
	for y := int(mario.Corner.TopRight.Y); y <= int(mario.Corner.BottomRight.Y); y++ {
		tile := s.config.GetCollider(common.Position{X: float64(nextX), Y: float64(y)})
		if tile == constants.GROUND {
			nextPixelGroundCount++
		}
	}
	if nextPixelGroundCount > 2 {
		remainder := nextX % constants.TILE_SILE
		distance := float64(nextX-remainder) - (mario.Position.X + mario.Width)
		mario.Velocity.X = distance
	}
}

func (s *mario) MoveLeft(mario *marioEntity.Mario) {
	mario.Velocity.X = -5

	currentPixelGroundCount := 0
	for y := int(mario.Corner.TopLeft.Y); y <= int(mario.Corner.BottomLeft.Y); y++ {
		tile := s.config.GetCollider(common.Position{X: mario.Corner.TopLeft.X, Y: float64(y)})
		if tile == constants.GROUND {
			currentPixelGroundCount++
		}
	}
	if currentPixelGroundCount > 2 {
		mario.Velocity.X = 0
	}

	prevX := int(mario.Corner.TopLeft.X + mario.Velocity.X)
	nextPixelGroundCount := 0
	for y := int(mario.Corner.TopLeft.Y); y <= int(mario.Corner.BottomLeft.Y); y++ {
		tile := s.config.GetCollider(common.Position{X: float64(prevX), Y: float64(y)})
		if tile == constants.GROUND {
			nextPixelGroundCount++
		}
	}
	if nextPixelGroundCount > 2 {
		remainder := (prevX % constants.TILE_SILE)
		switch remainder {
		case 0:
			mario.Velocity.X = float64(prevX) - mario.Position.X
		default:
			mario.Velocity.X = float64(prevX+(constants.TILE_SILE-remainder)) - mario.Position.X
		}
	}
}
