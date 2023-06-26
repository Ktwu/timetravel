package service

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/temelpa/timetravel/entity"
)

// Test basic read + write functionality
func TestSanity(t *testing.T) {
	service := NewInMemoryRecordService()

	testEntity := entity.Record{
		ID: 42,
		Version: 1,
		Data: map[string]string{
			"hello": "world",
		},
	}

	testEntityUpdate := entity.Record{
		ID: testEntity.ID,
		Version: 2,
		Data: map[string]string{
			"hello": "world",
			"goodbye": "world",
		},
	}

	testEntityUpdate2 := entity.Record{
		ID: testEntity.ID,
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
		t.Errorf("Should have failed grabbing nonexistant entry")
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
		t.Errorf("Unable to fetch record for %v", testEntity)
	} else if !cmp.Equal(r, testEntity) {
		t.Errorf("Fetched entry %v not the same as %v", r, testEntity)
	}

	// Test adding data
	testValue := "world"
	if r, err := service.UpdateRecord(ctx, testEntity.ID, map[string]*string {
		"goodbye": &testValue,
	}); err != nil {
		t.Errorf("Unable to update record %v", testEntity)
	} else if testValue = "unittest"; !cmp.Equal(r, testEntityUpdate) {
		t.Errorf("Update entry %v not the same as %v", r, testEntityUpdate)
	}

	// Test mutating and removing data
	if r, err := service.UpdateRecord(ctx, testEntity.ID, map[string]*string {
		"hello": nil,
		"goodbye": &testValue,
	}); err != nil {
		t.Errorf("Unable to update record %v", testEntity)
	} else if !cmp.Equal(r, testEntityUpdate2) {
		t.Errorf("Update entry %v not the same as %v", r, testEntityUpdate2)
	}
}