package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/fshttp"
	"github.com/rclone/rclone/fs/rc"
)

// consider using a configuration implementation here
// using file, env vars, args
var (
	noOutput  = false
	url       = "http://localhost:5572/"
	jsonInput = ""
	authUser  = ""
	authPass  = ""
	loopback  = false
	options   []string
	arguments []string
)

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

type ListConfig struct {
	Fs     string      `json:"fs"`
	Remote string      `json:"remote"`
	Opt    *ListOption `json:"opt"`
}

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

type ListResponse struct {
	List []ListResponseItem `json:"list"`
}

// Format an error and create a synthetic server return from it
func errorf(status int, path string, format string, arg ...any) (out rc.Params, err error) {
	err = fmt.Errorf(format, arg...)
	out = make(rc.Params)
	out["error"] = err.Error()
	out["path"] = path
	out["status"] = status
	return out, err
}

func ListJSON(ctx context.Context, fs string, remote string) error {
	path := "operations/list"
	config := ListConfig{
		Fs:     fs,
		Remote: remote,
		Opt: &ListOption{
			ShowHash:  true,
			HashTypes: []string{"dropbox", "md5"},
		},
	}
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to encode configuration: %w", err)
	}

	in := make(rc.Params)
	err = json.Unmarshal(configData, &in)
	if err != nil {
		return fmt.Errorf("failed to decode configuration: %w", err)
	}

	out, callErr := doCall(ctx, path, in)

	// Write the JSON blob to stdout if required
	if out != nil && !noOutput {
		err := rc.WriteJSON(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("failed to output JSON: %w", err)
		}
	}

	if out != nil {
		var lr ListResponse
		err := mapstructure.Decode(out, &lr)

		if err != nil {
			return fmt.Errorf("failed to map structure")
		}

		err = dumpListResponse(lr)
		if err != nil {
			return fmt.Errorf("failed to dump output JSON: %w", err)
		}
	}

	return callErr
}

func doCall(ctx context.Context, path string, in rc.Params) (out rc.Params, err error) {
	client := fshttp.NewClient(ctx)
	url += path
	data, err := json.Marshal(in)

	if err != nil {
		return errorf(http.StatusBadRequest, path, "failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return errorf(http.StatusInternalServerError, path, "failed to make request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authUser != "" || authPass != "" {
		req.SetBasicAuth(authUser, authPass)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errorf(http.StatusServiceUnavailable, path, "connection failed: %w", err)
	}
	defer fs.CheckClose(resp.Body, &err)

	// Read response
	var body []byte
	var bodyString string
	body, err = io.ReadAll(resp.Body)
	bodyString = strings.TrimSpace(string(body))
	if err != nil {
		return errorf(resp.StatusCode, "failed to read rc response: %s: %s", resp.Status, bodyString)
	}

	// Parse output
	out = make(rc.Params)
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(&out)
	if err != nil {
		return errorf(resp.StatusCode, path, "failed to decode response: %w: %s", err, bodyString)
	}

	// Check we got 200 OK
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("operation %q failed: %v", path, out["error"])
	}

	return out, err
}

func dumpListResponse(lr ListResponse) error {
	for _, item := range lr.List {
		fmt.Printf("Item Id %s name %s\n%s %d %t\n", item.ID, item.Name, item.Path, item.Size, item.IsDir)
	}
	return nil
}
