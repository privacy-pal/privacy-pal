import { Filter } from "firebase-admin/firestore";
import { FirestoreLocator } from "../../../../src";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";

export default function handleAccessFirestore(dataSubjectId: string, locator: FirestoreLocator, obj: any): Record<string, any> {
    switch (locator.dataType) {
        case 'user':
            return handleAccessUser(dataSubjectId, locator, obj as User);
        case 'groupChat':
            return handleAccessGroupChat(dataSubjectId, locator, obj as GroupChat);
        case 'message':
            return handleAccessMessage(dataSubjectId, locator, obj as Message);
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
            queries: [Filter.where('userID', '==', dataSubjectId)]
        } as FirestoreLocator
    };
};

function handleAccessUser(dataSubjectId: string, locator: FirestoreLocator, obj: User): Record<string, any> {
    return {
        name: obj.name,
        groupChats: obj.gcs.map((gc: string): FirestoreLocator => {
            return {
                dataType: 'groupChat',
                singleDocument: true,
                collectionPath: [FirestoreCollections.GroupChat],
                docIds: [gc],
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