# Privacy Pal

Privacy Pal is a tool that helps applications handle data access and deletion requests in compliance with privacy regulations like GDPR and CCPA. Currently, Privacy Pal only works with the Firestore database. 

## Overview

Privacy Pal provides two main functions for processing privacy requests:

- `ProcessAccessRequest` - Retrieve personal data for a data subject
- `ProcessDeletionRequest` - Delete or anonymize personal data for a data subject

To use PAL, applications need to:
1. Define structs for your application data and implement the `HandleAccess` and `HandleDeletion` methods. This allows PAL to interface with your application's data models. 
    ```
    type UserData struct {
        Name string
        Email string
        Posts []Post
    }

    func (u *UserData) HandleAccess(dataSubjectId string, currentDataNodeLocator Locator) map[string]interface{} {
        // return data to fulfill access request 
    }

    func (u *UserData) HandleDeletion(dataSubjectId string) (nodesToTraverse []Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {
        // return data to fulfill deletion request
    }
    ```

2. Initialize the pal client with your existing firestore client 

    ```
    client := pal.NewClient(firestoreClient)
    ```

3. Invoke `ProcessAccessRequest` and `ProcessDeletionRequest` functions. The pal client will handle recursively retrieving and redacting nested data based on the specifications outlined in `HandleAccess` and `HandleDeletion`.
    ```
    // Access request
    data, err := client.ProcessAccessRequest(
        dataSubjectLocator, 
        "user123" // data subject ID
    )

    // Deletion request 
    result, err := palClient.ProcessDeletionRequest(
        dataSubjectLocator, 
        "user123" // data subject ID  
    )
    ```

### Locator 
A Locator is a struct that contains information to locate a specific data node or collection in the database. It allows Privacy Pal to recursively retrieve data for access requests.

The fields in a Locator are:
- `Type`: Either `Document` or `Collection`
- `CollectionPath`: Collection path leading up to a particular Firestore document or collection
- `DocIDs`: List of document IDs in the order of collections
- `NewDataNode`: A function that returns a new instance of the data node struct (used to unmarshall the data fetched from Firestore)
- `Queries`: database queries (only applicable to collections)

Example:
```
pal.Locator{
		Type:           pal.Document,
		CollectionPath: []string{"users"},
		DocIDs:         []string{"123"},
		NewDataNode:    func() pal.DataNode { return &User{} },
	}
```

// TODO: more detailed explanation for collectionpath and DocIDs

### Implementing `HandleAccess`

Privacy Pal uses the `HandleAccess` method to retrieve the personal data for a data subject during an access request.

You should implement the method to return a map of the data to beinclude in the access report. The keys in the map will become the field names in the report.

There are a few ways to populate the map:

For fields that should be directly returned to the user, just use the field directly as the value to the map:

```
data["Name"] = u.Name
```

For fields that require reading other documents or collections from Firestore, you can put a locator or a list of locators as the value:
```
// list of locators
data["Groupchats"] = []Locator{...}

// locator

data["Messages"] = Locator{}
```

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
genpal -mode=typenames -input=User,Post,Comment -output=internal/chat/privacy.go

# YAML mode  
genpal -mode=yamlspec -input=internal/chat/privacypal.yaml -output=internal/chat/privacy.go
```

The generated code will live in the same package as the folder that contains the output file.