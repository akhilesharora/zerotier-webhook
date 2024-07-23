package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/zerotier/ztchooks"
	"zerotier-webhook/pkg/db"
	"zerotier-webhook/pkg/models"
)

var ErrUnhandledHook = errors.New("unhandled hook type")

type WebhookHandler struct {
	db db.Database
}

func NewWebhookHandler(database db.Database) *WebhookHandler {
	return &WebhookHandler{db: database}
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Skipping Signature verification for now
	// Need more time to understand it

	if err := h.processPayload(body); err != nil {
		if errors.Is(err, ErrUnhandledHook) {
			http.Error(w, "Unhandled hook type", http.StatusOK)
		} else {
			http.Error(w, "Error processing payload", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) processPayload(payload []byte) error {
	hType, err := ztchooks.GetHookType(payload)
	if err != nil {
		return err
	}

	event := models.NewEvent()

	switch hType {
	case ztchooks.MEMBER_CONFIG_CHANGED:
		var hook ztchooks.MemberConfigChanged
		if err := json.Unmarshal(payload, &hook); err != nil {
			return err
		}
		event.HookID = hook.HookID
		event.OrgID = hook.OrgID
		event.HookType = hook.HookType
		event.NetworkID = hook.NetworkID
		event.MemberID = hook.MemberID
		event.UserID = hook.UserID
		event.UserEmail = hook.UserEmail

		// Convert OldConfig and NewConfig to JSON strings
		oldConfig, err := json.Marshal(hook.OldConfig)
		if err != nil {
			return err
		}
		event.OldConfig = string(oldConfig)

		newConfig, err := json.Marshal(hook.NewConfig)
		if err != nil {
			return err
		}
		event.NewConfig = string(newConfig)

	// More cases can be added later

	default:
		return ErrUnhandledHook
	}

	return h.db.CreateEvent(event)
}

func (h *WebhookHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	networkID := query.Get("network_id")
	memberID := query.Get("member_id")
	userID := query.Get("user_id")

	events, err := h.db.GetEvents(networkID, memberID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}
