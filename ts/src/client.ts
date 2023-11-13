
import { Firestore } from 'firebase-admin/firestore';
import { MongoClient } from 'mongodb';
import Database from './database';
import { FieldsToUpdate, FirestoreLocator, HandleAccessFunc, HandleDeletionFunc, MongoLocator, isLocator, validateLocator } from './model';

class PrivacyPalClient<T extends FirestoreLocator | MongoLocator>{

    db: Database;

    constructor(client: Firestore | MongoClient) {
        this.db = new Database(client);
    }

    async processAccessRequest(handleAccess: HandleAccessFunc<T>, dataSubjectLocator: T, dataSubjectId: string) {
        console.log("Processing access request for data subject " + dataSubjectId);

        if (!dataSubjectLocator.singleDocument) {
            throw new Error("Data subject locator must be a single document");
        }

        const dataSubject = await this.db.getDocument(dataSubjectLocator);
        const data = await this.processAccessRequestHelper(handleAccess, dataSubject, dataSubjectId, dataSubjectLocator);
        return data;
        // TODO: maybe make all timestamp values human readable
    }

    private async processAccessRequestHelper(
        handleAccess: HandleAccessFunc<T>,
        dataNode: any,
        dataSubjectID: string,
        dataNodeLocator: T
    ): Promise<Record<string, any>> {

        const data = handleAccess(dataSubjectID, dataNodeLocator, dataNode);
        let report: Record<string, any> = {};

        for (const [key, value] of Object.entries(data)) {
            if (isLocator(value)) {
                // if locator, recursively process
                const retData = await this.processLocator(handleAccess, value as T, dataSubjectID);
                report[key] = retData;
            } else if (value instanceof Array && value.length > 0 && isLocator(value[0])) {
                // if locator slice, recursively process each locator
                report[key] = [];
                for (const loc of value) {
                    const retData = await this.processLocator(handleAccess, loc, dataSubjectID);
                    report[key].push(retData);
                }
            } else if (value instanceof Map && value.size > 0 && isLocator(value.values().next().value)) {
                // if map, recursively process each locator
                report[key] = new Map<string, any>();
                for (const [k, loc] of Object.entries(value)) {
                    const retData = await this.processLocator(handleAccess, loc, dataSubjectID);
                    report[key].set(k, retData);
                }
            } else {
                // else, directly add to the report
                report[key] = value;
            }
        }

        return report;
    }

    private async processLocator(handleAccess: HandleAccessFunc<T>, locator: T, dataSubjectID: string): Promise<Record<string, any>> {
        const err = validateLocator(locator);
        if (err) {
            throw err;
        }

        if (locator.singleDocument) {
            const dataNode = await this.db.getDocument(locator);
            const retData = await this.processAccessRequestHelper(handleAccess, dataNode, dataSubjectID, locator);
            return retData;
        } else {
            const dataNodes = await this.db.getDocuments(locator);
            const retData: Record<string, any>[] = [];
            for (var dataNode of dataNodes) {
                const currDataNodeData = await this.processAccessRequestHelper(handleAccess, dataNode, dataSubjectID, locator);
                retData.push(currDataNodeData);
            }
            return retData;
        }
    }

    async processDeletionRequest(handleDeletion: HandleDeletionFunc<T>, dataSubjectLocator: T, dataSubjectId: string, test: boolean) {
        const { fieldsToUpdate, nodesToDelete } = await this.processDeletionRequestHelper(handleDeletion, dataSubjectId, dataSubjectLocator);
        if (!test) {
            await this.db.updateAndDelete(fieldsToUpdate, nodesToDelete);
        }
    }

    private async processDeletionRequestHelper(
        handleDeletion: HandleDeletionFunc<T>,
        dataSubjectID: string,
        locator: T,
    ): Promise<{ fieldsToUpdate: FieldsToUpdate<T>[], nodesToDelete: T[] }> {
        let dataNodes: any[] = [];
        if (locator.singleDocument) {
            const node = await this.db.getDocument(locator);
            dataNodes.push(node);
        } else {
            const nodes = await this.db.getDocuments(locator);
            dataNodes = dataNodes.concat(nodes);
        }
        let allFieldsToUpdate: FieldsToUpdate<T>[] = [];
        let allNodesToDelete: T[] = [];

        for (const currentDataNode of dataNodes) {
            const { nodesToTraverse, deleteNode, fieldsToUpdate } = handleDeletion(dataSubjectID, locator, currentDataNode)
            // 1. first recursively process nested nodes
            if (nodesToTraverse.length > 0) {
                for (const nodeLocator of nodesToTraverse) {
                    const { fieldsToUpdate, nodesToDelete } = await this.processDeletionRequestHelper(handleDeletion, dataSubjectID, nodeLocator);
                    allFieldsToUpdate = allFieldsToUpdate.concat(fieldsToUpdate);
                    allNodesToDelete = allNodesToDelete.concat(nodesToDelete);
                }
            }

            // 2. delete current node if needed
            if (deleteNode) {
                allNodesToDelete.push(locator);
            } else if (fieldsToUpdate) {
                allFieldsToUpdate.push({ locator: locator as T, fieldsToUpdate: fieldsToUpdate });
            }
        }

        return { fieldsToUpdate: allFieldsToUpdate, nodesToDelete: allNodesToDelete };
    }

}

export default PrivacyPalClient;