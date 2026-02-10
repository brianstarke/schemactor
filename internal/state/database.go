package state

// DatabaseState represents the cumulative database state after all migrations
type DatabaseState struct {
	Domains map[string]*Domain
	Enums   map[string]*Enum
	Tables  map[string]*Table
	Views   map[string]*View

	// Track dropped objects to avoid recreating them
	DroppedTables  map[string]bool
	DroppedDomains map[string]bool
	DroppedEnums   map[string]bool
	DroppedViews   map[string]bool

	// Track indexes separately for later removal if table is modified
	Indexes map[string]*Index
}

// NewDatabaseState creates a new empty database state
func NewDatabaseState() *DatabaseState {
	return &DatabaseState{
		Domains:        make(map[string]*Domain),
		Enums:          make(map[string]*Enum),
		Tables:         make(map[string]*Table),
		Views:          make(map[string]*View),
		DroppedTables:  make(map[string]bool),
		DroppedDomains: make(map[string]bool),
		DroppedEnums:   make(map[string]bool),
		DroppedViews:   make(map[string]bool),
		Indexes:        make(map[string]*Index),
	}
}

// AddOrUpdateTable adds or updates a table
func (ds *DatabaseState) AddOrUpdateTable(table *Table) {
	ds.Tables[table.Name] = table
	delete(ds.DroppedTables, table.Name)
}

// GetTable returns a table by name
func (ds *DatabaseState) GetTable(name string) (*Table, bool) {
	table, ok := ds.Tables[name]
	return table, ok
}

// DropTable marks a table as dropped
func (ds *DatabaseState) DropTable(name string) {
	delete(ds.Tables, name)
	ds.DroppedTables[name] = true
}

// AddOrUpdateDomain adds or updates a domain
func (ds *DatabaseState) AddOrUpdateDomain(domain *Domain) {
	ds.Domains[domain.Name] = domain
	delete(ds.DroppedDomains, domain.Name)
}

// GetDomain returns a domain by name
func (ds *DatabaseState) GetDomain(name string) (*Domain, bool) {
	domain, ok := ds.Domains[name]
	return domain, ok
}

// DropDomain marks a domain as dropped
func (ds *DatabaseState) DropDomain(name string) {
	delete(ds.Domains, name)
	ds.DroppedDomains[name] = true
}

// AddOrUpdateEnum adds or updates an enum
func (ds *DatabaseState) AddOrUpdateEnum(enum *Enum) {
	ds.Enums[enum.Name] = enum
	delete(ds.DroppedEnums, enum.Name)
}

// GetEnum returns an enum by name
func (ds *DatabaseState) GetEnum(name string) (*Enum, bool) {
	enum, ok := ds.Enums[name]
	return enum, ok
}

// DropEnum marks an enum as dropped
func (ds *DatabaseState) DropEnum(name string) {
	delete(ds.Enums, name)
	ds.DroppedEnums[name] = true
}

// AddOrUpdateView adds or updates a view
func (ds *DatabaseState) AddOrUpdateView(view *View) {
	// Check if view already exists (for versioning)
	if existingView, ok := ds.Views[view.Name]; ok {
		view.Version = existingView.Version + 1
	}

	ds.Views[view.Name] = view
	delete(ds.DroppedViews, view.Name)
}

// GetView returns a view by name
func (ds *DatabaseState) GetView(name string) (*View, bool) {
	view, ok := ds.Views[name]
	return view, ok
}

// DropView marks a view as dropped
func (ds *DatabaseState) DropView(name string) {
	delete(ds.Views, name)
	ds.DroppedViews[name] = true
}

// AddIndex adds an index to the state
func (ds *DatabaseState) AddIndex(idx *Index) {
	ds.Indexes[idx.Name] = idx
}

// GetIndex returns an index by name
func (ds *DatabaseState) GetIndex(name string) (*Index, bool) {
	idx, ok := ds.Indexes[name]
	return idx, ok
}

// DropIndex removes an index from the state
func (ds *DatabaseState) DropIndex(name string) {
	delete(ds.Indexes, name)
}
