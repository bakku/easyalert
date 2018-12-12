package postgres_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bakku/easyalert"

	"github.com/bakku/easyalert/postgres"
	"github.com/stretchr/testify/require"
)

func createUserWithId(t *testing.T, db *sql.DB, id uint) {
	_, err := db.Exec(`
		INSERT INTO users(id, email, password_digest,
			token, created_at, updated_at)
		VALUES ($1, 'test@mail.com', '1234',
			'1234', NOW(), NOW())
	`, id)
	require.Nil(t, err)
}

func createAlert(t *testing.T, db *sql.DB, id uint, subject string, status uint, sentAt *time.Time, userID uint) {
	_, err := db.Exec(`
		INSERT INTO alerts(id, subject, status,
			sent_at, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		id, subject, status, sentAt, userID)

	require.Nil(t, err)
}

func TestFindAlert_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	sentAt := time.Now().AddDate(-1, 0, 0)

	createAlert(t, db, 1, "Testing", 0, &sentAt, 1)

	repo := postgres.AlertRepository{DB: db}

	alert, err := repo.FindAlert("WHERE id = $1", 1)
	require.Nil(t, err)

	var defaultTime time.Time

	require.Equal(t, uint(1), alert.ID)
	require.Equal(t, "Testing", alert.Subject)
	require.Equal(t, "pending", alert.HumanStatus())
	require.NotNil(t, alert.SentAt)
	require.NotEqual(t, defaultTime, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, defaultTime, alert.CreatedAt)
	require.NotEqual(t, defaultTime, alert.UpdatedAt)
}

func TestFindAlert_NullSentAt(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	createAlert(t, db, 1, "Testing", 0, nil, 1)

	repo := postgres.AlertRepository{DB: db}

	alert, err := repo.FindAlert("WHERE id = $1", 1)
	require.Nil(t, err)

	var defaultTime time.Time

	require.Equal(t, uint(1), alert.ID)
	require.Equal(t, "Testing", alert.Subject)
	require.Equal(t, "pending", alert.HumanStatus())
	require.Nil(t, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, defaultTime, alert.CreatedAt)
	require.NotEqual(t, defaultTime, alert.UpdatedAt)
}

func TestFindAlert_NotExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	sentAt := time.Now().AddDate(-1, 0, 0)

	createAlert(t, db, 1, "Testing", 0, &sentAt, 1)

	repo := postgres.AlertRepository{DB: db}

	_, err = repo.FindAlert("WHERE id = $1", 2)
	require.Equal(t, easyalert.ErrRecordDoesNotExist, err)
}

func TestFindAlerts_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	sentAt := time.Now().AddDate(-1, 0, 0)

	createAlert(t, db, 1, "Testing #1", 0, &sentAt, 1)
	createAlert(t, db, 2, "Testing #2", 1, &sentAt, 1)

	repo := postgres.AlertRepository{DB: db}

	alerts, err := repo.FindAlerts("WHERE user_id = $1", 1)
	require.Nil(t, err)

	require.Len(t, alerts, 2)

	var defaultTime time.Time

	alert := alerts[0]

	require.Equal(t, uint(1), alert.ID)
	require.Equal(t, "Testing #1", alert.Subject)
	require.Equal(t, "pending", alert.HumanStatus())
	require.NotNil(t, alert.SentAt)
	require.NotEqual(t, defaultTime, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, defaultTime, alert.CreatedAt)
	require.NotEqual(t, defaultTime, alert.UpdatedAt)

	alert = alerts[1]

	require.Equal(t, uint(2), alert.ID)
	require.Equal(t, "Testing #2", alert.Subject)
	require.Equal(t, "sent", alert.HumanStatus())
	require.NotNil(t, alert.SentAt)
	require.NotEqual(t, defaultTime, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, defaultTime, alert.CreatedAt)
	require.NotEqual(t, defaultTime, alert.UpdatedAt)
}

func TestFindAlerts_NonExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	sentAt := time.Now().AddDate(-1, 0, 0)

	createAlert(t, db, 1, "Testing #1", 0, &sentAt, 1)
	createAlert(t, db, 2, "Testing #2", 1, &sentAt, 1)

	repo := postgres.AlertRepository{DB: db}

	alerts, err := repo.FindAlerts("WHERE user_id = $1", 2)
	require.Nil(t, err)

	require.Len(t, alerts, 0)
}

func TestCreateAlert_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	repo := postgres.AlertRepository{DB: db}

	alert := easyalert.Alert{Subject: "Testing", Status: 0, SentAt: nil, UserID: 1}

	alert, err = repo.CreateAlert(alert)
	require.Nil(t, err)

	var defaultTime time.Time

	require.NotEqual(t, uint(0), alert.ID)
	require.Equal(t, "Testing", alert.Subject)
	require.Equal(t, "pending", alert.HumanStatus())
	require.Nil(t, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, defaultTime, alert.CreatedAt)
	require.NotEqual(t, defaultTime, alert.UpdatedAt)

	var exists bool

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM alerts WHERE id = $1)", alert.ID).Scan(&exists)
	require.Nil(t, err)
	require.True(t, exists)
}

func TestUpdateAlert_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	row := db.QueryRow(`
		INSERT INTO alerts(id, subject, status,
			sent_at, user_id, created_at, updated_at)
		VALUES (1, 'Test', 0, NULL, 1, NOW(), NOW())
		RETURNING updated_at
	`)

	var oldUpdatedAt time.Time
	err = row.Scan(&oldUpdatedAt)
	require.Nil(t, err)

	sentAt := time.Now()

	alert := easyalert.Alert{
		ID:      1,
		Subject: "Test",
		Status:  1,
		SentAt:  &sentAt,
		UserID:  1,
	}

	repo := postgres.AlertRepository{DB: db}

	alert, err = repo.UpdateAlert(alert)
	require.Nil(t, err)

	var defaultTime time.Time
	require.Equal(t, uint(1), alert.ID)
	require.Equal(t, "Test", alert.Subject)
	require.Equal(t, "sent", alert.HumanStatus())
	require.NotEqual(t, defaultTime, alert.SentAt)
	require.Equal(t, uint(1), alert.UserID)
	require.NotEqual(t, oldUpdatedAt, alert.UpdatedAt)
}

func TestUpdateAlert_NotExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	alert := easyalert.Alert{
		ID:      1,
		Subject: "Test",
		Status:  1,
		SentAt:  nil,
		UserID:  1,
	}

	repo := postgres.AlertRepository{DB: db}

	alert, err = repo.UpdateAlert(alert)
	require.Equal(t, easyalert.ErrRecordDoesNotExist, err)
}

func TestDeleteAlert_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	createUserWithId(t, db, 1)

	_, err = db.Exec(`
		INSERT INTO alerts(id, subject, status,
			sent_at, user_id, created_at, updated_at)
		VALUES (1, 'Test', 0, NULL, 1, NOW(), NOW())
	`)
	require.Nil(t, err)

	repo := postgres.AlertRepository{DB: db}

	alert := easyalert.Alert{ID: 1}

	err = repo.DeleteAlert(alert)
	require.Nil(t, err)

	row := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM alerts WHERE id = 1)
	`)

	var exists bool
	err = row.Scan(&exists)
	require.Nil(t, err)

	require.False(t, exists)
}
