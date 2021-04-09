package main

import (
	"encoding/json"
	"fmt"
	"log"
	"server/mario"
	rd "server/redis"
	"strings"

	bg_entity "server/background/entities"
	bg_service "server/background/services"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
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
	bgService := bg_service.New()
	redis := rd.New()
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
		level, err := bgService.GetBackground()
		if err != nil {
			log.Fatal(err)
		}
		s.Emit("setup", level.Backgrounds)
		err = bgService.Setup(level.Backgrounds, s)
		if err != nil {
			log.Fatal(err)
		}

	})

	mario := mario.Model{X: 276, Y: 44, Width: 16, Height: 16, Position: bg_entity.Position{X: 0, Y: 0}}
	server.OnEvent("/", "mario", func(s socketio.Conn, msg string) {
		b, err := json.Marshal(bg_entity.TileCollider{Level: "1-1", Tile: bg_entity.Position{X: mario.Position.X, Y: mario.Position.Y + 16}})
		if err != nil {
			log.Fatal(err)
		}
		result, _ := redis.Get(string(b))
		if strings.Compare(result, "ground") != 0 {
			s.Emit("mario", mario)
			mario.Position.Y++
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
