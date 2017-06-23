package schema

// Schema is the metadata container for a schema definition
type Schema struct {
	Tables map[string]*Table `json:"Tables"`
}

// Table is the metadata container for a SQL table definition
type Table struct {
	MultiKey bool   `json:"MultiKey"` // Use Primary or Primaries
	Primary  string `json:"Primary"`

	Primaries []string `json:"Primaries"`

	// Fields is the column definitions for the SQL table
	Fields map[string]*Field `json:"Fields"`

	EssentialFields []string `json:"EssentialFields"`

	Children map[string]*ChildTable `json:"Children"`

	// YAGNI?
	// TODO: ChildrenInsertionOrder?
	// TODO: DeletionOrder?
}

// Field represents a single column in a SQL table
type Field struct {
	Name         string `json:"Name"`
	AllowNull    bool   `json:"AllowNull"`
	DefaultValue string `json:"DefaultValue"` // Converts to integer if IsNumber is set
	IsNumber     bool   `json:"IsNumber"`
	IsIdentity   bool   `json:"IsIdentity"`
	DBType       string `json:"DBType"`
	Length       int    `json:"Length"`

	// TODO: IsUnique bool `json:"IsUnique"` // probably necessary because of say, GUID..
}

// ChildTable represents a relationship between a parent table
// and a child table
type ChildTable struct {
	ParentTable string `json:"ParentTable"`

	MultiKey     bool   `json:"MultiKey"`
	LocalField   string `json:"LocalField"`
	ForeignField string `json:"ForeignField"`

	LocalFields   []string `json:"LocalFields"`
	ForeignFields []string `json:"ForeignFields"`
}
