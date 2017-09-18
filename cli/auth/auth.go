package auth

import (
	"context"

	"github.com/antihax/goesi"
	"github.com/pkg/errors"

	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/log"
	"github.com/motki/motkid/model"
	"github.com/motki/motkid/model/client"
)

var ErrBadCredentials = errors.New("cli/auth: invalid username or password")
var ErrNotAuthenticated = client.ErrNotAuthenticated

type sessionKey *string

type Session struct {
	client client.Client
	api    *eveapi.EveAPI
	model  *model.Manager
	logger log.Logger

	sessionKey sessionKey
}

func NewSession(cl client.Client, mdl *model.Manager, api *eveapi.EveAPI, l log.Logger) *Session {
	return &Session{
		client: cl,
		api:    api,
		model:  mdl,
		logger: l,
	}
}

func (s *Session) Authenticate(user, password string) (*model.User, error) {
	key, err := s.client.Authenticate(user, password)
	if err != nil {
		s.logger.Warnf("unable to authenticate: %s", err.Error())
		return nil, ErrBadCredentials
	}
	s.sessionKey = &key
	return nil, nil
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
