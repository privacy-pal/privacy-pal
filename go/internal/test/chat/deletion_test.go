package chat

import (
	"encoding/json"
	"testing"

	"github.com/privacy-pal/privacy-pal/go/internal/test"
	pal "github.com/privacy-pal/privacy-pal/go/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeletionWithGroupChatFirestore(t *testing.T) {
	test.InitFirestoreClient()

	// create user 1
	user1, err := CreateUserFirestore("user1")
	if err != nil {
		panic(err)
	}

	// creater user 2
	user2, err := CreateUserFirestore("user2")
	if err != nil {
		panic(err)
	}

	// user 1 creates groupchat
	gc1, err := user1.CreateGroupChatFirestore()
	if err != nil {
		panic(err)
	}

	// user 2 joins groupchat
	err = user2.JoinOrQuitGroupChatFirestore(gc1.ID, JoinChat)
	if err != nil {
		panic(err)
	}

	// user 1 sends message to groupchat
	err = user1.SendMessageToGroupChatFirestore(gc1.ID, "hello")
	if err != nil {
		panic(err)
	}

	// user 2 sends message to groupchat
	err = user2.SendMessageToGroupChatFirestore(gc1.ID, "hi")
	if err != nil {
		panic(err)
	}

	// user 1 sends another message to groupchat
	err = user1.SendMessageToGroupChatFirestore(gc1.ID, "how are you?")
	if err != nil {
		panic(err)
	}

	// user 2 creates direct message with user 1
	dm1, err := user2.CreateDirectMessageFirestore(user1.ID)
	if err != nil {
		panic(err)
	}

	// user 2 sends message to direct message
	err = user2.SendMessageToDirectMessageFirestore(dm1.ID, "Hey! We are in direct message")
	if err != nil {
		panic(err)
	}

	// user 1 sends message to direct message
	err = user1.SendMessageToDirectMessageFirestore(dm1.ID, "Hello!")
	if err != nil {
		panic(err)
	}

	dataSubjectLocator := pal.Locator{
		LocatorType: pal.Document,
		DataType:    string(UserDataType),
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: []string{FirestoreUsersCollection},
			DocIDs:         []string{user1.ID},
		},
	}

	palClient := pal.NewClientWithFirestore(test.FirestoreClient)
	data, err := palClient.ProcessDeletionRequest(HandleDeletionFirestore, dataSubjectLocator, user1.ID, false)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	t.Log(string(json))
}

func TestDeletionWithGroupChatMongo(t *testing.T) {
	test.InitMongoClient()

	// create user 1
	user1, err := CreateUserMongo("user1")
	if err != nil {
		panic(err)
	}

	// creater user 2
	user2, err := CreateUserMongo("user2")
	if err != nil {
		panic(err)
	}

	// user 1 creates groupchat
	gc1, err := user1.CreateGroupChatMongo()
	if err != nil {
		panic(err)
	}

	// user 2 joins groupchat
	err = user2.JoinOrQuitGroupChatMongo(gc1.ID, JoinChat)
	if err != nil {
		panic(err)
	}

	// user 1 sends message to groupchat
	err = user1.SendMessageToGroupChatMongo(gc1.ID, "hello")
	if err != nil {
		panic(err)
	}

	// user 2 sends message to groupchat
	err = user2.SendMessageToGroupChatMongo(gc1.ID, "hi")
	if err != nil {
		panic(err)
	}

	// user 1 sends another message to groupchat
	err = user1.SendMessageToGroupChatMongo(gc1.ID, "how are you?")
	if err != nil {
		panic(err)
	}

	// user 2 creates direct message with user 1
	dm1, err := user2.CreateDirectMessageMongo(user1.ID)
	if err != nil {
		panic(err)
	}

	// user 2 sends message to direct message
	err = user2.SendMessageToDirectMessageMongo(dm1.ID, "Hey! We are in direct message")
	if err != nil {
		panic(err)
	}

	// user 1 sends message to direct message
	err = user1.SendMessageToDirectMessageMongo(dm1.ID, "Hello!")
	if err != nil {
		panic(err)
	}

	userID, err := primitive.ObjectIDFromHex(user1.ID)
	if err != nil {
		panic(err)
	}

	dataSubjectLocator := pal.Locator{
		LocatorType: pal.Document,
		DataType:    string(UserDataType),
		FirestoreLocator: pal.FirestoreLocator{
			CollectionPath: []string{FirestoreUsersCollection},
			DocIDs:         []string{user1.ID},
		},
		MongoLocator: pal.MongoLocator{
			Collection: "users",
			Filter:     bson.D{{Key: "_id", Value: userID}},
		},
	}

	palClient := pal.NewClientWithMongo(test.MongoDb)
	data, err := palClient.ProcessDeletionRequest(HandleDeletionMongo, dataSubjectLocator, user1.ID, false)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	t.Log(string(json))
}
