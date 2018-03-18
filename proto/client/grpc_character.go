package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// CharacterClient retrieves character, corporation, and alliance information.
type CharacterClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// CharacterForRole returns the current session's associated character for the given role.
func (c *CharacterClient) CharacterForRole(role model.Role) (*model.Character, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	var r proto.Role
	switch role {
	case model.RoleLogistics:
		r = proto.Role_LOGISTICS
	case model.RoleUser:
		r = proto.Role_USER
	default:
		//no op
	}
	res, err := service.GetCharacter(
		context.Background(),
		&proto.GetCharacterRequest{Token: &proto.Token{Identifier: c.token}, Role: r})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Character == nil {
		return nil, errors.New("expected grpc response to contain character, got nil")
	}
	return proto.ProtoToCharacter(res.Character), nil
}

// GetCharacter returns a populated Character for the given character ID.
func (c *CharacterClient) GetCharacter(charID int) (*model.Character, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetCharacter(
		context.Background(),
		&proto.GetCharacterRequest{Token: &proto.Token{Identifier: c.token}, CharacterId: int64(charID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Character == nil {
		return nil, errors.New("expected grpc response to contain character, got nil")
	}
	return proto.ProtoToCharacter(res.Character), nil
}

// GetCorporation returns a populated Corporation for the given corporation ID.
func (c *CharacterClient) GetCorporation(corpID int) (*model.Corporation, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetCorporation(
		context.Background(),
		&proto.GetCorporationRequest{Token: &proto.Token{Identifier: c.token}, CorporationId: int64(corpID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	corp := res.Corporation
	if corp == nil {
		return nil, errors.New("expected corporation in grpc response, got nil")
	}
	return proto.ProtoToCorporation(corp), nil
}

// GetAlliance returns a populated Alliance for the given alliance ID.
func (c *CharacterClient) GetAlliance(allianceID int) (*model.Alliance, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetAlliance(
		context.Background(),
		&proto.GetAllianceRequest{Token: &proto.Token{Identifier: c.token}, AllianceId: int64(allianceID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	alliance := res.Alliance
	if alliance == nil {
		return nil, errors.New("expected alliance in grpc response, got nil")
	}
	return proto.ProtoToAlliance(alliance), nil
}
