package eveapi

import (
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

func (api *EveAPI) GetCharacter(id int) (Character, error) {
	char := Character{}
	dat, _, err := api.client.V4.CharacterApi.GetCharactersCharacterId(int32(id), nil)
	if err != nil {
		return char, err
	}
	char.CharacterID = id
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
