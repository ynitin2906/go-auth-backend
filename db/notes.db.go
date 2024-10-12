package db

import (
	"context"
	"fmt"
	"golang-auth/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotesStore struct {
	collection *mongo.Collection
}

func (n *NotesStore) List(ctx context.Context, userId primitive.ObjectID) ([]*types.Notes, error) {
	filter := bson.M{"user_id": userId}
	cursor, err := n.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var notes []*types.Notes
	err = cursor.All(ctx, &notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (n *NotesStore) Get(ctx context.Context, id primitive.ObjectID) (*types.Notes, error) {
	var note *types.Notes
	filter := bson.M{"_id": id}
	err := n.collection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (n *NotesStore) Create(ctx context.Context, note *types.NotesCreate) (*types.Notes, error) {
	result, err := n.collection.InsertOne(ctx, note)
	if err != nil {
		return nil, err
	}
	newNote := types.Notes{
		Title:    note.Title,
		Category: note.Category,
		Note:     note.Note,
		UserID:   note.UserID,
		Id:       result.InsertedID.(primitive.ObjectID),
	}

	return &newNote, nil
}

func (n *NotesStore) Delete(ctx context.Context, id primitive.ObjectID) (*types.Notes, error) {
	var deletedNote *types.Notes
	err := n.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&deletedNote)
	if err != nil {
		return nil, err
	}
	_, err = n.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return deletedNote, nil
}

func (n *NotesStore) Update(ctx context.Context, id primitive.ObjectID, updatedData *types.NotesUpdate) (*types.Notes, error) {
	update := bson.M{
		"$set": updatedData,
	}
	result, err := n.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no note found")
	}

	var updatedNote *types.Notes
	err = n.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedNote)
	if err != nil {
		return nil, err
	}
	return updatedNote, nil
}
