package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	// This library uses cgo, which can complicate the
	// build environment and portability of the code, but
	// offers better performance at larger scales. For lack
	// of knowning exactly how many records this service will
	// manage, using this now is a slight bit of future-proofing.

	"github.com/mattn/go-sqlite3"
	"github.com/temelpa/timetravel/data"
	"github.com/temelpa/timetravel/entity"

	"os"
	"path/filepath"
)

// SQLiteRecordService is an SQLite-backed record service that
// persists data between runs of the server.
type SQLiteRecordService struct {
	db *sql.DB
}

type SQLiteRecordServiceSettings struct {
	// When the server is started, should the backing database
	// be purged?
	resetOnStart bool
}

func logError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func NewSQLiteRecordService(
	sqlDirectory string,
	settings SQLiteRecordServiceSettings,
) (SQLiteRecordService, error) {
	dbPath := filepath.Join(sqlDirectory, data.TIMETRAVEL_DB)
	if settings.resetOnStart {
		if err := os.RemoveAll(dbPath); err != nil {
			logError(err)
			return SQLiteRecordService{}, err
		}
	}

	// File permissions for the DB directory: R+W for owner, read only otherwise
	if err := os.MkdirAll(sqlDirectory, 0644); err != nil {
		logError(err)
		return SQLiteRecordService{}, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logError(err)
		return SQLiteRecordService{}, err
	}

	if err := db.Ping(); err != nil {
		logError(err)
		db.Close()
		return SQLiteRecordService{}, err
	}

	for _, sqlStatement := range []string{
		data.CREATE_RECORDS_TABLE,
		data.CREATE_RECORD_DELTAS_TABLE,
	} {
		statement, err := db.Prepare(sqlStatement)
		if err == nil {
			_, err = statement.Exec()
		}
		if err != nil {
			logError(err)
			db.Close()
			return SQLiteRecordService{}, err
		}
	}

	return SQLiteRecordService{db}, nil
}

func (s *SQLiteRecordService) GetRecord(
	ctx context.Context,
	id int,
) (entity.Record, error) {
	statement, err := s.db.Prepare(data.QUERY_RECORD)
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}
	row := statement.QueryRow(id)

	var jsonString string
	var recordVersion int
	if err = row.Scan(&id, &recordVersion, &jsonString); err != nil {
		logError(err)
		if err == sql.ErrNoRows {
			err = ErrRecordDoesNotExist
		}
		return entity.Record{}, err
	}

	if id == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	var data map[string]string
	if err = json.Unmarshal([]byte(jsonString), &data); err != nil {
		logError(err)
		return entity.Record{}, err
	}

	return entity.Record{ID: id, Version: recordVersion, Data: data}, nil
}

func (s *SQLiteRecordService) CreateRecord(
	ctx context.Context,
	record entity.Record,
) error {
	if record.ID <= 0 {
		return ErrRecordIDInvalid
	}

	statement, err := s.db.Prepare(data.INSERT_RECORD)
	if err != nil {
		logError(err)
		return err
	}

	jsonBytes, err := json.Marshal(record.Data)
	if err != nil {
		logError(err)
		return err
	}

	_, err = statement.Exec(record.ID, string(jsonBytes))
	if err != nil {
		logError(err)
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.Code == sqlite3.ErrConstraint {
			err = ErrRecordAlreadyExists
		}
		return err
	}

	return nil
}

func (s *SQLiteRecordService) UpdateRecord(
	ctx context.Context,
	id int,
	updates map[string]*string,
) (entity.Record, error) {
	entry, err := s.GetRecord(ctx, id)
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}

	updateInverse := entry.InverseUpdate(updates)
	if !entry.ApplyUpdate(updates) {
		return entry, nil
	}

	if statement, err := s.db.Prepare(data.INSERT_RECORD_DELTA); err == nil {
		if jsonBytes, err := json.Marshal(updateInverse); err == nil {
			_, err = statement.Exec(id, entry.Version, string(jsonBytes))
		}
	}
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}

	entry.Version += 1
	if statement, err := s.db.Prepare(data.UPDATE_RECORD); err == nil {
		if jsonBytes, err := json.Marshal(entry.Data); err == nil {
			_, err = statement.Exec(entry.Version, string(jsonBytes), id)
		}
	}
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}

	return entry.Copy(), nil
}

func (s *SQLiteRecordService) GetVersionedRecord(
	ctx context.Context,
	id int,
	version int,
) (entity.Record, error) {
	entry, err := s.GetRecord(ctx, id)
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}

	if version == 0 || version == entry.Version {
		return entry, nil
	}

	if entry.Version < version {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	statement, err := s.db.Prepare(data.QUERY_RECORD_DELTAS)
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}
	rows, err := statement.Query(version, id)
	if err != nil {
		logError(err)
		return entity.Record{}, err
	}

	for rows.Next() {
		var jsonString string
		var versionBeforeUpdate int

		if err = rows.Scan(&id, &versionBeforeUpdate, &jsonString); err != nil {
			logError(err)
			if err == sql.ErrNoRows {
				err = ErrRecordDoesNotExist
			}
			return entity.Record{}, err
		}

		var data map[string]*string
		if err = json.Unmarshal([]byte(jsonString), &data); err != nil {
			logError(err)
			return entity.Record{}, err
		}

		entry.ApplyUpdate(data)
		entry.Version = versionBeforeUpdate
	}

	return entry.Copy(), nil
}
