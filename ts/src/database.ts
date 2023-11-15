import { Firestore } from "firebase-admin/firestore";
import { MongoClient, Db } from "mongodb";
import { DocumentUpdates, FirestoreLocator, Locator, MongoLocator } from "./model";
import { executeTransactionInFirestore, getDocumentFromFirestore, getDocumentsFromFirestore } from "./firestore";
import { executeTransactionInMongo, getDocumentFromMongo, getDocumentsFromMongo } from "./mongodb";

class Database {
    type: "firestore" | "mongo";
    client: Firestore | MongoClient;
    mongoDb: Db;

    constructor(client: Firestore | MongoClient, mongoDb?: Db) {
        // check type matches
        if (client instanceof Firestore) {
            this.type = "firestore";
        } else if (client instanceof MongoClient && mongoDb) {
            this.type = "mongo";
            this.mongoDb = mongoDb;
        } else {
            throw new Error("The client argument must be either a Firestore or MongoClient. If client is MongoClient, mongoDb must be provided.");
        }
        this.client = client;
    }

    async getDocument(locator: Locator): Promise<any> {
        switch (this.type) {
            case "firestore":
                return getDocumentFromFirestore(this.client as Firestore, locator as FirestoreLocator);
            case "mongo":
                return getDocumentFromMongo(this.mongoDb, locator as MongoLocator);
        }
    }

    async getDocuments(locator: Locator): Promise<any[]> {
        switch (this.type) {
            case "firestore":
                return getDocumentsFromFirestore(this.client as Firestore, locator as FirestoreLocator);
            case "mongo":
                return getDocumentsFromMongo(this.mongoDb, locator as MongoLocator);
        }
    }

    async updateAndDelete(fieldsToUpdate: DocumentUpdates<MongoLocator | FirestoreLocator>[], nodesToDelete: Locator[]): Promise<void> {
        switch (this.type) {
            case "firestore":
                return executeTransactionInFirestore(this.client as Firestore, fieldsToUpdate as DocumentUpdates<FirestoreLocator>[], nodesToDelete as FirestoreLocator[])
            case "mongo":
                return executeTransactionInMongo(this.client as MongoClient, this.mongoDb, fieldsToUpdate as DocumentUpdates<MongoLocator>[], nodesToDelete as MongoLocator[])

        }
    }
}

export default Database;