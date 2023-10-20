package main

type User struct {
	ID   string            `firestore:"id,omitempty"`
	Name string            `firestore:"name"`
	GCs  []string          `firestore:"gcs"`
	DMs  map[string]string `firestore:"dms"` // map from other user id to dm id
}
