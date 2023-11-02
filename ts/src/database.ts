import { Firestore } from "firebase-admin/firestore";
import { MongoClient } from "mongodb";
import { FirestoreLocator, Locator, MongoLocator } from "./model";
import { getDocumentFromFirestore, getDocumentsFromFirestore } from "./firestore";
import { getDocumentFromMongo, getDocumentsFromMongo } from "./mongodb";

class Database {
    type: "firestore" | "mongo";
    client: Firestore | MongoClient;

    constructor(client: Firestore | MongoClient) {
        // check type matches
        if (client instanceof Firestore) {
            this.type = "firestore";
        } else if (client instanceof MongoClient) {
            this.type = "mongo";
        } else {
            throw new Error("Client must be either a Firestore or MongoClient");
        }
        this.client = client;
    }

    async getDocument(locator: Locator): Promise<any> {
        if (this.type == "firestore") {
            return getDocumentFromFirestore(this.client as Firestore, locator as FirestoreLocator);
        } else {
            return getDocumentFromMongo(this.client as MongoClient, locator as MongoLocator);
        }
    }

    async getDocuments(locator: Locator): Promise<any[]> {
        if (this.type == "firestore") {
            return getDocumentsFromFirestore(this.client as Firestore, locator as FirestoreLocator);
        } else {
            return getDocumentsFromMongo(this.client as MongoClient, locator as MongoLocator);
        }
    }
}

export default Database;