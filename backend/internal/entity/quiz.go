package entity

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Quiz struct {
	Id        bson.ObjectID  `json:"id" bson:"_id"` // Asegúrate de que esté como ObjectID
	Name      string         `json:"name"`
	Questions []QuizQuestion `json:"questions"`
}

type QuizQuestion struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Time    int          `json:"time"`
	Choices []QuizChoice `json:"choices"`
}

type QuizChoice struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}
