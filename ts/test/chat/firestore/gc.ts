import { db } from "../../firestore";
import { FirestoreCollections, doesNotExistError } from "../model/shared";
import GroupChat from "../model/gc";

export async function GetGroupChat(ID: string): Promise<GroupChat | null> {
    try {
        const doc = await db.collection(FirestoreCollections.GroupChat).doc(ID).get();

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
    } catch (err) {
        throw new Error(`Error getting group chat: ${err}`);
    }
}