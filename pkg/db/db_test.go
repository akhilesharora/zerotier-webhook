package db

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"zerotier-webhook/pkg/models"
)

func NewTestDatabase() (*GormDatabase, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Event{}); err != nil {
		return nil, err
	}

	return &GormDatabase{db: db}, nil
}

func clearDatabase(db *gorm.DB) error {
	return db.Exec("DELETE FROM events").Error
}

func TestCreateEvent(t *testing.T) {
	db, err := NewTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := clearDatabase(db.db); err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}

	event := &models.Event{
		ID:        uuid.New().String(),
		HookID:    "test_hook_id",
		NetworkID: "test_network",
		MemberID:  "test_member",
		UserID:    "test_user",
	}

	if err := db.CreateEvent(event); err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	var storedEvent models.Event
	if err := db.db.First(&storedEvent, "id = ?", event.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve stored event: %v", err)
	}

	if storedEvent.HookID != event.HookID || storedEvent.NetworkID != event.NetworkID {
		t.Errorf("Stored event does not match created event")
	}
}

func TestGetEventsWithFilters(t *testing.T) {
	db, err := NewTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := clearDatabase(db.db); err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}

	event1 := &models.Event{
		ID:        uuid.New().String(),
		HookID:    "test_hook_id_1",
		NetworkID: "test_network",
		MemberID:  "test_member1",
		UserID:    "test_user1",
	}

	event2 := &models.Event{
		ID:        uuid.New().String(),
		HookID:    "test_hook_id_2",
		NetworkID: "test_network",
		MemberID:  "test_member2",
		UserID:    "test_user2",
	}

	if err := db.CreateEvent(event1); err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	if err := db.CreateEvent(event2); err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	tests := []struct {
		networkID string
		memberID  string
		userID    string
		expected  int
	}{
		{"test_network", "", "", 2},
		{"test_network", "test_member1", "", 1},
		{"test_network", "", "test_user2", 1},
	}

	for _, tt := range tests {
		events, err := db.GetEvents(tt.networkID, tt.memberID, tt.userID)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		if len(events) != tt.expected {
			t.Errorf("Expected %d events, got %d", tt.expected, len(events))
		}
	}
}
