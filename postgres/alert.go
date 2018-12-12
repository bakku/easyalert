package postgres

import (
	"database/sql"

	"github.com/bakku/easyalert"
)

// AlertRepository is a postgres implementation of the AlertRepository interface
type AlertRepository struct {
	DB *sql.DB
}

// FindAlert fetches an alert using the query passed as a string and returns it. If the alert does not exist it will return easyalert.ErrRecordDoesNotExist.
func (repo AlertRepository) FindAlert(query string, params ...interface{}) (easyalert.Alert, error) {
	var alert easyalert.Alert

	baseQuery := `
		SELECT id, subject, status, sent_at, user_id, created_at, updated_at
		FROM alerts
	`

	row := repo.DB.QueryRow(baseQuery+query, params...)

	err := row.Scan(&alert.ID, &alert.Subject, &alert.Status, &alert.SentAt, &alert.UserID, &alert.CreatedAt, &alert.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return easyalert.Alert{}, easyalert.ErrRecordDoesNotExist
		}

		return easyalert.Alert{}, err
	}

	return alert, nil
}

// FindAlerts fetches all alerts based on the query and returns them.
func (repo AlertRepository) FindAlerts(query string, params ...interface{}) ([]easyalert.Alert, error) {
	var alerts []easyalert.Alert

	baseQuery := `
		SELECT id, subject, status, sent_at, user_id, created_at, updated_at
		FROM alerts
	`

	rows, err := repo.DB.Query(baseQuery+query, params...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a easyalert.Alert

		if err := rows.Scan(&a.ID, &a.Subject, &a.Status, &a.SentAt, &a.UserID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}

// CreateAlert creates a new alert in the Postgres database and returns it with ID and created_at/updated_at filled.
func (repo AlertRepository) CreateAlert(alert easyalert.Alert) (easyalert.Alert, error) {
	row := repo.DB.QueryRow(`
		INSERT INTO alerts(subject, status, sent_at, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`, alert.Subject, alert.Status, alert.SentAt, alert.UserID)

	err := row.Scan(&alert.ID, &alert.CreatedAt, &alert.UpdatedAt)

	if err != nil {
		return easyalert.Alert{}, err
	}

	return alert, nil
}

// UpdateAlert updates an existing alert in the Postgres database and returns it with updated_at updated.
func (repo AlertRepository) UpdateAlert(alert easyalert.Alert) (easyalert.Alert, error) {
	row := repo.DB.QueryRow(`
			UPDATE alerts
			SET subject = $1, status = $2, sent_at = $3,
			updated_at = NOW()
			WHERE alerts.id = $4
			RETURNING updated_at
		`, alert.Subject, alert.Status, alert.SentAt, alert.ID)

	err := row.Scan(&alert.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return easyalert.Alert{}, easyalert.ErrRecordDoesNotExist
		}

		return easyalert.Alert{}, err
	}

	return alert, nil
}

// DeleteAlert deletes the alert given as a parameter by using the ID.
func (repo AlertRepository) DeleteAlert(alert easyalert.Alert) error {
	_, err := repo.DB.Exec(`
			DELETE FROM alerts
			WHERE id = $1
		`, alert.ID)

	return err
}
