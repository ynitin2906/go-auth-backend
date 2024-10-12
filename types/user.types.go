package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" `
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	Role           string             `json:"role"`
	Notes          []*Notes           `json:"notes"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	SocialMedia    SocialMedia        `json:"social_media"    bson:"social_media"`
}
type UserUpdate struct {
	Name           string      `json:"name" `
	Email          string      `json:"email"`
	ProfilePicture string      `json:"profile_picture" bson:"profile_picture"`
	SocialMedia    SocialMedia `json:"social_media"    bson:"social_media"`
}
type UserResponse struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" `
	Email          string             `json:"email"`
	Notes          []*Notes           `json:"notes"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	SocialMedia    SocialMedia        `json:"social_media"    bson:"social_media"`
}
type UserRequest struct {
	Name     string `json:"name" `
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserCreate struct {
	Name     string `json:"name" `
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	// ProfilePicture string      `json:"profile_picture" bson:"profile_picture"`
	// SocialMedia    SocialMedia `json:"social_media"    bson:"social_media"`
}

type SocialMedia struct {
	Twitter   string `json:"twitter"   bson:"twitter"`
	LinkedIn  string `json:"linkedin"  bson:"linkedin"`
	GitHub    string `json:"github"    bson:"github"`
	Facebook  string `json:"facebook"  bson:"facebook"`
	Instagram string `json:"instagram" bson:"instagram"`
	Snapchat  string `json:"snapchat"  bson:"snapchat"`
	YouTube   string `json:"youtube"   bson:"youtube"`
	Pinterest string `json:"pinterest" bson:"pinterest"`
	Discord   string `json:"discord"   bson:"discord"`
	Website   string `json:"website"   bson:"website"`
}
