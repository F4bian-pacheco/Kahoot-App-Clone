package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"quiz.com/quiz/internal/entity"
)

type Player struct {
	Name       string
	Connection *websocket.Conn
}

type Game struct {
	Id      uuid.UUID
	Quiz    entity.Quiz
	Code    string
	Players []Player

	Host *websocket.Conn
}

func generateCode() string {
	return strconv.Itoa((100000 + rand.Intn(900000)))
}

func New(quiz entity.Quiz, host *websocket.Conn) Game {
	return Game{
		Id:      uuid.New(),
		Quiz:    quiz,
		Code:    generateCode(),
		Players: []Player{},
		Host:    host,
	}
}

func (g *Game) Start() {

	go func() {
		for {
			g.Tick()
			time.Sleep(time.Second)
		}
	}()

}

func (g *Game) Tick() {}

func (g *Game) OnPlayerJoin(name string, connection *websocket.Conn) {
	fmt.Println(name, " joined the game")
	g.Players = append(g.Players, Player{
		Name:       name,
		Connection: connection,
	})
}
