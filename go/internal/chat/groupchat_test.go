package chat

import (
	"encoding/json"
	"testing"

	pal "github.com/privacy-pal/privacy-pal/pkg"
)

func Test1(t *testing.T) {
	// create user 1
	user1, err := CreateUser("user1")
	if err != nil {
		panic(err)
	}

	// creater user 2
	user2, err := CreateUser("user2")
	if err != nil {
		panic(err)
	}

	// user 1 creates groupchat
	gc1, err := user1.CreateGroupChat()
	if err != nil {
		panic(err)
	}

	// user 2 joins groupchat
	err = user2.JoinOrQuitGroupChat(gc1.ID, JoinChat)
	if err != nil {
		panic(err)
	}

	// user 1 sends message to groupchat
	err = user1.SendMessageToGroupChat(gc1.ID, "hello")
	if err != nil {
		panic(err)
	}

	// user 2 sends message to groupchat
	err = user2.SendMessageToGroupChat(gc1.ID, "hi")
	if err != nil {
		panic(err)
	}

	// user 1 sends another message to groupchat
	err = user1.SendMessageToGroupChat(gc1.ID, "how are you?")
	if err != nil {
		panic(err)
	}

	// user 2 creates direct message with user 1
	dm1, err := user2.CreateDirectMessage(user1.ID)
	if err != nil {
		panic(err)
	}

	// user 2 sends message to direct message
	err = user2.SendMessageToDirectMessage(dm1.ID, "Hey! We are in direct message")
	if err != nil {
		panic(err)
	}

	// user 1 sends message to direct message
	err = user1.SendMessageToDirectMessage(dm1.ID, "Hello!")
	if err != nil {
		panic(err)
	}

	dataSubjectLocator := pal.Locator{
		LocatorType:    pal.Document,
		DataType:       string(UserDataType),
		CollectionPath: []string{FirestoreUsersCollection},
		DocIDs:         []string{user1.ID},
	}

	client := pal.NewClient(firestoreClient)
	data, err := client.ProcessAccessRequest(HandleAccess, dataSubjectLocator, user1.ID)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	t.Log(string(json))
}
