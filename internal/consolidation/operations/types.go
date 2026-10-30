package operations

import "context"

// OperationType represents the type of consolidation operation
type OperationType string

const (
	OperationArchive OperationType = "ARCHIVE"
	OperationDelete  OperationType = "DELETE"
	OperationMove    OperationType = "MOVE"
)

// OperationResult represents the result of a single consolidation operation
type OperationResult struct {
	FileID    int
	Operation OperationType
	Success   bool
	Error     string
	Details   string // Additional info like destination path
}

// ConsolidationFile represents a file to be consolidated
type ConsolidationFile struct {
	ID       int
	Path     string
	Remote   string
	FileSize int64
	Hash     string
}

// OperationExecutor defines the interface for executing consolidation operations
type OperationExecutor interface {
	// Archive copies a file to an archive destination
	Archive(ctx context.Context, file ConsolidationFile, destination string) *OperationResult

	// Delete removes a file from its remote location
	Delete(ctx context.Context, file ConsolidationFile) *OperationResult

	// Move relocates a file to a new path
	Move(ctx context.Context, file ConsolidationFile, destination string) *OperationResult

	// IsDryRun returns whether executor is in dry-run mode
	IsDryRun() bool

	// SetDryRun sets dry-run mode
	SetDryRun(dryRun bool)
}

// RcloneExecutor implements OperationExecutor using rclone RC
type RcloneExecutor struct {
	dryRun bool
}

// NewRcloneExecutor creates a new rclone-based operation executor
func NewRcloneExecutor() *RcloneExecutor {
	return &RcloneExecutor{
		dryRun: false,
	}
}

// IsDryRun returns whether executor is in dry-run mode
func (re *RcloneExecutor) IsDryRun() bool {
	return re.dryRun
}

// SetDryRun sets dry-run mode
func (re *RcloneExecutor) SetDryRun(dryRun bool) {
	re.dryRun = dryRun
}
