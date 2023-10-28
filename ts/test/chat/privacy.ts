import { Locator, LocatorType } from "../../src/model";
import GroupChat from "./model/gc";
import Message from "./model/message";
import { FirestoreCollections } from "./model/shared";
import User from "./model/user";

export default function handleAccess(dataSubjectId: string, locator: Locator, obj: any): Record<string, any> {
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

function handleAccessGroupChat(dataSubjectId: string, locator: Locator, obj: GroupChat): Record<string, any> {
    return {
        messages: {
            locatorType: LocatorType.Collection,
            dataType: 'message',
            collectionPath: [...locator.collectionPath, FirestoreCollections.Messages],
            docIds: locator.docIds,
            queries: [{
                fieldPath: 'userID',
                opStr: '==',
                value: dataSubjectId
            }]
        }
    };
};

function handleAccessUser(dataSubjectId: string, locator: Locator, obj: User): Record<string, any> {
    return {
        name: obj.name,
        groupChats: obj.gcs.map((gc: string): Locator => {
            return {
                locatorType: LocatorType.Document,
                dataType: 'groupChat',
                collectionPath: [FirestoreCollections.GroupChat],
                docIds: [gc],
            }
        })
    };
}

function handleAccessMessage(dataSubjectId: string, locator: Locator, obj: Message): Record<string, any> {
    return {
        content: obj.content,
        timestamp: obj.timestamp
    };
};