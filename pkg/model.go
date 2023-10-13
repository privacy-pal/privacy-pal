package pal

import (
	"fmt"

	"cloud.google.com/go/firestore"
)

type Locator struct {
	Type LocatorType
	// Collection path leading up to the document in Firestore
	// E.g. "[users]", "[courses,sections]"
	CollectionPath []string
	// List of document IDs in the order of collections
	DocIDs []string
	// List of queries to be applied to a collection. Ignored if type is document.
	Queries []Query
	// An empty struct that will be populated with the data from Firestore
	NewDataNode func() DataNode
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

type DataNode interface {
	HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{}
	// Only one of deleteNode and fieldsToUpdate should be set.
	// If deleteNode is set, fieldsToUpdate will be ignored (node will be deleted).
	HandleDeletion(dataSubjectId string) (nodesToTraverse []Locator, deleteNode bool, fieldsToUpdate []firestore.Update)
}

type Client struct {
	firestoreClient *firestore.Client
}

func NewClient(firestoreClient *firestore.Client) *Client {
	return &Client{firestoreClient: firestoreClient}
}

func validateLocator(loc Locator) error {
	if len(loc.CollectionPath) == 0 {
		return fmt.Errorf("collection path must have at least one element")
	}
	if loc.Type == Document && len(loc.DocIDs) != len(loc.CollectionPath) {
		return fmt.Errorf("document locator must have the same number of docIDs as collection path elements")
	}
	if loc.Type == Collection && len(loc.DocIDs) != len(loc.CollectionPath)-1 {
		return fmt.Errorf("collection locator must have one less docID than collection path elements")
	}
	return nil
}
