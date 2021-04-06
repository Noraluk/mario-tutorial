package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

type Ranges struct {
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Background struct {
	Tile     string   `json:"tile"`
	Position Position `json:"position"`
	Ranges   []Ranges `json:"ranges"`
}

type Level struct {
	Backgrounds []Background `json:"backgrounds"`
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

	oneToOne, err := os.Open("static/levels/1-1.json")
	if err != nil {
		log.Fatal(err)
	}
	defer oneToOne.Close()

	byteValue, err := ioutil.ReadAll(oneToOne)
	if err != nil {
		log.Fatal(err)
	}

	var level Level

	err = json.Unmarshal(byteValue, &level)
	if err != nil {
		log.Fatal(err)
	}

	server.OnEvent("/", "setup", func(s socketio.Conn, msg string) {
		log.Println(level.Backgrounds)
		for _, bg := range level.Backgrounds {
			s.Emit("setup", bg)
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
