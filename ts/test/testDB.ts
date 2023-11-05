import { configDotenv } from "dotenv";
import { cert, initializeApp } from "firebase-admin/app";
import { Firestore, getFirestore } from "firebase-admin/firestore";
import { MongoClient } from "mongodb";

export default class TestDatabase {
    public static mongoClient: MongoClient;
    public static firestoreClient: Firestore;
    public static database: "mongo" | "firestore";

    public static async initializeDB(database: "mongo" | "firestore") {
        configDotenv();
        if (database === "mongo") {
            const client = new MongoClient(process.env.MONGO_DB_URI as string);
            await client.connect()
            TestDatabase.mongoClient = client;
        } else if (database === "firestore") {
            initializeApp({
                credential: cert({
                    projectId: process.env.FIREBASE_PROJECT_ID,
                    clientEmail: process.env.FIREBASE_CLIENT_EMAIL,
                    privateKey: process.env.FIREBASE_PRIVATE_KEY
                }),
            });

            TestDatabase.firestoreClient = getFirestore();
        } else {
            throw new Error("Invalid database");
        }

        TestDatabase.database = database;
    }
}