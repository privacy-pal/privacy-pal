import TestDatabase from "../../testDB";
import DirectMessage from "../model/dm";
import { FirestoreCollections, doesNotExistError } from "../model/shared";

export async function GetDirectMessage(ID: string): Promise<DirectMessage | null> {
    try {
        const doc = await TestDatabase.firestoreClient.collection(FirestoreCollections.DirectMessages).doc(ID).get();

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
    } catch (err) {
        throw new Error(`Error getting direct message: ${err}`);
    }
}