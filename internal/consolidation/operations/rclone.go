package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/rclone/rclone/fs/rc"
	"github.com/haggishunk/filesprawl/internal/rclone"
)

// Archive copies a file to an archive destination
func (re *RcloneExecutor) Archive(ctx context.Context, file ConsolidationFile, destination string) *OperationResult {
	result := &OperationResult{
		FileID:    file.ID,
		Operation: OperationArchive,
		Success:   false,
	}

	if re.dryRun {
		result.Success = true
		result.Details = fmt.Sprintf("DRY-RUN: Would archive %s:%s to %s", file.Remote, file.Path, destination)
		return result
	}

	// Parse destination format (e.g., "archive:/path/to/archive")
	parts := strings.SplitN(destination, ":", 2)
	if len(parts) != 2 {
		result.Error = fmt.Sprintf("invalid destination format: %s (use remote:path)", destination)
		return result
	}

	destRemote := parts[0]
	destPath := strings.TrimLeft(parts[1], "/")

	// Construct full destination path
	fileName := extractFileName(file.Path)
	fullDest := fmt.Sprintf("%s/%s", strings.TrimRight(destPath, "/"), fileName)

	// Call rclone rc copyfile operation
	params := rc.Params{
		"srcFs":  fmt.Sprintf("%s:", file.Remote),
		"srcRemote": file.Path,
		"dstFs":  fmt.Sprintf("%s:", destRemote),
		"dstRemote": fullDest,
	}

	if err := rclone.CallRC(ctx, "operations/copyfile", params); err != nil {
		result.Error = fmt.Sprintf("rclone copyfile failed: %v", err)
		return result
	}

	result.Success = true
	result.Details = fmt.Sprintf("Archived %s:%s to %s:%s", file.Remote, file.Path, destRemote, fullDest)
	return result
}

// Delete removes a file from its remote location
func (re *RcloneExecutor) Delete(ctx context.Context, file ConsolidationFile) *OperationResult {
	result := &OperationResult{
		FileID:    file.ID,
		Operation: OperationDelete,
		Success:   false,
	}

	if re.dryRun {
		result.Success = true
		result.Details = fmt.Sprintf("DRY-RUN: Would delete %s:%s", file.Remote, file.Path)
		return result
	}

	// Call rclone rc delete operation
	params := rc.Params{
		"fs":     fmt.Sprintf("%s:", file.Remote),
		"remote": file.Path,
	}

	if err := rclone.CallRC(ctx, "operations/delete", params); err != nil {
		result.Error = fmt.Sprintf("rclone delete failed: %v", err)
		return result
	}

	result.Success = true
	result.Details = fmt.Sprintf("Deleted %s:%s", file.Remote, file.Path)
	return result
}

// Move relocates a file to a new path
func (re *RcloneExecutor) Move(ctx context.Context, file ConsolidationFile, destination string) *OperationResult {
	result := &OperationResult{
		FileID:    file.ID,
		Operation: OperationMove,
		Success:   false,
	}

	if re.dryRun {
		result.Success = true
		result.Details = fmt.Sprintf("DRY-RUN: Would move %s:%s to %s", file.Remote, file.Path, destination)
		return result
	}

	// Parse destination format (e.g., "remote:/path/to/file")
	parts := strings.SplitN(destination, ":", 2)
	if len(parts) != 2 {
		result.Error = fmt.Sprintf("invalid destination format: %s (use remote:path)", destination)
		return result
	}

	destRemote := parts[0]
	destPath := strings.TrimLeft(parts[1], "/")

	// Call rclone rc movefile operation
	params := rc.Params{
		"srcFs":  fmt.Sprintf("%s:", file.Remote),
		"srcRemote": file.Path,
		"dstFs":  fmt.Sprintf("%s:", destRemote),
		"dstRemote": destPath,
	}

	if err := rclone.CallRC(ctx, "operations/movefile", params); err != nil {
		result.Error = fmt.Sprintf("rclone movefile failed: %v", err)
		return result
	}

	result.Success = true
	result.Details = fmt.Sprintf("Moved %s:%s to %s:%s", file.Remote, file.Path, destRemote, destPath)
	return result
}

// extractFileName extracts the filename from a full path
func extractFileName(path string) string {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}
