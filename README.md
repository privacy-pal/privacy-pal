# Privacy Pal

Privacy Pal is a framework that helps developers handle data subject access and deletion requests for applications using NoSQL databases. This repo contains the implementation of a Privacy Pal client that supports applications using MongoDB or Firestore. Privacy Pal is available as a Go module and a TypeScript npm package.
Privacy Pal is a framework that helps developers handle data subject access and deletion requests for applications using NoSQL databases. This repo contains the implementation of a Privacy Pal client that supports applications using MongoDB or Firestore. Privacy Pal is available as a Go module and a TypeScript npm package.

## Setup

### Go module
### Go module

To install:
To install:

```bash
go get github.com/privacy-pal/privacy-pal
```

To import:
To import:

```golang
import pal "github.com/privacy-pal/privacy-pal/go/pkg"
```


### TypeScript npm package

To install:

```bash
npm i privacy-pal
```

To import:

```typescript
import * from 'privacy-pal';
```

## Overview

Privacy Pal helps developers focus on the core data ownership logic for fulfilling privacy requests as it handles the logic to retrieve, update, and delete data.

To use Privacy Pal, follow these steps.

1. Define `HandleAccess` and `HandleDeletion` functions, specifying what data should be returned, modified, or deleted when fulfilling access or deletion requests. See more details in [Implementing HandleAccess](#implementing-handleaccess) and [Implementing HandleDeletion](#implementing-handledeletion).

2. Initialize the Privacy Pal client with existing database client (MongoDB or Firestore).

3. Invoke `ProcessAccessRequest` and `ProcessDeletionRequest` methods of the Privacy Pal client, passing in the corresponding handler functions (`HandleAccess` or `HandleDeletion`) and a data subject locator, to fulfill access and deletion requests respectively.
   - A data subject locator is a [locator](#locator) object used to retrieve the initial document storing user data. This document is the starting point for handling personal data stored in other documents.

## Core Concepts

### Locator

Privacy Pal uses locators to locate documents in the database. When you invoke the `ProcessXXXRequest` methods and implement your `HandleXXX` functions, you will need to use Locators.

A locator needs to contain the following basic information:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| LocatorType / IsSingleDocument   | Enum / Boolean      | Whether the locator points to a single document (`Document`) or a collection of documents (`Collection`). |
| DataType       | String     | The type of data represented by the document retrieved using this locator. Can be used in `HandleAccess` and `HandleDeletion` to determine what type of document is being processed.  <br/><br/>  We recommend choosing a name that is consistent with the naming of your application's structs/types. E.g. if the locator points to a "user" document, the `DataType` can be set to "user". |
|Context (Optional)  | Any    | Arbitrary metadata that can be used to pass additional information to `HandleAccess` and `HandleDeletion`. See an example in [HandleDeletion](#implementing-handledeletion)|

Since MongoDB and Firestore support different ways of locating documents, the Locator can contain either a MongoLocator or a FirestoreLocator, **but not both**.

MongoLocator uses the following information to locate:

| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| Collection | String | The name of the collection containing the target document(s). |
| Filter         | MongoDB Filter | The filter used to find the target document(s) from the `collection`. |

Example of a locator that locates a user with id "123" in MongoDB:

```golang
pal.Locator{
    LocatorType:         pal.Document,
    DataType:            "user",
    MongoLocator:        &pal.MongoLocator{
        Collection: "users",
        Filter: bson.M{
            "_id": bson.M{
                "$eq": "123",
            },
        },
    },
}
```

FirestoreLocator uses the following information to locate:

| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| CollectionPath | List of Strings | The path of collections leading up to a particular Firestore document or collection. E.g. `["courses]` represents the courses collection and `["courses", "sections"]` represents the sections subcollection under a document in the courses collection |
| DocIDs         | List of Strings | Document IDs in the order they appear within collections. <br/><br/> To locate a single document, `DocIDs` should have the same length as `CollectionPath`. To locate a collection or a set of documents, `DocsIDs` should be one item shorter than the length of `CollectionPath`. |
| Filters        | Optional List of Firestore Filters | The Firestore filters that can be applied to filter out set of documents from a collection of documents. Hence, it's not applicable to the single document locator type. |

Example of a locator that locates a user with id "123" in Firestore:

```golang
pal.Locator{
    LocatorType:         pal.Document,
    DataType:            "user",
    FirestoreLocator:    &pal.FirestoreLocator{
        CollectionPath: []string{"users"},
        DocIDs:         []string{"123"},
    },
}
```

### Implementing HandleAccess

`ProcessAccessRequest` uses the `HandleAccess` function you define to retrieve all personal data for a particular user. 

`HandleAccess` has access to 3 values through its parameters:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| dataSubjectID    | String       | The ID of the data subject we are building the access report for. You will pass this value in when `ProcessAccessRequest` is invoked. <br/> This ID can be used to determine whether a piece of data belongs to the data subject. |
| locator    | Locator       | The current [locator](#locator) used to retreive document(s) from the database, containing information such as the type of data retrieved. |
| dbObj    | DatabaseObject      | The document retrieved from the database using the locator, augmented with an `_id` field that contains a string |

You should specify (for each object/document type), (1) which data fields are personal data belonging to the data subject and (2) more documents in the database that contain user data (using locators). The Privacy Pal client will use this function to perform recursive lookups and construct the privacy report, which follows the structure specified in your function. Here's an example of what a `HandleAccess` function might look like for a simple application with data types `user` and `post`:

```typescript
function handleAccess(dataSubjectId: string, locator: Locator, dbObj: any): map<string, any> {
  switch (locator.dataType) {
     case 'user':
        return {
           name: dbObj.name,
           email: dbObj.email,
           // construct locator for each post
           posts: dbObj.posts.map((postId) => ({
                dataType: 'post',
                isSingleDocument: true,
                collection: 'posts',
                filter: {
                    _id: new ObjectId(postId)
                }
           } as MongoLocator))
        }
     case 'post':
        return {
           content: dbObj.content,
           …
        }
 }
}
```
The user’s name and email are directly included in the report. For each post id stored in the user document, the client looks up the referenced post documents and returns its content. Based on this `handleAccess` function, the following privacy report will be generated:

```json
{
 "name": "User1",
  "email": "user1@gmail.com",
  "posts": [
          { "content": "Hey!"   }, 
          { "content": "Hello!" }
         ]
}
```

Below is the same example in Go:

```go
func HandleAccess(dataSubjectID string, locator Locator, dbObj DatabaseObject) (data map[string]interface{}, err error) {
    switch currentDbObjLocator.DataType {
    case "user":
        data["name"] = dbObj["name"]
        data["email"] = dbObj["email"]
        data["posts"] = make([]pal.Locator, 0)
        // construct locator for each post
        for _, id := range dbObj["posts"].([]interface{}) {
            data["posts"] = append(data["posts"], pal.Locator{
                LocatorType: pal.Document,
                DataType:    "post",
                FirestoreLocator: pal.FirestoreLocator{
                    CollectionPath: []string{"posts"},
                    DocIDs:         []string{id.(string)},
                },
		    })
        }
        return data, nil
    case "post":
        data["content"] = dbObj["content"]
        return data, nil
    default:
        return nil, fmt.Errorf("unknown data type %s", currentDbObjLocator.DataType)
    }
}
```

Note that in Go, `dbObj` is of type `map[string]interface{}`. In order to properly operate on its fields, you need to cast the values into an appropriate type. In the above example, in order to operate on the list of postIds, you need to first cast the `posts` field to a list of interfaces and then cast each id to a string.

Alternatively, you can (and are encouraged to) marshall the database document into your own struct to access the fields more easily (this will require you to put struct tags on your structs):

```go
jsonStr, err := json.Marshal(dbObj)
// User struct must have struct tags for json
user := &User{}
err := json.Unmarshal(jsonStr, &user);

// Can access the fields directly
for _, id := range user.Posts {
   ...
}
```

### Implementing HandleDeletion

`HandleDeletion` should specify (for each object/document type): (a) whether the current document should be deleted, (b) if not, any fields that should be updated in the current document, and (c) more documents containing personal data. 

Privacy Pal client uses the function to read from the database and keep track of documents to delete and update. When there is no more document to look up, it sends a transaction request to the database to delete and update all documents together. 

The deletion logic can also be customized for each application depending on user's privacy preferences. In the below example, a user can choose to delete or anonymize their posts. When deleting a user, the developer specifies in the data subject locator’s context field whether to anonymize their posts and passes this information to the post locators. When `HandleDeletion` is later invoked on a post, it has access to the post locator and can perform the corresponding operation. 

```TypeScript
function handleDeletion(dataSubjectId: string, locator: Locator, dbObj: any): map<string, any> {
  switch (locator.dataType) {
     case 'user':
        return {
           nodesToTraverse: [
              {
                  dataType: "post",
                  singleDocument: false,
                  collection: "posts",
                  filter: { author: dbObj._id },
                  context: locator.context
              }
            ],
           deleteNode: true
        }
     case 'post':
        if (locator.context.anonymize) {
           return {
              nodesToTraverse: [],
              deleteNode: false,
              fieldsToUpdate: {
                $set: {
                    postedBy: 'anonymous'
                }
              }
           }
        } else {
           return {
              nodesToTraverse: [],
              deleteNode: true
           }
        }
 }
}
```

For easier debugging, Privacy Pal offers a trial run mode. When invoking `processDeletionRequest`, if `writeToDatabase` is set to false, the function performs a dry run traversal and returns all documents that would be deleted or updated without committing to the database. You can inspect this report to validate your `HandleDeletion` logic before enabling database changes.

## Genpal - Stub Generator for Handle functions

`HandleAccess` and `HandleDeletion` involve a lot of boilterplate code due to the modularized nature of the functions. To get started, we provide the GenPal code generation tool for automatically producing function stubs.

To use genpal,
```bash
# TypeScript
genpal -t=user,post -output=./privacy.ts -mode=[access/deletion]

# Go
# The generated code will live in the same package as the folder that contains the output file.
genpal -input=user,post -output=./privacy.go -mode=[access/deletion]
```

Here's an example of the generated code for `HandleAccess`:

```TypeScript
function handleAccess(dataSubjectId: string, locator: Loctaor, dbObj: any): map<string, any> {
  switch (locator.dataType) {
     case 'user':
        return handleAccessUser(dataSubjectId, locator, dbObj)
     case 'post':
        return handleAccessPost(dataSubjectId, locator, dbObj)
   }
 }

function handleAccessUser(dataSubjectId: string, locator: Loctaor, dbObj: any): {
   // personal data stored in user document
   return {}
}

function handleAccessPost(dataSubjectId: string, locator: Loctaor, dbObj: any): {
   // personal data stored in post document
   return {}
}

```
