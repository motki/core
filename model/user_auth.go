package model

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func (m *Manager) AuthenticateUser(name, password string) (*User, string, error) {
	var emptyKey = ""
	db, err := m.pool.Open()
	if err != nil {
		return nil, emptyKey, err
	}
	defer m.pool.Release(db)
	u := &User{}
	var p []byte
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
	defer m.pool.Release(db)
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
	b, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(
		`INSERT INTO app.user_authorizations
			 (
			     user_id
		 	   , character_id
			   , role
			   , token
			 )
		       VALUES($1, $2, $3, $4)
			 ON CONFLICT
		       ON CONSTRAINT "user_authorizations_pkey"
			 DO UPDATE
			   SET character_id = EXCLUDED.character_id
			     , token = EXCLUDED.token`,
		u.UserID,
		characterID,
		int(r),
		b,
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
	defer m.pool.Release(db)
	a := &Authorization{}
	token := &oAuth2Token{}
	var b []byte
	ri := 0
	row := db.QueryRow(
		`SELECT
			  user_id
			, character_id
			, "role"
			, token
		    FROM app.user_authorizations
		    WHERE user_id = $1
			AND "role" = $2`,
		user.UserID,
		role)
	err = row.Scan(&a.UserID, &a.CharacterID, &ri, &b)
	a.Role = Role(ri)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("not authorized")
		}
		return nil, err
	}
	err = token.Scan(b)
	if err != nil {
		return nil, err
	}
	a.Token = (*oauth2.Token)(token)
	source, err := m.eveapi.TokenSource(a.Token)
	if err != nil {
		return nil, err
	}
	info, err := m.eveapi.Verify(source)
	if err != nil {
		return nil, err
	}
	if int(info.CharacterID) != a.CharacterID {
		return nil, errors.New("expected character IDs to match!")
	}
	a.source = source
	// Force retrieval of current char info from the API
	char, err := m.getCharacterFromAPI(a.CharacterID)
	if err != nil {
		return nil, err
	}
	a.CorporationID = char.CorporationID
	return a, nil
}

func (m *Manager) RemoveAuthorization(user *User, role Role) error {
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(
		`DELETE
			 FROM app.user_authorizations
			 WHERE user_id = $1 AND "role" = $2`,
		user.UserID,
		int(role))
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
	UserID        int
	CharacterID   int
	CorporationID int
	Role          Role
	Token         *oauth2.Token
	source        oauth2.TokenSource
}

func (a *Authorization) Context() context.Context {
	return newAuthContext(a)
}
