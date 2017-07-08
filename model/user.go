package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type Role int

const (
	RoleAnon Role = iota
	RoleUser
	RoleMember
	RoleDirector
	RoleAdmin
)

func (r Role) Value() (driver.Value, error) {
	return int(r), nil
}

func (r *Role) Scan(src interface{}) error {
	i, ok := src.(int)
	if !ok {
		return fmt.Errorf("invalid value for role: %v", src)
	}
	switch Role(i) {
	case RoleAnon:
		*r = RoleAnon
	case RoleUser:
		*r = RoleUser
	case RoleMember:
		*r = RoleMember
	case RoleDirector:
		*r = RoleDirector
	case RoleAdmin:
		*r = RoleAdmin
	default:
		return fmt.Errorf("invalid value for role: %v", i)
	}
	return nil
}

var ErrUserExists = errors.New("user already exists")

type User struct {
	UserID int
	Name   string
	Email  string
}

func (m *Manager) NewUser(name, email, password string) (*User, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow(`SELECT COUNT(1) FROM app.users WHERE username = $1 OR email = $2`, name, email)
	i := 0
	err = row.Scan(&i)
	if err != nil {
		if err != sql.ErrNoRows {
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
	res, err := db.Exec("INSERT INTO app.users (id, username, email, password) VALUES(DEFAULT, $1, $2, $3)", name, email, p)
	if err != nil {
		return nil, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		UserID: int(lid),
		Name:   name,
		Email:  email,
	}, nil
}

func (m *Manager) AuthenticateUser(name, password string) (*User, string, error) {
	var emptyKey = ""
	db, err := m.pool.Open()
	if err != nil {
		return nil, emptyKey, err
	}
	defer db.Close()
	u := &User{}
	p := []byte{}
	row := db.QueryRow(`SELECT id, username, email, password FROM app.users WHERE username = $1`, name)
	err = row.Scan(&u.UserID, &u.Name, &u.Email, &p)
	if err != nil {
		return nil, emptyKey, err
	}
	err = bcrypt.CompareHashAndPassword(p, []byte(password))
	if err != nil {
		return nil, emptyKey, err
	}
	bk := make([]byte, 32)
	n, err := rand.Read(bk)
	if err != nil || n != len(bk) {
		return nil, emptyKey, errors.New("unable to securely generate user session key")
	}
	key := base64.RawURLEncoding.EncodeToString(bk)
	_, err = db.Exec(`INSERT INTO app.user_sessions (user_id, key) VALUES($1, $2)
						ON CONFLICT ON CONSTRAINT "user_sessions_pkey" DO
						UPDATE SET key = EXCLUDED.key,
							     last_seen_at = DEFAULT,
							     created_at = DEFAULT`, u.UserID, key)
	if err != nil {
		return nil, emptyKey, err
	}
	return u, key, nil
}

func (m *Manager) GetUserBySessionKey(key string) (*User, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	u := &User{}
	row := db.QueryRow(`UPDATE app.user_sessions us
					    SET last_seen_at = NOW()
					    FROM (
					    	SELECT u.id, u.username, u.email
					    	FROM app.users u
					    	  JOIN app.user_sessions s ON s.user_id = u.id
					    	WHERE s.key = $1
					    	  AND s.last_seen_at >= NOW() - INTERVAL '30 mins'
					    ) u
					    WHERE us.user_id = u.id
					    RETURNING u.id, u.username, u.email`, key)
	err = row.Scan(&u.UserID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *Manager) AuthorizeUser(user *User, role Role) (context.Context, error) {
	// Check database
	return nil, nil
}

type oAuth2Token oauth2.Token

func (r *oAuth2Token) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *oAuth2Token) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("invalid value for token: %v", src)
	}
	return json.Unmarshal([]byte(s), &r)
}

type Authorization struct {
	UserID      int
	CharacterID int
	Role        Role

	token *oAuth2Token
}
