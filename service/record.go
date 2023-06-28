package service

import (
	"context"
	"errors"
	"sync"

	"github.com/temelpa/timetravel/entity"
)

var ErrRecordDoesNotExist = errors.New("record with that id does not exist")
var ErrRecordIDInvalid = errors.New("record id must >= 0")
var ErrRecordAlreadyExists = errors.New("record already exists")

type RecordServiceBase interface {
	// TODO: It seems awkward for a rwlock to be used mostly outside the
	// struct that owns it. Ideally the service ought to be the one locking.
	// Try to make things structured so that the lock lives elsewhere,
	// or the service is the true user of the lock.
	GetRWLockForAPI() *sync.RWMutex
}

// Implements method to get, create, and update record data.
type RecordServiceV1 interface {
	RecordServiceBase

	// GetRecord will retrieve an record.
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// CreateRecord will insert a new record.
	//
	// If it a record with that id already exists it will fail.
	CreateRecord(ctx context.Context, record entity.Record) error

	// UpdateRecord will change the internal `Map` values of the record if they exist.
	// if the update[key] is null it will delete that key from the record's Map.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error)
}

// Introduce the concept of record versions. Versions will start at 1 and increment per
// later versions that exist
type RecordServiceV2 interface {
	RecordServiceBase
	RecordServiceV1

	// Retrieve a record. If `version` is nil or 0, return the latest version that
	// exists.
	GetVersionedRecord(ctx context.Context, id int, version int) (entity.Record, error)

	// Retrieves all versions of a record.
	GetAllRecordVersions(ctx context.Context, id int) ([]entity.Record, error)
}

type RecordService interface {
	// The current supported max API level
	RecordServiceV2
}
