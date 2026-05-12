package rclone

import (
	"encoding/json"
	"fmt"

	"github.com/rclone/rclone/fs/rc"
)

type ListOptionFunc func(*ListOption)
type ListConfigFunc func(*ListConfig)

// ListOption structifies options related to rc operations/list
type ListOption struct {
	DirsOnly      bool     `json:"dirsOnly"`
	FilesOnly     bool     `json:"filesOnly"`
	HashTypes     []string `json:"hashTypes"`
	Metadata      bool     `json:"metadata"`
	NoMimeType    bool     `json:"noMimeType"`
	NoModTime     bool     `json:"noModTime"`
	Recurse       bool     `json:"recurse"`
	ShowEncrypted bool     `json:"showEncrypted"`
	ShowHash      bool     `json:"showHash"`
	ShowOrigIDs   bool     `json:"showOrigIDs"`
}

// ListConfig structifies configuration related to rc operations/list
// note: rclone implementation for configuring operations is less strict
// and uses rc.Params interface{}
type ListConfig struct {
	Fs     string      `json:"fs"`     // the remote name is the Fs
	Remote string      `json:"remote"` // the path is the Remote
	Opt    *ListOption `json:"opt"`
}

func ListOptionDirsOnly() ListOptionFunc {
	return func(lo *ListOption) {
		lo.DirsOnly = true
	}
}

func ListOptionFilesOnly() ListOptionFunc {
	return func(lo *ListOption) {
		lo.FilesOnly = true
	}
}

func NewListOption(options ...ListOptionFunc) ListOption {
	// operation configured using some inputs
	lo := ListOption{
		DirsOnly:      false,
		FilesOnly:     false,
		HashTypes:     []string{},
		NoMimeType:    true,
		NoModTime:     true,
		Recurse:       false,
		ShowEncrypted: false,
		ShowHash:      true,
		ShowOrigIDs:   true,
	}

	for _, opt := range options {
		opt(&lo)
	}

	return lo
}

// NewListConfig creates a struct to be used in a ListJSON operation
func NewListConfig(fs string, r string, o *ListOption) ListConfig {
	return ListConfig{
		Fs:     fs,
		Remote: r,
		Opt:    o,
	}
}

// Decode converts a ListConfig into an generic rclone Params
// for use with rc calls
func (lr *ListConfig) Decode() (rc.Params, error) {
	data, err := json.Marshal(lr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode configuration: %w", err)
	}

	decoded := make(rc.Params)
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	return decoded, nil
}
