package saga

const (
	UserBlockedSubject     = "saga.user.blocked"
	UserBlockFailedSubject = "saga.user.block.failed"
)

type UserBlockedEvent struct {
	UserID string `json:"user_id"`
}

type UserBlockFailedEvent struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}
