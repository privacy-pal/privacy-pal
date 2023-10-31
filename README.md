# Privacy Pal

Privacy Pal is a tool that helps applications handle data access and deletion requests in compliance with privacy regulations like GDPR and CCPA. Currently, Privacy Pal only works with the Firestore database. 

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
A Locator is an object that contains information to locate a document or a set of documents in the database. It allows Privacy Pal to retrieve, update, or delete data in the database.

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

You need to implement a `HandleAccess`` function, which Privacy Pal will call to build the access report. The function should return a map containing personal data related to the data subject, for **all possible data types**. The keys in the map will become the field names in the report. For example,

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

You will be passing this function to the `ProcessAccessRequest` to generate the report.

HandleAccess exposes you to 3 values:
| Field          | Type       | Description                                           |
| -------------- | ---------- | ----------------------------------------------------- |
| dataSubjectID    | String       | The ID of the data subject we are building the access report for. Passed in when `ProcessAccessRequest` is invoked. |
| locator    | Locator       | The current locator used to retreive document(s) from the database, <br/> containing information about the type of data retrieved, the collection path, and docIDs leading up to the document. |
| dbObj    | DatabaseObject, <br/> an alias for `map[string]interface{}`       | The document retrieved from the database using the locator, augmented with an `_id` field that contains a string |

Note that `dbObj` is of type `map[string]interface{}`. In order to properly operate on its fields, you need to cast the values into an appropriate type. For example

// TODO: change the example
```
if dbObj["_id"].(string) != dataSubjectId {
		data["Name"] = dbObj["name"]
		return data
}
```
Note: if you are using a typed language (E.g. Go), you are encouraged to marshall the database document into your own struct for type safety. For example,
```
```

--------------------
NOT FINISHED BELOW

There are a few ways to populate the map:

- For fields that should be directly returned to the user, just use the field directly as the value to the map:

        data["Name"] = dbObj["name"]

- For fields that require reading other documents or collections from Firestore, you can put a locator or a list of locators as the value:

        // list of locators
        data["Groupchats"] = []Locator{...}

        // locator
        data["Messages"] = Locator{}

    Privacy Pal will use these locators to traverse down the database and return all data relevant to the data subject.

You can see a full example here
// TODO: add link to example

### Implementing `HandleDeletion`



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
