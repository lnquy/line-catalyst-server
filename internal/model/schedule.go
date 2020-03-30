package model

type Schedule struct {
	Name    string `json:"name" bson:"_id"`
	Cron    string `json:"cron" bson:"cron"`
	Message string `json:"message" bson:"message"`
	ReplyTo string `json:"reply_to" bson:"reply_to"`
	IsDone  bool   `json:"is_done" bson:"is_done"`
}
