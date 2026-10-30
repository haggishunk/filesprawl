package flags

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// FlagType represents the type of consolidation decision
type FlagType string

const (
	FlagTypeKeep   FlagType = "KEEP"
	FlagTypeArchive FlagType = "ARCHIVE"
	FlagTypeDelete  FlagType = "DELETE"
)

// Flag represents a consolidation flag decision
type Flag struct {
	ID          int
	FileID      int
	UserID      int
	FlagType    FlagType
	FlaggedAt   time.Time
	Reason      *string
	ArchiveDest *string
}

// Review represents a review decision in the audit trail
type Review struct {
	ID           int
	FileID       int
	OriginalFlag *FlagType
	OriginalUser *int
	OriginalAt   *time.Time
	ReviewedBy   int
	ReviewedAt   time.Time
	NewFlag      *FlagType
	ReviewReason *string
}

// Repository interface for flag and review operations
type Repository interface {
	// CreateFlag stores a new flag decision
	CreateFlag(ctx context.Context, fileID int, userID int, flagType FlagType, reason *string, archiveDest *string) (*Flag, error)
	
	// GetFlag retrieves a flag by ID
	GetFlag(ctx context.Context, id int) (*Flag, error)
	
	// GetFlagsByFile retrieves all flags for a file
	GetFlagsByFile(ctx context.Context, fileID int) ([]Flag, error)
	
	// GetFlagsByUser retrieves all flags created by a user
	GetFlagsByUser(ctx context.Context, userID int) ([]Flag, error)
	
	// UpdateFlag updates a flag decision
	UpdateFlag(ctx context.Context, flagID int, flagType FlagType, reason *string, archiveDest *string) (*Flag, error)
	
	// DeleteFlag removes a flag
	DeleteFlag(ctx context.Context, flagID int) error
	
	// CreateReview stores a review of a flag decision
	CreateReview(ctx context.Context, fileID int, originalFlag *FlagType, originalUser *int, originalAt *time.Time, reviewedBy int, newFlag *FlagType, reviewReason *string) (*Review, error)
	
	// GetReviewsByFile retrieves all reviews for a file
	GetReviewsByFile(ctx context.Context, fileID int) ([]Review, error)
	
	// GetReviewsByUser retrieves all reviews made by a user
	GetReviewsByUser(ctx context.Context, userID int) ([]Review, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL-backed flag repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateFlag creates a new flag
func (r *PostgresRepository) CreateFlag(ctx context.Context, fileID int, userID int, flagType FlagType, reason *string, archiveDest *string) (*Flag, error) {
	if fileID <= 0 || userID <= 0 {
		return nil, errors.New("fileID and userID must be positive")
	}

	flag := &Flag{
		FileID:      fileID,
		UserID:      userID,
		FlagType:    flagType,
		FlaggedAt:   time.Now(),
		Reason:      reason,
		ArchiveDest: archiveDest,
	}

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO consolidation_flags (file_id, user_id, flag_type, flagged_at, reason, archive_dest)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, flagged_at`,
		fileID, userID, string(flagType), flag.FlaggedAt, reason, archiveDest,
	).Scan(&flag.ID, &flag.FlaggedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create flag: %w", err)
	}

	return flag, nil
}

// GetFlag retrieves a flag by ID
func (r *PostgresRepository) GetFlag(ctx context.Context, id int) (*Flag, error) {
	flag := &Flag{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, file_id, user_id, flag_type, flagged_at, reason, archive_dest
		 FROM consolidation_flags WHERE id = $1`,
		id,
	).Scan(&flag.ID, &flag.FileID, &flag.UserID, (*string)(&flag.FlagType), &flag.FlaggedAt, &flag.Reason, &flag.ArchiveDest)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("flag not found: id=%d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get flag: %w", err)
	}

	return flag, nil
}

// GetFlagsByFile retrieves all flags for a file
func (r *PostgresRepository) GetFlagsByFile(ctx context.Context, fileID int) ([]Flag, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, file_id, user_id, flag_type, flagged_at, reason, archive_dest
		 FROM consolidation_flags WHERE file_id = $1
		 ORDER BY flagged_at DESC`,
		fileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get flags: %w", err)
	}
	defer rows.Close()

	var flags []Flag
	for rows.Next() {
		var flag Flag
		var flagTypeStr string
		if err := rows.Scan(&flag.ID, &flag.FileID, &flag.UserID, &flagTypeStr, &flag.FlaggedAt, &flag.Reason, &flag.ArchiveDest); err != nil {
			return nil, fmt.Errorf("failed to scan flag: %w", err)
		}
		flag.FlagType = FlagType(flagTypeStr)
		flags = append(flags, flag)
	}

	return flags, rows.Err()
}

// GetFlagsByUser retrieves all flags by a user
func (r *PostgresRepository) GetFlagsByUser(ctx context.Context, userID int) ([]Flag, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, file_id, user_id, flag_type, flagged_at, reason, archive_dest
		 FROM consolidation_flags WHERE user_id = $1
		 ORDER BY flagged_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get flags: %w", err)
	}
	defer rows.Close()

	var flags []Flag
	for rows.Next() {
		var flag Flag
		var flagTypeStr string
		if err := rows.Scan(&flag.ID, &flag.FileID, &flag.UserID, &flagTypeStr, &flag.FlaggedAt, &flag.Reason, &flag.ArchiveDest); err != nil {
			return nil, fmt.Errorf("failed to scan flag: %w", err)
		}
		flag.FlagType = FlagType(flagTypeStr)
		flags = append(flags, flag)
	}

	return flags, rows.Err()
}

// UpdateFlag updates a flag
func (r *PostgresRepository) UpdateFlag(ctx context.Context, flagID int, flagType FlagType, reason *string, archiveDest *string) (*Flag, error) {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE consolidation_flags SET flag_type = $1, reason = $2, archive_dest = $3 WHERE id = $4`,
		string(flagType), reason, archiveDest, flagID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update flag: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return nil, fmt.Errorf("flag not found: id=%d", flagID)
	}

	return r.GetFlag(ctx, flagID)
}

// DeleteFlag removes a flag
func (r *PostgresRepository) DeleteFlag(ctx context.Context, flagID int) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM consolidation_flags WHERE id = $1`,
		flagID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete flag: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("flag not found: id=%d", flagID)
	}

	return nil
}

// CreateReview creates a new review
func (r *PostgresRepository) CreateReview(ctx context.Context, fileID int, originalFlag *FlagType, originalUser *int, originalAt *time.Time, reviewedBy int, newFlag *FlagType, reviewReason *string) (*Review, error) {
	if fileID <= 0 || reviewedBy <= 0 {
		return nil, errors.New("fileID and reviewedBy must be positive")
	}

	review := &Review{
		FileID:       fileID,
		OriginalFlag: originalFlag,
		OriginalUser: originalUser,
		OriginalAt:   originalAt,
		ReviewedBy:   reviewedBy,
		ReviewedAt:   time.Now(),
		NewFlag:      newFlag,
		ReviewReason: reviewReason,
	}

	var originalFlagStr, newFlagStr *string
	if originalFlag != nil {
		s := string(*originalFlag)
		originalFlagStr = &s
	}
	if newFlag != nil {
		s := string(*newFlag)
		newFlagStr = &s
	}

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO consolidation_reviews (file_id, original_flag, original_user, original_at, reviewed_by, reviewed_at, new_flag, review_reason)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, reviewed_at`,
		fileID, originalFlagStr, originalUser, originalAt, reviewedBy, review.ReviewedAt, newFlagStr, reviewReason,
	).Scan(&review.ID, &review.ReviewedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

// GetReviewsByFile retrieves all reviews for a file
func (r *PostgresRepository) GetReviewsByFile(ctx context.Context, fileID int) ([]Review, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, file_id, original_flag, original_user, original_at, reviewed_by, reviewed_at, new_flag, review_reason
		 FROM consolidation_reviews WHERE file_id = $1
		 ORDER BY reviewed_at DESC`,
		fileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		var origFlagStr, newFlagStr *string
		if err := rows.Scan(&review.ID, &review.FileID, &origFlagStr, &review.OriginalUser, &review.OriginalAt, &review.ReviewedBy, &review.ReviewedAt, &newFlagStr, &review.ReviewReason); err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		if origFlagStr != nil {
			f := FlagType(*origFlagStr)
			review.OriginalFlag = &f
		}
		if newFlagStr != nil {
			f := FlagType(*newFlagStr)
			review.NewFlag = &f
		}
		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}

// GetReviewsByUser retrieves all reviews by a user
func (r *PostgresRepository) GetReviewsByUser(ctx context.Context, userID int) ([]Review, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, file_id, original_flag, original_user, original_at, reviewed_by, reviewed_at, new_flag, review_reason
		 FROM consolidation_reviews WHERE reviewed_by = $1
		 ORDER BY reviewed_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		var origFlagStr, newFlagStr *string
		if err := rows.Scan(&review.ID, &review.FileID, &origFlagStr, &review.OriginalUser, &review.OriginalAt, &review.ReviewedBy, &review.ReviewedAt, &newFlagStr, &review.ReviewReason); err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		if origFlagStr != nil {
			f := FlagType(*origFlagStr)
			review.OriginalFlag = &f
		}
		if newFlagStr != nil {
			f := FlagType(*newFlagStr)
			review.NewFlag = &f
		}
		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}
