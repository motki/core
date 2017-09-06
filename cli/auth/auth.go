package auth

import (
	"context"

	"github.com/antihax/goesi"
	"github.com/pkg/errors"

	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/model"
)

var ErrBadCredentials = errors.New("cli/auth: invalid username or password")
var ErrNotAuthenticated = errors.New("cli/auth: not authenticated")

type sessionKey *string

type Session struct {
	model *model.Manager
	api   *eveapi.EveAPI

	sessionKey sessionKey
}

func NewSession(model *model.Manager, api *eveapi.EveAPI) *Session {
	return &Session{model, api, nil}
}

func (s *Session) Authenticate(user, password string) (*model.User, error) {
	u, key, err := s.model.AuthenticateUser(user, password)
	if err != nil {
		return nil, ErrBadCredentials
	}
	s.sessionKey = &key
	return u, nil
}

func (s *Session) User() (*model.User, error) {
	if s.sessionKey == nil {
		return nil, ErrNotAuthenticated
	}
	return s.model.GetUserBySessionKey(*s.sessionKey)
}

// AuthorizedContext returns a context containing credentials and the corresponding characterID.
func (s *Session) AuthorizedContext(role model.Role) (context.Context, int, error) {
	user, err := s.User()
	if err != nil {
		return nil, 0, err
	}
	a, err := s.model.GetAuthorization(user, role)
	if err != nil {
		return nil, 0, err
	}
	source, err := s.api.TokenSource((*goesi.CRESTToken)(a.Token))
	if err != nil {
		return nil, 0, err
	}
	info, err := s.api.Verify(source)
	if err != nil {
		return nil, 0, err
	}
	t, err := source.Token()
	if err != nil {
		return nil, 0, err
	}
	if err = s.model.SaveAuthorization(user, role, int(info.CharacterID), t); err != nil {
		return nil, 0, err
	}
	return context.WithValue(context.Background(), goesi.ContextOAuth2, source), int(info.CharacterID), nil
}
