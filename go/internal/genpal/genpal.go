package genpal

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ExportedHandleAccessFuncName = "HandleAccess"
	InternalHandleAccessFuncName = "handleAccess"
)

func GenerateStubs(typenames []string) (ret string) {
	ret += generateHandleAccess(typenames)

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
	ret += "func " + ExportedHandleAccessFuncName + "(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	ret += "switch currentDbObjLocator.DataType {\n"
	for _, typename := range typenames {
		ret += fmt.Sprintf("case %sDataType:\n", toCamelCase(typename))
		ret += fmt.Sprintf("return %s%s(dataSubjectId, currentDbObjLocator, dbObj)\n", InternalHandleAccessFuncName, toCamelCase(typename))
	}
	ret += "default:\n"
	ret += "return nil\n"
	ret += "}\n"
	ret += "}\n\n"

	return
}

func generateHandleAccessForType(typename string) (ret string) {
	ret += "func " + InternalHandleAccessFuncName + toCamelCase(typename) + "(dataSubjectId string, currentDbObjLocator pal.Locator, dbObj pal.DatabaseObject) map[string]interface{} {\n"
	// only generate function headers
	ret += "return nil\n"
	ret += "}\n\n"
	return
}

// example: 'group_chat' yields 'GroupChat'
func toCamelCase(s string) string {
	lst := strings.Split(s, "_")
	for i, s := range lst {
		lst[i] = cases.Title(language.English).String(s)
	}
	return strings.Join(lst, "")
}
