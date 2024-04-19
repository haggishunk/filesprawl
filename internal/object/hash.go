package object

import "fmt"

type Hash struct {
	Id        int
	Hash      string
	Type      string
	Persisted bool
}

func (h Hash) String() string {
	return fmt.Sprintf("{hash: %s, type: %s}", h.Hash, h.Type)
}

type HashOption func(*Hash)

func WithHashId(i int) HashOption {
	return func(h *Hash) {
		h.Id = i
	}
}

func NewHash(h string, t string, ho ...HashOption) Hash {
	n := Hash{
		Id:        -1,
		Hash:      h,
		Type:      t,
		Persisted: false,
	}
	for _, f := range ho {
		f(&n)
	}
	return n
}
