
import { CollectionReference, Firestore, Query, WhereFilterOp } from 'firebase-admin/firestore';
import { Locator, validateLocator, isLocator } from "./model/locator";
import { DataNode } from './model/datanode';

export class PrivacyPalClient {

    db: Firestore;

    constructor(db: Firestore) {
        this.db = db;
    }

    processAccessRequest(dataSubjectLocator: Locator, dataSubjectId: string) {
        console.log("Processing access request for data subject " + dataSubjectId);

        if (dataSubjectLocator.type != "document") {
            throw new Error("Data subject locator type must be document");
        }

        const dataSubject = this.getDocumentFromFirestore(dataSubjectLocator);
        const data = this.processAccessRequestHelper(dataSubject, dataSubjectId, dataSubjectLocator);
        return data;
    }

    processDeletionRequest(dataSubjectLocator: Locator, dataSubjectId: string) {

    }

    private getDocumentFromFirestore(locator: Locator): DataNode {
        let docRef = this.db.collection(locator.collectionPath[0]).doc(locator.docIds[0]);

        for (let i = 1; i < locator.collectionPath.length; i++) {
            docRef = docRef.collection(locator.collectionPath[i]).doc(locator.docIds[i]);
        }

        let dataNode: DataNode = new locator.newDataNode();
        docRef.get()
            .then((doc) => {
                if (!doc.exists) {
                    throw new Error("Document does not exist");
                }
                dataNode = doc.data() as (typeof dataNode);
            })
            .catch((err) => {
                throw new Error('Error getting document: ' + err);
            });

        return dataNode;
    }

    private getDocumentsFromFirestore(locator: Locator): DataNode[] {
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

        let dataNodes: DataNode[] = [];

        query.get()
            .then((snapshot) => {
                snapshot.forEach((doc) => {
                    // TODO: verify
                    const newDataNode = new locator.newDataNode();
                    dataNodes.push(doc.data() as (typeof newDataNode));
                });
            })
            .catch((err) => {
                throw new Error('Error getting documents: ' + err);
            });

        return dataNodes;
    }

    private processAccessRequestHelper(dataNode: DataNode, dataSubjectID: string, dataNodeLocator: Locator): Record<string, any> {
        const data = dataNode.handleAccess(dataSubjectID, dataNodeLocator);
        let report: Record<string, any> = {};

        Object.entries(data).forEach(([key, value]) => {
            if (isLocator(value)) {
                // if locator, recursively process
                const retData = this.processLocator(report, value, dataSubjectID);
                report[key] = retData;
            } else if (value instanceof Array) {
                // if locator slice, recursively process each locator
                report[key] = [];
                value.forEach((loc) => {
                    const retData = this.processLocator(report, loc, dataSubjectID);
                    report[key].push(retData);
                });
            } else if (value instanceof Map) {
                // if map, recursively process each locator
                report[key] = new Map<string, any>();
                Object.entries(value).forEach(([k, loc]) => {
                    const retData = this.processLocator(report, loc, dataSubjectID);
                    report[key].set(k, retData);
                });
            } else {
                // else, directly add to report
                report[key] = value;
            }
        });

        return report;

    }

    private processLocator(report: Record<string, any>, locator: Locator, dataSubjectID: string): Record<string, any> {
        const err = validateLocator(locator);
        if (err) {
            throw err;
        }

        if (locator.type == "document") {
            const dataNode = this.getDocumentFromFirestore(locator);
            const retData = this.processAccessRequestHelper(dataNode, dataSubjectID, locator);
            return retData;
        }

        if (locator.type == "collection") {
            const dataNodes = this.getDocumentsFromFirestore(locator);
            const retData: Record<string, any>[] = [];
            dataNodes.forEach((dataNode) => {
                const currDataNodeData = this.processAccessRequestHelper(dataNode, dataSubjectID, locator);
                retData.push(currDataNodeData);
            });
            return retData;
        }

        throw new Error("Invalid locator type");
    }
}