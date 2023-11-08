package chat

import (
	pal "github.com/privacy-pal/privacy-pal/pkg"
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
		return HandleAccessUser(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return HandleAccessGroupChat(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return HandleAccessMessage(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return HandleAccessDirectMessage(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		// TODO: should return error
		return nil
	}
}

func HandleAccessUser(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
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
		data["Groupchats"] = append(data["Groupchats"].([]pal.Locator), pal.Locator{
			LocatorType: pal.Document,
			DataType:    GroupChatDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"gcs"},
				DocIDs:         []string{id},
			},
		})
	}
	data["DirectMessages"] = make([]pal.Locator, 0)
	// TODO: support this in yaml
	for _, DMId := range dbObj["dms"].(map[string]interface{}) {
		DMId := DMId.(string)
		data["DirectMessages"] = append(data["DirectMessages"].([]pal.Locator), pal.Locator{
			LocatorType: pal.Document,
			DataType:    DirectMessageDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"dms"},
				DocIDs:         []string{DMId},
			},
		})
	}

	return data
}

func HandleAccessGroupChat(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Messages"] = pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: append(currentDbObjLocator.FirestoreLocator.CollectionPath, "messages"),
			DocIDs:         currentDbObjLocator.DocIDs,
			Queries: []pal.Query{
				{
					Path:  "userId",
					Op:    "==",
					Value: dataSubjectId,
				},
			},
		},
	}

	return data
}

func HandleAccessMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Content"] = dbObj["content"]
	data["Timestamp"] = dbObj["timestamp"]

	return data
}

func HandleAccessDirectMessage(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	var otherUserId string
	if dbObj["user1"].(string) == dataSubjectId {
		otherUserId = dbObj["user2"].(string)
	} else {
		otherUserId = dbObj["user1"].(string)
	}

	data["Other User"] = pal.Locator{
		LocatorType: pal.Document,
		DataType:    UserDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: []string{"users"},
			DocIDs:         []string{otherUserId},
		},
	}
	data["Messages"] = pal.Locator{
		LocatorType: pal.Collection,
		DataType:    MessageDataType,
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: append(currentDbObjLocator.FirestoreLocator.CollectionPath, "messages"),
			DocIDs:         currentDbObjLocator.DocIDs,
			Queries: []pal.Query{
				{
					Path:  "userId",
					Op:    "==",
					Value: dataSubjectId,
				},
			},
		},
	}

	return data
}
