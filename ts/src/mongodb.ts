import { MongoClient, ObjectId } from "mongodb";
import { MongoLocator } from "./model";

export async function getDocumentFromMongo(db: MongoClient, locator: MongoLocator): Promise<any> {
    return db.db().collection(locator.collection).findOne(locator.filter)
}

export async function getDocumentsFromMongo(db: MongoClient, locator: MongoLocator): Promise<any> {
    return db.db().collection(locator.collection).find(locator.filter).toArray()
}