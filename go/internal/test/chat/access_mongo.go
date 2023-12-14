package chat

import (
	"fmt"

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

func HandleAccessMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleAccessUserMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleAccessGroupChatMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleAccessMessageMongo(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleAccessDirectMessageMongo(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		err = fmt.Errorf("invalid data type")
		return
	}
}

func handleAccessUserMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	id, ok := dbObj["_id"].(string)
	if !ok {
		err = fmt.Errorf("invalid id")
		return
	}
	if id != dataSubjectId {
		data["Name"] = dbObj["name"]
		return
	}
	// if currentDbObjLocator.DocIDs[len(currentDbObjLocator.DocIDs)-1] != dataSubjectId {
	// 	data["Name"] = dbObj["name"]
	// 	return
	// }

	data["Name"] = dbObj["name"]
	groupchatLocators := make([]pal.Locator, 0)
	gcs, ok := dbObj["gcs"].([]interface{})
	if !ok {
		err = fmt.Errorf("invalid gcs")
		return
	}
	for _, id := range gcs {
		id, ok := id.(string)
		if !ok {
			err = fmt.Errorf("invalid id")
			return
		}
		var objectID primitive.ObjectID
		objectID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return
		}

		groupchatLocators = append(groupchatLocators, pal.Locator{
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
	data["Groupchats"] = groupchatLocators

	directMessageLocators := make([]pal.Locator, 0)
	dms, ok := dbObj["dms"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("invalid dms")
		return
	}
	for _, dmID := range dms {
		dmID, ok := dmID.(string)
		if !ok {
			err = fmt.Errorf("invalid id")
			return
		}
		var dmIDObj primitive.ObjectID
		dmIDObj, err = primitive.ObjectIDFromHex(dmID)
		if err != nil {
			return
		}

		directMessageLocators = append(directMessageLocators, pal.Locator{
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
	data["DirectMessages"] = directMessageLocators

	return
}

func handleAccessGroupChatMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

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

	return
}

func handleAccessMessageMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	data["Content"] = dbObj["content"]
	data["Timestamp"] = dbObj["timestamp"]

	return
}

func handleAccessDirectMessageMongo(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	user1, ok := dbObj["user1"].(string)
	if !ok {
		err = fmt.Errorf("invalid user1")
		return
	}
	user2, ok := dbObj["user2"].(string)
	if !ok {
		err = fmt.Errorf("invalid user2")
		return
	}
	var otherUserId string
	if user1 == dataSubjectId {
		otherUserId = user2
	} else {
		otherUserId = user1
	}
	var otherUserIDObj primitive.ObjectID
	otherUserIDObj, err = primitive.ObjectIDFromHex(otherUserId)
	if err != nil {
		return
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

	return
}
