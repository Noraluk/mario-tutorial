package bg_service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	bg_entity "server/background/entities"
	rd "server/redis"

	socketio "github.com/googollee/go-socket.io"
)

const tileSize int = 16

type Background interface {
	Setup(backgrounds []bg_entity.Background, conn socketio.Conn) error
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

func (s *background) Setup(backgrounds []bg_entity.Background, conn socketio.Conn) error {
	for _, bg := range backgrounds {
		for _, val := range bg.Ranges {
			for x := val.X1; x < val.X2; x++ {
				for y := val.Y1; y < val.Y2; y++ {
					b, err := json.Marshal(bg_entity.TileCollider{Level: "1-1", Tile: bg_entity.Position{X: (x - 1) * tileSize, Y: (y - 1) * tileSize}})
					if err != nil {
						return err
					}
					_, err = s.redis.Set(string(b), bg.Tile)
					if err != nil {
						return err
					}
					if bg.Tile == "sky" {
						continue
					}
					conn.Emit("collider", bg_entity.Position{X: x, Y: y})
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
