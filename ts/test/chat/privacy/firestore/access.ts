import { Filter } from "firebase-admin/firestore";
import { FirestoreLocator } from "../../../../src";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";
import DirectMessage from "../../model/dm";

export default function handleAccessFirestore(dataSubjectId: string, locator: FirestoreLocator, obj: any): Record<string, any> {
    switch (locator.dataType) {
        case 'user':
            const user = obj as User;
            user.id = obj._id.toString();
            return handleAccessUser(dataSubjectId, locator, user);
        case 'groupChat':
            const groupChat = obj as GroupChat;
            groupChat.id = obj._id.toString();
            return handleAccessGroupChat(dataSubjectId, locator, groupChat);
        case 'directMessage':
            const directMessage = obj as DirectMessage;
            directMessage.id = obj._id.toString();
            return handleAccessDirectMessage(dataSubjectId, locator, directMessage);
        case 'message':
            const message = obj as Message;
            message.id = obj._id.toString();
            return handleAccessMessage(dataSubjectId, locator, message);
        default:
            throw new Error(`Data type ${locator.dataType} does not have handleAccess function implemented`);
    }
}

function handleAccessGroupChat(dataSubjectId: string, locator: FirestoreLocator, obj: GroupChat): Record<string, any> {
    return {
        messages: {
            dataType: 'message',
            singleDocument: false,
            collectionPath: [...locator.collectionPath, FirestoreCollections.Messages],
            docIds: locator.docIds,
            filters: [Filter.where('userID', '==', dataSubjectId)]
        } as FirestoreLocator
    };
};

function handleAccessUser(dataSubjectId: string, locator: FirestoreLocator, obj: User): Record<string, any> {
    if (dataSubjectId !== obj.id) {
        return {
            name: obj.name,
        }
    }
    return {
        name: obj.name,
        groupChats: obj.gcs.map((gc: string): FirestoreLocator => {
            return {
                dataType: 'groupChat',
                singleDocument: true,
                collectionPath: [FirestoreCollections.GroupChat],
                docIds: [gc],
            }
        }),
        directMessages: Object.values(obj.dms).map((dm: string): FirestoreLocator => {
            return {
                dataType: 'directMessage',
                singleDocument: true,
                collectionPath: [FirestoreCollections.DirectMessages],
                docIds: [dm],
            }
        })
    };
}

function handleAccessMessage(dataSubjectId: string, locator: FirestoreLocator, obj: Message): Record<string, any> {
    return {
        content: obj.content,
        timestamp: obj.timestamp
    };
};

function handleAccessDirectMessage(dataSubjectId: string, locator: FirestoreLocator, obj: DirectMessage): Record<string, any> {
    let otherUserId: string;
    if (obj.user1 === dataSubjectId) {
        otherUserId = obj.user2;
    } else {
        otherUserId = obj.user1;
    }

    return {
        otherUser: {
            dataType: 'user',
            singleDocument: true,
            collectionPath: [FirestoreCollections.Users],
            docIds: [otherUserId]
        } as FirestoreLocator,
        messages: {
            dataType: 'message',
            singleDocument: false,
            collectionPath: [...locator.collectionPath, FirestoreCollections.Messages],
            docIds: locator.docIds,
            filters: [Filter.where('userID', '==', dataSubjectId)]
        } as FirestoreLocator
    };
};