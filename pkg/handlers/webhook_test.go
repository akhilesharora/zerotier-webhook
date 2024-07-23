package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zerotier/ztchooks"
	"zerotier-webhook/pkg/models"
)

type MockDatabase struct {
	events []models.Event
}

func (m *MockDatabase) CreateEvent(event *models.Event) error {
	m.events = append(m.events, *event)
	return nil
}

func (m *MockDatabase) GetEvents(networkID, memberID, userID string) ([]models.Event, error) {
	var result []models.Event
	for _, event := range m.events {
		if (networkID == "" || event.NetworkID == networkID) &&
			(memberID == "" || event.MemberID == memberID) &&
			(userID == "" || event.UserID == userID) {
			result = append(result, event)
		}
	}
	return result, nil
}

func TestHandleWebhook(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := NewWebhookHandler(mockDB)

	payload := ztchooks.MemberConfigChanged{
		HookBase: ztchooks.HookBase{
			HookID:   "test-hook-id",
			OrgID:    "test-org-id",
			HookType: ztchooks.MEMBER_CONFIG_CHANGED,
		},
		NetworkID: "test-network-id",
		MemberID:  "test-member-id",
		UserID:    "test-user-id",
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer(payloadBytes))
	rr := httptest.NewRecorder()

	handler.HandleWebhook(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if len(mockDB.events) != 1 {
		t.Errorf("Expected 1 event to be created, got %d", len(mockDB.events))
	}

	createdEvent := mockDB.events[0]
	if createdEvent.HookID != payload.HookID {
		t.Errorf("Expected HookID %s, got %s", payload.HookID, createdEvent.HookID)
	}
}

func TestHandleSearch(t *testing.T) {
	mockDB := &MockDatabase{
		events: []models.Event{
			{NetworkID: "network1", MemberID: "member1", UserID: "user1"},
			{NetworkID: "network1", MemberID: "member2", UserID: "user2"},
		},
	}
	handler := NewWebhookHandler(mockDB)

	req, _ := http.NewRequest("GET", "/search?network_id=network1", nil)
	rr := httptest.NewRecorder()

	handler.HandleSearch(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result []models.Event
	json.Unmarshal(rr.Body.Bytes(), &result)

	if len(result) != 2 {
		t.Errorf("Expected 2 events, got %d", len(result))
	}
}
