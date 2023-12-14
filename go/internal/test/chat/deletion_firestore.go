package chat

import (
	"fmt"

	"cloud.google.com/go/firestore"
	pal "github.com/privacy-pal/privacy-pal/go/pkg"
)

func HandleDeletionFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	switch currentDbObjLocator.DataType {
	case UserDataType:
		return handleDeletionUserFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case GroupChatDataType:
		return handleDeletionGroupChatFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case MessageDataType:
		return handleDeletionMessageFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	case DirectMessageDataType:
		return handleDeletionDirectMessageFirestore(dataSubjectId, currentDbObjLocator, dbObj)
	default:
		err = fmt.Errorf("invalid data type: %s", currentDbObjLocator.DataType)
		return
	}
}

func handleDeletionUserFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
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

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    GroupChatDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"gcs"},
				DocIDs:         []string{id},
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

		nodesToTraverse = append(nodesToTraverse, pal.Locator{
			LocatorType: pal.Document,
			DataType:    DirectMessageDataType,
			FirestoreLocator: pal.FirestoreLocator{
				CollectionPath: []string{"dms"},
				DocIDs:         []string{dmID},
			},
		})
	}

	deleteNode = true
	return
}

func handleDeletionGroupChatFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
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
	})

	deleteNode = false

	firestoreUpdates := []firestore.Update{}

	users, ok := dbObj["users"].([]interface{})
	if !ok {
		err = fmt.Errorf("invalid users field")
		return
	}
	if (dbObj["owner"] == dataSubjectId) && (len(users) > 1) {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  "owner",
			Value: users[0],
		})
	} else {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  "users",
			Value: firestore.ArrayRemove(dataSubjectId),
		})
	}

	fieldsToUpdate = pal.FieldUpdates{
		FirestoreUpdates: firestoreUpdates,
	}

	return
}

func handleDeletionMessageFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
	deleteNode = true
	return
}

func handleDeletionDirectMessageFirestore(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate pal.FieldUpdates, err error) {
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
	})

	deleteNode = false

	return
}
