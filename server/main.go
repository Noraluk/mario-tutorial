package main

import (
	"fmt"
	"log"

	bgEntity "server/background/entities"
	bgService "server/background/services"
	"server/common"
	"server/config"
	"server/constants"
	marioEntity "server/mario/entities"
	marioService "server/mario/services"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Camera struct {
	Size     common.Size     `json:"size"`
	Position common.Position `json:"position"`
}

type Screen struct {
	Background []bgEntity.Background `json:"backgrounds"`
	Camera     Camera                `json:"camera"`
	Mario      *marioEntity.Mario    `json:"mario"`
}

func getScreen(config config.Config, backgroundService bgService.Background, mario *marioEntity.Mario) Screen {
	level, err := backgroundService.GetBackground()
	if err != nil {
		log.Fatal(err)
	}

	camera := Camera{Position: common.Position{X: 0, Y: 0}, Size: common.Size{Width: 256, Height: 256}}
	if mario.Position.X >= constants.HALF_SCREEN {
		camera.Position.X = mario.Position.X - constants.HALF_SCREEN
	}

	cameraStart := camera.Position.X
	cameraEnd := int(camera.Position.X+constants.TILE_SILE) + camera.Size.Width

	for i, bg := range level.Backgrounds {
		newRanges := []bgEntity.Ranges{}
		for _, val := range bg.Ranges {
			x1 := val.X1 * constants.TILE_SILE
			x2 := val.X2 * constants.TILE_SILE

			if (cameraStart > float64(x1) && cameraStart > float64(x2)) || (x1 > cameraEnd && x2 > cameraEnd) {
				continue
			}
			newRange := bgEntity.Ranges{X1: x1 / constants.TILE_SILE, X2: x2 / constants.TILE_SILE, Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: int(constants.TILE_SILE), Height: int(constants.TILE_SILE)}}
			if cameraStart > float64(x1) && float64(x2) > cameraStart {
				newRange = bgEntity.Ranges{X1: int(cameraStart / constants.TILE_SILE), X2: x2 / int(constants.TILE_SILE), Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: (val.X2 - int(cameraStart)) / constants.TILE_SILE, Height: int(constants.TILE_SILE)}}
			} else if cameraEnd > x1 && x2 > cameraEnd {
				newRange = bgEntity.Ranges{X1: x1 / constants.TILE_SILE, X2: cameraEnd / constants.TILE_SILE, Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: (cameraEnd - val.X1) / constants.TILE_SILE, Height: int(constants.TILE_SILE)}}
			}
			newRanges = append(newRanges, newRange)
		}
		level.Backgrounds[i].Ranges = newRanges
	}

	screen := Screen{Background: level.Backgrounds, Mario: mario, Camera: camera}
	return screen
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func main() {
	config := config.New()
	backgroundService := bgService.New(config)
	marioService := marioService.New(config)
	router := gin.New()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	var mario *marioEntity.Mario
	server.OnEvent("/", "setup", func(s socketio.Conn, msg string) {
		mario = &marioEntity.Mario{X: 276, Y: 44, Width: 16, Height: 16, Position: common.Position{X: 0, Y: 0}, Velocity: marioEntity.Velocity{X: 0, Y: 0.1}}
		err = backgroundService.Setup()
		if err != nil {
			log.Fatal("setup error ", err.Error())
		}

		s.Emit("draw", getScreen(config, backgroundService, mario))
		s.Emit("drawMario", mario)
	})

	server.OnEvent("/", "draw", func(s socketio.Conn, msg string) {
		mario.SetCorner(config)

		if msg == "right" {
			marioService.MoveRight(mario)
		} else if msg == "left" {
			marioService.MoveLeft(mario)
		} else if msg == "up" {
			canFall := marioService.CanFall(mario)
			if !canFall {
				mario.Velocity.Y = -2.5
				mario.Action = "jump"
			}
		}

		mario.Velocity.Y += 0.05
		if mario.Velocity.Y > 0 {
			marioService.CanFall(mario)
			mario.Action = ""
		} else {
			marioService.IsCeiling(mario)
		}

		mario.Position.X += mario.Velocity.X
		mario.Position.Y += mario.Velocity.Y

		if mario.Velocity.X != 0 || mario.Velocity.Y != 0 || mario.Action == "jump" {
			screen := getScreen(config, backgroundService, mario)

			s.Emit("draw", screen)
			s.Emit("drawMario", screen.Mario)
		}

		mario.Velocity.X = 0
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:3000"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	router.Run()
}
