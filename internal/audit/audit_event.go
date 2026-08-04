package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditEvent представляет событие аудита.
type AuditEvent struct {
	Timestamp int64     `json:"ts"`
	Action    Action    `json:"action"`
	UserID    uuid.UUID `json:"user_id"`
	URL       string    `json:"url"`
}

// NewEvent создает новое событие аудита.
func NewEvent(action Action, userID uuid.UUID, url string) AuditEvent {
	return AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    action,
		UserID:    userID,
		URL:       url,
	}
}

func (e AuditEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
