package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/entity"
	"quiz.com/quiz/internal/game"
)

type NetService struct {
	quizService *QuizService
	games       []*game.Game

	// games []*Game
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*game.Game{},
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
	}
	return nil
}

func (c *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QuestionShowPacket:
		{
			return 2, nil
		}
	}

	return 0, errors.New("invalid packet type")

}

func (c *NetService) geGameByCode(code string) *game.Game {
	for _, game := range c.games {
		if game.Code == code {
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

			newGame := game.New(*quiz, con)
			fmt.Println(newGame.Code, newGame.Id)
			c.games = append(c.games, &newGame)

			// go func() {
			// 	time.Sleep(time.Second * 5)
			// 	c.SendPacket(con, QuestionShowPacket{
			// 		Question: entity.QuizQuestion{
			// 			Name: "Whats is 2+2?",
			// 			Choices: []entity.QuizChoice{
			// 				{
			// 					Name: "4",
			// 				},
			// 				{
			// 					Name: "9",
			// 				},
			// 				{
			// 					Name: "hola",
			// 				},
			// 				{
			// 					Name: "chao",
			// 				},
			// 			},
			// 		},
			// 	})

			// }()
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
