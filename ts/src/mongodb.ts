import { MongoClient, ObjectId, UpdateFilter } from "mongodb";
import { FieldsToUpdate, Locator, MongoLocator } from "./model";
import { UpdateData } from "firebase-admin/firestore";

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
            // delete nodes
            for (const nodeLocator of nodesToDelete) {
                const locator = nodeLocator as MongoLocator;
                await db.db().collection(locator.collection).deleteOne(locator.filter, { session });
            }

            // update nodes
            for (const { locator, fieldsToUpdate } of toUpdate) {
                await db.db().collection(locator.collection).updateOne(locator.filter, fieldsToUpdate, { session });
            }
        });
    } catch (err) {
        console.log(err);
    } finally {
        await session.endSession();
        await db.close();
    }
}