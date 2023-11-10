import { CollectionReference, FieldPath, Firestore, Query, UpdateData, WhereFilterOp } from "firebase-admin/firestore";
import { FieldsToUpdate, FirestoreLocator } from "./model";

export async function getDocumentFromFirestore(db: Firestore, locator: FirestoreLocator): Promise<any> {
    let docRef = db.collection(locator.collectionPath[0]).doc(locator.docIds[0]);

    for (let i = 1; i < locator.collectionPath.length; i++) {
        docRef = docRef.collection(locator.collectionPath[i]).doc(locator.docIds[i]);
    }

    return docRef.get()
        .then((doc) => {
            if (!doc.exists) {
                throw new Error("Document does not exist");
            }
            return doc.data();
        })
        .catch((err) => {
            throw new Error('Error getting document: ' + err);
        });
}

export async function getDocumentsFromFirestore(db: Firestore, locator: FirestoreLocator): Promise<any[]> {
    let docRef: CollectionReference = db.collection(locator.collectionPath[0]);

    for (let i = 1; i < locator.collectionPath.length; i++) {
        docRef = docRef.doc(locator.docIds[i - 1]).collection(locator.collectionPath[i]);
    }

    let query: Query = docRef;
    if (locator.queries?.length) {
        query = query.where(locator.queries[0]);
        for (let i = 1; i < locator.queries.length; i++) {
            query = query.where(locator.queries[i]);
        }
    }


    return query.get()
        .then((snapshot) => {
            let dataNodes: any[] = [];
            snapshot.forEach((doc) => {
                // TODO: verify
                dataNodes.push(doc.data());
            });
            return dataNodes;
        })
        .catch((err) => {
            throw new Error('Error getting documents: ' + err);
        });
}

export async function executeTransactionInFirestore(
    db: Firestore,
    toUpdate: FieldsToUpdate<FirestoreLocator>[],
    nodesToDelete: FirestoreLocator[]
) {
    try {
        await db.runTransaction(async (t) => {
            // delete nodes
            let promises = [];
            for (const nodeLocator of nodesToDelete) {
                let docRef = db.collection(nodeLocator.collectionPath[0]).doc(nodeLocator.docIds[0]);

                for (let i = 1; i < nodeLocator.collectionPath.length; i++) {
                    docRef = docRef.collection(nodeLocator.collectionPath[i]).doc(nodeLocator.docIds[i]);
                }

                promises.push(t.delete(docRef));
            }

            // update nodes
            for (const { locator, fieldsToUpdate } of toUpdate) {
                let docRef = db.collection(locator.collectionPath[0]).doc(locator.docIds[0]);

                for (let i = 1; i < locator.collectionPath.length; i++) {
                    docRef = docRef.collection(locator.collectionPath[i]).doc(locator.docIds[i]);
                }

                promises.push(t.update(docRef, fieldsToUpdate));
            }

            await Promise.all(promises);
        });
        console.log("Privacy Pal: successfully updated and deleted data")
    } catch (err) {
        console.log(err)
    }
}