package bg_service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	bgEntity "server/background/entities"
	"server/common"
	"server/config"
	"server/constants"
)

type Background interface {
	Setup() error
	SetBackground() error
	GetBackground() bgEntity.Level
}

type background struct {
	config config.Config
	level  bgEntity.Level
}

func New(config config.Config) Background {
	return &background{
		config: config,
	}
}

func (s *background) Setup() error {
	level := s.GetBackground()

	for _, bg := range level.Backgrounds {
		for _, val := range bg.Ranges {
			for x := val.X1 * constants.TILE_SILE; x <= val.X2*constants.TILE_SILE; x++ {
				for y := val.Y1 * constants.TILE_SILE; y <= val.Y2*constants.TILE_SILE; y++ {
					position := common.Position{X: float64(x), Y: float64(y)}
					s.config.SetCollider(position, bg.IsCollide)
				}
			}
		}
	}
	return nil
}

func (s *background) SetBackground() error {
	oneToOne, err := os.Open("static/levels/1-1.json")
	if err != nil {
		return err
	}
	defer oneToOne.Close()

	b, err := ioutil.ReadAll(oneToOne)
	if err != nil {
		return err
	}

	var level bgEntity.Level
	err = json.Unmarshal(b, &level)
	if err != nil {
		return err
	}

	s.level = level
	return nil
}

func (s *background) GetBackground() bgEntity.Level {
	return s.level
}
