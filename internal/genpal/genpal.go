package genpal

import (
	"fmt"
	"strings"
)

func GenerateWithTypenameMode(typenames []string) (ret string) {
	for _, typename := range typenames {
		ret += generateHandleAccess(typename, nil)
		ret += generateHandleDeletion(typename)
	}
	return
}

func GenerateWithYamlspecMode(data []map[string]DataNodeFields) (ret string) {
	for _, obj := range data {
		// only has one key, which is the typename (e.g. "User")
		for typename, fields := range obj {
			ret += generateHandleAccess(typename, fields)
			ret += generateHandleDeletion(typename)
		}
	}
	return
}

func generateHandleAccess(typename string, fields DataNodeFields) (ret string) {
	ret += "func (" + strings.ToLower(typename[0:1]) + " *" + typename + ") HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{} {\n"
	// only generate function headers
	if fields == nil {
		ret += "return nil\n"
		ret += "}\n\n"
		return
	}

	// generate function body
	ret += "data := make(map[string]interface{})\n\n"
	for fieldName, val := range fields {
		switch val.(type) {
		case string:
			if val == "direct_field" {
				ret += fmt.Sprintf(`data["%s"] = c.%s`+"\n", fieldName, fieldName)
			}
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
