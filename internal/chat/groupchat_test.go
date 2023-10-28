package chat

import (
	"testing"

	"github.com/stretchr/testify/suite"

	pal "github.com/privacy-pal/privacy-pal/pkg"
)

type BasicChatTestSuite struct {
	suite.Suite
	user1 *User
	user2 *User
	gc1   *GroupChat
}

func (suite *BasicChatTestSuite) SetupTest() {
	suite.user1, _ = CreateUser("user1")
	suite.user2, _ = CreateUser("user2")
	suite.gc1, _ = suite.user1.CreateGroupChat()
	suite.user2.JoinOrQuitGroupChat(suite.gc1.ID, JoinChat)
	suite.user1.SendMessageToGroupChat(suite.gc1.ID, "hello")
	suite.user2.SendMessageToGroupChat(suite.gc1.ID, "hi")
	suite.user1.SendMessageToGroupChat(suite.gc1.ID, "how are you?")
}

func (suite *BasicChatTestSuite) TestProcessAccessRequestUser1() {
	dataSubjectLocator := pal.Locator{
		Type:           pal.Document,
		CollectionPath: []string{FirestoreUsersCollection},
		DocIDs:         []string{suite.user1.ID},
		NewDataNode:    func() pal.DataNode { return &User{} },
	}

	client := pal.NewClient(firestoreClient)
	data, err := client.ProcessAccessRequest(dataSubjectLocator, suite.user1.ID)
	suite.NoError(err)
	suite.NotNil(data)
	// data map[Groupchats:[map[Messages:[map[Content:hello Timestamp:2023-10-28 19:58:55.407095 +0000 UTC] map[Content:how are you? Timestamp:2023-10-28 19:58:55.677507 +0000 UTC]]]] Name:user1]
	suite.Equal(suite.user1.Name, data["Name"])
	suite.Equal(1, len(data["Groupchats"].([]interface{})))
	suite.Equal(2, len(data["Groupchats"].([]interface{})[0].(map[string]interface{})["Messages"].([]interface{})))
	suite.Equal("hello", data["Groupchats"].([]interface{})[0].(map[string]interface{})["Messages"].([]interface{})[0].(map[string]interface{})["Content"])
	suite.Equal("how are you?", data["Groupchats"].([]interface{})[0].(map[string]interface{})["Messages"].([]interface{})[1].(map[string]interface{})["Content"])
}

func (suite *BasicChatTestSuite) TestProcessAccessRequestUser2() {
	dataSubjectLocator := pal.Locator{
		Type:           pal.Document,
		CollectionPath: []string{FirestoreUsersCollection},
		DocIDs:         []string{suite.user2.ID},
		NewDataNode:    func() pal.DataNode { return &User{} },
	}

	client := pal.NewClient(firestoreClient)
	data, err := client.ProcessAccessRequest(dataSubjectLocator, suite.user2.ID)
	suite.NoError(err)
	suite.NotNil(data)

	suite.Equal(suite.user2.Name, data["Name"])
	suite.Equal(1, len(data["Groupchats"].([]interface{})))
	suite.Equal(1, len(data["Groupchats"].([]interface{})[0].(map[string]interface{})["Messages"].([]interface{})))
	suite.Equal("hi", data["Groupchats"].([]interface{})[0].(map[string]interface{})["Messages"].([]interface{})[0].(map[string]interface{})["Content"])
}

func TestBasicChatTestSuite(t *testing.T) {
	suite.Run(t, new(BasicChatTestSuite))
}
