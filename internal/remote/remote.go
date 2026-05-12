package remote

type Remote struct {
	Id        int
	Hostname  string
	Name      string // different hosts may have same remote with different name
	Type      string
	Persisted bool
}

type RemoteOption func(*Remote)

func WithRemoteID(id int) RemoteOption {
	return func(r *Remote) {
		r.Id = id
	}
}

func WithRemoteType(kind string) RemoteOption {
	return func(r *Remote) {
		r.Type = kind
	}
}

func NewRemote(hostname string, name string, options ...RemoteOption) Remote {
	r := Remote{
		Id:        -1,
		Hostname:  hostname,
		Name:      name,
		Type:      "rclone",
		Persisted: false,
	}

	for _, opt := range options {
		opt(&r)
	}

	return r
}
