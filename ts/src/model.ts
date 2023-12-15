import { UpdateData } from 'firebase-admin/firestore';
import { Filter, UpdateFilter } from 'mongodb';

export type HandleAccessFunc<T extends Locator> =
    (dataSubjectId: string, locator: T, databaseObject: any) => Record<string, any>

export type HandleDeletionFunc<T extends Locator> =
    (dataSubjectId: string, locator: T, databaseObject: any) => {
        nodesToTraverse: T[],
        deleteNode: boolean,
        fieldsToUpdate?: T extends MongoLocator ? UpdateFilter<any> | Partial<any> : UpdateData<any>
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
export interface DocumentUpdates<T extends Locator> {
    locator: T;
    fieldsToUpdate: T extends MongoLocator ? UpdateFilter<any> : UpdateData<any>;
}

interface LocatorBase {
    dataType: string;
    singleDocument: boolean; // whether the output document(s) should be nested in an array
    context?: any; // any additional information to be passed to the handleAccess or handleDeletion function
}

export function validateLocator(locator: Locator): Error | null {
    if (isFirestoreLocator(locator)) {
        if (!locator.collectionPath || locator.collectionPath.length === 0) {
            return new Error("Locator must have a collectionPath");
        }

        if (locator.singleDocument) {
            if (locator.docIds.length != locator.collectionPath.length) {
                return new Error("Single Document locator must have the same number of docIds as collectionPath elements");
            }
        } else {
            if (locator.docIds.length != locator.collectionPath.length - 1) {
                return new Error("Multi Document locator must have one less docId than collectionPath elements");
            }
        }
    }
    return null;
}

export function isLocator<T extends Locator>(obj: any): obj is T {
    return obj && obj.dataType && (obj.singleDocument !== undefined);
}

function isFirestoreLocator(obj: any): obj is FirestoreLocator {
    return obj.collectionPath && obj.docIds;
}