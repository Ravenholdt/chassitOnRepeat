package service

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal().Str("tag", "mongo").Err(err).Msg("Can't setup MongoDB connection")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Str("tag", "mongo").Err(err).Msg("Can't connect to MongoDB server")
	}
	return client
}
