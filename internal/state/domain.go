package state

// Domain represents a PostgreSQL domain type
type Domain struct {
	Name       string
	BaseType   string
	Default    string
	Constraint string
	Comment    string
	CreatedIn  int
}

// NewDomain creates a new domain
func NewDomain(name string) *Domain {
	return &Domain{
		Name: name,
	}
}
