package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/entity"
)

type NetService struct {
	quizService *QuizService
	games       []*Game

	// games []*Game
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
	}
}

type ConnectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostNamePacket struct {
	QuizId string `json:"quizid"`
}

type QuestionShowPacket struct {
	Question entity.QuizQuestion `json:"question"`
}

type ChangeGameStatePacket struct {
	State GameState `json:"state"`
}

type PlayerJoinPacket struct {
	Player Player `json:"player"`
}

type StartGamePacket struct {
}

type TickPacket struct {
	Tick int `json:"tick"`
}

func (c *NetService) packetIdToPacket(packetId uint8) any {
	switch packetId {
	case 0:
		{
			return &ConnectPacket{}
		}
	case 1:
		{
			return &HostNamePacket{}
		}
	case 5:
		{

		}
	}
	return nil
}

func (c *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QuestionShowPacket:
		{
			return 2, nil
		}
	case ChangeGameStatePacket:
		{
			return 3, nil
		}
	case PlayerJoinPacket:
		{
			return 4, nil
		}
	case TickPacket:
		{
			return 6, nil
		}
	}

	return 0, errors.New("invalid packet type")

}

func (c *NetService) geGameByCode(code string) *Game {
	for _, game := range c.games {
		if game.Code == code {
			return game
		}
	}
	return nil
}

func (c *NetService) getGameByHost(host *websocket.Conn) *Game {
	for _, game := range c.games {
		if game.Host == host {
			return game
		}
	}
	return nil
}

func (c *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {

	if len(msg) < 2 {
		return
	}

	packetId := msg[0]
	data := msg[1:]

	packet := c.packetIdToPacket(packetId)
	if packet == nil {
		return
	}

	err := json.Unmarshal(data, &packet)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch data := packet.(type) {
	case *ConnectPacket:
		{
			// fmt.Println(data.Name, "wants to join game ", data.Code)
			game := c.geGameByCode(data.Code)
			if game == nil {
				return
			}
			game.OnPlayerJoin(data.Name, con)
			break
		}
	case *HostNamePacket:
		{
			// fmt.Println("User wants to host quiz", data.QuizId)
			quizId, err := primitive.ObjectIDFromHex(data.QuizId)
			if err != nil {
				fmt.Println("error al convertir el Id", err)
				return
			}
			quiz, err := c.quizService.quizCollection.GetQuizById(quizId)
			if err != nil {
				fmt.Println("Error al obtener quiz:", err)
				return
			}

			if quiz == nil {
				fmt.Println("Quiz es nil aunque no hubo error")
				return
			}

			game := newGame(*quiz, con, c)
			fmt.Println("new game", game.Code)
			c.games = append(c.games, &game)

			c.SendPacket(con, ChangeGameStatePacket{
				State: LobbyState,
			})

			break
		}
	case *StartGamePacket:
		{
			game := c.getGameByHost(con)
			if game == nil {
				return
			}
			game.Start()
			break
		}
	}

}

func (c *NetService) SendPacket(connection *websocket.Conn, packet any) error {
	bytes, err := c.PackageToBytes(packet)
	if err != nil {
		return err
	}

	return connection.WriteMessage(websocket.BinaryMessage, bytes)
}

func (c *NetService) PackageToBytes(packet any) ([]byte, error) {
	packetId, err := c.packetToPacketId(packet)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	final := append([]byte{packetId}, bytes...)
	return final, nil
}
