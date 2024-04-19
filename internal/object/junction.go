package object

import (
	"fmt"
	"time"
)

type MetaHashJunction struct {
	Id        int
	ScanTime  time.Time
	MetaId    int
	HashId    int
	Persisted bool
}

func (o MetaHashJunction) String() string {
	return fmt.Sprintf("{MetaId: %d, HashId: %d}", o.MetaId, o.HashId)
}

type MetaHashJunctionOption func(*MetaHashJunction)

func WithMetaHashJunctionId(i int) MetaHashJunctionOption {
	return func(o *MetaHashJunction) {
		o.Id = i
	}
}

func NewMetaHashJunction(mid int, hid int, oo ...MetaHashJunctionOption) MetaHashJunction {
	o := MetaHashJunction{
		Id:        -1,
		ScanTime:  time.Now(),
		MetaId:    mid,
		HashId:    hid,
		Persisted: false,
	}
	for _, of := range oo {
		of(&o)
	}
	return o
}
