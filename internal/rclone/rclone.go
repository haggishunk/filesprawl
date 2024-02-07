package rclone

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/rclone/rclone/fs/rc"
)

// ListOption structifies options related to rc operations/list
type ListOption struct {
	Recurse       bool     `json:"recurse"`
	NoModTime     bool     `json:"noModTime"`
	ShowEncrypted bool     `json:"showEncrypted"`
	ShowOrigIDs   bool     `json:"showOrigIDs"`
	ShowHash      bool     `json:"showHash"`
	NoMimeType    bool     `json:"noMimeType"`
	DirsOnly      bool     `json:"dirsOnly"`
	FilesOnly     bool     `json:"filesOnly"`
	Metadata      bool     `json:"metadata"`
	HashTypes     []string `json:"hashTypes"`
}

// ListConfig structifies configuration related to rc operations/list
type ListConfig struct {
	Fs     string      `json:"fs"`
	Remote string      `json:"remote"`
	Opt    *ListOption `json:"opt"`
}

// ListResponseItem structifies the response items from rc operations/list
type ListResponseItem struct {
	ID       string            `json:"ID"`
	IsDir    bool              `json:"IsDir"`
	MimeType string            `json:"MimeType"`
	ModTime  string            `json:"ModTime"`
	Name     string            `json:"Name"`
	Path     string            `json:"Path"`
	Size     int64             `json:"Size"`
	Hashes   map[string]string `json:"Hashes,omitempty"`
}

// ListResponse structifies the response from rc operations/list
type ListResponse struct {
	List []ListResponseItem `json:"list"`
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

// EncodeListResponse takes a generic return from rc calls and
// encodes into a ListReponse struct
func EncodeListResponse(rcp rc.Params, lr *ListResponse) error {
	err := mapstructure.Decode(rcp, &lr)
	if err != nil {
		return fmt.Errorf("failed to map structure")
	}
	return nil
}

// ListJSON queries an rc server for objects
//
// options for handling this expensive call:
// - use it to walk the remote dirs and listing objects
// - pass in a handler for it to sink objects to
// - pass in a channel ref for it to sink objects to
//
// err := clients.ListJSON(ctx, "remote:", "path/in/remote")
func ListJSON(ctx context.Context, fs string, remote string) error {
	// fixed path means only one operation type (ie `list`)
	path := "operations/list"

	// operation configured using some inputs
	config := ListConfig{
		Fs:     fs,
		Remote: remote,
		Opt: &ListOption{
			ShowHash:  true,
			HashTypes: []string{"dropbox", "md5"},
		},
	}
	configDecoded, err := config.Decode()
	if err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	// call is made to lower level "imported" call to rc
	out, callErr := doCall(ctx, path, configDecoded)

	// outputs and handling

	// Write the JSON blob to stdout if required
	if out != nil && !noOutput {
		err := rc.WriteJSON(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("failed to output JSON: %w", err)
		}
	}

	// Convert generic output to a list response and do something with it
	if out != nil {
		var lr = ListResponse{}
		err := EncodeListResponse(out, &lr)
		if err != nil {
			return fmt.Errorf("failed to encode list response: %w", err)
		}
		err = dumpListResponse(lr)
		if err != nil {
			return fmt.Errorf("failed to dump output JSON: %w", err)
		}
	}

	return callErr
}

// dumpListResponse is a toy to print a ListResponse to stdout
func dumpListResponse(lr ListResponse) error {
	for _, item := range lr.List {
		fmt.Printf("Item Id %s name %s\n%s %d %t\n", item.ID, item.Name, item.Path, item.Size, item.IsDir)
	}
	return nil
}
