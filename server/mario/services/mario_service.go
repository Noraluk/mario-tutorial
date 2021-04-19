package marioService

import (
	bgEntity "server/background/entities"
	bgService "server/background/services"
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
		tile := s.config.GetCollider(bgEntity.Position{X: float64(x), Y: float64(int(mario.Corner.BottomLeft.Y))})
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
		tile := s.config.GetCollider(bgEntity.Position{X: float64(x), Y: float64(nextY)})
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
		tile := s.config.GetCollider(bgEntity.Position{X: float64(x), Y: mario.Corner.TopLeft.Y})
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
	nextX := int(mario.Corner.CenterRight.X + mario.Velocity.X)
	nextP := bgEntity.Position{X: float64(nextX), Y: mario.Corner.CenterRight.Y}
	nextTile := s.config.GetCollider(nextP)
	if mario.Corner.CenterRightTile == constants.SKY && nextTile == constants.GROUND {
		mario.Velocity.X = (nextP.X - float64(int(nextP.X)%constants.TILE_SILE)) - (mario.Position.X + mario.Width)
	} else if mario.Corner.CenterRightTile == constants.GROUND {
		mario.Velocity.X = 0
	}
}

func (s *mario) MoveLeft(mario *marioEntity.Mario) {
	mario.Velocity.X = -5
	prevX := int(mario.Corner.CenterLeft.X + mario.Velocity.X)
	prevP := bgEntity.Position{X: float64(prevX), Y: mario.Corner.CenterLeft.Y}
	prevTile := s.config.GetCollider(prevP)
	if mario.Corner.CenterLeftTile == constants.SKY && prevTile == constants.GROUND {
		if int(prevP.X)%constants.TILE_SILE == 0 {
			mario.Velocity.X = prevP.X - mario.Position.X
		} else {
			mario.Velocity.X = (prevP.X + (constants.TILE_SILE - float64(int(prevP.X)%constants.TILE_SILE))) - mario.Position.X
		}
	} else if mario.Corner.CenterLeftTile == constants.GROUND {
		mario.Velocity.X = 0
	}
}
