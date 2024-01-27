package object

import "path"

type Object struct {
	Name string
	Path string
}

func (o *Object) Base() string {
	return path.Base(o.Path)
}
