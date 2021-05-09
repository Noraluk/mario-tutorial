package main

import (
	"fmt"
	"log"
	"strings"

	bgService "server/background/services"
	"server/config"
	marioEntity "server/mario/entities"
	marioService "server/mario/services"
	screenEntity "server/screen/entities"
	screenService "server/screen/services"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var (
	screen       *screenEntity.Screen
	actionIndex  = 0
	marioActions = marioEntity.NewActions()
)

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
	screenService := screenService.New(backgroundService)
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

	server.OnEvent("/", "setup", func(s socketio.Conn, msg string) {
		err := backgroundService.SetBackground()
		if err != nil {
			log.Fatal(err)
		}

		mario := marioEntity.NewMario(marioActions[0])
		err = backgroundService.Setup()
		if err != nil {
			log.Fatal("setup error ", err.Error())
		}

		camera := screenEntity.NewCamera()
		screen = screenService.GetScreen(camera, mario)
		s.Emit("draw", screen)
		s.Emit("drawMario", mario)
	})

	server.OnEvent("/", "fall", func(s socketio.Conn, msg string) {
		screen.Mario.SetCorner(config)

		screen.Mario.Velocity.Y += 0.05
		if screen.Mario.Velocity.Y > 0 {
			marioService.CanFall(screen.Mario)
		} else {
			marioService.IsCeiling(screen.Mario)
		}

		if screen.Mario.Velocity.Y == 0 && strings.Contains(screen.Mario.Movement, "jump") {
			screen.Mario.Action = marioActions[0]
			screen.Mario.Movement = strings.Replace(screen.Mario.Movement, "jump", "", -1)
		} else if screen.Mario.Velocity.Y != 0 {
			screen.Mario.Action = marioActions[3]
		}

		screen.Mario.Position.X += screen.Mario.Velocity.X
		screen.Mario.Position.Y += screen.Mario.Velocity.Y

		screen := screenService.GetScreen(screen.Camera, screen.Mario)

		s.Emit("draw", screen)
		s.Emit("drawMario", screen.Mario)

		screen.Mario.Velocity.X = 0
	})

	server.OnEvent("/", "right", func(s socketio.Conn, msg string) {
		screen.Mario.SetCorner(config)

		marioService.MoveRight(screen.Mario)

		if strings.Contains(screen.Mario.Movement, "jump") {
			screen.Mario.Movement = "rightjump"
		} else {
			screen.Mario.Movement = "right"
		}

		if screen.Mario.Velocity.X > 0 && !strings.Contains(screen.Mario.Movement, "jump") {
			actionIndex++
			screen.Mario.Action = marioActions[actionIndex%4]
		}
	})

	server.OnEvent("/", "left", func(s socketio.Conn, msg string) {
		screen.Mario.SetCorner(config)

		marioService.MoveLeft(screen.Mario)

		if strings.Contains(screen.Mario.Movement, "jump") {
			screen.Mario.Movement = "leftjump"
		} else {
			screen.Mario.Movement = "left"
		}

		if screen.Mario.Velocity.X < 0 && !strings.Contains(screen.Mario.Movement, "jump") {
			actionIndex++
			screen.Mario.Action = marioActions[actionIndex%4]
		}
	})

	server.OnEvent("/", "jump", func(s socketio.Conn, msg string) {
		screen.Mario.SetCorner(config)

		canFall := marioService.CanFall(screen.Mario)
		if !canFall {
			screen.Mario.Velocity.Y = -2.7
			screen.Mario.Movement = fmt.Sprintf("%sjump", screen.Mario.Movement)

			screen.Mario.Action = marioActions[3]
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
