import { UpdateData } from 'firebase-admin/firestore';
import { Filter, UpdateFilter } from 'mongodb';

export type HandleAccessFunc<T extends FirestoreLocator | MongoLocator> =
    (dataSubjectId: string, locator: T, obj: any) => Record<string, any>

export type HandleDeletionFunc<T extends FirestoreLocator | MongoLocator> = (dataSubjectId: string, locator: T) => {
    nodesToTraverse: T[],
    deleteNode: boolean,
    fieldsToUpdate?: T extends MongoLocator ? UpdateFilter<any>[] : UpdateData<any>[]
};

export type Locator = FirestoreLocator | MongoLocator;

export interface MongoLocator extends LocatorBase {
    collection: string;
    filter: Filter<any>;
}

export interface FirestoreLocator extends LocatorBase {
    collectionPath: string[];
    docIds: string[];
    queries?: FirebaseFirestore.Filter[];
}

// internal 
export interface FieldsToUpdate<T extends FirestoreLocator | MongoLocator> {
    locator: T;
    fieldsToUpdate: T extends MongoLocator ? UpdateFilter<any> : UpdateData<any>;
}

interface LocatorBase {
    dataType: string;
    singleDocument: boolean; // whether the output document(s) should be nested in an array
}

export function validateLocator(locator: FirestoreLocator | MongoLocator): Error | null {
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