package eveapi

import (
	"context"
	"time"
)

type Character struct {
	CharacterID   int
	Name          string
	BloodlineID   int
	RaceID        int
	AncestryID    int
	CorporationID int
	AllianceID    int
	BirthDate     time.Time
	Description   string
}

func (api *EveAPI) GetCharacter(characterID int) (char *Character, err error) {
	dat, _, err := api.client.ESI.CharacterApi.GetCharactersCharacterId(context.Background(), int32(characterID), nil)
	if err != nil {
		return char, err
	}
	char = &Character{}
	char.CharacterID = characterID
	char.Name = dat.Name
	char.Description = dat.Description
	char.RaceID = int(dat.RaceId)
	char.AncestryID = int(dat.RaceId)
	char.BloodlineID = int(dat.BloodlineId)
	char.BirthDate = dat.Birthday
	char.CorporationID = int(dat.CorporationId)
	char.AllianceID = int(dat.AllianceId)
	return char, nil
}
