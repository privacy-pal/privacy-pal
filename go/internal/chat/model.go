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
	ID   string            `firestore:"id,omitempty"`
	Name string            `firestore:"name"`
	GCs  []string          `firestore:"gcs"`
	DMs  map[string]string `firestore:"dms"` // map from other user id to dm id
}

type GroupChat struct {
	ID       string    `firestore:"id,omitempty"`
	Owner    string    `firestore:"owner"`
	Users    []string  `firestore:"users"`
	Messages []Message `firestore:"-"` // subcollection
}

type DirectMessage struct {
	ID       string    `firestore:"id,omitempty"`
	User1    string    `firestore:"user1"`
	User2    string    `firestore:"user2"`
	Messages []Message `firestore:"-"` // subcollection
}

type Message struct {
	ID        string    `firestore:"id,omitempty"`
	UserID    string    `firestore:"userId"`
	Content   string    `firestore:"content"`
	Timestamp time.Time `firestore:"timestamp"`
}

type DataType string

const (
	UserDataType      DataType = "user"
	GroupChatDataType DataType = "groupchat"
	DirectMessageTyep DataType = "directmessage"
	MessageDataType   DataType = "message"
)
