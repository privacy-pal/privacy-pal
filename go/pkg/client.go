package pal

import (
	"cloud.google.com/go/firestore"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandleAccessFunc func(dataSubjectId string, currentDbObjLocator Locator, dbObj DatabaseObject) map[string]interface{}

// Only one of deleteNode and fieldsToUpdate should be set.
// If deleteNode is set, fieldsToUpdate will be ignored (node will be deleted).
type HandleDeletionFunc func(dataSubjectId string, currentDbObjLocator Locator, dbObj DatabaseObject) (nodesToTraverse []Locator, deleteNode bool, fieldsToUpdate FieldUpdates)

type Client struct {
	dbClient databaseClient
}

func NewClientWithFirestore(firestoreClient *firestore.Client) *Client {
	return &Client{dbClient: newDbClientForFirestore(firestoreClient)}
}

func NewClientWithMongo(mongoDb *mongo.Database) *Client {
	return &Client{dbClient: newDbClientForMongo(mongoDb)}
}
