package pal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoClient struct {
	db *mongo.Database
}

func newDbClientForMongo(mongoDb *mongo.Database) databaseClient {
	return &mongoClient{db: mongoDb}
}

func (c *mongoClient) getDocument(loc Locator) (DatabaseObject, error) {
	// Get a single result based on the collection and filter supplied in the locator
	collection := c.db.Collection(loc.MongoLocator.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonResult := bson.M{}
	if err := collection.FindOne(ctx, loc.MongoLocator.Filter).Decode(&bsonResult); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s document does not exist", GET_DOCUMENT_ERROR)
		}
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	// Convert bson.M to DatabaseObject
	tempBytes, err := bson.MarshalExtJSON(bsonResult, true, true)
	if err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	result := DatabaseObject{}
	err = json.Unmarshal(tempBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}
	result["_id"] = bsonResult["_id"].(primitive.ObjectID).Hex()

	return result, nil
}

func (c *mongoClient) getDocuments(loc Locator) ([]DatabaseObject, error) {
	// Get a list of results based on the collection and filter supplied in the locator
	collection := c.db.Collection(loc.MongoLocator.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, loc.MongoLocator.Filter)
	if err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	bsonResults := []bson.M{}
	if err = cursor.All(ctx, &bsonResults); err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	// Convert bson.M to DatabaseObject
	results := []DatabaseObject{}
	for _, result := range bsonResults {
		tempBytes, err := bson.MarshalExtJSON(result, true, true)
		if err != nil {
			return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
		}
		convertedResult := DatabaseObject{}
		err = json.Unmarshal(tempBytes, &convertedResult)
		if err != nil {
			return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
		}
		convertedResult["_id"] = result["_id"].(primitive.ObjectID).Hex()
		results = append(results, convertedResult)
	}

	return results, nil
}
