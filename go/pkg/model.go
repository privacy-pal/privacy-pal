package pal

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"go.mongodb.org/mongo-driver/mongo"
)

type Locator struct {
	LocatorType LocatorType
	DataType    string
	// Collection path leading up to the document in Firestore
	// E.g. "[users]", "[courses,sections]"
	CollectionPath []string
	// List of document IDs in the order of collections
	DocIDs []string
	// List of queries to be applied to a collection. Ignored if type is document.
	Queries []Query
}

type Query struct {
	Path  string
	Op    string
	Value interface{}
}

type LocatorType string

const (
	Document   LocatorType = "document"
	Collection LocatorType = "collection"
)

type HandleAccessFunc func(dataSubjectId string, currentDbObjLocator Locator, dbObj DatabaseObject) map[string]interface{}

// Only one of deleteNode and fieldsToUpdate should be set.
// If deleteNode is set, fieldsToUpdate will be ignored (node will be deleted).
type HandleDeletionFunc func(dataSubjectId string, obj interface{}) (nodesToTraverse []Locator, deleteNode bool, fieldsToUpdate []firestore.Update)

type Client struct {
	FirestoreClient *firestore.Client
	DbClient        DatabaseClient
}

func NewClientWithFirestore(firestoreClient *firestore.Client) *Client {
	return &Client{DbClient: newDbClientForFirestore(firestoreClient)}
}

func NewClientWithMongo(mongoClient *mongo.Client) *Client {
	return &Client{DbClient: newDbClientForMongo(mongoClient)}
}

func validateLocator(loc Locator) error {
	if len(loc.CollectionPath) == 0 {
		return fmt.Errorf("collection path must have at least one element")
	}
	if loc.LocatorType == Document && len(loc.DocIDs) != len(loc.CollectionPath) {
		return fmt.Errorf("document locator must have the same number of docIDs as collection path elements")
	}
	if loc.LocatorType == Collection && len(loc.DocIDs) != len(loc.CollectionPath)-1 {
		return fmt.Errorf("collection locator must have one less docID than collection path elements")
	}
	return nil
}
