package rclone

import (
	"context"
	"fmt"

	"github.com/rclone/rclone/fs/rc"
)

// ListJSON queries an rc server for objects
//
// options for handling this expensive call:
// - use it to walk the remote dirs and listing objects
// - pass in a handler for it to sink objects to
// - pass in a channel ref for it to sink objects to
func ListJSON(ctx context.Context, lc ListConfig) (*ListResponse, error) {
	configDecoded, err := lc.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// fixed path means only one operation type (ie `list`)
	path := "operations/list"
	// call is made to lower level "imported" call to rc
	out, callErr := doCall(ctx, path, configDecoded)
	if callErr != nil {
		return nil, fmt.Errorf("failed to make rclone rc call: %w", callErr)
	}
	var lr = ListResponse{}
	if out == nil {
		return &lr, nil
	}
	err = EncodeListResponse(out, &lr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode list response: %w", err)
	}
	return &lr, nil
}

// ListRemotes queries an rc server for configured remotes.
func ListRemotes(ctx context.Context) (*ListRemotesResponse, error) {
	out, callErr := doCall(ctx, "config/listremotes", rc.Params{})
	if callErr != nil {
		return nil, fmt.Errorf("failed to list remotes: %w", callErr)
	}

	var lr = ListRemotesResponse{}
	if out == nil {
		return &lr, nil
	}
	err := EncodeListRemotesResponse(out, &lr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode remote list response: %w", err)
	}
	return &lr, nil
}
