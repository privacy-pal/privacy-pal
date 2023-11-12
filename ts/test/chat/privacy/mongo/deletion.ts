import { ObjectId, UpdateFilter } from "mongodb";
import { MongoLocator } from "../../../../src";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";
import DirectMessage from "../../model/dm";

export default function handleDeletionMongo(dataSubjectId: string, locator: MongoLocator, obj: any): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    switch (locator.dataType) {
        case 'user':
            const user = obj as User;
            user.id = obj._id.toString();
            return handleDeletionUser(dataSubjectId, locator, user);
        case 'groupChat':
            const groupChat = obj as GroupChat;
            groupChat.id = obj._id.toString();
            return handleDeletionGroupChat(dataSubjectId, locator, groupChat);
        case 'directMessage':
            const directMessage = obj as DirectMessage;
            directMessage.id = obj._id.toString();
            return handleDeletionDirectMessage(dataSubjectId, locator, directMessage);
        case 'message':
            const message = obj as Message;
            message.id = obj._id.toString();
            return handleDeletionMessage(dataSubjectId, locator, message);
        default:
            throw new Error(`Data type ${locator.dataType} does not have handleDeletion function implemented`);
    }
}

function handleDeletionGroupChat(dataSubjectId: string, locator: MongoLocator, obj: GroupChat): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    return {
        nodesToTraverse: [{
            dataType: 'message',
            singleDocument: false,
            collection: FirestoreCollections.Messages,
            filter: {
                userID: dataSubjectId,
                chatID: obj.id
            },
            context: {
                anonymize: true
            }
        }],
        deleteNode: obj.owner === dataSubjectId,
        fieldsToUpdate: obj.owner === dataSubjectId ? undefined :
            // remove dataSubjectId from users field
            {
                $pull: {
                    users: dataSubjectId
                }
            }
    };
};

function handleDeletionUser(dataSubjectId: string, locator: MongoLocator, obj: User): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    if (obj.id !== dataSubjectId) {
        throw new Error(`User ${dataSubjectId} cannot delete another user ${obj.id}`);
    }

    const GCsToTraverse = obj.gcs.map((gc: string): MongoLocator => {
        return {
            dataType: 'groupChat',
            singleDocument: true,
            collection: FirestoreCollections.GroupChat,
            filter: {
                _id: new ObjectId(gc)
            }
        }
    })

    const DMsToTraverse = Object.values(obj.dms).map((dm: string): MongoLocator => {
        return {
            dataType: 'directMessage',
            singleDocument: true,
            collection: FirestoreCollections.DirectMessages,
            filter: {
                _id: new ObjectId(dm)
            }
        }
    })

    return {
        nodesToTraverse: [...GCsToTraverse, ...DMsToTraverse],
        deleteNode: true,
    };
}

function handleDeletionDirectMessage(dataSubjectId: string, locator: MongoLocator, obj: DirectMessage): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    let otherUserId: string;
    let thisUserField: string;
    if (obj.user1 === dataSubjectId) {
        thisUserField = 'user1';
        otherUserId = obj.user2;
    } else if (obj.user2 === dataSubjectId) {
        thisUserField = 'user2';
        otherUserId = obj.user1;
    } else {
        throw new Error(`User ${dataSubjectId} does not have access to direct message ${obj.id}`);
    }

    return {
        // TODO: locator should probably have some sort of context to be passed to the next object
        nodesToTraverse: [{
            dataType: 'message',
            singleDocument: false,
            collection: FirestoreCollections.Messages,
            filter: {
                userID: dataSubjectId,
                chatID: obj.id
            },
            context: {
                anonymize: true
            }
        }],
        deleteNode: otherUserId === "anonymous",
        fieldsToUpdate: {
            $set: {
                [thisUserField]: 'anonymous'
            }
        }
    };
}

function handleDeletionMessage(dataSubjectId: string, locator: MongoLocator, obj: Message): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    const fieldsToUpdate = (obj.userID === dataSubjectId && locator.context.anonymize) ?
        {
            $set: {
                userID: 'anonymous'
            }
        } : undefined;
    return {
        nodesToTraverse: [],
        deleteNode: obj.userID === dataSubjectId && !locator.context.anonymize,
        fieldsToUpdate: fieldsToUpdate
    };
};