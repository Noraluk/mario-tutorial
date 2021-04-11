package bg_service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	bgEntity "server/background/entities"
	"server/config"
	rd "server/redis"
)

const tileSize int = 16

var (
	colliders []bgEntity.TileCollider
	col       map[bgEntity.Position]string
)

type Background interface {
	Setup() error
	GetBackground() (*bgEntity.Level, error)
	// GetPositions() []bgEntity.Position
	// SetPositions(colliders []bgEntity.Position)
	GetColliders() []bgEntity.TileCollider
	SetColliders(colliders []bgEntity.TileCollider)
	// Get(position bgEntity.Position) string
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

	newColliders := []bgEntity.TileCollider{}
	for _, bg := range level.Backgrounds {
		for _, val := range bg.Ranges {
			for x := val.X1; x < val.X2; x++ {
				for y := val.Y1; y < val.Y2; y++ {
					position := bgEntity.Position{X: (x) * tileSize, Y: (y) * tileSize}
					collider := bgEntity.TileCollider{Name: bg.Tile, Tile: position}

					newColliders = append(newColliders, collider)
					s.config.SetCollider(position, bg.Tile)
				}
			}
		}
	}
	s.SetColliders(newColliders)
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
