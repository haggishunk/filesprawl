package rclone

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rclone/rclone/fs/rc"
)

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

func (lri ListResponseItem) String() string {
	return fmt.Sprintf("Name: %s, Path: %s", lri.Name, lri.Path)
}

// ListResponse structifies the response from rc operations/list
type ListResponse struct {
	List []ListResponseItem `json:"list"`
}

// ListRemotesResponse structifies the response from rc config/listremotes.
type ListRemotesResponse struct {
	Remotes []string `json:"remotes"`
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

// EncodeListRemotesResponse takes a generic return from rc calls and
// encodes it into a ListRemotesResponse struct.
func EncodeListRemotesResponse(rcp rc.Params, lr *ListRemotesResponse) error {
	err := mapstructure.Decode(rcp, &lr)
	if err != nil {
		return fmt.Errorf("failed to map structure")
	}
	return nil
}
