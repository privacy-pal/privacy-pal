package genpal

import (
	"fmt"
	"strings"

	pal "github.com/privacy-pal/privacy-pal/pkg"
)

var typenames []string
var collectionPaths map[string][]string

func init() {
	typenames = make([]string, 0)
	collectionPaths = make(map[string][]string)
}

func GenerateWithTypenameMode(typenames []string) (ret string) {
	for _, typename := range typenames {
		ret += generateHandleAccessForType(typename, nil)
	}
	return
}

func GenerateWithYamlspecMode(data map[string]DataNodeProperty) (ret string) {
	// Store collection paths for each typename
	for typename, property := range data {
		typenames = append(typenames, typename)
		collectionPaths[typename] = property.CollectionPath
	}

	ret += generateHandleAccess(typenames)

	// Generate interface methods for each typename, in sequence of typenames in yaml spec
	for _, typename := range typenames {
		property := data[typename]
		ret += generateHandleAccessForType(typename, &property)
		// ret += generateHandleDeletion(typename)
	}
	return
}

func generateHandleAccess(typenames []string) (ret string) {
	// constant data type names
	ret += "const (\n"
	for _, typename := range typenames {
		ret += fmt.Sprintf("%sDataType = \"%s\"\n", toCamelCase(typename), typename)
	}
	ret += ")\n\n"

	// generate handle access function
	ret += "func HandleAccess(dataSubjectId string, currentDataNodeLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	ret += "switch currentDataNodeLocator.DataType {\n"
	for _, typename := range typenames {
		ret += fmt.Sprintf("case %sDataType:\n", toCamelCase(typename))
		ret += fmt.Sprintf("return HandleAccess%s(dataSubjectId, currentDataNodeLocator, dbObj)\n", toCamelCase(typename))
	}
	ret += "default:\n"
	ret += "return nil\n"
	ret += "}\n"
	ret += "}\n\n"

	return
}

func generateHandleAccessForType(typename string, dataNodeProperty *DataNodeProperty) (ret string) {
	ret += "func HandleAccess" + toCamelCase(typename) + "(dataSubjectId string, currentDataNodeLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	// only generate function headers
	if dataNodeProperty == nil {
		ret += "return nil\n"
		ret += "}\n\n"
		return
	}

	// generate function body
	ret += "data := make(map[string]interface{})\n\n"

	// generate code for data subject
	if dataNodeProperty.IsDataSubject {
		ret += fmt.Sprintf(
			`if dbObj["_id"].(string) != dataSubjectId {
				// current database object is not the data subject
				return data
			}` + "\n\n",
		)
	}

	// generate code for direct fields
	for _, field := range dataNodeProperty.DirectFields {
		ret += fmt.Sprintf(`data["%s"] = dbObj["%s"]`+"\n", field, field)
	}

	// generate code for indirect fields
	for _, field := range dataNodeProperty.IndirectFields {
		// Parse type
		locatorType, list, dataType, err := parseIndirectFieldType(field.Type)
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
					LocatorType:           %s,
					DataType:       "%s",
					CollectionPath: %s,
					DocIDs:         %s,
				}`,
				locatorTypeStr,
				dataType,
				fmt.Sprintf("%#v", collectionPaths[dataType]),
				strings.ReplaceAll(fmt.Sprintf("%#v", []string{"id"}), `"`, ""),
			)

			if list {
				ret += fmt.Sprintf(`data["%s"] = make([]pal.Locator, 0)`+"\n", field.ExportedName)
				ret += fmt.Sprintf(
					`for _, id := range dbObj["%s"].([]interface{}) {
					id := id.(string)
					data["%s"] = append(data["%s"].([]pal.Locator), %s)
				}`+"\n",
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
					LocatorType:           %s,
					DataType:       "%s",
					CollectionPath: append(currentDataNodeLocator.CollectionPath, "%s"),
					DocIDs:         currentDataNodeLocator.DocIDs,`+"\n",
				field.ExportedName,
				locatorTypeStr,
				dataType,
				collectionPaths[dataType][len(collectionPaths[dataType])-1],
			)

			if field.Queries != nil && len(*field.Queries) > 0 {
				queryStr := ""
				for _, query := range *field.Queries {
					queryVal := fmt.Sprintf("\"%s\"", query.Value)
					if strings.HasPrefix(query.Value.(string), "${") && strings.HasSuffix(query.Value.(string), "}") {
						queryVal = queryVal[3 : len(queryVal)-2]
					}

					queryStr += fmt.Sprintf(
						`{
							Path:  "%s",
							Op:    "%s",
							Value: %s,
						},`+"\n",
						query.Path,
						query.Op,
						queryVal,
					)
				}
				ret += fmt.Sprintf(
					`Queries: []pal.Query{
						%s
					},`+"\n",
					queryStr,
				)
			}

			ret += "}\n"
		}
	}

	ret += "\nreturn data\n"
	ret += "}\n\n"
	return
}

func parseIndirectFieldType(fieldtype string) (locatorType pal.LocatorType, list bool, dataType string, err error) {
	// options include list<ID<xxx>>, ID<xxx>, subcollection<xxx>
	// if ID<xxx>, list is false, dataType is xxx, locatorType is Document
	// if list<ID<xxx>>, list is true, dataType is xxx, locatorType is Document
	// if subcollection<xxx>, list is false, dataType is xxx, locatorType is Collection

	if strings.HasPrefix(fieldtype, "list<") {
		list = true
		fieldtype = fieldtype[5 : len(fieldtype)-1]
	} else {
		list = false
	}

	if strings.HasPrefix(fieldtype, "ID<") {
		dataType = fieldtype[3 : len(fieldtype)-1]
		locatorType = pal.Document
	} else if strings.HasPrefix(fieldtype, "subcollection<") {
		dataType = fieldtype[14 : len(fieldtype)-1]
		locatorType = pal.Collection
	} else {
		err = fmt.Errorf("invalid indirect field type: %s", fieldtype)
	}

	return
}

// example input: group_chat; output: GroupChat
func toCamelCase(s string) string {
	lst := strings.Split(s, "_")
	for i, s := range lst {
		lst[i] = strings.Title(s)
	}
	return strings.Join(lst, "")
}
