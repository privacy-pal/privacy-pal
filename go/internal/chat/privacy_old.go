package chat

import (
	"github.com/mitchellh/mapstructure"
	pal "github.com/privacy-pal/privacy-pal/pkg"
)

func HandleAccessOld(dataSubjectId string, currentDbObjLocator pal.Locator, obj pal.DatabaseObject) map[string]interface{} {
	switch currentDbObjLocator.DataType {
	case string(UserDataType):
		user := &User{}
		err := mapstructure.Decode(obj, user)
		if err != nil {
			panic(err)
		}
		return user.HandleAccess(dataSubjectId, currentDbObjLocator)
	case string(GroupChatDataType):
		groupChat := &GroupChat{}
		err := mapstructure.Decode(obj, groupChat)
		if err != nil {
			panic(err)
		}
		return groupChat.HandleAccess(dataSubjectId, currentDbObjLocator)
	case string(MessageDataType):
		message := &Message{}
		err := mapstructure.Decode(obj, message)
		if err != nil {
			panic(err)
		}
		return message.HandleAccess(dataSubjectId, currentDbObjLocator)
	default:
		return nil
	}
}

// TODO: in documentation
// You can also make your structs implement the handle access function, which allows you to access the data in the object
func (u *User) HandleAccess(dataSubjectId string, currentDbObjLocator pal.Locator) map[string]interface{} {
	data := make(map[string]interface{})

	data["Name"] = u.Name
	data["Groupchats"] = make([]pal.Locator, 0)
	for _, id := range u.GCs {
		data["Groupchats"] = append(data["Groupchats"].([]pal.Locator), pal.Locator{
			LocatorType:    pal.Document,
			DataType:       string(GroupChatDataType),
			CollectionPath: []string{"gcs"},
			DocIDs:         []string{id},
		})
	}

	return data
}

func (g *GroupChat) HandleAccess(dataSubjectId string, currentDbObjLocator pal.Locator) map[string]interface{} {
	data := make(map[string]interface{})

	data["Messages"] = pal.Locator{
		LocatorType:    pal.Collection,
		DataType:       string(MessageDataType),
		CollectionPath: append(currentDbObjLocator.CollectionPath, "messages"),
		DocIDs:         currentDbObjLocator.DocIDs,
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

func (m *Message) HandleAccess(dataSubjectId string, currentDbObjLocator pal.Locator) map[string]interface{} {
	data := make(map[string]interface{})

	data["Content"] = m.Content
	data["Timestamp"] = m.Timestamp

	return data
}
