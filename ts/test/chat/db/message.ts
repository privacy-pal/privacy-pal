import { ObjectId } from "mongodb";
import TestDatabase from "../../testDB";
import DirectMessage from "../model/dm";
import { FirestoreCollections, doesNotExistError } from "../model/shared";

export async function GetMessage(ID: string): Promise<DirectMessage | null> {
    try {
        if (TestDatabase.database === "firestore") {
            const doc = await TestDatabase.firestoreClient.collection(FirestoreCollections.Messages).doc(ID).get();

            if (!doc.exists) {
                throw new Error(doesNotExistError);
            }

            const dm = new DirectMessage('', '');

            dm.id = doc.id;
            const data = doc.data();

            if (data) {
                dm.user1 = data.user1 || '';
                dm.user2 = data.user2 || '';
            }

            return dm;

        } else if (TestDatabase.database === "mongo") {
            const doc = await TestDatabase.mongoDb.collection(FirestoreCollections.Messages).findOne({ _id: new ObjectId(ID) });

            if (!doc) {
                throw new Error(doesNotExistError);
            }

            const dm = new DirectMessage('', '');

            dm.id = doc._id.toString();
            dm.user1 = doc.user1 || '';
            dm.user2 = doc.user2 || '';

            return dm;
        } else {
            throw new Error("Database not initialized");
        }
    } catch (err) {
        throw new Error(`Error getting direct message: ${err}`);
    }
}