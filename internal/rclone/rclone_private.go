package rclone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

// Format an error and create a synthetic server return from it
func errorf(status int, path string, format string, arg ...any) (out rc.Params, err error) {
	err = fmt.Errorf(format, arg...)
	out = make(rc.Params)
	out["error"] = err.Error()
	out["path"] = path
	out["status"] = status
	return out, err
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
