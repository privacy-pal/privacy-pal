
import { CollectionReference, Firestore, Query, WhereFilterOp } from 'firebase-admin/firestore';
import { Locator, validateLocator, isLocator } from "./model/locator";
import { DataNode } from './model/datanode';

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

        const dataSubject = await this.getDocumentFromFirestore(dataSubjectLocator);
        const data = await this.processAccessRequestHelper(dataSubject, dataSubjectId, dataSubjectLocator);
        return data;
    }

    processDeletionRequest(dataSubjectLocator: Locator, dataSubjectId: string) {

    }

    private async getDocumentFromFirestore(locator: Locator): Promise<DataNode> {
        let docRef = this.db.collection(locator.collectionPath[0]).doc(locator.docIds[0]);

        for (let i = 1; i < locator.collectionPath.length; i++) {
            docRef = docRef.collection(locator.collectionPath[i]).doc(locator.docIds[i]);
        }

        return docRef.get()
            .then((doc) => {
                if (!doc.exists) {
                    throw new Error("Document does not exist");
                }
                return doc.data() as DataNode;
            })
            .catch((err) => {
                throw new Error('Error getting document: ' + err);
            });
    }

    private async getDocumentsFromFirestore(locator: Locator): Promise<DataNode[]> {
        let docRef: CollectionReference = this.db.collection(locator.collectionPath[0]);

        for (let i = 1; i < locator.collectionPath.length; i++) {
            docRef = docRef.doc(locator.docIds[i - 1]).collection(locator.collectionPath[i]);
        }

        let query: Query = docRef;
        if (locator.queries?.length) {
            query = query.where(locator.queries[0].path, locator.queries[0].op as WhereFilterOp, locator.queries[0].value);
            for (let i = 1; i < locator.queries.length; i++) {
                query = query.where(locator.queries[i].path, locator.queries[i].op as WhereFilterOp, locator.queries[i].value);
            }
        }

        
        return query.get()
            .then((snapshot) => {
                let dataNodes: DataNode[] = [];
                snapshot.forEach((doc) => {
                    // TODO: verify
                    dataNodes.push(doc.data() as DataNode);
                });
                return dataNodes;
            })
            .catch((err) => {
                throw new Error('Error getting documents: ' + err);
            });
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
            const dataNode = await this.getDocumentFromFirestore(locator);
            const retData = await this.processAccessRequestHelper(dataNode, dataSubjectID, locator);
            return retData;
        }

        if (locator.type == "collection") {
            const dataNodes = await this.getDocumentsFromFirestore(locator);
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