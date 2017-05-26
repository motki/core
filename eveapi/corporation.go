package eveapi

import (
	"time"
)

type Corporation struct {
	CorporationID int
	Name          string
	LogoURL       string
	AllianceID    int
	CreationDate  time.Time
	Description   string
	Ticker        string
}

type Alliance struct {
	AllianceID  int
	Name        string
	LogoURL     string
	DateFounded time.Time
	Ticker      string
}

func (api *EveAPI) GetCorporation(corpID int) (corp Corporation, err error) {
	dat, _, err := api.client.V3.CorporationApi.GetCorporationsCorporationId(int32(corpID), nil)
	if err != nil {
		return corp, err
	}
	corp.CorporationID = corpID
	corp.Name = dat.CorporationName
	corp.Description = dat.CorporationDescription
	corp.Ticker = dat.Ticker
	corp.CreationDate = dat.CreationDate
	corp.AllianceID = int(dat.AllianceId)
	corp.LogoURL = ImageURL(ImageCorpLogo, corpID, 128)
	return corp, nil
}

func (api *EveAPI) GetAlliance(allianceID int) (alliance Alliance, err error) {
	dat, _, err := api.client.V2.AllianceApi.GetAlliancesAllianceId(int32(allianceID), nil)
	if err != nil {
		return alliance, err
	}
	alliance.AllianceID = allianceID
	alliance.Name = dat.AllianceName
	alliance.Ticker = dat.Ticker
	alliance.DateFounded = dat.DateFounded
	alliance.LogoURL = ImageURL(ImageAllianceLogo, allianceID, 128)
	return alliance, nil
}
