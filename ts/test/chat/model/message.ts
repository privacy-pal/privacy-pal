import { UpdateData } from "firebase-admin/firestore";
import { DataNode } from "../../../model/datanode";
import { Locator } from "../../../model/locator";

export default class Message implements DataNode{
    id: string;
    userID: string;
    content: string;
    timestamp: Date;

    constructor(userID: string, content: string, timestamp: Date) {
        this.userID = userID;
        this.content = content;
        this.timestamp = timestamp;
    }

    handleAccess(dataSubjectId: string, locator: Locator): Record<string, any> {
        return {
            content: this.content,
            timestamp: this.timestamp
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