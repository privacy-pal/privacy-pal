import { UpdateData } from "firebase-admin/firestore";
import { DataNode } from "../../../model/datanode";
import { Locator, LocatorType } from "../../../model/locator";
import Message from "./message"
import { FirestoreCollections } from "./shared";

export default class GroupChat implements DataNode{
    id: string;
    users: string[];
    messages: Message[];
    owner: string;

    constructor(owner: string, users: string[]) {
        this.owner = owner;
        this.users = users;
        this.messages = [];
    }

    handleAccess(dataSubjectId: string, locator: Locator): Record<string, any> {
        return {
            messages: {
                    type: LocatorType.Collection,
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

    handleDeletion(dataSubjectId: string): {
        nodesToTraverse: Locator[],
        deleteNode: boolean,
        updateData?: UpdateData<any>
    } {
        return {nodesToTraverse: [], deleteNode: false};
    }
}