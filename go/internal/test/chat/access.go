package chat

import (
	"log"

	pal "github.com/privacy-pal/privacy-pal/go/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserDataType          = "user"
	GroupChatDataType     = "groupchat"
	DirectMessageDataType = "directmessage"
	MessageDataType       = "message"
)

func HandleAccess(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleAccessUser(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleAccessGroupChat(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleAccessMessage(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleAccessDirectMessage(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		// TODO: should return error
		return nil
	}
}

func handleAccessUser(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	// TODO: include in documentation: you can access the id in 2 ways
	if dbObj["_id"].(string) != dataSubjectId {
		data["Name"] = dbObj["name"]
		return data
	}
	// if currentDbObjLocator.DocIDs[len(currentDbObjLocator.DocIDs)-1] != dataSubjectId {
	// 	data["Name"] = dbObj["name"]
	// 	return data
	// }

	data["Name"] = dbObj["name"]
	data["Groupchats"] = make([]pal.Locator, 0)
	for _, id := range dbObj["gcs"].([]interface{}) {
		id := id.(string)
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println(err)
		}

		data["Groupchats"] = append(data["Groupchats"].([]pal.Locator), pal.Locator{
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
	data["DirectMessages"] = make([]pal.Locator, 0)
	// TODO: support this in yaml
	for _, dmID := range dbObj["dms"].(map[string]interface{}) {
		dmID := dmID.(string)
		dmIDObj, err := primitive.ObjectIDFromHex(dmID)
		if err != nil {
			log.Println(err)
		}

		data["DirectMessages"] = append(data["DirectMessages"].([]pal.Locator), pal.Locator{
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

	return data
}

func handleAccessGroupChat(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Messages"] = pal.Locator{
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
	}

	return data
}

func handleAccessMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Content"] = dbObj["content"]
	data["Timestamp"] = dbObj["timestamp"]

	return data
}

func handleAccessDirectMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	var otherUserId string
	if dbObj["user1"].(string) == dataSubjectId {
		otherUserId = dbObj["user2"].(string)
	} else {
		otherUserId = dbObj["user1"].(string)
	}
	otherUserIDObj, err := primitive.ObjectIDFromHex(otherUserId)
	if err != nil {
		log.Println(err)
	}

	data["Other User"] = pal.Locator{
		LocatorType: pal.Document,
		DataType:    UserDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: []string{"users"},
			DocIDs:         []string{otherUserId},
		},
		MongoLocator: pal.MongoLocator{
			Collection: "users",
			Filter:     bson.D{{Key: "_id", Value: otherUserIDObj}},
		},
	}
	data["Messages"] = pal.Locator{
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
	}

	return data
}
