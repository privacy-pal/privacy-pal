import { ObjectId, UpdateFilter } from "mongodb";
import { MongoLocator } from "../../../../src";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";

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
    const deleteNode = obj.owner === dataSubjectId;
    let fieldsToUpdate: UpdateFilter<any> | Partial<any> = {};
    if (!deleteNode) {
        // remove dataSubjectId from users field
        fieldsToUpdate = {
            $pull: {
                users: dataSubjectId
            }
        }
    }
    return {
        nodesToTraverse: [{
            dataType: 'message',
            singleDocument: false,
            collection: FirestoreCollections.Messages,
            filter: {
                userID: dataSubjectId,
                chatID: obj.id
            }
        }],
        deleteNode: deleteNode,
        fieldsToUpdate: fieldsToUpdate
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
    return {
        nodesToTraverse: obj.gcs.map((gc: string): MongoLocator => {
            return {
                dataType: 'groupChat',
                singleDocument: true,
                collection: FirestoreCollections.GroupChat,
                filter: {
                    _id: new ObjectId(gc)
                }
            }
        }),
        deleteNode: true,
    };
}

function handleDeletionMessage(dataSubjectId: string, locator: MongoLocator, obj: Message): {
    nodesToTraverse: MongoLocator[],
    deleteNode: boolean,
    fieldsToUpdate?: UpdateFilter<any> | Partial<any>
} {
    return {
        nodesToTraverse: [],
        deleteNode: obj.userID === dataSubjectId,
    };
};