package pal

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

type firestoreClient struct {
	client *firestore.Client
}

func newDbClientForFirestore(client *firestore.Client) databaseClient {
	return &firestoreClient{client: client}
}

func (c *firestoreClient) getDocument(loc Locator) (DatabaseObject, error) {
	docRef := c.client.Collection(loc.FirestoreLocator.CollectionPath[0]).Doc(loc.DocIDs[0])

	for i := 1; i < len(loc.FirestoreLocator.CollectionPath); i++ {
		docRef = docRef.Collection(loc.FirestoreLocator.CollectionPath[i]).Doc(loc.DocIDs[i])
	}

	doc, err := docRef.Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}
	if !doc.Exists() {
		return nil, fmt.Errorf("%s document does not exist", GET_DOCUMENT_ERROR)
	}

	data := doc.Data()
	data["_id"] = doc.Ref.ID
	return data, nil
}

func (c *firestoreClient) getDocuments(loc Locator) ([]DatabaseObject, error) {
	docRef := c.client.Collection(loc.FirestoreLocator.CollectionPath[0])

	for i := 1; i < len(loc.FirestoreLocator.CollectionPath); i++ {
		docRef = docRef.Doc(loc.DocIDs[i-1]).Collection(loc.FirestoreLocator.CollectionPath[i])
	}

	var query firestore.Query = docRef.Query
	if len(loc.Filters) > 0 {
		query = docRef.Where(loc.Filters[0].Path, loc.Filters[0].Op, loc.Filters[0].Value)
		for i := 1; i < len(loc.Filters); i++ {
			query = query.Where(loc.Filters[i].Path, loc.Filters[i].Op, loc.Filters[i].Value)
		}
	}

	doc, err := query.Documents(context.Background()).GetAll()

	if err != nil {
		return nil, fmt.Errorf("%s %w", GET_DOCUMENT_ERROR, err)
	}

	dataNodes := make([]DatabaseObject, len(doc))
	for i, d := range doc {
		data := d.Data()
		data["_id"] = d.Ref.ID
		dataNodes[i] = data
	}
	return dataNodes, nil
}

func (c *firestoreClient) updateAndDelete(documentsToUpdate []documentUpdates, nodesToDelete []Locator) {
	err := c.client.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {
		// delete nodes
		for _, nodeLocator := range nodesToDelete {
			docRef := c.client.Collection(nodeLocator.FirestoreLocator.CollectionPath[0]).Doc(nodeLocator.DocIDs[0])

			for i := 1; i < len(nodeLocator.FirestoreLocator.CollectionPath); i++ {
				docRef = docRef.Collection(nodeLocator.FirestoreLocator.CollectionPath[i]).Doc(nodeLocator.DocIDs[i])
			}

			err := t.Delete(docRef)
			if err != nil {
				return err
			}
		}

		// update nodes
		for _, update := range documentsToUpdate {
			docRef := c.client.Collection(update.Locator.FirestoreLocator.CollectionPath[0]).Doc(update.Locator.DocIDs[0])

			for i := 1; i < len(update.Locator.FirestoreLocator.CollectionPath); i++ {
				docRef = docRef.Collection(update.Locator.FirestoreLocator.CollectionPath[i]).Doc(update.Locator.DocIDs[i])
			}

			err := t.Update(docRef, update.FieldsToUpdate.FirestoreUpdates)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error updating and deleting data: %v", err)
	}
}
