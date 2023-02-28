package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ MigrationHistoriesStore = (*migrationHistories)(nil)

// MigrationHistories is the default instance of the MigrationHistoriesStore.
var MigrationHistories MigrationHistoriesStore

// MigrationHistoriesStore is the persistent interface for migrationHistories.
type MigrationHistoriesStore interface {
	// GetById returns a migrationHistory with the given id.
	GetById(ctx context.Context, id int64) (*MigrationHistory, error)
	// GetBySequence returns a migrationHistory with the given sequence.
	GetBySequence(ctx context.Context, sequence int64) (*MigrationHistory, error)
	// GetByVersion returns a migrationHistory with the given version.
	GetByVersion(ctx context.Context, version string) (*MigrationHistory, error)
	// GetByType returns a migrationHistory with the given type.
	GetByType(ctx context.Context, typ string) (*MigrationHistory, error)
	// GetByCreatedTs returns a migrationHistory with the given createdTs.
	GetByCreatedTs(ctx context.Context, createdTs int64) (*MigrationHistory, error)
}

// NewMigrationHistoriesStore returns a MigrationHistoriesStore instance with the given database connection.
func NewMigrationHistoriesStore(db *gorm.DB) MigrationHistoriesStore {
	return &migrationHistories{db}
}

// MigrationHistory represents the migrationHistories.
type MigrationHistory struct {
	gorm.Model

	CreatedBy           string
	CreatedTs           int64
	UpdatedBy           string
	UpdatedTs           int64
	ReleaseVersion      string
	Namespace           string
	Sequence            int64
	Source              string
	Type                string
	Status              string
	Version             string
	Description         string
	Statement           string
	Schema              string
	SchemaPrev          string
	ExecutionDurationNs int64
	IssueID             string
	Payload             string
}

type migrationHistories struct {
	*gorm.DB
}

var ErrMigrationHistoryNotExists = errors.New("migrationHistory dose not exist")

func (db *migrationHistories) GetById(ctx context.Context, id int64) (*MigrationHistory, error) {
	var migrationHistory MigrationHistory
	if err := db.WithContext(ctx).Model(&MigrationHistory{}).Where("id = ?", id).First(&migrationHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMigrationHistoryNotExists
		}
	}
	return &migrationHistory, nil
}

func (db *migrationHistories) GetBySequence(ctx context.Context, sequence int64) (*MigrationHistory, error) {
	var migrationHistory MigrationHistory
	if err := db.WithContext(ctx).Model(&MigrationHistory{}).Where("sequence = ?", sequence).First(&migrationHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMigrationHistoryNotExists
		}
	}
	return &migrationHistory, nil
}

func (db *migrationHistories) GetByVersion(ctx context.Context, version string) (*MigrationHistory, error) {
	var migrationHistory MigrationHistory
	if err := db.WithContext(ctx).Model(&MigrationHistory{}).Where("version = ?", version).First(&migrationHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMigrationHistoryNotExists
		}
	}
	return &migrationHistory, nil
}

func (db *migrationHistories) GetByType(ctx context.Context, typ string) (*MigrationHistory, error) {
	var migrationHistory MigrationHistory
	if err := db.WithContext(ctx).Model(&MigrationHistory{}).Where("type = ?", typ).First(&migrationHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMigrationHistoryNotExists
		}
	}
	return &migrationHistory, nil
}

func (db *migrationHistories) GetByCreatedTs(ctx context.Context, createdTs int64) (*MigrationHistory, error) {
	var migrationHistory MigrationHistory
	if err := db.WithContext(ctx).Model(&MigrationHistory{}).Where("createdTs = ?", createdTs).First(&migrationHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMigrationHistoryNotExists
		}
	}
	return &migrationHistory, nil
}
