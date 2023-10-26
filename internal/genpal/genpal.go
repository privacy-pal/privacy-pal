package genpal

import (
	"fmt"
	"strings"

	pal "github.com/privacy-pal/privacy-pal/pkg"
)

var collectionPaths map[string][]string

func init() {
	collectionPaths = make(map[string][]string)
}

func GenerateWithTypenameMode(typenames []string) (ret string) {
	for _, typename := range typenames {
		ret += generateHandleAccess(typename, nil)
		ret += generateHandleDeletion(typename)
	}
	return
}

func GenerateWithYamlspecMode(data map[string]DataNodeProperty) (ret string) {
	// Store collection paths for each typename
	for typename, property := range data {
		collectionPaths[typename] = property.CollectionPath
	}

	// Generate interface methods for each typename
	for typename, property := range data {
		ret += generateHandleAccess(typename, &property)
		ret += generateHandleDeletion(typename)
	}
	return
}

func generateHandleAccess(typename string, dataNodeProperty *DataNodeProperty) (ret string) {
	obj := strings.ToLower(typename[0:1])
	ret += "func (" + obj + " *" + typename + ") HandleAccess(dataSubjectId string, currentDataNodeLocator pal.Locator) map[string]interface{} {\n"
	// only generate function headers
	if dataNodeProperty == nil {
		ret += "return nil\n"
		ret += "}\n\n"
		return
	}

	// generate function body
	ret += "data := make(map[string]interface{})\n\n"

	// generate code for direct fields
	for _, field := range dataNodeProperty.DirectFields {
		ret += fmt.Sprintf(`data["%s"] = %s.%s`+"\n", field, obj, field)
	}

	// generate code for indirect fields
	for _, field := range dataNodeProperty.IndirectFields {
		// Parse type
		locatorType, list, dataNode, err := parseIndirectFieldType(field.Type)
		if err != nil {
			// TODO: handle error
			panic(err)
		}

		var locatorTypeStr string
		if locatorType == pal.Document {
			locatorTypeStr = "pal.Document"
		} else if locatorType == pal.Collection {
			locatorTypeStr = "pal.Collection"
		}

		var locatorStr string
		if locatorType == pal.Document {
			locatorStr = fmt.Sprintf(
				`pal.Locator{
					Type:           %s,
					CollectionPath: %s,
					DocIDs:         %s,
					NewDataNode:    func() pal.DataNode { return &%s{} },
				}`,
				locatorTypeStr,
				fmt.Sprintf("%#v", collectionPaths[dataNode]),
				strings.ReplaceAll(fmt.Sprintf("%#v", []string{"id"}), `"`, ""),
				dataNode,
			)

			if list {
				ret += fmt.Sprintf(`data["%s"] = make([]pal.Locator, 0)`+"\n", field.ExportedName)
				ret += fmt.Sprintf(
					`for _, id := range %s.%s {
					data["%s"] = append(data["%s"].([]pal.Locator), %s)
				}`+"\n",
					obj,
					field.FieldName,
					field.ExportedName,
					field.ExportedName,
					locatorStr,
				)
			} else {
				ret += fmt.Sprintf(`data["%s"] = %s`+"\n", field.ExportedName, locatorStr)
			}
		}

		if locatorType == pal.Collection {
			ret += fmt.Sprintf(
				`data["%s"] = pal.Locator{
					Type:           %s,
					CollectionPath: append(currentDataNodeLocator.CollectionPath, "%s"),
					DocIDs:         currentDataNodeLocator.DocIDs,
					NewDataNode:    func() pal.DataNode { return &%s{} },
				}`+"\n",
				field.ExportedName,
				locatorTypeStr,
				collectionPaths[dataNode][len(collectionPaths[dataNode])-1],
				dataNode,
			)
		}
	}

	ret += "\nreturn data\n"
	ret += "}\n\n"
	return
}

func generateHandleDeletion(typename string) (ret string) {
	ret += "func (" + strings.ToLower(typename[0:1]) + " *" + typename + ") HandleDeletion(dataSubjectId string) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {\n"
	ret += "return nil, false, nil\n"
	ret += "}\n\n"
	return
}

func parseIndirectFieldType(fieldtype string) (locatorType pal.LocatorType, list bool, dataNode string, err error) {
	// options include list<ID<xxx>>, ID<xxx>, subcollection<xxx>
	// if ID<xxx>, list is false, dataNode is xxx, locatorType is Document
	// if list<ID<xxx>>, list is true, dataNode is xxx, locatorType is Document
	// if subcollection<xxx>, list is false, dataNode is xxx, locatorType is Collection

	if strings.HasPrefix(fieldtype, "list<") {
		list = true
		fieldtype = fieldtype[5 : len(fieldtype)-1]
	} else {
		list = false
	}

	if strings.HasPrefix(fieldtype, "ID<") {
		dataNode = fieldtype[3 : len(fieldtype)-1]
		locatorType = pal.Document
	} else if strings.HasPrefix(fieldtype, "subcollection<") {
		dataNode = fieldtype[14 : len(fieldtype)-1]
		locatorType = pal.Collection
	} else {
		err = fmt.Errorf("invalid indirect field type: %s", fieldtype)
	}

	return

}
