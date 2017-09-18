package model

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
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
	defer db.Close()
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
	defer db.Close()
	res, err := db.Exec("UPDATE app.users u SET verified = 1 FROM (SELECT user_id FROM app.user_verifications JOIN app.users ON user_id = id WHERE verified = 0 AND email = $1 AND hash = $2) uv WHERE uv.user_id = u.id", email, hash)
	if err != nil {
		return false, err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return r == 1, nil
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
	row := db.QueryRow(`SELECT id, username, email, password FROM app.users WHERE username = $1 AND verified = 1 AND disabled <> 1`, name)
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

func (m *Manager) SaveAuthorization(u *User, r Role, characterID int, tok *oauth2.Token) error {
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(
		`INSERT INTO app.user_authorizations (user_id, character_id, role, token)
		  	   VALUES($1, $2, $3, $4)
			   ON CONFLICT ON CONSTRAINT "user_authorizations_pkey"  DO
				UPDATE SET character_id = EXCLUDED.character_id,
				  token = EXCLUDED.token`,
		u.UserID,
		characterID,
		r,
		(*oAuth2Token)(tok),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetAuthorization(user *User, role Role) (*Authorization, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	a := &Authorization{}
	token := &oAuth2Token{}
	row := db.QueryRow(`SELECT user_id, character_id, "role", token
					    FROM app.user_authorizations
					    WHERE user_id = $1
						AND "role" = $2`, user.UserID, role)
	err = row.Scan(&a.UserID, &a.CharacterID, &a.Role, &token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not authorized")
		}
		return nil, err
	}
	a.Token = (*oauth2.Token)(token)
	return a, nil
}

func (m *Manager) RemoveAuthorization(user *User, role Role) error {
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`DELETE FROM app.user_authorizations WHERE user_id = $1 AND "role" = $2`, user.UserID, role)
	return err
}

type oAuth2Token oauth2.Token

func (r *oAuth2Token) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *oAuth2Token) Scan(src interface{}) error {
	s, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid value for token: %v", src)
	}
	return json.Unmarshal(s, &r)
}

type Authorization struct {
	UserID      int
	CharacterID int
	Role        Role
	Token       *oauth2.Token
}
