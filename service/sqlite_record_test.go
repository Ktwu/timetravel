package service

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/temelpa/timetravel/entity"
)

// Test basic read + write functionality
func TestSanitySQL(t *testing.T) {
	service, err := NewSQLiteRecordService(
		"testdata",
		SQLiteRecordServiceSettings{ResetOnStart: true},
	)
	if err != nil {
		t.Errorf("Unable to create testing database, error %v", err)
	}
	defer func() {
		os.RemoveAll("testdata")
	}()

	testEntity := entity.Record{
		ID:      42,
		Version: 1,
		Data: map[string]string{
			"hello": "world",
		},
	}

	testEntityUpdate := entity.Record{
		ID:      testEntity.ID,
		Version: 2,
		Data: map[string]string{
			"hello":   "world",
			"goodbye": "world",
		},
	}

	testEntityUpdate2 := entity.Record{
		ID:      testEntity.ID,
		Version: 3,
		Data: map[string]string{
			"goodbye": "unittest",
		},
	}

	// TODO: read up on contexts to understand what should
	// actually go here
	ctx := context.Background()

	// Make sure the system is empty
	if _, err := service.GetRecord(ctx, testEntity.ID); err != ErrRecordDoesNotExist {
		t.Errorf("Should have failed grabbing nonexistant entry, got error %v", err)
	}
	if service.CreateRecord(ctx, testEntity) != nil {
		t.Errorf("Unable to create record for %v", testEntity)
	}

	// Attempting to create the same record again ought to fail
	// as well as not mutate the original data
	if service.CreateRecord(ctx, testEntityUpdate) != ErrRecordAlreadyExists {
		t.Errorf("Erroneously created second conflicting record %v", testEntityUpdate)
	}

	// Test fetching the data and data integrity
	if r, err := service.GetRecord(ctx, testEntity.ID); err != nil {
		t.Errorf("Unable to fetch record for %v, got error %v", testEntity, err)
	} else if !cmp.Equal(r, testEntity) {
		t.Errorf("Fetched entry %v not the same as %v", r, testEntity)
	}

	// Test adding data
	testValue := "world"
	if r, err := service.UpdateRecord(ctx, testEntity.ID, map[string]*string{
		"goodbye": &testValue,
	}); err != nil {
		t.Errorf("Unable to update record %v, got error %v", testEntity, err)
	} else if !cmp.Equal(r, testEntityUpdate) {
		t.Errorf("Update entry %v not the same as %v", r, testEntityUpdate)
	}

	// Test mutating and removing data
	testValue = "unittest"
	if r, err := service.UpdateRecord(ctx, testEntity.ID, map[string]*string{
		"hello":   nil,
		"goodbye": &testValue,
	}); err != nil {
		t.Errorf("Unable to update record %v, got err %v", testEntity, err)
	} else if !cmp.Equal(r, testEntityUpdate2) {
		t.Errorf("Update entry %v not the same as %v", r, testEntityUpdate2)
	}

	if _, err := service.GetVersionedRecord(ctx, testEntity.ID, 4); err != ErrRecordDoesNotExist {
		t.Errorf("Should have failed grabbing entry for nonexistant version, error %v", err)
	}

	if r, err := service.GetVersionedRecord(ctx, testEntity.ID, 0); err != nil {
		t.Errorf("Error grabbing newest version record, error %v", err)
	} else if !cmp.Equal(r, testEntityUpdate2) {
		t.Errorf("Failed to grab newest version, got %v, expected %v", r, testEntityUpdate2)
	}

	if r, err := service.GetVersionedRecord(ctx, testEntity.ID, 3); err != nil {
		t.Errorf("Error grabbing newest version record, error %v", err)
	} else if !cmp.Equal(r, testEntityUpdate2) {
		t.Errorf("Failed to grab newest version, got %v, expected %v", r, testEntityUpdate2)
	}

	if r, err := service.GetVersionedRecord(ctx, testEntity.ID, 2); err != nil {
		t.Errorf("Error grabbing versioned record, error %v", err)
	} else if !cmp.Equal(r, testEntityUpdate) {
		t.Errorf("Failed to grab version, got %v, expected %v", r, testEntityUpdate)
	}

	if r, err := service.GetVersionedRecord(ctx, testEntity.ID, 1); err != nil {
		t.Errorf("Error grabbing versioned record, error %v", err)
	} else if !cmp.Equal(r, testEntity) {
		t.Errorf("Failed to grab version, got %v, expected %v", r, testEntity)
	}

	if rs, err := service.GetAllRecordVersions(ctx, testEntity.ID); err != nil {
		t.Errorf("Error grabbing versions of record, error %v", err)
	} else if expected := []entity.Record{testEntity, testEntityUpdate, testEntityUpdate2}; !cmp.Equal(rs, expected) {
		t.Errorf("Failed to grab all versions, got %v, expected %v", rs, expected)
	}

	service, err = NewSQLiteRecordService(
		"testdata",
		SQLiteRecordServiceSettings{ResetOnStart: false},
	)
	if err != nil {
		t.Errorf("Unable to create testing database, got err %v", err)
	}

	// Test fetching the data and data integrity
	if r, err := service.GetRecord(ctx, testEntity.ID); err != nil {
		t.Errorf("Unable to fetch record for %v, got err %v", testEntityUpdate2, err)
	} else if !cmp.Equal(r, testEntityUpdate2) {
		t.Errorf("Fetched entry %v not the same as %v", r, testEntityUpdate2)
	}
}

// Test creating an inverse update on a map for basic add, delete, and mutate ops
func TestUpdateInverse(t *testing.T) {
	basicMap := map[string]string{
		"hello": "world",
		"basic": "data",
	}

	worldText := "world"
	update := map[string]*string{
		"goodbye": &worldText,
		"hello":   nil,
	}

	expectedInverse := map[string]*string{
		"hello":   &worldText,
		"goodbye": nil,
	}

	testEntry := entity.Record{Data: basicMap}
	inverse := testEntry.InverseUpdate(update)
	if !cmp.Equal(inverse, expectedInverse) {
		t.Errorf("Expected %v, got %v", expectedInverse, inverse)
	}
}
