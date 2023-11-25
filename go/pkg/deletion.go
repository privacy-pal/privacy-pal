package pal

import "cloud.google.com/go/firestore"

type DocumentUpdates struct {
	Locator        Locator
	FieldsToUpdate []firestore.Update
}

func (pal *Client) ProcessDeletionRequest(handleDeletion HandleDeletionFunc, dataSubjectLocator Locator, dataSubjectID string) string {
	documentsToUpdate, nodesToDelete, err := pal.processDeletionRequest(handleDeletion, dataSubjectLocator, dataSubjectID)
	if err != nil {
		return err.Error()
	}
	pal.dbClient.updateAndDelete(documentsToUpdate, nodesToDelete)
	return ""
}

func (pal *Client) processDeletionRequest(
	handleDelection HandleDeletionFunc,
	locator Locator,
	dataSubjectID string,
) (documentsToUpdate []DocumentUpdates, nodesToDelete []Locator, err error) {
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

	allDocumentsToUpdate := make([]DocumentUpdates, 0)
	allNodesToDelete := make([]Locator, 0)
	for _, currentDataNode := range dataNodes {
		nodesToTraverse, deleteNode, fieldsToUpdate := handleDelection(dataSubjectID, locator, currentDataNode)
		// 1. first recursively process nested nodes
		if len(nodesToTraverse) > 0 {
			for _, nodeLocator := range nodesToTraverse {
				documentsToUpdate, nodesToDelete, err := pal.processDeletionRequest(handleDelection, nodeLocator, dataSubjectID)
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
		} else if fieldsToUpdate != nil {
			allDocumentsToUpdate = append(allDocumentsToUpdate, DocumentUpdates{Locator: locator, FieldsToUpdate: fieldsToUpdate})
		}
	}

	return allDocumentsToUpdate, allNodesToDelete, nil
}
