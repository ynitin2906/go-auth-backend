package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tasks struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title         string             `json:"title" `
	Category      string             `json:"category"`
	Task          string             `json:"task"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	StatusHistory []*Status          `json:"status_history"`
}

type TasksUpdate struct {
	Title         string    `json:"title" `
	Category      string    `json:"category"`
	Task          string    `json:"task"`
	StatusHistory []*Status `json:"status_history"`
}

type TasksCreate struct {
	Title         string             `json:"title" `
	Category      string             `json:"category"`
	Task          string             `json:"task"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	StatusHistory []*Status          `json:"status_history"`
}
type TasksRequest struct {
	Title    string `json:"title" `
	Category string `json:"category"`
	Task     string `json:"task"`
	Status   string `json:"status"`
}

type Status struct {
	Status string `json:"status"`
	UserId string `json:"user_id" bson:"user_id"`
}
