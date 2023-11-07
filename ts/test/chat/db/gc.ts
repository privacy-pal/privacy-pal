import { FirestoreCollections, doesNotExistError } from "../model/shared";
import GroupChat from "../model/gc";
import TestDatabase from "../../testDB";
import { ObjectId } from "mongodb";

export async function GetGroupChat(ID: string): Promise<GroupChat | null> {
    try {
        if (TestDatabase.database === "firestore") {
            const doc = await TestDatabase.firestoreClient.collection(FirestoreCollections.GroupChat).doc(ID).get();

            if (!doc.exists) {
                throw new Error(doesNotExistError);
            }

            const chat = new GroupChat('', []);

            chat.id = doc.id;
            const data = doc.data();

            if (data) {
                chat.owner = data.owner || '';
                chat.users = data.users || [];
                chat.messages = data.messages || [];
            }

            return chat;
        } else if (TestDatabase.database === "mongo") {
            const doc = await TestDatabase.mongoClient.db().collection(FirestoreCollections.GroupChat).findOne({ _id: new ObjectId(ID) });

            if (!doc) {
                throw new Error(doesNotExistError);
            }

            const chat = new GroupChat('', []);

            chat.id = doc._id.toString();
            chat.owner = doc.owner || '';
            chat.users = doc.users || [];
            chat.messages = doc.messages || [];

            return chat;
        } else {
            throw new Error("Database not initialized");
        }
    } catch (err) {
        throw new Error(`Error getting group chat: ${err}`);
    }
}