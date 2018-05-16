package v1

import "github.com/ByteFlinger/tecos/backend"

type ModuleList struct {
	Meta    PaginationMeta `json:"meta,omitempty"`
	Modules []ModuleInfo   `json:"modules,omitempty"`
}

// PaginationMeta represents pagination in API responses
type PaginationMeta struct {
	Limit          int `json:"limit,omitempty"`
	CurrentOffset  int `json:"current_offset,omitempty"`
	NextOffset     int `json:"next_offset,omitempty"`
	PreviousOffset int `json:"previous_offset,omitempty"`
}

type ModuleInfo struct {
	ID          string `json:"id,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	Downloads   int    `json:"downloads,omitempty"`
	Verified    bool   `json:"verified,omitempty"`
}

// FromData returns a ModuleInfo based on the given backend.ModuleData
func FromData(m backend.ModuleData) ModuleInfo {
	return ModuleInfo{
		Name:        m.Name,
		Namespace:   m.Namespace,
		Description: m.Description,
		PublishedAt: m.PublishedAt,
		Version:     m.Version,
	}
}
