package main

import (
	"fmt"
	"log"

	bgEntity "server/background/entities"
	bgService "server/background/services"
	"server/config"
	"server/constants"
	marioEntity "server/mario/entities"
	marioService "server/mario/services"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Screen struct {
	Background []bgEntity.Background   `json:"backgrounds"`
	Colliders  []bgEntity.TileCollider `json:"colliders"`
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
		mario = &marioEntity.Mario{X: 276, Y: 44, Width: 16, Height: 16, Position: bgEntity.Position{X: 0, Y: 0}, Velocity: marioEntity.Velocity{X: 5, Y: 1}}
		err = backgroundService.Setup()
		if err != nil {
			log.Fatal("setup error ", err.Error())
		}
	})

	server.OnEvent("/", "draw", func(s socketio.Conn, msg string) {
		level, err := backgroundService.GetBackground()
		if err != nil {
			log.Fatal("err ", err.Error())
		}
		screen := Screen{Background: level.Backgrounds, Colliders: backgroundService.GetColliders()}

		canFall := marioService.CanFall(mario)
		mario.Position.Y += mario.Velocity.Y
		if !canFall {
			mario.Position.Y = mario.Position.Y - (mario.Position.Y % constants.TILE_SILE)
		}

		s.Emit("draw", screen)
		s.Emit("drawMario", mario)
	})

	server.OnEvent("/", "right", func(s socketio.Conn, msg string) {
		mario.Position.X += mario.Velocity.X
		marioService.MoveHorizontal(mario)
	})

	server.OnEvent("/", "left", func(s socketio.Conn, msg string) {
		mario.Position.X -= mario.Velocity.X
		marioService.MoveHorizontal(mario)
	})

	server.OnEvent("/", "up", func(s socketio.Conn, msg string) {
		canFall := marioService.CanFall(mario)
		if !canFall {
			mario.Position.Y -= 40
		}
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
