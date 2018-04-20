package backend

// ModuleInfo describes information about a module
// At the moment this is an exact copy of v1.ModuleInfo
// but kept separate to allow the backend data model to
// change independently of the API
type ModuleData struct {
	ID          string
	Owner       string
	Namespace   string
	Name        string
	Version     string
	Provider    string
	Description string
	Source      string
	PublishedAt string
	Downloads   int
	Verified    bool
}
