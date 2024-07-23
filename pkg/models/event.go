package models

import (
	"github.com/google/uuid"
	"github.com/zerotier/ztchooks"
)

type Event struct {
	ID        string            `json:"id" gorm:"primaryKey"`
	HookID    string            `json:"hook_id"`
	OrgID     string            `json:"org_id"`
	HookType  ztchooks.HookType `json:"hook_type"`
	NetworkID string            `json:"network_id"`
	MemberID  string            `json:"member_id"`
	UserID    string            `json:"user_id"`
	UserEmail string            `json:"user_email"`
	OldConfig string            `json:"old_config"`
	NewConfig string            `json:"new_config"`
}

func NewEvent() *Event {
	return &Event{
		ID: uuid.New().String(),
	}
}
