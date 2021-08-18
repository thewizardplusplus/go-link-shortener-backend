package entities

// Link ...
type Link struct {
	ServerID string `json:",omitempty"`
	Code     string
	URL      string
}
