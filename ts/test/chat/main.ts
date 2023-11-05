import { ObjectId } from "mongodb";
import PrivacyPalClient from "../../src/client";
import { FirestoreLocator, MongoLocator } from "../../src/model";
import TestDatabase from "../testDB";
import { GetGroupChat } from "./firestore/gc";
import { CreateUser } from "./firestore/user";
import { FirestoreCollections, JoinQuitAction } from "./model/shared";
import handleAccessMongo from "./privacyMongo";
import handleAccessFirestore from "./privacyFirestore";

async function testMongo() {
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

    console.log("Starting to test privacy pal")

    const privacyPalClient = new PrivacyPalClient<MongoLocator>(TestDatabase.mongoClient);
    const dataSubjectLocator: MongoLocator = {
        dataType: 'user',
        singleDocument: true,
        collection: FirestoreCollections.Users,
        filter: {
            _id: new ObjectId(user1.id)
        }
    }

    const res = await privacyPalClient.processAccessRequest(handleAccessMongo, dataSubjectLocator, user1.id)
    console.log(JSON.stringify(res))
}

async function testFirestore() {
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

testMongo().then(() => TestDatabase.cleanupDB())
// testFirestore().then(() => TestDatabase.cleanupDB())

