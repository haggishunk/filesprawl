package object

import (
	"fmt"
	"path"
)

type Meta struct {
	Id   int
	Name string
	Path string
	// BasePath represents the object store paths investigated
	// keep these as metadata for scan result object
	// and for freshness of scan directory targeting
	BasePath  string
	MimeType  string
	Persisted bool
}

func (o Meta) String() string {
	return fmt.Sprintf("{Name: %s, Path: %s}", o.Name, o.Path)
}

type MetaOption func(*Meta)

func WithMetaId(i int) MetaOption {
	return func(o *Meta) {
		o.Id = i
	}
}

func NewMeta(n string, p string, m string, oo ...MetaOption) Meta {
	o := Meta{
		Id:        -1,
		Name:      n,
		Path:      p,
		BasePath:  path.Base(p),
		MimeType:  m,
		Persisted: false,
	}
	for _, of := range oo {
		of(&o)
	}
	return o
}
