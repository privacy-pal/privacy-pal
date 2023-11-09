package pal

type databaseClient interface {
	getDocument(loc Locator) (DatabaseObject, error)
	getDocuments(loc Locator) ([]DatabaseObject, error)
}

type DatabaseObject map[string]interface{}
