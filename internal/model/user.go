package model

type User struct {
	ID            string `json:"id" bson:"_id"`
	Name          string `json:"name" bson:"name"`
	PictureURL    string `json:"picture_url" bson:"picture_url"`
	StatusMessage string `json:"status_message" bson:"status_message"`
}
