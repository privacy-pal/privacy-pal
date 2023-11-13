import { ObjectId } from "mongodb";
import { MongoLocator } from "../../../../src/model";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";
import DirectMessage from "../../model/dm";

export default function handleAccessMongo(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {
    switch (locator.dataType) {
        case 'user':
            const user = obj as User;
            // TODO: should we always convert _id to string after reading from database?
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

function handleAccessGroupChat(dataSubjectId: string, locator: MongoLocator, obj: GroupChat): Record<string, any> {
    return {
        messages: {
            dataType: 'message',
            singleDocument: false,
            collection: FirestoreCollections.Messages,
            filter: {
                userID: dataSubjectId,
                chatID: obj.id
            }
        } as MongoLocator
    };
};

function handleAccessUser(dataSubjectId: string, locator: MongoLocator, obj: User): Record<string, any> {
    if (dataSubjectId !== obj.id) {
        return {
            name: obj.name,
        }
    }

    return {
        name: obj.name,
        groupChats: obj.gcs.map((gc: string): MongoLocator => {
            return {
                dataType: 'groupChat',
                singleDocument: true,
                collection: FirestoreCollections.GroupChat,
                filter: {
                    _id: new ObjectId(gc)
                }
            }
        }),
        directMessages: Object.values(obj.dms).map((dm: string): MongoLocator => {
            return {
                dataType: 'directMessage',
                singleDocument: true,
                collection: FirestoreCollections.DirectMessages,
                filter: {
                    _id: new ObjectId(dm)
                }
            }
        })
    };
}

function handleAccessMessage(dataSubjectId: string, locator: MongoLocator, obj: Message): Record<string, any> {
    return {
        content: obj.content,
        timestamp: obj.timestamp
    };
};

function handleAccessDirectMessage(dataSubjectId: string, locator: MongoLocator, obj: DirectMessage): Record<string, any> {
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
            collection: FirestoreCollections.Users,
            filter: {
                _id: new ObjectId(otherUserId)
            }
        } as MongoLocator,
        messages: {
            dataType: 'message',
            singleDocument: false,
            collection: FirestoreCollections.Messages,
            filter: {
                userID: dataSubjectId,
                chatID: obj.id
            }
        } as MongoLocator
    };
}