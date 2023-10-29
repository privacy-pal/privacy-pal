package chat

import (
	pal "github.com/privacy-pal/privacy-pal/pkg"
)

func HandleAccess(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	switch currentDataNodeLocator.DataType {
	case string(UserDataType):
		return HandleAccessUser(dataSubjectId, currentDataNodeLocator, obj)
	case string(GroupChatDataType):
		return HandleAccessGroupChat(dataSubjectId, currentDataNodeLocator, obj)
	case string(MessageDataType):
		return HandleAccessMessage(dataSubjectId, currentDataNodeLocator, obj)
	default:
		return nil
	}
}

func HandleAccessUser(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Name"] = obj["name"]
	data["Groupchats"] = make([]pal.Locator, 0)
	for _, id := range obj["gcs"].([]interface{}) {
		id := id.(string)
		data["Groupchats"] = append(data["Groupchats"].([]pal.Locator), pal.Locator{
			LocatorType:    pal.Document,
			DataType:       string(GroupChatDataType),
			CollectionPath: []string{"gcs"},
			DocIDs:         []string{id},
		})
	}
	// data["DirectMessages"] = make(map[string]pal.Locator)
	// for id, _ := range obj["dms"].(map[string]interface{}) {

	return data
}

func HandleAccessGroupChat(dataSubjectId string, currentDataNodeLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	data := make(map[string]interface{})

	data["Messages"] = pal.Locator{
		LocatorType:    pal.Collection,
		DataType:       string(MessageDataType),
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
