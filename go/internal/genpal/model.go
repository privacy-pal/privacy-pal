package genpal

import pal "github.com/privacy-pal/privacy-pal/pkg"

type Mode string

const (
	ModeTypenames Mode = "typenames"
	ModeYamlspec  Mode = "yamlspec"
)

type DataNodeProperty struct {
	CollectionPath []string        `yaml:"collection_path,omitempty"`
	DirectFields   []string        `yaml:"direct_fields,omitempty"`
	IndirectFields []IndirectField `yaml:"indirect_fields,omitempty"`
}

type IndirectField struct {
	Type         string       `yaml:"type"`
	FieldName    string       `yaml:"field_name,omitempty"`
	ExportedName string       `yaml:"exported_name"`
	Queries      *[]pal.Query `yaml:"queries,omitempty"`
}
