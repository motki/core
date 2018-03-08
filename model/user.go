package model

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	RoleAnon Role = iota
	RoleUser
	RoleMember
	RoleLogistics
	RoleDirector
	RoleAdmin
)

func (r Role) Value() (driver.Value, error) {
	return int64(r), nil
}

func (r *Role) Scan(src interface{}) error {
	i, ok := src.(int32)
	if !ok {
		return fmt.Errorf("invalid %t for role: %v", src, src)
	}
	switch Role(i) {
	case RoleAnon:
		*r = RoleAnon
	case RoleUser:
		*r = RoleUser
	case RoleMember:
		*r = RoleMember
	case RoleLogistics:
		*r = RoleLogistics
	case RoleDirector:
		*r = RoleDirector
	case RoleAdmin:
		*r = RoleAdmin
	default:
		return fmt.Errorf("invalid value for role: %v", i)
	}
	return nil
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrMissingField = errors.New("missing field")
)

type User struct {
	UserID int
	Name   string
	Email  string
}

func (m *Manager) NewUser(name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		return nil, ErrMissingField
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email address")
	}
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	row := db.QueryRow(`SELECT COUNT(1) FROM app.users WHERE username = $1 OR email = $2`, name, email)
	i := 0
	err = row.Scan(&i)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
	}
	if i != 0 {
		return nil, ErrUserExists
	}
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	row = db.QueryRow("INSERT INTO app.users (id, username, email, password) VALUES(DEFAULT, $1, $2, $3) RETURNING id", name, email, p)
	lid := 0
	err = row.Scan(&lid)
	if err != nil {
		return nil, err
	}
	if lid == 0 {
		return nil, errors.New("invalid last insert id")
	}
	return &User{
		UserID: int(lid),
		Name:   name,
		Email:  email,
	}, nil
}

func (m *Manager) CreateUserVerificationHash(user *User) ([]byte, error) {
	if user == nil {
		return nil, errors.New("cannot get hash for nil user")
	}
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	h := sha256.New()
	b := make([]byte, 24)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	h.Write(b)
	hash := h.Sum(nil)
	_, err = db.Exec(`INSERT INTO app.user_verifications (user_id, hash) VALUES($1, $2)
						ON CONFLICT ON CONSTRAINT "user_verifications_pkey" DO
						UPDATE SET hash = EXCLUDED.hash, expires_at = DEFAULT`, user.UserID, hash)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (m *Manager) VerifyUserEmail(email string, hash []byte) (bool, error) {
	if !strings.Contains(email, "@") {
		return false, errors.New("invalid email address")
	}
	db, err := m.pool.Open()
	if err != nil {
		return false, err
	}
	defer m.pool.Release(db)
	res, err := db.Exec("UPDATE app.users u SET verified = 1 FROM (SELECT user_id FROM app.user_verifications JOIN app.users ON user_id = id WHERE verified = 0 AND email = $1 AND hash = $2) uv WHERE uv.user_id = u.id", email, hash)
	if err != nil {
		return false, err
	}
	r := res.RowsAffected()
	return r == 1, nil
}
