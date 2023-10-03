package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID       uint
	UUID     uuid.UUID
	Username string
	Email    string
	Password string // this will be the hashed value, not plaintext
}

func New(username string, email string) *User {
	return &User{
		UUID:     uuid.New(),
		Username: username,
		Email:    email,
	}
}

func (u *User) WithPasswordHash(pw string) *User {
	u.Password = pw
	return u
}

type UserModel struct {
	DB *sql.DB
}

type Repository interface {
	GetUserByEmailOrUsername(email string, username string) (*User, error)
	CreateUser(user *User) error
  GetUserById(userId uint) (*User, error)
}

func (um *UserModel) GetUserByEmailOrUsername(
	email string,
	username string,
) (*User, error) {
	var user User
	res := um.DB.QueryRow(`
    SELECT
      id,
      uuid,
      email,
      username,
      password
    FROM users WHERE email = ? OR username = ?;
  `, email, username)
	err := res.Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	return &user, err
}

func (um *UserModel) CreateUser(user *User) error {
	_, err := um.DB.Exec(`
    INSERT INTO users (
      uuid,
      username,
      email,
      password
    ) VALUES (?, ?, ?, ?);
  `,
		user.UUID,
		user.Username,
		user.Email,
		user.Password,
	)

	return err
}

func (um *UserModel) GetUserById(userId uint) (*User, error) {
	var user User
	res := um.DB.QueryRow(`
    SELECT 
      id,
      uuid,
      username,
      email,
      password
    FROM users WHERE id = ?
  `, userId)

	err := res.Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	return &user, err
}
