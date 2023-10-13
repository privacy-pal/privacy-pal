package pal

import (
	"fmt"

	"cloud.google.com/go/firestore"
)

type DeletionResult struct {
	NumDocumentsDeleted int
	NumDocumentsUpdated int
}

func (pal *Client) ProcessDeletionRequest(dataSubjectLocator Locator, dataSubjectID string) (DeletionResult, error) {
	fmt.Printf("Processing deletion request for data subject %s\n", dataSubjectID)
	// TODO: check data subject is valid
	batch := pal.firestoreClient.Batch()

	delResult := DeletionResult{}
	err := pal.processDeletionRequest(dataSubjectLocator, dataSubjectID, batch, &delResult)
	if err != nil {
		return DeletionResult{}, fmt.Errorf("%s %w", DELETION_REQUEST_ERROR, err)
	}

	err = pal.commitBatch(batch)
	if err != nil {
		return DeletionResult{}, fmt.Errorf("%s %w", DELETION_REQUEST_ERROR, err)
	}

	fmt.Printf("Deleted %d documents\n", delResult.NumDocumentsDeleted)
	fmt.Printf("Updated %d documents\n", delResult.NumDocumentsUpdated)

	return delResult, nil
}

func (pal *Client) processDeletionRequest(
	locator Locator, dataSubjectID string, batch *firestore.WriteBatch, delResult *DeletionResult,
) error {

	dataNode, err := pal.getDocumentFromFirestore(locator)
	if err != nil {
		return err
	}
	nodesToTraverse, deleteNode, fieldsToUpdate := dataNode.HandleDeletion(dataSubjectID)

	// 1. first recursively process nested nodes
	if len(nodesToTraverse) > 0 {
		for _, loc := range nodesToTraverse {
			pal.processDeletionRequest(loc, dataSubjectID, batch, delResult)
		}
	}

	// 2. delete current node if needed
	if deleteNode {
		pal.addDeletionOperationToBatch(
			batch,
			Locator{
				CollectionPath: locator.CollectionPath,
				DocIDs:         locator.DocIDs,
			},
		)
		delResult.NumDocumentsDeleted++

	} else if len(fieldsToUpdate) > 0 {
		// 3. override fields if only nested data needs to be deleted
		pal.addUpdateOperationToBatch(
			batch,
			Locator{
				CollectionPath: locator.CollectionPath,
				DocIDs:         locator.DocIDs,
			},
			fieldsToUpdate,
		)
		delResult.NumDocumentsUpdated++
	}

	return nil
}
