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
	marioActions      []marioEntity.Action
}

func New(config config.Config, marioActions []marioEntity.Action) Mario {
	return &mario{
		backgroundService: bgService.New(config),
		config:            config,
		marioActions:      marioActions,
	}
}

func (s *mario) CanFall(mario *marioEntity.Mario) bool {
	currentPixelGroundCount := 0
	for x := int(mario.Corner.BottomLeft.X); x <= int(mario.Corner.BottomRight.X); x++ {
		isCollide := s.config.GetCollider(common.Position{X: float64(x), Y: float64(int(mario.Corner.BottomLeft.Y))})
		if !isCollide {
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
		isCollide := s.config.GetCollider(common.Position{X: float64(x), Y: float64(nextY)})
		if !isCollide {
			nextPixelGroundCount++
		}
	}
	if nextPixelGroundCount < 15 {
		remainder := nextY % constants.TILE_SILE
		distance := float64(nextY-remainder) - (mario.Position.Y + float64(mario.Action.Size.Height))
		mario.Velocity.Y = distance
		return true
	}

	return true
}

func (s *mario) IsCeiling(mario *marioEntity.Mario) bool {
	pixelGroundCount := 0
	for x := int(mario.Corner.TopLeft.X); x <= int(mario.Corner.TopRight.X); x++ {
		isCollide := s.config.GetCollider(common.Position{X: float64(x), Y: mario.Corner.TopLeft.Y})
		if isCollide {
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
	if mario.Velocity.X < 0 {
		mario.Action = s.marioActions[4]
	}
	mario.Velocity.X += 0.2
	if mario.Velocity.X > 2 {
		mario.Velocity.X = 2
	}

	currentPixelGroundCount := 0
	for y := int(mario.Corner.TopRight.Y); y <= int(mario.Corner.BottomRight.Y); y++ {
		isCollide := s.config.GetCollider(common.Position{X: mario.Corner.TopRight.X, Y: float64(y)})
		if isCollide {
			currentPixelGroundCount++
		}
	}

	nextX := int(mario.Corner.TopRight.X + mario.Velocity.X)
	nextPixelGroundCount := 0
	for y := int(mario.Corner.TopRight.Y); y <= int(mario.Corner.BottomRight.Y); y++ {
		isCollide := s.config.GetCollider(common.Position{X: float64(nextX), Y: float64(y)})
		if isCollide {
			nextPixelGroundCount++
		}
	}

	if currentPixelGroundCount > 2 || nextPixelGroundCount > 2 {
		mario.Velocity.X = 0
	}
}

func (s *mario) MoveLeft(mario *marioEntity.Mario) {
	if mario.Velocity.X > 0 {
		mario.Action = s.marioActions[4]
	}
	mario.Velocity.X -= 0.2
	if mario.Velocity.X < -2 {
		mario.Velocity.X = -2
	}

	currentPixelGroundCount := 0
	for y := int(mario.Corner.TopLeft.Y); y <= int(mario.Corner.BottomLeft.Y); y++ {
		isCollide := s.config.GetCollider(common.Position{X: mario.Corner.TopLeft.X, Y: float64(y)})
		if isCollide {
			currentPixelGroundCount++
		}
	}

	prevX := int(mario.Corner.TopLeft.X + mario.Velocity.X)
	nextPixelGroundCount := 0
	for y := int(mario.Corner.TopLeft.Y); y <= int(mario.Corner.BottomLeft.Y); y++ {
		isCollide := s.config.GetCollider(common.Position{X: float64(prevX), Y: float64(y)})
		if isCollide {
			nextPixelGroundCount++
		}
	}

	if currentPixelGroundCount > 2 || nextPixelGroundCount > 2 {
		mario.Velocity.X = 0
	}
}
