package pal

type databaseClient interface {
	getDocument(loc Locator) (DatabaseObject, error)
	getDocuments(loc Locator) ([]DatabaseObject, error)
	updateAndDelete(documentsToUpdate []DocumentUpdates, nodesToDelete []Locator)
}

type DatabaseObject map[string]interface{}
