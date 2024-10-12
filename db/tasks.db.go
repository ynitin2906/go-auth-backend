package db

import (
	"context"
	"fmt"
	"golang-auth/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TasksStore struct {
	collection *mongo.Collection
}

// List retrieves all tasks for a specific user
func (n *TasksStore) List(ctx context.Context, userId primitive.ObjectID) ([]*types.Tasks, error) {
	filter := bson.M{"user_id": userId}
	cursor, err := n.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var tasks []*types.Tasks
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Get retrieves a single task by ID
func (n *TasksStore) Get(ctx context.Context, id primitive.ObjectID) (*types.Tasks, error) {
	var task *types.Tasks
	filter := bson.M{"_id": id}
	err := n.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Create inserts a new task into the database
func (n *TasksStore) Create(ctx context.Context, task *types.TasksCreate) (*types.Tasks, error) {
	result, err := n.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	newTask := types.Tasks{
		Id:            result.InsertedID.(primitive.ObjectID),
		Title:         task.Title,
		Category:      task.Category,
		Task:          task.Task,
		UserID:        task.UserID,
		StatusHistory: task.StatusHistory,
	}

	return &newTask, nil
}

// Delete removes a task by ID and returns the deleted task
func (n *TasksStore) Delete(ctx context.Context, id primitive.ObjectID) (*types.Tasks, error) {
	var deletedTask *types.Tasks
	err := n.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&deletedTask)
	if err != nil {
		return nil, err
	}
	_, err = n.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return deletedTask, nil
}

// Update modifies an existing task based on its ID
func (n *TasksStore) Update(ctx context.Context, id primitive.ObjectID, updatedData *types.TasksUpdate) (*types.Tasks, error) {
	update := bson.M{
		"$set": updatedData,
	}
	result, err := n.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no task found")
	}

	// Fetch and return the updated task
	var updatedTask *types.Tasks
	err = n.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedTask)
	if err != nil {
		return nil, err
	}
	return updatedTask, nil
}
