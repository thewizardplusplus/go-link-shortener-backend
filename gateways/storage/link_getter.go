package storage

// LinkGetter ...
type LinkGetter struct {
	Client     Client
	Database   string
	Collection string
	KeyField   string
}
