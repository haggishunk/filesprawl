package operation

import (
	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/remote"
)

// scanning:
// recursively descending through directories
// retrieving objects and persisting to a database
// persist object aspects in relationships with scan timestamp

// synthesize a list response with the targeted remote
type ScanResult struct {
	Remote *remote.Remote
	Hash   *object.Hash
	Object *object.Object
	// Timestamp
}

// ScanRemote begins the scanning of a remote
func ScanRemote() {}

// save new remotes, hashes, objects or retrieve existing ids
// save relationships betweeen these items with timestamp
func PersistScanResult() {}
