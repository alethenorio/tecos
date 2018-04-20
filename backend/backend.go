package backend

// Storage is an API for retrieving modules from
// some backend
type Storage interface {
	ListModules() []ModuleData
}
