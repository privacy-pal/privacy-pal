package genpal

import (
	"fmt"
	"strings"
)

var dataNodeProperties map[string]DataNodeProperty

func init() {
	dataNodeProperties = make(map[string]DataNodeProperty)
}

func GenerateWithTypenameMode(typenames []string) (ret string) {
	for _, typename := range typenames {
		ret += generateHandleAccessForType(typename)
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
	ret += "func HandleAccess(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	ret += "switch currentDbObjLocator.DataType {\n"
	for _, typename := range typenames {
		ret += fmt.Sprintf("case %sDataType:\n", toCamelCase(typename))
		ret += fmt.Sprintf("return HandleAccess%s(dataSubjectId, currentDbObjLocator, dbObj)\n", toCamelCase(typename))
	}
	ret += "default:\n"
	ret += "return nil\n"
	ret += "}\n"
	ret += "}\n\n"

	return
}

func generateHandleAccessForType(typename string) (ret string) {
	ret += "func HandleAccess" + toCamelCase(typename) + "(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	// only generate function headers
	ret += "return nil\n"
	ret += "}\n\n"
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
