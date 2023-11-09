package chat

import "time"

const (
	FirestoreUsersCollection          = "users"
	FirestoreGroupChatCollection      = "gcs"
	FirestoreDirectMessagesCollection = "dms"
	FirestoreMessagesCollection       = "messages"
)

type joinQuitAction string

const (
	JoinChat joinQuitAction = "join"
	QuitChat joinQuitAction = "quit"
)

const (
	doesNotExistError string = "does not exist"
)

type User struct {
	ID   string            `firestore:"id,omitempty" bson:"_id,omitempty"`
	Name string            `firestore:"name" bson:"name"`
	GCs  []string          `firestore:"gcs" bson:"gcs"`
	DMs  map[string]string `firestore:"dms" bson:"dms"` // map from other user id to dm id
}

type GroupChat struct {
	ID       string    `firestore:"id,omitempty" bson:"_id,omitempty"`
	Owner    string    `firestore:"owner" bson:"owner"`
	Users    []string  `firestore:"users" bson:"users"`
	Messages []Message `firestore:"-" bson:"messages"` // subcollection
}

type DirectMessage struct {
	ID       string    `firestore:"id,omitempty" bson:"_id,omitempty"`
	User1    string    `firestore:"user1" bson:"user1"`
	User2    string    `firestore:"user2" bson:"user2"`
	Messages []Message `firestore:"-" bson:"messages"` // subcollection
}

type Message struct {
	ID        string    `firestore:"id,omitempty" bson:"_id,omitempty"`
	ChatID    string    `firestore:"-" bson:"chatId"`
	UserID    string    `firestore:"userId" bson:"userId"`
	Content   string    `firestore:"content" bson:"content"`
	Timestamp time.Time `firestore:"timestamp" bson:"timestamp"`
}
