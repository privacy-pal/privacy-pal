
import { Firestore } from 'firebase-admin/firestore';
import { Locator, validateLocator, isLocator } from "./model/locator";
import { DataNode } from './model/datanode';
import { getDocumentFromFirestore, getDocumentsFromFirestore } from "./firestore";

export class PrivacyPalClient {

    db: Firestore;

    constructor(db: Firestore) {
        this.db = db;
    }

    async processAccessRequest(dataSubjectLocator: Locator, dataSubjectId: string) {
        console.log("Processing access request for data subject " + dataSubjectId);

        if (dataSubjectLocator.type != "document") {
            throw new Error("Data subject locator type must be document");
        }

        const dataSubject = await getDocumentFromFirestore(dataSubjectLocator);
        const data = await this.processAccessRequestHelper(dataSubject, dataSubjectId, dataSubjectLocator);
        return data;
    }

    processDeletionRequest(dataSubjectLocator: Locator, dataSubjectId: string) {

    }

    private async processAccessRequestHelper(dataNode: DataNode, dataSubjectID: string, dataNodeLocator: Locator): Promise<Record<string, any>> {
        const data = dataNode.handleAccess(dataSubjectID, dataNodeLocator);
        let report: Record<string, any> = {};

        for (const [key, value] of Object.entries(data)) {
            if (isLocator(value)) {
                // if locator, recursively process
                const retData = await this.processLocator(value, dataSubjectID);
                report[key] = retData;
            } else if (value instanceof Array) {
                // if locator slice, recursively process each locator
                report[key] = [];
                for (const loc of value) {
                    const retData = await this.processLocator(loc, dataSubjectID);
                    report[key].push(retData);
                }
            } else if (value instanceof Map) {
                // if map, recursively process each locator
                report[key] = new Map<string, any>();
                for (const [k, loc] of Object.entries(value)) {
                    const retData = await this.processLocator(loc, dataSubjectID);
                    report[key].set(k, retData);
                }
            } else {
                // else, directly add to the report
                report[key] = value;
            }
        }

        return report;
    }

    private async processLocator(locator: Locator, dataSubjectID: string): Promise<Record<string, any>> {
        const err = validateLocator(locator);
        if (err) {
            throw err;
        }

        if (locator.type == "document") {
            const dataNode = await getDocumentFromFirestore(locator);
            const retData = await this.processAccessRequestHelper(dataNode, dataSubjectID, locator);
            return retData;
        }

        if (locator.type == "collection") {
            const dataNodes = await getDocumentsFromFirestore(locator);
            const retData: Record<string, any>[] = [];
            for (var dataNode of dataNodes) {
                const currDataNodeData = await this.processAccessRequestHelper(dataNode, dataSubjectID, locator);
                retData.push(currDataNodeData);
            }
            return retData;
        }

        throw new Error("Invalid locator type");
    }
}