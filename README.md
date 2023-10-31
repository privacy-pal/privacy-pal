# Privacy Pal

Privacy Pal is a tool that helps applications handle data access and deletion requests in compliance with privacy regulations like GDPR and CCPA. Currently, Privacy Pal only works with the Firestore database. It is available as a Go module and an npm package.

## Installation

```
// Golang
go get github.com/privacy-pal/privacy-pal

// TypeScript
npm i privacy-pal
```

## Overview

Privacy Pal provides two main functions for processing privacy requests:

- `ProcessAccessRequest` - Retrieve personal data for a data subject
- `ProcessDeletionRequest` - Delete or anonymize personal data for a data subject

To use Privacy Pal, applications need to:
1. Define `HandleAccess` and `HandleDeletion` functions to specify what data should be returned, modified, or deleted when fulfilling access or deletion requests (Explained in detail below).

```
func HandleAccess(dataSubjectID string, locator Locator, dbObj DatabaseObject) 
```

2. Initialize the pal client with your existing firestore client 

    ```
    client := pal.NewClient(firestoreClient)
    ```

3. Invoke `ProcessAccessRequest` and `ProcessDeletionRequest` functions. The pal client will handle recursively retrieving and redacting nested data based on the specifications outlined in `HandleAccess` and `HandleDeletion`.
    ```
    // access request
    data, err := client.ProcessAccessRequest(
        HandleAccess, 
        dataSubjectLocator, 
        "user123" // userID
    )
    ```

### Locator 
Privacy Pal uses Locators to retrieve, update, or delete data in the database. When you invoke the `ProcessXXXRequest` functions and implement your `HandleXXX` functions, you will need to use Locators, which are objects that contain information to locate a document or a set of documents in the database. 

The fields in a Locator are:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| LocatorType    | Enum       | Specifies whether the locator points to a single document (`Document`) or a collection of documents (`Collection`). |
| DataType       | String     | Indicates the type of data represented by the database object retrieved using this locator. Can be used in `HandleAccess` to decide what data needs to be returned in the object.  <br/><br/>  We recommend choosing a name that is consistent with the naming of your application's structs/types. E.g. if the locator points to a "user" document, the DataType can be set to "user". |
| CollectionPath | List of Strings | Represents the collection path leading up to a particular Firestore document or collection. E.g. `["course]`, `["course", "sections"]` |
| DocIDs         | List of Strings | Document IDs in the order they appear within collections. <br/><br/> If locating a document, should have the same length as CollectionPath. If locating a collection, should be one item shorter than the length of CollectionPath. |
| Queries        | Query Object | Represents the database queries that can be applied to filter and retrieve specific data within a collection. Only applicable to `Collection` locators. |

Example:
```
pal.Locator{
    LocatorType:         pal.Document,
    DataType:            "user",
    CollectionPath:      []string{"users"},
    DocIDs:              []string{"123"},
    NewDataNode:         func() pal.DataNode { return &User{} },
}
```

### Implementing HandleAccess

When you invoke `ProcessAccessRequest` to retrieve all personal data for a particular user, you need to pass in a `HandleAccess` function, which Privacy Pal will use to build the access report. The function should return a map containing personal data related to the data subject, for **all possible data types**. The keys in the map will become the field names in the report. For example,

```
func HandleAccess(dataSubjectID string, locator Locator, dbObj DatabaseObject) {
    switch currentDbObjLocator.DataType {
	case "user":
		// handle data in user document
	case "groupchat":
		// handle data in groupchat document
    ...
	default:
		return nil
	}
}
```

HandleAccess exposes you to 3 values:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| dataSubjectID    | String       | The ID of the data subject we are building the access report for. You will pass this value in when `ProcessAccessRequest` is invoked. |
| locator    | Locator       | The current locator used to retreive document(s) from the database, containing information about the type of data retrieved, the collection path, and docIDs leading up to the document. |
| dbObj    | DatabaseObject, <br/> an alias for `map[string]interface{}`       | The document retrieved from the database using the locator, augmented with an `_id` field that contains a string |

Note that `dbObj` is of type `map[string]interface{}`. In order to properly operate on its fields, you need to cast the values into an appropriate type. For example, if you want to compare whether the "owner" field of a particular document matches with the dataSubjectId, you need to cast the owner field to a string in order to perform the comparison:

```
if dbObj["owner"].(string) != dataSubjectId {
    ...
}
```
If you are using a typed language (E.g. Go), you are encouraged to marshall the database document into your own struct to access the fields more easily (this will require you to put struct tags on your structs). For example, 
```
jsonStr, err := json.Marshal(dbObj)
message := &Message{}
err := json.Unmarshal(jsonStr, &message);

// Then you can access the "owner" field via user.Owner
if message.Owner != dataSubjectId {
   ...
}
```

When you populate the map,
- For fields that should be directly returned to the user, just use the field directly as the value to the map:

        data["Name"] = dbObj["name"]

- For fields that require reading other documents or collections from Firestore, you can put a locator, a list of locators, or a map from string to locators as the value:

        // locator
        data["Messages"] = Locator{}

        // list of locators
        data["Groupchats"] = []Locator{...}

    Privacy Pal will use these locators to look up the document or set of documents and recursively return the data relevant to the data subject based on the `HandleAccess` function.

You can see a full example here
// TODO: add link to example

### Implementing HandleDeletion



## Generate Handle Methods with Genpal

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
