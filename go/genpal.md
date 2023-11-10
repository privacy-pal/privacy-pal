# Genpal

## Testing
1. Under root, run the following to generate the code stub:
```
go run cmd/genpal.go -mode=yamlspec -input=internal/test/chat/privacypal.yaml -output=internal/test/chat/privacy.go
```
2. Run tests on the chat application
```
cd internal/test/chat
go test -v ./...
```


## YAML Schema Specification

The YAML schema file defines the data model so genpal can generate more complete method stubs.

Here is an example [specification](./internal/chat/privacypal.yaml) and the corresponding [generated code](./internal/chat/privacy.go).

### Format

The file should define a YAML mapping (dictionary) where the keys are the type names. The types must be present in the same package as the output path and must be spelled identically.

Under each type, specify:
- `collection_path` (required): collection path leading up to a particular document of the specified type
- `direct_fields` - list of fields to be directly returned to the data subject. Each field must exist in the type and be spelled identically.
- `indirect_fields` - List of fields that require reading additional documents or collections from database.

For example:
```
User:
  collection_path: 
    - users
  
  direct_fields:
    - Name
    - Email

  indirect_fields:
    - type: ID<Post>
      #...
```

### indirect_fields
Each indirect field definition contains:

- `type` (required) Options:
    -  `ID<TypeName>` - Single document reference
    - `list<ID<TypeName>>` - List of document references
    - `subcollection<TypeName>` - A list of all or a subset of documents in a subcollection
- `field_name`: Field name on the type/struct. Must be spelled identically. Not applicable for subcollections.
- `exported_name` (required): The name to use in the exported map data.
- `queries`: List of queries to filter a subcollection. Only applicable for subcollections.

### queries
You can additionally specify a list of queries to filter a subcollection. Every query requires 3 fields:
- `path`: Field path to query. Path must be identical to the path in database.
- `op`: Comparison operator in string (==, !=, < etc.)
- `value` - Value to compare against. Can use `${dataSubjectId}` to reference the id of the data subject we are performing the access request for

For example, the below specification will generate a query that returns the list of messages where the field `userId` is equal to the data subject ID.
```
GroupChat:
    collection_path:
        - gcs
    indirect_fields:
        - type: subcollection<Message>
            exported_name: Messages 
            queries:
                - path: userId
                op: ==
                value: ${dataSubjectId}
```