package chat

import (
	"context"
	"os"

	firebaseSDK "firebase.google.com/go"
	"google.golang.org/api/option"
)

// App is a global variable to hold the initialized Firebase App object
var App *firebaseSDK.App
var Context context.Context

func initializeFirebaseApp() {
	ctx := context.Background()
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CONFIG")))
	app, err := firebaseSDK.NewApp(ctx, nil, opt)
	if err != nil {
		panic(err.Error())
	}

	App = app
	Context = ctx
}
