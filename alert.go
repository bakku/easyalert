package easyalert

import "time"

type AlertRepository interface {
	FindAlert(query string, params ...interface{}) (Alert, error)
	FindAlerts(query string, params ...interface{}) ([]Alert, error)
	CreateAlert(alert Alert) (Alert, error)
	UpdateAlert(alert Alert) (Alert, error)
	DeleteAlert(alert Alert) error
}

const (
	AlertStatusPending = iota
	AlertStatusSent
	AlertStatusFailed
)

type Alert struct {
	ID        uint
	Subject   string
	Status    uint
	SentAt    *time.Time
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Alert) HumanStatus() string {
	switch a.Status {
	case AlertStatusPending:
		return "pending"
	case AlertStatusSent:
		return "sent"
	case AlertStatusFailed:
		return "failed"
	default:
		return "invalid status"
	}
}
