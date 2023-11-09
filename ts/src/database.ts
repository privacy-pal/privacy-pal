import { Firestore, UpdateData } from "firebase-admin/firestore";
import { MongoClient, UpdateFilter } from "mongodb";
import { FieldsToUpdate, FirestoreLocator, Locator, MongoLocator } from "./model";
import { executeTransactionInFirestore, getDocumentFromFirestore, getDocumentsFromFirestore } from "./firestore";
import { executeTransactionInMongo, getDocumentFromMongo, getDocumentsFromMongo } from "./mongodb";

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
        switch (this.type) {
            case "firestore":
                return getDocumentFromFirestore(this.client as Firestore, locator as FirestoreLocator);
            case "mongo":
                return getDocumentFromMongo(this.client as MongoClient, locator as MongoLocator);
        }
    }

    async getDocuments(locator: Locator): Promise<any[]> {
        switch (this.type) {
            case "firestore":
                return getDocumentsFromFirestore(this.client as Firestore, locator as FirestoreLocator);
            case "mongo":
                return getDocumentsFromMongo(this.client as MongoClient, locator as MongoLocator);
        }
    }

    async updateAndDelete(fieldsToUpdate: FieldsToUpdate<MongoLocator | FirestoreLocator>[], nodesToDelete: Locator[]): Promise<void> {
        switch (this.type) {
            case "firestore":
                executeTransactionInFirestore(this.client as Firestore, fieldsToUpdate as FieldsToUpdate<FirestoreLocator>[], nodesToDelete as FirestoreLocator[])
            case "mongo":
                executeTransactionInMongo(this.client as MongoClient, fieldsToUpdate as FieldsToUpdate<MongoLocator>[], nodesToDelete as MongoLocator[])

        }
    }
}

export default Database;