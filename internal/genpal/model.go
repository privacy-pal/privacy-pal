package genpal

type Mode string

const (
	ModeTypenames Mode = "typenames"
	ModeYamlspec  Mode = "yamlspec"
)

type DataNodeFields map[string]interface{}
