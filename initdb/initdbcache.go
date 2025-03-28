package initdb

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Package-level context.
var CTX = context.Background()

var MongoClient *mongo.Client // Global MongoDB client

var rxdurl string = os.Getenv("REDIS_URL")
var rxdopts, _ = redis.ParseURL(rxdurl)
var RedisClient = redis.NewClient(rxdopts)

// init initializes MongoDB client.
func init() {
	errx := godotenv.Load()
	if errx != nil {
		log.Fatal("Error loading .env file")
	}

	var mongo_url = os.Getenv("MONGODB_URI")

	if _, err := RedisClient.Ping(CTX).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize MongoDB client.
	var err error
	MongoClient, err = mongo.Connect(CTX, options.Client().ApplyURI(mongo_url))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err = MongoClient.Ping(CTX, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
}
