package chat

import (
	"fmt"
	"time"

	"github.com/privacy-pal/privacy-pal/go/internal/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	newChat.ID = result.InsertedID.(primitive.ObjectID).Hex()

	// add the group chat to the user
	userID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.M{"$push": bson.M{"gcs": newChat.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return &newChat, nil
}

func GetGroupChatMongo(ID string) (chat *GroupChat, err error) {
	gcID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}

	chat = &GroupChat{}
	if err = test.MongoDb.Collection(FirestoreGroupChatCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: gcID}}).Decode(chat); err != nil {
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
	newUser.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return &newUser, nil
}

func GetUserMongo(ID string) (user *User, err error) {
	userID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}

	user = &User{}
	if err = test.MongoDb.Collection(FirestoreUsersCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: userID}}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(doesNotExistError)
		}
		return nil, err
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

	userID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return fmt.Errorf("error converting ID to ObjectID: %v", err)
	}
	chatIDObj, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("error converting ID to ObjectID: %v", err)
	}

	// add or remove user from mongo group chat users based on action
	if action == JoinChat {
		// Update group chat
		filter := bson.D{{Key: "_id", Value: chatIDObj}}
		update := bson.M{
			"$addToSet": bson.M{
				"users": u.ID,
			},
		}
		_, err = test.MongoDb.Collection(FirestoreGroupChatCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating group chat: %v", err)
		}

		// Update user
		filter = bson.D{{Key: "_id", Value: userID}}
		update = bson.M{
			"$addToSet": bson.M{
				"gcs": chatID,
			},
		}
		_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
	} else if action == QuitChat {
		// Update group chat
		filter := bson.D{{Key: "_id", Value: chatIDObj}}
		update := bson.M{
			"$pull": bson.M{
				"users": u.ID,
			},
		}
		_, err = test.MongoDb.Collection(FirestoreGroupChatCollection).UpdateOne(test.Context, filter, update)
		if err != nil {
			return fmt.Errorf("error updating group chat: %v", err)
		}

		// Update user
		filter = bson.D{{Key: "_id", Value: userID}}
		update = bson.M{
			"$pull": bson.M{
				"gcs": chatID,
			},
		}
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
		User1:    u.ID,
		User2:    user2ID,
		Messages: []Message{},
	}

	result, err := test.MongoDb.Collection(FirestoreDirectMessagesCollection).InsertOne(test.Context, newDM)
	if err != nil {
		return nil, fmt.Errorf("error creating direct message: %v", err)
	}
	newDM.ID = result.InsertedID.(primitive.ObjectID).Hex()

	// add the DM to both users
	userID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.M{"$set": bson.M{fmt.Sprintf("dms.%s", user2ID): newDM.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	user2IDObj, err := primitive.ObjectIDFromHex(user2ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}
	filter = bson.D{{Key: "_id", Value: user2IDObj}}
	update = bson.M{"$set": bson.M{fmt.Sprintf("dms.%s", u.ID): newDM.ID}}
	_, err = test.MongoDb.Collection(FirestoreUsersCollection).UpdateOne(test.Context, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return &newDM, nil
}

func GetDirectMessageMongo(ID string) (chat *DirectMessage, err error) {
	dmID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("error converting ID to ObjectID: %v", err)
	}

	chat = &DirectMessage{}
	if err = test.MongoDb.Collection(FirestoreDirectMessagesCollection).FindOne(test.Context, bson.D{{Key: "_id", Value: dmID}}).Decode(chat); err != nil {
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
	gc, err := GetGroupChatMongo(chatID)
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
		ChatID:    chatID,
		Content:   message,
		Timestamp: time.Now(),
	}

	_, err = test.MongoDb.Collection(FirestoreMessagesCollection).InsertOne(test.Context, newMessage)
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
		ChatID:    chatID,
		Content:   message,
		Timestamp: time.Now(),
	}

	// write to mongo message collection
	_, err = test.MongoDb.Collection(FirestoreMessagesCollection).InsertOne(test.Context, newMessage)
	if err != nil {
		return fmt.Errorf("error creating message: %v", err)
	}

	return nil
}
