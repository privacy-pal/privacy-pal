package chat

import (
	"fmt"
	"time"

	"github.com/privacy-pal/privacy-pal/internal/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *User) CreateGroupChatMongo() (chat *GroupChat, err error) {
	newChat := GroupChat{
		Owner:    u.ID,
		Users:    []string{},
		Messages: []Message{},
	}

	result, err := test.MongoDb.Collection(FirestoreGroupChatCollection).InsertOne(test.Context, newChat)
	if err != nil {
		return nil, fmt.Errorf("error creating group chat: %v", err)
	}
	newChat.ID = result.InsertedID.(string)

	// add the group chat to the user
	filter := bson.D{{Key: "_id", Value: u.ID}}
	update := bson.M{"$push": bson.M{"user.$.gcs": newChat.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return &newChat, nil
}

func GetGroupChatMongo(ID string) (chat *GroupChat, err error) {
	if err = test.MongoDb.Collection(FirestoreGroupChatCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: ID}}).Decode(chat); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(doesNotExistError)
		}
		return nil, fmt.Errorf("error getting group chat: %v", err)
	}
	chat.ID = ID

	return chat, nil
}

func CreateUserMongo(name string) (user *User, err error) {
	newUser := User{
		Name: name,
		GCs:  []string{},
		DMs:  map[string]string{},
	}

	result, err := test.MongoDb.Collection(FirestoreUsersCollection).InsertOne(test.Context, newUser)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	newUser.ID = result.InsertedID.(string)

	return &newUser, nil
}

func GetUserMongo(ID string) (user *User, err error) {
	if err = test.MongoDb.Collection(FirestoreUsersCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: ID}}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(doesNotExistError)
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	user.ID = ID

	return user, nil
}

func (u *User) JoinOrQuitGroupChatMongo(chatID string, action joinQuitAction) (err error) {
	_, err = GetUserMongo(u.ID)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	_, err = GetGroupChatMongo(chatID)
	if err != nil {
		return fmt.Errorf("error getting group chat: %v", err)
	}

	// add or remove user from mongo group chat users based on action
	if action == JoinChat {
		// Update group chat
		filter := bson.D{{Key: "_id", Value: chatID}}
		update := bson.M{"$set": bson.M{"gc.$.users": bson.M{"setUnion": bson.M{"gc.$.users": u.ID}}}}
		_, err = test.MongoDb.Collection(FirestoreGroupChatCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating group chat: %v", err)
		}

		// Update user
		filter = bson.D{{Key: "_id", Value: u.ID}}
		update = bson.M{"$set": bson.M{"user.$.gcs": bson.M{"setUnion": bson.M{"user.$.gcs": chatID}}}}
		_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
	} else if action == QuitChat {
		// Update group chat
		filter := bson.D{{Key: "_id", Value: chatID}}
		update := bson.M{"$set": bson.M{"gc.$.users": bson.M{"setDifference": bson.M{"gc.$.users": u.ID}}}}
		_, err = test.MongoDb.Collection(FirestoreGroupChatCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating group chat: %v", err)
		}

		// Update user
		filter = bson.D{{Key: "_id", Value: u.ID}}
		update = bson.M{"$set": bson.M{"user.$.gcs": bson.M{"setDifference": bson.M{"user.$.gcs": chatID}}}}
		_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
	}

	return nil
}

func (u *User) CreateDirectMessageMongo(user2ID string) (chat *DirectMessage, err error) {
	// check if user exists and if DM already exists
	user2, err := GetUserMongo(user2ID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	if _, ok := user2.DMs[u.ID]; ok {
		return nil, fmt.Errorf("direct message already exists")
	}
	if _, ok := u.DMs[user2ID]; ok {
		return nil, fmt.Errorf("direct message already exists")
	}

	newDM := DirectMessage{
		User1: u.ID,
		User2: user2ID,
	}

	result, err := test.MongoDb.Collection(FirestoreDirectMessagesCollection).InsertOne(test.Context, newDM)
	if err != nil {
		return nil, fmt.Errorf("error creating direct message: %v", err)
	}
	newDM.ID = result.InsertedID.(string)

	// add the DM to both users
	filter := bson.D{{Key: "_id", Value: u.ID}}
	update := bson.M{"$set": bson.M{fmt.Sprintf("dm.%s", user2ID): newDM.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	filter = bson.D{{Key: "_id", Value: user2ID}}
	update = bson.M{"$set": bson.M{fmt.Sprintf("dm.%s", u.ID): newDM.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return &newDM, nil
}

func GetDirectMessageMongo(ID string) (chat *DirectMessage, err error) {
	if err = test.MongoDb.Collection(FirestoreDirectMessagesCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: ID}}).Decode(chat); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(doesNotExistError)
		}
		return nil, fmt.Errorf("error getting direct message: %v", err)
	}
	chat.ID = ID

	return chat, nil
}

func (u *User) SendMessageToGroupChatMongo(chatID string, message string) error {
	// get the group chat
	gc, err := GetGroupChatFirestore(chatID)
	if err != nil {
		return err
	}

	// check if user is in chat
	if (!stringInSlice(u.ID, gc.Users)) && (gc.Owner != u.ID) {
		return fmt.Errorf("user is not in group chat")
	}

	// create message
	newMessage := Message{
		UserID:    u.ID,
		Content:   message,
		Timestamp: time.Now(),
	}

	// write to mongo group chat collection
	filter := bson.D{{Key: "_id", Value: chatID}}
	update := bson.M{"$push": bson.M{"messages": newMessage}}
	_, err = test.MongoDb.Collection(FirestoreGroupChatCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return fmt.Errorf("error creating message: %v", err)
	}

	return nil
}

func (u *User) SendMessageToDirectMessageMongo(chatID string, message string) error {
	// get the direct message
	dm, err := GetDirectMessageMongo(chatID)
	if err != nil {
		return err
	}

	// check if user is in dm
	if dm.User1 != u.ID && dm.User2 != u.ID {
		return fmt.Errorf("user is not in direct message")
	}

	// create message
	newMessage := Message{
		UserID:    u.ID,
		Content:   message,
		Timestamp: time.Now(),
	}

	// write to mongo direct message collection
	filter := bson.D{{Key: "_id", Value: chatID}}
	update := bson.M{"$push": bson.M{"messages": newMessage}}
	_, err = test.MongoDb.Collection(FirestoreDirectMessagesCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return fmt.Errorf("error creating message: %v", err)
	}

	return nil
}
