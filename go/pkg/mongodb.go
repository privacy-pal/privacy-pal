package pal

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoClient struct {
	Client *mongo.Client
}

func newDbClientForMongo(client *mongo.Client) DatabaseClient {
	return &mongoClient{Client: client}
}

func (c *mongoClient) getDocument(loc Locator) (DatabaseObject, error) {
	return nil, nil
}

func (c *mongoClient) getDocuments(loc Locator) ([]DatabaseObject, error) {
	return nil, nil
}
