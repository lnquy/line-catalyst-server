package model

import "time"

const (
	MessageTypeUser  MessageType = "user"
	MessageTypeGroup MessageType = "group"
)

type (
	MessageType string

	Message struct {
		ID             string      `bson:"_id"`
		Timestamp      time.Time   `bson:"timestamp" json:"timestamp"`
		Type           MessageType `json:"type" bson:"type"`
		MessageID      string      `json:"message_id" bson:"message_id"`
		UserID         string      `json:"user_id" bson:"user_id"`
		GroupID        string      `json:"group_id" bson:"group_id"`
		Text           string      `json:"text" bson:"text"`
		ExpirationDate time.Time   `json:"expiration_date" bson:"expiration_date"`
	}
)
