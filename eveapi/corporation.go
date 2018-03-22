package eveapi

import (
	"strconv"
	"time"

	"github.com/antihax/goesi/esi"
	"golang.org/x/net/context"
)

type Corporation struct {
	CorporationID int
	Name          string
	AllianceID    int
	CreationDate  time.Time
	Description   string
	Ticker        string
}

type Alliance struct {
	AllianceID  int
	Name        string
	DateFounded time.Time
	Ticker      string
}

func (api *EveAPI) GetCorporation(corpID int) (corp *Corporation, err error) {
	dat, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationId(context.Background(), int32(corpID), nil)
	if err != nil {
		return corp, err
	}
	corp = &Corporation{}
	corp.CorporationID = corpID
	corp.Name = dat.Name
	corp.Description = dat.Description
	corp.Ticker = dat.Ticker
	corp.CreationDate = dat.DateFounded
	corp.AllianceID = int(dat.AllianceId)
	return corp, nil
}

func (api *EveAPI) GetAlliance(allianceID int) (alliance *Alliance, err error) {
	dat, _, err := api.client.ESI.AllianceApi.GetAlliancesAllianceId(context.Background(), int32(allianceID), nil)
	if err != nil {
		return alliance, err
	}
	alliance = &Alliance{}
	alliance.AllianceID = allianceID
	alliance.Name = dat.Name
	alliance.Ticker = dat.Ticker
	alliance.DateFounded = dat.DateFounded
	return alliance, nil
}

type CorporationSheet struct {
	CorporationID   int
	CorporationName string
	Ticker          string
	CEOID           int
	CEOName         string
	StationID       int
	StationName     string
	Description     string
	AllianceID      int
	AllianceName    string
	FactionID       int
	URL             string
	MemberCount     int
	Shares          int
}

func corpResponseToSheet(corpID int, sheet esi.GetCorporationsCorporationIdOk) CorporationSheet {
	return CorporationSheet{
		CorporationID:   corpID,
		CorporationName: sheet.Name,
		Ticker:          sheet.Ticker,
		CEOID:           int(sheet.CeoId),
		StationID:       int(sheet.HomeStationId),
		Description:     sheet.Description,
		AllianceID:      int(sheet.AllianceId),
		FactionID:       int(sheet.FactionId),
		URL:             sheet.Url,
		MemberCount:     int(sheet.MemberCount),
		Shares:          int(sheet.Shares),
	}
}

func (api *EveAPI) GetPublicCorporationSheet(corpID int) (*CorporationSheet, error) {
	sheet, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationId(context.Background(), int32(corpID), nil)
	if err != nil {
		return nil, err
	}
	res := corpResponseToSheet(corpID, sheet)
	return &res, nil
}

// Divisions is a map of division key to division name.
type Divisions map[string]string

func (d Divisions) GetName(idx int) (string, bool) {
	id := strconv.Itoa(idx)
	if v, ok := d[id]; ok {
		return v, true
	}
	return "", false
}

type CorporationSheetDetail struct {
	CorporationSheet
	Wallets Divisions
	Hangars Divisions
}

func (api *EveAPI) GetCorporationSheet(ctx context.Context, corpID int) (*CorporationSheetDetail, error) {
	_, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sheet, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationId(ctx, int32(corpID), nil)
	if err != nil {
		return nil, err
	}
	res := &CorporationSheetDetail{
		CorporationSheet: corpResponseToSheet(corpID, sheet),
		Wallets:          make(Divisions),
		Hangars:          make(Divisions),
	}
	divs, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationIdDivisions(ctx, int32(corpID), nil)
	if err != nil {
		return nil, err
	}
	for _, div := range divs.Hangar {
		res.Hangars[strconv.Itoa(int(div.Division))] = div.Name
	}
	for _, div := range divs.Wallet {
		key := strconv.Itoa(int(div.Division))
		if key == "1" {
			res.Wallets["1"] = "Master Wallet"
			continue
		}
		res.Wallets[key] = div.Name
	}
	return res, nil
}
