package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	User  UserStore
	Notes NotesStore
	Tasks TasksStore
}

// NewStore initializes the DB connection and returns a new Store
func NewStore() *Store {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URL")
	if mongoURI == "" {
		log.Fatal("MONGO_URL not set")
	}

	// Create a new context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB directly
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	userCollection := client.Database("go-lang-auth-db").Collection("user")
	notesCollection := client.Database("go-lang-auth-db").Collection("note")
	tasksCollection := client.Database("go-lang-auth-db").Collection("task")

	// Return the store containing the UserStore
	return &Store{
		User: UserStore{
			collection: userCollection,
		},
		Notes: NotesStore{
			collection: notesCollection,
		},
		Tasks: TasksStore{
			collection: tasksCollection,
		},
	}
}
