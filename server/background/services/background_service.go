package bg_service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	bg_entity "server/background/entities"
	rd "server/redis"
)

type Background interface {
	Setup(backgrounds []bg_entity.Background) error
	GetBackground() (*bg_entity.Level, error)
}

type background struct {
	redis rd.RedisStorage
}

func New() Background {
	return &background{
		redis: rd.New(),
	}
}

func (s *background) Setup(backgrounds []bg_entity.Background) error {
	for _, bg := range backgrounds {
		for _, val := range bg.Ranges {
			for x := val.X1; x < val.X2; x++ {
				for y := val.Y1; y < val.Y2; y++ {
					b, err := json.Marshal(bg_entity.TileCollider{Level: "1-1", Tile: bg_entity.Position{X: x, Y: y}})
					if err != nil {
						return err
					}
					_, err = s.redis.Set(string(b), bg.Tile)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *background) GetBackground() (*bg_entity.Level, error) {
	oneToOne, err := os.Open("static/levels/1-1.json")
	if err != nil {
		return nil, err
	}
	defer oneToOne.Close()

	b, err := ioutil.ReadAll(oneToOne)
	if err != nil {
		return nil, err
	}

	var level bg_entity.Level
	err = json.Unmarshal(b, &level)
	if err != nil {
		return nil, err
	}

	return &level, nil
}
