package pal

import (
	"encoding/json"

	"cloud.google.com/go/firestore"
)

type documentUpdates struct {
	Locator        Locator
	FieldsToUpdate FieldUpdates
}

type FieldUpdates struct {
	FirestoreUpdates []firestore.Update
	MongoUpdates     []interface{}
}

func (pal *Client) ProcessDeletionRequest(handleDeletion HandleDeletionFunc, dataSubjectLocator Locator, dataSubjectID string, writeToDatabase bool) (string, error) {
	documentsToUpdate, nodesToDelete, err := pal.processDeletionRequest(handleDeletion, dataSubjectLocator, dataSubjectID)
	if err != nil {
		return "", err
	}
	if writeToDatabase {
		pal.dbClient.updateAndDelete(documentsToUpdate, nodesToDelete)
	}

	result, err := json.Marshal(map[string]interface{}{
		"writeToDatabase":   writeToDatabase,
		"nodesToDelete":     nodesToDelete,
		"documentsToUpdate": documentsToUpdate,
	})
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (pal *Client) processDeletionRequest(
	handleDeletion HandleDeletionFunc,
	locator Locator,
	dataSubjectID string,
) (documentsToUpdate []documentUpdates, nodesToDelete []Locator, err error) {
	dataNodes := make([]DatabaseObject, 0)
	if locator.LocatorType == Document {
		node, err := pal.dbClient.getDocument(locator)
		if err != nil {
			return nil, nil, err
		}
		dataNodes = append(dataNodes, node)
	} else {
		nodes, err := pal.dbClient.getDocuments(locator)
		if err != nil {
			return nil, nil, err
		}
		dataNodes = append(dataNodes, nodes...)
	}

	allDocumentsToUpdate := make([]documentUpdates, 0)
	allNodesToDelete := make([]Locator, 0)
	for _, currentDataNode := range dataNodes {
		nodesToTraverse, deleteNode, fieldsToUpdate, err := handleDeletion(dataSubjectID, locator, currentDataNode)
		if err != nil {
			return nil, nil, err
		}

		// 1. first recursively process nested nodes
		if len(nodesToTraverse) > 0 {
			for _, nodeLocator := range nodesToTraverse {
				documentsToUpdate, nodesToDelete, err := pal.processDeletionRequest(handleDeletion, nodeLocator, dataSubjectID)
				if err != nil {
					return nil, nil, err
				}
				allDocumentsToUpdate = append(allDocumentsToUpdate, documentsToUpdate...)
				allNodesToDelete = append(allNodesToDelete, nodesToDelete...)
			}
		}

		// 2. delete current node if needed
		if deleteNode {
			allNodesToDelete = append(allNodesToDelete, locator)
		} else {
			allDocumentsToUpdate = append(allDocumentsToUpdate, documentUpdates{Locator: locator, FieldsToUpdate: fieldsToUpdate})
		}
	}

	return allDocumentsToUpdate, allNodesToDelete, nil
}
