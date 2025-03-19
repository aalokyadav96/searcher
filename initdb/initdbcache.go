package initdb

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Package-level context.
var ctx = context.Background()
var CTX = ctx

var RedisClient *redis.Client
var Client *mongo.Client // Global MongoDB client

// init initializes Redis and MongoDB clients.
func init() {
	// Initialize Redis client.
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Update if needed.
	})
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize MongoDB client.
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err = Client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
}
