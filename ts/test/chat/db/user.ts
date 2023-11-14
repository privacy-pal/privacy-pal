import TestDatabase from "../../testDB";
import { FirestoreCollections, doesNotExistError } from "../model/shared";
import User from "../model/user";
import { ObjectId } from "mongodb";

export async function GetUser(ID: string): Promise<User | undefined> {
    try {
        if (TestDatabase.database === "firestore") {
            const doc = await TestDatabase.firestoreClient.collection(FirestoreCollections.Users).doc(ID).get();

            if (!doc.exists) {
                throw new Error(doesNotExistError);
            }

            const user = new User(ID);
            const data = doc.data();
            if (data) {
                user.name = data.name || '';
                user.gcs = data.gcs || [];
                user.dms = data.dms || {};

                return user;
            }
        } else if (TestDatabase.database === "mongo") {
            const doc = await TestDatabase.mongoDb.collection(FirestoreCollections.Users).findOne({ _id: new ObjectId(ID) });

            if (!doc) {
                throw new Error(doesNotExistError);
            }

            const user = new User(ID);
            user.name = doc.name || '';
            user.gcs = doc.gcs || [];
            user.dms = doc.dms || {};

            return user;
        } else {
            throw new Error("Database not initialized");
        }
    } catch (err) {
        throw new Error(`Error getting user: ${err}`);
    }
}

export async function CreateUser(name: string): Promise<User> {
    try {
        const user = new User(name);
        if (TestDatabase.database === "firestore") {
            const docRef = await TestDatabase.firestoreClient.collection(FirestoreCollections.Users).add(Object.assign({}, user));
            user.id = docRef.id;
            return user;
        } else if (TestDatabase.database === "mongo") {
            const doc = await TestDatabase.mongoDb.collection(FirestoreCollections.Users).insertOne(Object.assign({}, user));
            user.id = doc.insertedId.toString();
            return user;
        } else {
            throw new Error("Database not initialized");
        }
    } catch (err) {
        throw new Error(`Error creating user: ${err}`);
    }
}