import { MongoClient, Db } from "mongodb";
import { DocumentUpdates, Locator, MongoLocator } from "./model";

export async function getDocumentFromMongo(db: Db, locator: MongoLocator): Promise<any> {
    return db.collection(locator.collection).findOne(locator.filter)
}

export async function getDocumentsFromMongo(db: Db, locator: MongoLocator): Promise<any> {
    return db.collection(locator.collection).find(locator.filter).toArray()
}

export async function executeTransactionInMongo(
    client: MongoClient,
    db: Db,
    toUpdate: DocumentUpdates<MongoLocator>[],
    nodesToDelete: Locator[]
) {
    const session = client.startSession();
    return session.withTransaction(async () => {
        let promises = [];
        // delete nodes
        for (const nodeLocator of nodesToDelete) {
            const locator = nodeLocator as MongoLocator;
            promises.push(db.collection(locator.collection).deleteOne(locator.filter, { session }));
        }

        // update nodes
        for (const { locator, fieldsToUpdate } of toUpdate) {
            promises.push(db.collection(locator.collection).updateOne(locator.filter, fieldsToUpdate, { session }));
        }

        await Promise.all(promises);
    }).finally(() => session.endSession());
}