import { Firestore, UpdateData } from 'firebase-admin/firestore';
import { Filter } from 'mongodb';

// export type FirestoreDb = "firestore";
// export type MongoDb = "mongo";

export type HandleAccessFunc<T extends FirestoreLocator | MongoLocator> =
    (dataSubjectId: string, locator: T, obj: any) => Record<string, any>

// export type HandleAccessFunc<T extends FirestoreDb | MongoDb> =
//     (
//         dataSubjectId: string,
//         locator: T extends FirestoreDb ? FirestoreLocator : MongoLocator,
//         obj: any
//     ) => Record<string, any>

export type HandleDeletionFunc = (dataSubjectId: string) => {
    nodesToTraverse: Locator[],
    deleteNode: boolean,
    updateData?: UpdateData<any>
};


export enum LocatorType {
    Document = "document",
    Collection = "collection"
}

export type Locator = MongoLocator | FirestoreLocator;

interface LocatorBase {
    dataType: string;
    singleDocument: boolean; // whether the output document(s) should be nested in an array
    // TODO: enforce this
}

export interface MongoLocator extends LocatorBase {
    collection: string;
    filter: Filter<any>;
}

export interface FirestoreLocator extends LocatorBase {
    locatorType: LocatorType; // TODO: remove this 
    collectionPath: string[];
    docIds: string[];
    queries?: FirebaseFirestore.Filter[];
}

export function validateLocator(locator: Locator): Error | null {
    // if (!locator.collectionPath || locator.collectionPath.length === 0) {
    //     return new Error("Locator must have a collectionPath");
    // }

    // if (locator.locatorType == LocatorType.Document && locator.docIds.length != locator.collectionPath.length) {
    //     return new Error("Document locator must have the same number of docIds as collectionPath elements");
    // }

    // if (locator.locatorType == LocatorType.Collection && locator.docIds.length != locator.collectionPath.length - 1) {
    //     return new Error("Collection locator must have one less docId than collectionPath elements");
    // }

    return null;
}

export function isLocator<T extends FirestoreLocator | MongoLocator>(obj: any): obj is T {
    return obj.dataType && (obj.singleDocument !== undefined);
}