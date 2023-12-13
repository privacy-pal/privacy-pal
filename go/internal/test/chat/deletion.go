package chat

import (
	"log"

	"cloud.google.com/go/firestore"
	pal "github.com/privacy-pal/privacy-pal/go/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleDeletion(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates) {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleDeletionUser(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleDeletionGroupChat(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleDeletionMessage(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleDeletionDirectMessage(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		// TODO: should return error
		return nil, false, pal.FieldUpdates{}
	}
}

func handleDeletionUser(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates) {
	for _, id := range dbObj["gcs"].([]interface{}) {
		id := id.(string)
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println(err)
		}

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    GroupChatDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"gcs"},
				DocIDs:         []string{id},
			},
			MongoLocator: pal.MongoLocator{
				Collection: "gcs",
				Filter:     bson.D{{Key: "_id", Value: objectID}},
			},
		})
	}

	for _, dmID := range dbObj["dms"].(map[string]interface{}) {
		dmID := dmID.(string)
		dmIDObj, err := primitive.ObjectIDFromHex(dmID)
		if err != nil {
			log.Println(err)
		}

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    DirectMessageDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"dms"},
				DocIDs:         []string{dmID},
			},
			MongoLocator: pal.MongoLocator{
				Collection: "dms",
				Filter:     bson.D{{Key: "_id", Value: dmIDObj}},
			},
		})
	}

	deleteNode = true
	return
}

func handleDeletionGroupChat(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates) {
	nodesToTraverse = append(nodesToTraverse, pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: append(currentDbObjLocator.FirestoreLocator.CollectionPath, "messages"),
			DocIDs:         currentDbObjLocator.DocIDs,
			Filters: []pal.Filter{
				{
					Path:  "userId",
					Op:    "==",
					Value: dataSubjectId,
				},
			},
		},
		MongoLocator: pal.MongoLocator{
			Collection: "messages",
			Filter:     bson.D{{Key: "userId", Value: dataSubjectId}, {Key: "chatId", Value: dbObj["_id"]}},
		},
	})

	deleteNode = false

	firestoreUpdates := []firestore.Update{}
	mongoUpdates := []interface{}{}

	if (dbObj["owner"] == dataSubjectId) && (len(dbObj["users"].([]interface{})) > 1) {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  "owner",
			Value: dbObj["users"].([]interface{})[0],
		})
		mongoUpdates = append(mongoUpdates, bson.D{{Key: "owner", Value: dbObj["users"].([]interface{})[0]}})
	} else {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  "users",
			Value: firestore.ArrayRemove(dataSubjectId),
		})
		mongoUpdates = append(mongoUpdates, bson.D{{Key: "$pull", Value: bson.D{{Key: "users", Value: dataSubjectId}}}})
	}

	fieldsToUpdate = pal.FieldUpdates{
		FirestoreUpdates: firestoreUpdates,
		MongoUpdates:     mongoUpdates,
	}

	return
}

func handleDeletionMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates) {
	deleteNode = true
	return
}

func handleDeletionDirectMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates) {
	nodesToTraverse = append(nodesToTraverse, pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: append(currentDbObjLocator.FirestoreLocator.CollectionPath, "messages"),
			DocIDs:         currentDbObjLocator.DocIDs,
			Filters: []pal.Filter{
				{
					Path:  "userId",
					Op:    "==",
					Value: dataSubjectId,
				},
			},
		},
		MongoLocator: pal.MongoLocator{
			Collection: "messages",
			Filter:     bson.D{{Key: "userId", Value: dataSubjectId}, {Key: "chatId", Value: dbObj["_id"]}},
		},
	})

	deleteNode = false

	fieldsToUpdate = pal.FieldUpdates{
		MongoUpdates: []interface{}{
			bson.D{{Key: "$pull", Value: bson.D{{Key: "$pull", Value: bson.D{{Key: "messages", Value: bson.D{{Key: "userId", Value: dataSubjectId}, {Key: "chatId", Value: dbObj["_id"]}}}}}}}},
		},
	}

	return
}
