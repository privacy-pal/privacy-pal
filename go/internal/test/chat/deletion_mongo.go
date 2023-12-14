package chat

import (
	"fmt"
	"log"

	pal "github.com/privacy-pal/privacy-pal/go/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleDeletionMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleDeletionUserMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleDeletionGroupChatMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleDeletionMessageMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleDeletionDirectMessageMongo(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		err = fmt.Errorf("invalid data type: %s", currentDbObjLocator.DataType)
		return
	}
}

func handleDeletionUserMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	gcs, ok := dbObj["gcs"].([]interface{})
	if !ok {
		err = fmt.Errorf("invalid gcs field")
		return
	}
	for _, id := range gcs {
		id, ok := id.(string)
		if !ok {
			err = fmt.Errorf("invalid gcs field")
			return
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println(err)
		}

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    GroupChatDataType,
			MongoLocator: pal.MongoLocator{
				Collection: "gcs",
				Filter:     bson.D{{Key: "_id", Value: objectID}},
			},
		})
	}

	dms, ok := dbObj["dms"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("invalid dms field")
		return
	}
	for _, dmID := range dms {
		dmID, ok := dmID.(string)
		if !ok {
			err = fmt.Errorf("invalid dms field")
			return
		}
		dmIDObj, err := primitive.ObjectIDFromHex(dmID)
		if err != nil {
			log.Println(err)
		}

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    DirectMessageDataType,
			MongoLocator: pal.MongoLocator{
				Collection: "dms",
				Filter:     bson.D{{Key: "_id", Value: dmIDObj}},
			},
		})
	}

	deleteNode = true
	return
}

func handleDeletionGroupChatMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	nodesToTraverse = append(nodesToTraverse, pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
		MongoLocator: pal.MongoLocator{
			Collection: "messages",
			Filter:     bson.D{{Key: "userId", Value: dataSubjectId}, {Key: "chatId", Value: dbObj["_id"]}},
		},
	})

	deleteNode = false

	mongoUpdates := []interface{}{}

	users, ok := dbObj["users"].([]interface{})
	if !ok {
		err = fmt.Errorf("invalid users field")
		return
	}
	if (dbObj["owner"] == dataSubjectId) && (len(users) > 1) {
		mongoUpdates = append(mongoUpdates, bson.D{{Key: "owner", Value: users[0]}})
	} else {
		mongoUpdates = append(mongoUpdates, bson.D{{Key: "$pull", Value: bson.D{{Key: "users", Value: dataSubjectId}}}})
	}

	fieldsToUpdate = pal.FieldUpdates{
		MongoUpdates: mongoUpdates,
	}

	return
}

func handleDeletionMessageMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	deleteNode = true
	return
}

func handleDeletionDirectMessageMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	nodesToTraverse = append(nodesToTraverse, pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
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
