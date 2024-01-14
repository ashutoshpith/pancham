package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbConnectionPayload struct {
	Url string
}

func (dbConnPayload DbConnectionPayload) onConnectionEstablish() {
	log.Println("Connected to MongoDB!")
}

func (dbConnPayload DbConnectionPayload) onConnectionFailed(err error) {
	log.Fatal("Connection Failed to MongoDB!", err)
}

type DbConnection struct {
	DbConnectionPayload
}

func Connect(dbConn DbConnection) *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(dbConn.Url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		dbConn.onConnectionFailed(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		dbConn.onConnectionFailed(err)
	}

	dbConn.onConnectionEstablish()

	return client
}
