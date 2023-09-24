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

func CreateTable(db *sql.DB) error {
  createStr := `
    CREATE TABLE IF NOT EXISTS sessions (
      id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
      uuid TEXT NOT NULL UNIQUE,
      init_time DATETIME,
      expiration DATETIME NOT NULL,
      user_id INTEGER NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users(id)
    );
  `
  _, err := db.Exec(createStr)
  return err
}

func New() *Session {
  now := time.Now()
  return &Session{
    UUID: uuid.New(),
    InitTime: now,
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
}

func (sr *SessionModel) CreateSession(s *Session) error {
  return nil
}
