package remote

import (
	"github.com/haggishunk/filesprawl/internal/location"
)

type Remote struct {
	Id        int
	LocationId int // FK to location table
	Location  location.Location
	Name      string // path or remote name within the location
	Persisted bool
}

type RemoteOption func(*Remote)

func WithRemoteID(id int) RemoteOption {
	return func(r *Remote) {
		r.Id = id
	}
}

func WithLocation(loc location.Location) RemoteOption {
	return func(r *Remote) {
		r.Location = loc
		r.LocationId = loc.Id
	}
}

// NewRemote creates a remote for the given location and name.
// name: path/subdirectory for local scans or remote name for rclone
func NewRemote(loc location.Location, name string, options ...RemoteOption) Remote {
	r := Remote{
		Id:        -1,
		LocationId: loc.Id,
		Location:  loc,
		Name:      name,
		Persisted: false,
	}

	for _, opt := range options {
		opt(&r)
	}

	return r
}

// NewLocalRemote creates a remote for a local scan.
// hostname: the machine hostname
// name: the path/directory name
func NewLocalRemote(hostname string, name string, options ...RemoteOption) Remote {
	loc := location.LocalLocation(hostname)
	return NewRemote(loc, name, options...)
}

// NewRcloneRemote creates a remote for an rclone remote.
// remoteName: the rclone remote identifier
// path: the path within the remote
func NewRcloneRemote(remoteName string, path string, options ...RemoteOption) Remote {
	loc := location.RemoteLocation(remoteName)
	return NewRemote(loc, path, options...)
}
