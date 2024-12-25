package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"quiz.com/quiz/internal/entity"
)

type NetService struct {
	quizService *QuizService
	host        *websocket.Conn

	tick int

	// games []*Game
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		// games:       []*Game{},
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

	switch packet := packet.(type) {
	case *ConnectPacket:
		{
			fmt.Println(packet.Name, "wants to join game ", packet.Code)
			break
		}
	case *HostNamePacket:
		{
			fmt.Println("User wants to host quiz", packet.QuizId)
			break
		}
	}

}

func (c *NetService) SendPackage(connection *websocket.Conn, packet any) error {
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
