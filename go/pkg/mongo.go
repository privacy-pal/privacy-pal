package pal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func (c *mongoClient) getDocument(loc Locator) (locatorAndObject, error) {
	// Get a single result based on the collection and filter supplied in the locator
	collection := c.db.Collection(loc.MongoLocator.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonResult := bson.M{}
	if err := collection.FindOne(ctx, loc.MongoLocator.Filter).Decode(&bsonResult); err != nil {
		if err == mongo.ErrNoDocuments {
			return locatorAndObject{}, fmt.Errorf("%s document does not exist", GET_DOCUMENT_ERROR)
		}
		return locatorAndObject{}, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	// Convert bson.M to DatabaseObject
	tempBytes, err := bson.MarshalExtJSON(bsonResult, true, true)
	if err != nil {
		return locatorAndObject{}, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	result := DatabaseObject{}
	err = json.Unmarshal(tempBytes, &result)
	if err != nil {
		return locatorAndObject{}, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}
	result["_id"] = bsonResult["_id"].(primitive.ObjectID).Hex()

	loc.LocatorType = Document
	return locatorAndObject{Locator: loc, Object: result}, nil
}

func (c *mongoClient) getDocuments(loc Locator) ([]locatorAndObject, error) {
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

	loc.LocatorType = Document
	// Convert bson.M to DatabaseObject
	results := []locatorAndObject{}
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
		results = append(results, locatorAndObject{Locator: loc, Object: convertedResult})
	}

	return results, nil
}

func (c *mongoClient) updateAndDelete(documentsToUpdate []documentUpdates, nodesToDelete []Locator) {
	session, err := c.db.Client().StartSession()
	if err != nil {
		log.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		// delete nodes
		for _, nodeLocator := range nodesToDelete {
			collection := c.db.Collection(nodeLocator.MongoLocator.Collection)
			_, err := collection.DeleteOne(sessionContext, nodeLocator.MongoLocator.Filter)
			if err != nil {
				return nil, err
			}
		}

		// update nodes
		for _, update := range documentsToUpdate {
			collection := c.db.Collection(update.Locator.MongoLocator.Collection)
			_, err := collection.UpdateOne(sessionContext, update.Locator.MongoLocator.Filter, update.FieldsToUpdate)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback)
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}
}
