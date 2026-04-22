package entity

// Primitive represents a single versioned entry from a repository's index.json.
type Primitive struct {
	Type         string `json:"type"`
	Vendor       string `json:"vendor"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	SourceRepoID string `json:"sourceRepoId,omitempty"`
}
