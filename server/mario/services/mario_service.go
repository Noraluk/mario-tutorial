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
	MoveHorizontal(mario *marioEntity.Mario)
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
	bottomLeftPosition := bgEntity.Position{
		X: mario.Position.X - (mario.Position.X % constants.TILE_SILE),
		Y: mario.Position.Y - (mario.Position.Y % constants.TILE_SILE) + mario.Height,
	}
	bottomRightPosition := bgEntity.Position{
		X: mario.Position.X + mario.Width - (mario.Position.X % constants.TILE_SILE),
		Y: mario.Position.Y - (mario.Position.Y % constants.TILE_SILE) + mario.Height,
	}

	bottomLeftTile := s.config.GetCollider(bgEntity.Position{X: bottomLeftPosition.X, Y: bottomLeftPosition.Y})
	bottomRightTile := s.config.GetCollider(bgEntity.Position{X: bottomRightPosition.X, Y: bottomRightPosition.Y})
	if bottomLeftTile == constants.GROUND || bottomRightTile == constants.GROUND {
		return false
	}
	return true
}

func (s *mario) MoveHorizontal(mario *marioEntity.Mario) {
	topRightPosition := bgEntity.Position{
		X: mario.Position.X + mario.Width - (mario.Position.X % constants.TILE_SILE),
		Y: mario.Position.Y - (mario.Position.Y % constants.TILE_SILE),
	}

	topLeftPosition := bgEntity.Position{
		X: mario.Position.X - (mario.Position.X % constants.TILE_SILE),
		Y: mario.Position.Y - (mario.Position.Y % constants.TILE_SILE),
	}

	topRightTile := s.config.GetCollider(bgEntity.Position{X: topRightPosition.X, Y: topRightPosition.Y})
	if topRightTile == constants.GROUND {
		mario.Position.X = mario.Position.X - (mario.Position.X % constants.TILE_SILE)
		return
	}

	topLeftTile := s.config.GetCollider(bgEntity.Position{X: topLeftPosition.X, Y: topLeftPosition.Y})
	if topLeftTile == constants.GROUND {
		mario.Position.X = mario.Position.X - (mario.Position.X % constants.TILE_SILE) + mario.Width
		return
	}
}
