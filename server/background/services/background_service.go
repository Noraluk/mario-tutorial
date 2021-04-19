package bg_service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	bgEntity "server/background/entities"
	"server/config"
	"server/constants"
	rd "server/redis"
)

var (
	colliders []bgEntity.TileCollider
)

type Background interface {
	Setup() error
	GetBackground() (*bgEntity.Level, error)
	GetColliders() []bgEntity.TileCollider
	SetColliders(colliders []bgEntity.TileCollider)
}

type background struct {
	redis  rd.RedisStorage
	config config.Config
}

func New(config config.Config) Background {
	return &background{
		redis:  rd.New(),
		config: config,
	}
}

func (s *background) Setup() error {
	level, err := s.GetBackground()
	if err != nil {
		return err
	}

	for _, bg := range level.Backgrounds {
		for _, val := range bg.Ranges {
			for x := val.X1 * constants.TILE_SILE; x <= val.X2*constants.TILE_SILE; x++ {
				for y := val.Y1 * constants.TILE_SILE; y <= val.Y2*constants.TILE_SILE; y++ {
					position := bgEntity.Position{X: float64(x), Y: float64(y)}
					s.config.SetCollider(position, bg.Tile)
				}
			}
		}
	}
	return nil
}

func (s *background) GetBackground() (*bgEntity.Level, error) {
	oneToOne, err := os.Open("static/levels/1-1.json")
	if err != nil {
		return nil, err
	}
	defer oneToOne.Close()

	b, err := ioutil.ReadAll(oneToOne)
	if err != nil {
		return nil, err
	}

	var level bgEntity.Level
	err = json.Unmarshal(b, &level)
	if err != nil {
		return nil, err
	}

	return &level, nil
}

func (s *background) GetColliders() []bgEntity.TileCollider {
	return colliders
}

func (s *background) SetColliders(newColliders []bgEntity.TileCollider) {
	colliders = newColliders
}
