package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notes struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" `
	Category string             `json:"category"`
	Note     string             `json:"note"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type NotesUpdate struct {
	Title    string `json:"title" `
	Category string `json:"category"`
	Note     string `json:"note"`
}

type NotesCreate struct {
	Title    string             `json:"title" `
	Category string             `json:"category"`
	Note     string             `json:"note"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
}
type NotesRequest struct {
	Title    string `json:"title" `
	Category string `json:"category"`
	Note     string `json:"note"`
}
