package pal

type databaseClient interface {
	getDocument(loc Locator) (locatorAndObject, error)
	getDocuments(loc Locator) ([]locatorAndObject, error)
	updateAndDelete(documentsToUpdate []documentUpdates, nodesToDelete []Locator)
}

type DatabaseObject map[string]interface{}

type locatorAndObject struct {
	Locator Locator
	Object  DatabaseObject
}
