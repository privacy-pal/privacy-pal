import { ObjectId } from "mongodb";
import { MongoLocator } from "../../../../src/model";
import GroupChat from "../../model/gc";
import Message from "../../model/message";
import { FirestoreCollections } from "../../model/shared";
import User from "../../model/user";

export default function handleAccessMongo(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {
    switch (locator.dataType) {
        case 'user':
            const user = obj as User;
            user.id = obj._id.toString();
            return handleAccessUser(dataSubjectId, locator, user);
        case 'groupChat':
            const groupChat = obj as GroupChat;
            groupChat.id = obj._id.toString();
            return handleAccessGroupChat(dataSubjectId, locator, groupChat);
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
        })
    };
}

function handleAccessMessage(dataSubjectId: string, locator: MongoLocator, obj: Message): Record<string, any> {
    return {
        content: obj.content,
        timestamp: obj.timestamp
    };
};