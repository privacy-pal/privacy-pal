import { GetGroupChat } from "./firestore/gc";
import { CreateUser } from "./firestore/user";
import { JoinQuitAction } from "./model/shared";
import PrivacyPalClient from "../../src/client";
import { db } from "../firestore";
import { FirestoreLocator, Locator, LocatorType } from "../../src/model";
import handleAccess from "./privacy";

async function test1() {
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

    const privacyPalClient = new PrivacyPalClient<FirestoreLocator>(db);
    const dataSubjectLocator: FirestoreLocator = {
        dataType: 'user',
        singleDocument: true,
        locatorType: LocatorType.Document,
        collectionPath: ['users'],
        docIds: [user1.id]
    }

    const res = await privacyPalClient.processAccessRequest(handleAccess, dataSubjectLocator, user1.id)
    console.log(JSON.stringify(res))
}

test1();
