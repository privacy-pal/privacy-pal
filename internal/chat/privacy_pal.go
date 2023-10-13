package chat

import (
	"cloud.google.com/go/firestore"
	pal "github.com/tianrendong/privacy-pal/pkg"
)

func (u *User) HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{} {
	ret := make(map[string]interface{})

	ret["name"] = u.Name
	ret["group chat"] = make([]pal.Locator, 0)
	for _, gcId := range u.GCs {
		ret["group chat"] = append(ret["group chat"].([]pal.Locator), pal.Locator{
			Type:           pal.Document,
			CollectionPath: []string{FirestoreGroupChatCollection},
			DocIDs:         []string{gcId},
			NewDataNode:    func() pal.DataNode { return &GroupChat{} },
		})
	}

	return ret
}

func (u *User) HandleDeletion(dataSubjectId string) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {
	return
}

func (gc *GroupChat) HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["messages"] = []pal.Locator{{
		Type:           pal.Collection,
		CollectionPath: []string{FirestoreGroupChatCollection, FirestoreMessagesCollection},
		DocIDs:         []string{currentDocumentID},
		Queries: []pal.Query{{
			Path:  "userId",
			Op:    "==",
			Value: dataSubjectId,
		}},
		NewDataNode: func() pal.DataNode { return &Message{} },
	}}

	return ret
}

func (m *GroupChat) HandleDeletion(dataSubjectId string) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {
	return
}

func (m *Message) HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{} {
	ret := make(map[string]interface{})
	if m.UserID == dataSubjectId {
		ret["content"] = m.Content
		ret["timestamp"] = m.Timestamp
	}

	return ret
}

func (m *Message) HandleDeletion(dataSubjectId string) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {
	return
}
