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

func HandleAccess(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	switch currentDataNodeLocator.DataType {
	case UserDataType:
		return HandleAccessUser(dataSubjectId, currentDataNodeLocator, obj)
	case GroupChatDataType:
		return HandleAccessGroupChat(dataSubjectId, currentDataNodeLocator, obj)
	case MessageDataType:
		return HandleAccessMessage(dataSubjectId, currentDataNodeLocator, obj)
	case DirectMessageDataType:
		return HandleAccessDirectMessage(dataSubjectId, currentDataNodeLocator, obj)
	default:
		// TODO: should return error
		return nil
	}
}

func HandleAccessUser(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	// TODO: include in documentation: you can access the id in 2 ways
	if obj["_id"].(string) != dataSubjectId {
		data["Name"] = obj["name"]
		return data
	}
	// if currentDataNodeLocator.DocIDs[len(currentDataNodeLocator.DocIDs)-1] != dataSubjectId {
	// 	data["Name"] = obj["name"]
	// 	return data
	// }

	data["Name"] = obj["name"]
	data["Groupchats"] = make([]pal.Locator, 0)
	for _, id := range obj["gcs"].([]interface{}) {
		id := id.(string)
		data["Groupchats"] = append(data["Groupchats"].([]pal.Locator), pal.Locator{
			LocatorType:    pal.Document,
			DataType:       GroupChatDataType,
			CollectionPath: []string{"gcs"},
			DocIDs:         []string{id},
		})
	}
	data["DirectMessages"] = make([]pal.Locator, 0)
	for _, DMId := range obj["dms"].(map[string]interface{}) {
		DMId := DMId.(string)
		data["DirectMessages"] = append(data["DirectMessages"].([]pal.Locator), pal.Locator{
			LocatorType:    pal.Document,
			DataType:       DirectMessageDataType,
			CollectionPath: []string{"dms"},
			DocIDs:         []string{DMId},
		})
	}

	return data
}

func HandleAccessGroupChat(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Messages"] = pal.Locator{
		LocatorType:    pal.Collection,
		DataType:       MessageDataType,
		CollectionPath: append(currentDataNodeLocator.CollectionPath, "messages"),
		DocIDs:         currentDataNodeLocator.DocIDs,
		Queries: []pal.Query{
			{
				Path:  "userId",
				Op:    "==",
				Value: dataSubjectId,
			},
		},
	}

	return data
}

func HandleAccessMessage(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Content"] = obj["content"]
	data["Timestamp"] = obj["timestamp"]

	return data
}

func HandleAccessDirectMessage(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	var otherUserId string
	if obj["user1"].(string) == dataSubjectId {
		otherUserId = obj["user2"].(string)
	} else {
		otherUserId = obj["user1"].(string)
	}

	data["Other User"] = pal.Locator{
		LocatorType:    pal.Document,
		DataType:       UserDataType,
		CollectionPath: []string{"users"},
		DocIDs:         []string{otherUserId},
	}
	data["Messages"] = pal.Locator{
		LocatorType:    pal.Collection,
		DataType:       MessageDataType,
		CollectionPath: append(currentDataNodeLocator.CollectionPath, "messages"),
		DocIDs:         currentDataNodeLocator.DocIDs,
		Queries: []pal.Query{
			{
				Path:  "userId",
				Op:    "==",
				Value: dataSubjectId,
			},
		},
	}

	return data
}
