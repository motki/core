package model

import (
	"github.com/antihax/goesi"
	"golang.org/x/net/context"
)

type authContext interface {
	context.Context

	UserID() int

	Role() Role
	CharacterID() int
	CorporationID() int
}

type authContextImpl struct {
	context.Context
	a *Authorization
}

func (ctx authContextImpl) CorporationID() int {
	return ctx.a.CorporationID
}

func (ctx authContextImpl) CharacterID() int {
	return ctx.a.CharacterID
}

func (ctx authContextImpl) UserID() int {
	return ctx.a.UserID
}

func (ctx authContextImpl) Role() Role {
	return ctx.a.Role
}

func newAuthContext(a *Authorization) context.Context {
	return authContextImpl{
		Context: context.WithValue(context.Background(), goesi.ContextOAuth2, a.source),
		a:       a,
	}
}

func authContextFromContext(ctx context.Context) (authContext, bool) {
	if v, ok := ctx.(authContext); ok {
		return v, true
	}
	return nil, false
}
