package state

// Enum represents a PostgreSQL enum type
type Enum struct {
	Name        string
	Values      []string
	TypeComment string
	CreatedIn   int
	UsedBy      []string
}

// NewEnum creates a new enum
func NewEnum(name string) *Enum {
	return &Enum{
		Name:   name,
		Values: []string{},
		UsedBy: []string{},
	}
}

// AddValue adds a value to the enum
func (e *Enum) AddValue(value string) {
	// Check if value already exists
	for _, v := range e.Values {
		if v == value {
			return
		}
	}
	e.Values = append(e.Values, value)
}

// AddUsedBy records that a table uses this enum
func (e *Enum) AddUsedBy(tableName string) {
	for _, t := range e.UsedBy {
		if t == tableName {
			return
		}
	}
	e.UsedBy = append(e.UsedBy, tableName)
}
