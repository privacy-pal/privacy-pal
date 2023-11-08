package chat

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/privacy-pal/privacy-pal/internal/test"
)

func init() {
	test.InitFirestoreClient()
}

func (u *User) CreateGroupChat() (chat *GroupChat, err error) {
	newChat := GroupChat{
		Owner:    u.ID,
		Users:    []string{},
		Messages: []Message{},
	}

	ref, _, err := test.FirestoreClient.Collection(FirestoreGroupChatCollection).Add(test.Context, newChat)
	if err != nil {
		return nil, fmt.Errorf("error creating group chat: %v\n", err)
	}
	newChat.ID = ref.ID

	// add the group chat to the user
	_, err = test.FirestoreClient.Collection(FirestoreUsersCollection).Doc(u.ID).Set(test.Context, map[string]interface{}{
		"gcs": firestore.ArrayUnion(newChat.ID),
	}, firestore.MergeAll)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %v\n", err)
	}

	return &newChat, nil
}

func GetGroupChat(ID string) (chat *GroupChat, err error) {
	doc, err := test.FirestoreClient.Collection(FirestoreGroupChatCollection).Doc(ID).Get(test.Context)
	if err != nil {
		return nil, fmt.Errorf("error getting group chat: %v\n", err)
	}

	if !doc.Exists() {
		return nil, fmt.Errorf(doesNotExistError)
	}

	chat = &GroupChat{}
	err = doc.DataTo(chat)
	if err != nil {
		return nil, fmt.Errorf("error parsing group chat: %v\n", err)
	}
	chat.ID = doc.Ref.ID

	return chat, nil
}

func CreateUser(name string) (user *User, err error) {
	newUser := User{
		Name: name,
		GCs:  []string{},
		DMs:  map[string]string{},
	}

	ref, _, err := test.FirestoreClient.Collection(FirestoreUsersCollection).Add(test.Context, newUser)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error creating user: %v\n", err)
	}
	newUser.ID = ref.ID

	return &newUser, nil
}

func GetUser(ID string) (user *User, err error) {
	doc, err := test.FirestoreClient.Collection(FirestoreUsersCollection).Doc(ID).Get(test.Context)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v\n", err)
	}

	if !doc.Exists() {
		return nil, fmt.Errorf(doesNotExistError)
	}

	user = &User{}
	err = doc.DataTo(user)
	if err != nil {
		return nil, fmt.Errorf("error parsing user: %v\n", err)
	}
	user.ID = doc.Ref.ID

	return user, nil
}

func (u *User) JoinOrQuitGroupChat(chatID string, action joinQuitAction) (err error) {
	_, err = GetUser(u.ID)
	if err != nil {
		return fmt.Errorf("error getting user: %v\n", err)
	}

	_, err = GetGroupChat(chatID)
	if err != nil {
		return fmt.Errorf("error getting group chat: %v\n", err)
	}

	updates := []firestore.Update{}
	if action == JoinChat {
		updates = []firestore.Update{
			{Path: "users", Value: firestore.ArrayUnion(u.ID)},
		}
	} else if action == QuitChat {
		updates = []firestore.Update{
			{Path: "users", Value: firestore.ArrayRemove(u.ID)},
		}
	}
	_, err = test.FirestoreClient.Collection(FirestoreGroupChatCollection).Doc(chatID).Update(test.Context, updates)

	if err != nil {
		return fmt.Errorf("error updating group chat: %v\n", err)
	}

	// update user
	updates = []firestore.Update{}
	if action == JoinChat {
		updates = []firestore.Update{
			{Path: "gcs", Value: firestore.ArrayUnion(chatID)},
		}
	} else if action == QuitChat {
		updates = []firestore.Update{
			{Path: "gcs", Value: firestore.ArrayRemove(chatID)},
		}
	}

	_, err = test.FirestoreClient.Collection(FirestoreUsersCollection).Doc(u.ID).Update(test.Context, updates)
	if err != nil {
		return fmt.Errorf("error updating user: %v\n", err)
	}

	return nil
}

func (u *User) CreateDirectMessage(user2ID string) (chat *DirectMessage, err error) {
	// check if user exists and if DM already exists
	user2, err := GetUser(user2ID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v\n", err)
	}
	if _, ok := user2.DMs[u.ID]; ok {
		return nil, fmt.Errorf("direct message already exists\n")
	}
	if _, ok := u.DMs[user2ID]; ok {
		return nil, fmt.Errorf("direct message already exists\n")
	}

	newDM := DirectMessage{
		User1: u.ID,
		User2: user2ID,
	}

	ref, _, err := test.FirestoreClient.Collection(FirestoreDirectMessagesCollection).Add(test.Context, newDM)
	if err != nil {
		return nil, fmt.Errorf("error creating direct message: %v\n", err)
	}
	newDM.ID = ref.ID

	// add the DM to both users
	_, err = test.FirestoreClient.Collection(FirestoreUsersCollection).Doc(u.ID).Set(test.Context, map[string]interface{}{
		"dms": map[string]string{
			user2ID: newDM.ID,
		},
	}, firestore.MergeAll)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v\n", err)
	}

	_, err = test.FirestoreClient.Collection(FirestoreUsersCollection).Doc(user2ID).Set(test.Context, map[string]interface{}{
		"dms": map[string]string{
			u.ID: newDM.ID,
		},
	}, firestore.MergeAll)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v\n", err)
	}
	return &newDM, nil
}

func GetDirectMessage(ID string) (chat *DirectMessage, err error) {
	doc, err := test.FirestoreClient.Collection(FirestoreDirectMessagesCollection).Doc(ID).Get(test.Context)
	if err != nil {
		return nil, fmt.Errorf("error getting direct message: %v\n", err)
	}

	if !doc.Exists() {
		return nil, fmt.Errorf(doesNotExistError)
	}

	chat = &DirectMessage{}
	err = doc.DataTo(chat)
	if err != nil {
		return nil, fmt.Errorf("error parsing direct message: %v\n", err)
	}
	chat.ID = doc.Ref.ID

	return chat, nil
}

func (u *User) SendMessageToGroupChat(chatID string, message string) error {
	// get the group chat
	gc, err := GetGroupChat(chatID)
	if err != nil {
		return err
	}

	// check if user is in chat
	if (!stringInSlice(u.ID, gc.Users)) && (gc.Owner != u.ID) {
		return fmt.Errorf("user is not in group chat\n")
	}

	// create message
	newMessage := Message{
		UserID:    u.ID,
		Content:   message,
		Timestamp: time.Now(),
	}

	// write to firestore subcollection
	ref, _, err := test.FirestoreClient.Collection(FirestoreGroupChatCollection).Doc(chatID).
		Collection(FirestoreMessagesCollection).Add(test.Context, newMessage)
	if err != nil {
		return fmt.Errorf("error creating message: %v\n", err)
	}
	newMessage.ID = ref.ID

	return nil
}

func (u *User) SendMessageToDirectMessage(chatID string, message string) error {
	// get the direct message
	dm, err := GetDirectMessage(chatID)
	if err != nil {
		return err
	}

	// check if user is in dm
	if dm.User1 != u.ID && dm.User2 != u.ID {
		return fmt.Errorf("user is not in direct message\n")
	}

	// create message
	newMessage := Message{
		UserID:    u.ID,
		Content:   message,
		Timestamp: time.Now(),
	}

	// write to firestore subcollection
	ref, _, err := test.FirestoreClient.Collection(FirestoreDirectMessagesCollection).Doc(chatID).
		Collection(FirestoreMessagesCollection).Add(test.Context, newMessage)
	if err != nil {
		return fmt.Errorf("error creating message: %v\n", err)
	}
	newMessage.ID = ref.ID

	return nil
}
