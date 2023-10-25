package genpal

import "strings"

func GenerateWithTypenameMode(typenames []string) (ret string) {
	for _, typename := range typenames {
		ret += generateHandleAccess(typename)
		ret += generateHandleDeletion(typename)
	}
	return
}

func GenerateWithYamlspecMode(data []map[string]interface{}) (ret string) {
	for _, obj := range data {
		// Top level obj only has one key, which is the typename (e.g. "User")
		for key := range obj {
			typename := key
			ret += generateHandleAccess(typename)
			ret += generateHandleDeletion(typename)
		}
	}
	return
}

func generateHandleAccess(typename string) (ret string) {
	ret += "func (" + strings.ToLower(typename[0:1]) + " *" + typename + ") HandleAccess(dataSubjectId string, currentDocumentID string) map[string]interface{} {\n"
	ret += "return nil\n"
	ret += "}\n\n"
	return
}

func generateHandleDeletion(typename string) (ret string) {
	ret += "func (" + strings.ToLower(typename[0:1]) + " *" + typename + ") HandleDeletion(dataSubjectId string) (nodesToTraverse []pal.Locator, deleteNode bool, fieldsToUpdate []firestore.Update) {\n"
	ret += "return nil, false, nil\n"
	ret += "}\n\n"
	return
}
