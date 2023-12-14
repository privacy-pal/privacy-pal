package chat

import (
	"fmt"

	pal "github.com/privacy-pal/privacy-pal/go/pkg"
)

func HandleAccessFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleAccessUserFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleAccessGroupChatFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleAccessMessageFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleAccessDirectMessageFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		err = fmt.Errorf("invalid data type")
		return
	}
}

func handleAccessUserFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	if currentDbObjLocator.DocIDs[len(currentDbObjLocator.DocIDs)-1] != dataSubjectId {
		data["Name"] = dbObj["name"]
		return
	}

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

		groupchatLocators = append(groupchatLocators, pal.Locator{
			LocatorType: pal.Document,
			DataType:    GroupChatDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"gcs"},
				DocIDs:         []string{id},
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

		directMessageLocators = append(directMessageLocators, pal.Locator{
			LocatorType: pal.Document,
			DataType:    DirectMessageDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"dms"},
				DocIDs:         []string{dmID},
			},
		})
	}
	data["DirectMessages"] = directMessageLocators

	return
}

func handleAccessGroupChatFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
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
	}

	return
}

func handleAccessMessageFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	data["Content"] = dbObj["content"]
	data["Timestamp"] = dbObj["timestamp"]

	return
}

func handleAccessDirectMessageFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (data map[string]interface{}, err error) {
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
			Filters: []pal.Filter{
				{
					Path:  "userId",
					Op:    "==",
					Value: dataSubjectId,
				},
			},
		},
	}

	return
}
