import { ObjectId } from "mongodb";
import PrivacyPalClient from "../../src/client";
import { FirestoreLocator, MongoLocator } from "../../src/model";
import TestDatabase from "../testDB";
import { GetGroupChat } from "./db/gc";
import { CreateUser } from "./db/user";
import { FirestoreCollections, JoinQuitAction } from "./model/shared";
import handleAccessMongo from "./privacy/mongo/access";
import handleAccessFirestore from "./privacy/firestore/access";
import handleDeletionMongo from "./privacy/mongo/deletion";

const DELETION = false

async function testMongo(deletion: boolean) {
    await TestDatabase.initializeDB("mongo");

    // create user 1
    const user1 = await CreateUser('user1');
    console.log('Created User:', user1);

    // create user 2
    const user2 = await CreateUser('user2');
    console.log('Created User:', user2);

    // user1 create groupchat
    let groupChat = await user1.CreateGroupChat();
    console.log('Created Group Chat:', groupChat);

    if (!groupChat) {
        console.error('Group Chat not found');
        return;
    }

    // user2 joins groupchat
    await user2.JoinOrQuitGroupChat(groupChat.id, JoinQuitAction.JoinChat);
    console.log('Joined Group Chat');

    groupChat = await GetGroupChat(groupChat.id);
    if (!groupChat) {
        console.error('Group Chat not found');
        return;
    }

    // user 1 sends message to groupchat
    await user1.SendMessageToGroupChat(groupChat.id, 'hello');

    // user 2 sends message to groupchat
    await user2.SendMessageToGroupChat(groupChat.id, 'hi');

    // user 1 sends another message to groupchat
    await user1.SendMessageToGroupChat(groupChat.id, 'hello again');

    // user 2 creates direct message with user 1
    let dm = await user2.CreateDirectMessage(user1.id);

    if (!dm) {
        console.error('Direct Message not found');
        return;
    }
    // user 2 sends message to direct message
    await user2.SendMessageToDirectMessage(dm.id, "Hey! We are in direct message");

    // user 1 sends message to direct message
    await user1.SendMessageToDirectMessage(dm.id, "Hello!");

    console.log("Starting to test privacy pal")

    const privacyPalClient = new PrivacyPalClient<MongoLocator>(TestDatabase.mongoClient);
    const user1Locator: MongoLocator = {
        dataType: 'user',
        singleDocument: true,
        collection: FirestoreCollections.Users,
        filter: {
            _id: new ObjectId(user1.id)
        }
    }

    const res = await privacyPalClient.processAccessRequest(handleAccessMongo, user1Locator, user1.id)
    console.log(JSON.stringify(res))

    if (deletion) {
        const user2Locator: MongoLocator = {
            dataType: 'user',
            singleDocument: true,
            collection: FirestoreCollections.Users,
            filter: {
                _id: new ObjectId(user2.id)
            }
        }
        await privacyPalClient.processDeletionRequest(handleDeletionMongo, user2Locator, user2.id, false)

        console.log("finished deletion")
        // const res = await privacyPalClient.processAccessRequest(handleAccessMongo, dataSubjectLocator, user1.id)
        // console.log(JSON.stringify(res))
    }
}

async function testFirestore(deletion: boolean = false) {
    await TestDatabase.initializeDB("firestore");

    // create user 1
    const user1 = await CreateUser('user1');
    console.log('Created User:', user1);

    // create user 2
    const user2 = await CreateUser('user2');
    console.log('Created User:', user2);

    // user1 create groupchat
    let groupChat = await user1.CreateGroupChat();
    console.log('Created Group Chat:', groupChat);

    if (!groupChat) {
        console.error('Group Chat not found');
        return;
    }

    // user2 joins groupchat
    await user2.JoinOrQuitGroupChat(groupChat.id, JoinQuitAction.JoinChat);
    console.log('Joined Group Chat');

    console.log(groupChat.id)
    groupChat = await GetGroupChat(groupChat.id);
    if (!groupChat) {
        console.log('Group Chat not found');
        return;
    }

    // user 1 sends message to groupchat
    await user1.SendMessageToGroupChat(groupChat.id, 'hello');

    // user 2 sends message to groupchat
    await user2.SendMessageToGroupChat(groupChat.id, 'hi');

    // user 1 sends another message to groupchat
    await user1.SendMessageToGroupChat(groupChat.id, 'hello again');


    const privacyPalClient = new PrivacyPalClient<FirestoreLocator>(TestDatabase.firestoreClient);

    const dataSubjectLocator: FirestoreLocator = {
        dataType: 'user',
        singleDocument: true,
        collectionPath: ['users'],
        docIds: [user1.id]
    }

    const res = await privacyPalClient.processAccessRequest(handleAccessFirestore, dataSubjectLocator, user1.id)
    console.log(JSON.stringify(res))
}

// testMongo(DELETION).then(() => TestDatabase.cleanupDB())
testFirestore().then(() => TestDatabase.cleanupDB())

