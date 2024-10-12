package db

import (
	"context"
	"fmt"
	"golang-auth/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	collection *mongo.Collection
}

// List retrieves all users from the database

func (u *UserStore) FindByEmail(email string) (*types.User, error) {
	var user types.User
	err := u.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (u *UserStore) List(ctx context.Context) ([]*types.UserResponse, error) {
	cursor, err := u.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.UserResponse
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err

	}
	return users, nil
}

// Get retrieves a single user by ID and returns added user and error
func (u *UserStore) Get(ctx context.Context, id primitive.ObjectID) (*types.UserResponse, error) {
	var user types.UserResponse

	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {

		return nil, err
	}
	return &user, nil

}

// Create creates a new user and returns added user and error
func (u *UserStore) Create(ctx context.Context, user *types.UserCreate) (*types.UserResponse, error) {
	result, err := u.collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}
	userResponse := types.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		Id:    result.InsertedID.(primitive.ObjectID),
	}

	return &userResponse, nil
}

// Delete deletes a user and returns an deleted object and an error
func (u *UserStore) Delete(ctx context.Context, id primitive.ObjectID) (*types.UserResponse, error) {
	var user *types.UserResponse
	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	_, err = u.collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update modifies an existing user
// func (u *UserStore) Update(ctx context.Context, id primitive.ObjectID, updatedUser types.User) (*types.UserResponse, error) {
func (u *UserStore) Update(ctx context.Context, id primitive.ObjectID, updateData *types.UserUpdate) (*types.UserResponse, error) {
	// Prepare the update document, using $set to update only specific fields
	update := bson.M{
		"$set": updateData,
	}

	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	// Check if any document was matched (i.e., if the user exists)
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no user found")
	}

	var updatedUser types.UserResponse
	err = u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
	// return u.Get(ctx, id)

}
