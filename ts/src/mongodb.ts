import { MongoClient } from "mongodb";
import { FieldsToUpdate, Locator, MongoLocator } from "./model";

export async function getDocumentFromMongo(db: MongoClient, locator: MongoLocator): Promise<any> {
    return db.db().collection(locator.collection).findOne(locator.filter)
}

export async function getDocumentsFromMongo(db: MongoClient, locator: MongoLocator): Promise<any> {
    return db.db().collection(locator.collection).find(locator.filter).toArray()
}

export async function executeTransactionInMongo(
    db: MongoClient,
    toUpdate: FieldsToUpdate<MongoLocator>[],
    nodesToDelete: Locator[]
) {
    const session = db.startSession();
    try {
        await session.withTransaction(async () => {
            let promises = [];
            // delete nodes
            for (const nodeLocator of nodesToDelete) {
                const locator = nodeLocator as MongoLocator;
                promises.push(db.db().collection(locator.collection).deleteOne(locator.filter, { session }));
            }

            // update nodes
            for (const { locator, fieldsToUpdate } of toUpdate) {
                promises.push(db.db().collection(locator.collection).updateOne(locator.filter, fieldsToUpdate, { session }));
            }

            await Promise.all(promises);
        });
    } catch (err) {
        throw "Transaction aborted: " + err;
    } finally {
        await session.endSession();
        await db.close();
    }
}