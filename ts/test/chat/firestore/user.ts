import { db } from "../../firestore";
import { FirestoreCollections, doesNotExistError } from "../model/shared";
import User from "../model/user";

export async function GetUser(ID: string): Promise<User | undefined> {
    try {
        const doc = await db.collection(FirestoreCollections.Users).doc(ID).get();

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
    } catch (err) {
        throw new Error(`Error getting user: ${err}`);
    }
}

export async function CreateUser(name: string): Promise<User> {
    try {
        const user = new User(name);
        const docRef = await db.collection(FirestoreCollections.Users).add(Object.assign({}, user));
        user.id = docRef.id;
        return user;
    } catch (err) {
        throw new Error(`Error creating user: ${err}`);
    }
}