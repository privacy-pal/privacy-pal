package pal

type DatabaseClient interface {
	getDocument(loc Locator) (DatabaseObject, error)
	getDocuments(loc Locator) ([]DatabaseObject, error)
}

type DatabaseObject map[string]interface{}
