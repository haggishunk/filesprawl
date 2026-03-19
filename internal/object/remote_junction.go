package object

type ObjectRemoteJunction struct {
	Id        int
	MetaId    int
	RemoteId  int
	Persisted bool
}

type ObjectRemoteJunctionOption func(*ObjectRemoteJunction)

func WithObjectRemoteJunctionID(id int) ObjectRemoteJunctionOption {
	return func(o *ObjectRemoteJunction) {
		o.Id = id
	}
}

func NewObjectRemoteJunction(metaID int, remoteID int, options ...ObjectRemoteJunctionOption) ObjectRemoteJunction {
	o := ObjectRemoteJunction{
		Id:        -1,
		MetaId:    metaID,
		RemoteId:  remoteID,
		Persisted: false,
	}

	for _, opt := range options {
		opt(&o)
	}

	return o
}
