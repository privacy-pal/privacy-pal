import { FieldValue, UpdateData } from "firebase-admin/firestore";
import { db } from "../../firestore";
import { GetDirectMessage } from "../firestore/dm";
import { GetGroupChat } from "../firestore/gc";
import { GetUser } from "../firestore/user";
import DirectMessage from "./dm";
import GroupChat from "./gc";
import Message from "./message";
import { FirestoreCollections, JoinQuitAction } from "./shared";
import { DataNode } from "../../../model/datanode";
import { Locator, LocatorType } from "../../../model/locator";

export default class User implements DataNode {
    id: string;
    name: string;
    gcs: string[];
    dms: { [key: string]: string };

    constructor(name: string) {
        this.name = name;
        this.gcs = [];
        this.dms = {};
    }

    handleAccess(dataSubjectId: string, locator: Locator): Record<string, any> {
        return {
            name: this.name,
            groupChats: this.gcs.map((gc: string): Locator => {
                return {
                    type: LocatorType.Document,
                    collectionPath: [FirestoreCollections.GroupChat],
                    docIds: [gc],
                }
            })
        };
    }

    handleDeletion(dataSubjectId: string): {
        nodesToTraverse: Locator[],
        deleteNode: boolean,
        updateData?: UpdateData<any>
    } {
        return {nodesToTraverse: [], deleteNode: false};
    }

    async CreateGroupChat(): Promise<GroupChat | null> {
        const newChat = new GroupChat(this.id, []);

        try {
            const ref = await db.collection(FirestoreCollections.GroupChat).add(Object.assign({}, newChat));
            newChat.id = ref.id;

            // Add the group chat to the user
            await db.collection(FirestoreCollections.Users).doc(this.id).set(
                {
                    gcs: FieldValue.arrayUnion(newChat.id),
                },
                { merge: true }
            );

            this.gcs.push(newChat.id);
            return newChat;
        } catch (err) {
            throw new Error(`Error creating group chat: ${err}`);
        }
    }

    async JoinOrQuitGroupChat(chatID: string, action: JoinQuitAction): Promise<void> {
        try {
            // Check if the user and group chat exist
            await GetUser(this.id);
            await GetGroupChat(chatID);

            const updates: UpdateData<{[x: string]: any}> = {};

            if (action === JoinQuitAction.JoinChat) {
                updates.users = FieldValue.arrayUnion(this.id)
            } else if (action === JoinQuitAction.QuitChat) {
                updates.users = FieldValue.arrayRemove(this.id)
            }

            // Update the group chat
            await db.collection(FirestoreCollections.GroupChat).doc(chatID).update(updates);

            // Update the user
            updates.length = 0; // Clear the updates array

            if (action === JoinQuitAction.JoinChat) {
                updates.gcs = FieldValue.arrayUnion(chatID);
            } else if (action === JoinQuitAction.QuitChat) {
                updates.gcs = FieldValue.arrayRemove(chatID);
            }

            await db.collection(FirestoreCollections.Users).doc(this.id).update(updates);
            
        } catch (err) {
            throw new Error(`Error updating user or group chat: ${err}`);
        }
    }

    async CreateDirectMessage(user2ID: string): Promise<DirectMessage | null> {
        try {
            // Check if the user exists and if the direct message already exists
            const user2 = await GetUser(user2ID);

            if (!user2) {
                throw new Error('User does not exist');
            }

            if (user2.dms && user2.dms[this.id]) {
                throw new Error('Direct message already exists');
            }

            if (this.dms && this.dms[user2ID]) {
                throw new Error('Direct message already exists');
            }

            const newDM = new DirectMessage(this.id, user2ID);

            const ref = await db.collection(FirestoreCollections.DirectMessages).add(Object.assign({}, newDM));
            newDM.id = ref.id;

            // Add the DM to both users
            await db.collection(FirestoreCollections.Users).doc(this.id).set(
                {
                    dms: {
                        [user2ID]: newDM.id,
                    },
                },
                { merge: true }
            );

            await db.collection(FirestoreCollections.Users).doc(user2ID).set(
                {
                    dms: {
                        [this.id]: newDM.id,
                    },
                },
                { merge: true }
            );

            return newDM;
        } catch (err) {
            throw new Error(`Error creating direct message: ${err}`);
        }
    }

    async SendMessageToGroupChat(chatID: string, message: string): Promise<void> {
        try {
            // Get the group chat
            const gc = await GetGroupChat(chatID);

            if (!gc) {
                throw new Error('Group chat does not exist');
            }

            // Check if the user is in the chat
            if (!(gc.users?.includes(this.id)) && gc.owner !== this.id) {
                throw new Error('User is not in the group chat');
            }

            // Create a message
            const newMessage = new Message(this.id, message, new Date());

            // Write the message to the Firestore subcollection
            const ref = await db.collection(FirestoreCollections.GroupChat)
                .doc(chatID)
                .collection(FirestoreCollections.Messages)
                .add(Object.assign({}, newMessage));

            newMessage.id = ref.id;
        } catch (err) {
            throw new Error(`Error creating message: ${err}`);
        }
    }

    async SendMessageToDirectMessage(chatID: string, message: string): Promise<void> {
        try {
            // Get the direct message
            const dm = await GetDirectMessage(chatID);

            if (!dm) {
                throw new Error('Direct message does not exist');
            }

            // Check if the user is in the direct message
            if (dm.user1 !== this.id && dm.user2 !== this.id) {
                throw new Error('User is not in the direct message');
            }

            // Create a message
            const newMessage = new Message(this.id, message, new Date());

            // Write the message to the Firestore subcollection
            const ref = await db.collection(FirestoreCollections.DirectMessages)
                .doc(chatID)
                .collection(FirestoreCollections.Messages)
                .add(Object.assign({}, newMessage));

            newMessage.id = ref.id;
        } catch (err) {
            throw new Error(`Error creating message: ${err}`);
        }
    }
}