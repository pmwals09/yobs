package session

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pmwals09/yobs/internal/models/user"
)

type Session struct {
	ID         uint
	UUID       uuid.UUID
	InitTime   time.Time
	Expiration time.Time
	UserID     uint
}

func New() *Session {
	now := time.Now()
	return &Session{
		UUID:       uuid.New(),
		InitTime:   now,
		Expiration: now.Add(time.Minute * 30),
	}
}

func (s *Session) WithExpiration(t time.Time) *Session {
	s.Expiration = t
	return s
}

func (s *Session) WithUser(u *user.User) *Session {
	s.UserID = u.ID
	return s
}

type SessionModel struct {
	DB *sql.DB
}

type Repository interface {
	CreateSession(s *Session) error
	GetSessionByUUID(uuid uuid.UUID) (*Session, error)
	UpdateSession(s *Session) error
	DeleteSessionByUUID(uuid uuid.UUID) error
}

func (sm *SessionModel) CreateSession(s *Session) error {
	_, err := sm.DB.Exec(`
    INSERT INTO sessions (
      uuid,
      init_time,
      expiration,
      user_id
    ) VALUES ( ?, ?, ?, ? );
  `,
		s.UUID,
		s.InitTime,
		s.Expiration,
		s.UserID,
	)

	return err
}

func (sm *SessionModel) GetSessionByUUID(uuid uuid.UUID) (*Session, error) {
	row := sm.DB.QueryRow(`
    SELECT 
      id,
      uuid,
      init_time,
      expiration,
      user_id
    FROM sessions WHERE uuid = ?;
  `,
		uuid)
	var session Session
	err := row.Scan(
		&session.ID,
		&session.UUID,
		&session.InitTime,
		&session.Expiration,
		&session.UserID,
	)

	return &session, err
}

func (sm *SessionModel) UpdateSession(s *Session) error {
	// NOTE: the only thing that should change after a session is created is the
	// expiration time
	_, err := sm.DB.Exec(`
    UPDATE sessions
    SET expiration = ?
    WHERE id = ?
  `, s.Expiration, s.ID)

	return err
}

func (sm *SessionModel) DeleteSessionByUUID(uuid uuid.UUID) error {
	_, err := sm.DB.Exec(`DELETE FROM sessions WHERE uuid = ?`, uuid)
	return err
}
