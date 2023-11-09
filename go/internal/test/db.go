package test

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebaseSDK "firebase.google.com/go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

var Context context.Context
var FirestoreClient *firestore.Client
var MongoDb *mongo.Database

func init() {
	err := godotenv.Load("../../../.env")
	fmt.Println("Loaded .env file")
	if err != nil {
		panic(err)
	}

	Context = context.TODO()
}

func InitFirestoreClient() {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CONFIG")))
	app, err := firebaseSDK.NewApp(Context, nil, opt)
	if err != nil {
		panic(err.Error())
	}

	FirestoreClient, err = app.Firestore(Context)
	if err != nil {
		panic(fmt.Errorf("firestore client error: %v", err))
	}
}

func InitMongoClient() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(Context, opts)
	if err != nil {
		panic(fmt.Errorf("mongo client error: %v", err))
	}

	MongoDb = client.Database(os.Getenv("MONGO_DB_NAME"))
}
