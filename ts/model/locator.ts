export interface Locator {
    locatorType: LocatorType;
    dataType: string;
    collectionPath: string[];
    docIds: string[];
    queries?: Query[];
}

export interface Query {
    path: string;
    op: string;
    value: any;
}

export enum LocatorType {
    Document = "document",
    Collection = "collection"
}

export function validateLocator(locator: Locator): Error | null {
    if (!locator.collectionPath || locator.collectionPath.length === 0) {
        return new Error("Locator must have a collectionPath");
    }

    if (locator.locatorType == LocatorType.Document && locator.docIds.length != locator.collectionPath.length) {
        return new Error("Document locator must have the same number of docIds as collectionPath elements");
    }

    if (locator.locatorType == LocatorType.Collection && locator.docIds.length != locator.collectionPath.length - 1) {
        return new Error("Collection locator must have one less docId than collectionPath elements");
    }

    return null;
}

export function isLocator(obj: any): obj is Locator {
    return obj.type && obj.collectionPath && obj.docIds && obj.newDataNode;
}