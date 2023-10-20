package main

type GroupChat struct {
	ID    string   `firestore:"id,omitempty"`
	Owner string   `firestore:"owner"`
	Users []string `firestore:"users"`
}
