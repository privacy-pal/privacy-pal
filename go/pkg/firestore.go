package pal

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type firestoreClient struct {
	client *firestore.Client
}

func newDbClientForFirestore(client *firestore.Client) databaseClient {
	return &firestoreClient{client: client}
}

func (c *firestoreClient) getDocument(loc Locator) (DatabaseObject, error) {
	docRef := c.client.Collection(loc.CollectionPath[0]).Doc(loc.DocIDs[0])

	for i := 1; i < len(loc.CollectionPath); i++ {
		docRef = docRef.Collection(loc.CollectionPath[i]).Doc(loc.DocIDs[i])
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
	docRef := c.client.Collection(loc.CollectionPath[0])

	for i := 1; i < len(loc.CollectionPath); i++ {
		docRef = docRef.Doc(loc.DocIDs[i-1]).Collection(loc.CollectionPath[i])
	}

	var query firestore.Query = docRef.Query
	if len(loc.Queries) > 0 {
		query = docRef.Where(loc.Queries[0].Path, loc.Queries[0].Op, loc.Queries[0].Value)
		for i := 1; i < len(loc.Queries); i++ {
			query = query.Where(loc.Queries[i].Path, loc.Queries[i].Op, loc.Queries[i].Value)
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

// func (pal *Client) addDeletionOperationToBatch(batch *firestore.WriteBatch, loc Locator) {
// 	docRef := pal.dbClient.Collection(loc.CollectionPath[0]).Doc(loc.DocIDs[0])

// 	for i := 1; i < len(loc.CollectionPath); i++ {
// 		docRef = docRef.Collection(loc.CollectionPath[i]).Doc(loc.DocIDs[i])
// 	}
// 	batch.Delete(docRef)
// }

// func (pal *Client) addUpdateOperationToBatch(batch *firestore.WriteBatch, loc Locator, fieldsToUpdate []firestore.Update) {
// 	docRef := pal.dbClient.Collection(loc.CollectionPath[0]).Doc(loc.DocIDs[0])

// 	for i := 1; i < len(loc.CollectionPath); i++ {
// 		docRef = docRef.Collection(loc.CollectionPath[i]).Doc(loc.DocIDs[i])
// 	}
// 	batch.Update(docRef, fieldsToUpdate)
// }

// func (pal *Client) commitBatch(batch *firestore.WriteBatch) error {
// 	_, err := batch.Commit(context.Background())
// 	if err != nil {
// 		return fmt.Errorf("%s %w", WRITE_BATCH_ERROR, err)
// 	}
// 	return nil
// }
