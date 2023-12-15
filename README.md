# Privacy Pal

Privacy Pal is a framework that helps developers handle data subject access and deletion requests for applications using NoSQL databases. This repo contains the implementation of a Privacy Pal client that supports applications using MongoDB or Firestore as data store. The Privacy Pal client can be used in both Golang and TypeScript applications.

## Setup

### Installation

To use privacy-pal in your Golang application:

```bash
go get github.com/privacy-pal/privacy-pal
```

To use privacy-pal in your TypeScript application:

```bash
npm i privacy-pal
```

### Importing

To import the Privacy Pal client in your Golang application:

```golang
import pal "github.com/privacy-pal/privacy-pal/go/pkg"
```

To import the Privacy Pal client in your TypeScript application:

```typescript
import * from 'privacy-pal';
```

## Overview

Privacy Pal helps developers focus on the core data ownership logic for fulfilling privacy requests since it has all the boilerplate code for retrieving, updating, and deleting data from the database.

To use Privacy Pal, follow these steps.

1. Define `HandleAccess` and `HandleDeletion` functions, specifying what data should be returned, modified, or deleted when fulfilling access or deletion requests. See more details in [Implementing HandleAccess](#implementing-handleaccess) and [Implementing HandleDeletion](#implementing-handledeletion).

2. Initialize the Privacy Pal client with existing database client (MongoDB or Firestore).

3. Invoke `ProcessAccessRequest` and `ProcessDeletionRequest` methods of the Privacy Pal client, passing in arguments like `HandleAccess` and `HandleDeletion`, to fulfill access and deletion requests respectively.

## Core Concepts

It is important to understand the following concepts since they are core to the Privacy Pal framework.

### Data Subject

A data subject is the person whose data is being accessed or deleted. For example, if a user wants to access or delete their personal data from a chat application, the user is the data subject.

A data subject locator is used to locate the document for the data subject. This document is used as the starting point for retrieving all personal data for the data subject.

A data subject ID is passed down the data ownership logic since it is often needed to determine whether a particular piece of data belongs to the data subject. For example, data subject ID might be used to check whether a given message nested under a group chat is sent by the data subject.

### Database Object

A database object is the data stored in a database document. The database object is passed to the `HanldeAccess` and `HandleDeletion` functions you define. You can then use the data to determine what data to return, modify, or delete, and identify documents to further traverse.

### Locator

Privacy Pal uses Locators locate documents in the database. When you invoke the `ProcessXXXRequest` functions and implement your `HandleXXX` functions, you will need to use Locators.

A locator needs to contain the following basic information:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| LocatorType / IsSingleDocument   | Enum / Boolean      | Whether the locator points to a single document (`Document`) or a collection of documents (`Collection`). |
| DataType       | String     | The type of data represented by the database object retrieved using this locator. Can be used in `HandleAccess` and `HandleDeletion` to determine what type of document is being processed.  <br/><br/>  We recommend choosing a name that is consistent with the naming of your application's structs/types. E.g. if the locator points to a "user" document, the `DataType` can be set to "user". |

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

Use `ProcessAccessRequest` to retrieve all personal data for a particular user, passing in a `HandleAccess` function you defined. The `HandleAccess` should specify how to construct the access report returned to the user. To do so, it should return a map containing personal data related to the data subject, for **all possible data types**. The keys in the map will become the field names in the report. Here's an example of what a `HandleAccess` function might look like for a simple chat application:

```golang
func HandleAccess(dataSubjectID string, locator Locator, dbObj DatabaseObject) (data map[string]interface{}, err error) {
    switch currentDbObjLocator.DataType {
    case "user":
        // handle data in user document
    case "groupchat":
        // handle data in groupchat document
    case "directmessage":
        // handle data in directmessage document
    default:
        return nil, fmt.Errorf("unknown data type %s", currentDbObjLocator.DataType)
    }
}
```

HandleAccess has access to 3 values through its parameters:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| dataSubjectID    | String       | The ID of the data subject we are building the access report for. You will pass this value in when `ProcessAccessRequest` is invoked. <br/> This ID can be used to determine whether a piece of data belongs to the data subject. |
| locator    | Locator       | The current locator used to retreive document(s) from the database, containing information such as the type of data retrieved. See [Locator section](#locator) for information contained in a locator. |
| dbObj    | DatabaseObject, <br/> an alias for `map[string]interface{}`       | The document retrieved from the database using the locator, augmented with an `_id` field that contains a string |

Note that `dbObj` is of type `map[string]interface{}`. In order to properly operate on its fields, you need to cast the values into an appropriate type. For example, if you want to compare whether the "owner" field of a particular document matches with the dataSubjectId, you need to cast the owner field to a string in order to perform the comparison:

```golang
// Ensure that the "owner" field is a string
owner, ok := dbObj["owner"].(string)
if !ok {
    return nil, fmt.Errorf("owner field is not a string")
}

// Compare the owner field with the dataSubjectId
if dbObj["owner"].(string) != dataSubjectId {
    ...
}
```

If you are using a typed language (E.g. Golang), you are encouraged to marshall the database document into your own struct to access the fields more easily (this will require you to put struct tags on your structs):

```golang
jsonStr, err := json.Marshal(dbObj)
// Message struct must have struct tags for json
message := &Message{}
err := json.Unmarshal(jsonStr, &message);

// Can access the "owner" via message.Owner
if message.Owner != dataSubjectId {
   ...
}
```

When populating the map to return, there are two types of fields to consider:

- For fields that should be directly added to the report returned to the user, add the field directly as the value into the map:

    ```golang
    data["Name"] = dbObj["name"]
    ```

- For fields that require reading other documents or collections from database, you can put a locator, a list of locators, or a map from string to locators as the value:

    ```golang
    // locator
    data["Messages"] = Locator{}

    // list of locators
    data["Groupchats"] = []Locator{...}
    ```

    Privacy Pal will use these locators to retrieve the document or set of documents. Then it will invoke `HandleAccess` on the retrieved documents to recursively build the report.

You can see a full example here
TODO: add link to example

### Implementing HandleDeletion

TODO: update HandleDeletion documentation

## Generate Handle Methods with Genpal

TODO: update genpal documentation

Implementing the `HandleAccess` and `HandleDeletion` methods manually can be tedious. The genpal tool can generate stub implementations to get started.

`genpal` accepts 3 parameters:
- **mode (required):** one of the following 2 generation modes
    - `typenames`: Only generates method headers based on a list of type names passed in.
    - `yamlspec`: Generates more complete method stubs based on the yaml specification file. See the [genpal documentation](./genpal.md) for more details on the YAML format and options.
- **input (required):** The input for the given mode. For `typenames` this is a comma-separated list of type names. For `yamlspec` it is the path to a YAML schema file.
- **output (optional):** The path of the generated Go file. Defaults to `./privacy.go`

For example,
```
# Typenames mode
genpal -mode=typenames -input=User,GroupChat,Message -output=internal/chat/privacy.go

# YAML mode  
genpal -mode=yamlspec -input=internal/chat/privacypal.yaml -output=internal/chat/privacy.go
```

The generated code will live in the same package as the folder that contains the output file.
