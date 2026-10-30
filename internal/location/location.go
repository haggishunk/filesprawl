package location

// Location identifies the source of indexed data.
// For local scans: name=hostname, type='local'
// For rclone remotes: name=remote identifier, type='remote'
type Location struct {
	Id        int
	Name      string // hostname (local) or remote identifier (remote)
	Type      string // 'local' or 'remote'
	Persisted bool
}

type LocationOption func(*Location)

func WithLocationID(id int) LocationOption {
	return func(loc *Location) {
		loc.Id = id
	}
}

// NewLocation creates a location for identifying a data source.
// name: hostname for local scans, remote identifier for rclone remotes
// locationType: 'local' or 'remote'
func NewLocation(name, locationType string, options ...LocationOption) Location {
	loc := Location{
		Id:        -1,
		Name:      name,
		Type:      locationType,
		Persisted: false,
	}
	
	for _, opt := range options {
		opt(&loc)
	}
	
	return loc
}

// LocalLocation creates a location for a local scan.
func LocalLocation(hostname string) Location {
	return NewLocation(hostname, "local")
}

// RemoteLocation creates a location for a rclone remote.
func RemoteLocation(remoteName string) Location {
	return NewLocation(remoteName, "remote")
}

// IsLocal returns true if this location represents a local scan.
func (loc *Location) IsLocal() bool {
	return loc.Type == "local"
}

// IsRemote returns true if this location represents a rclone remote.
func (loc *Location) IsRemote() bool {
	return loc.Type == "remote"
}
