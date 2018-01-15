package eveapi

import (
	"strconv"
	"time"

	"github.com/antihax/goesi/eveapi"
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

func xmlToCorporationSheet(sheet *eveapi.CorporationSheetXML) CorporationSheet {
	return CorporationSheet{
		CorporationID:   int(sheet.CorporationID),
		CorporationName: sheet.CorporationName,
		Ticker:          sheet.Ticker,
		CEOID:           int(sheet.CEOID),
		CEOName:         sheet.CEOName,
		StationID:       int(sheet.StationID),
		StationName:     sheet.StationName,
		Description:     sheet.Description,
		AllianceID:      int(sheet.AllianceID),
		AllianceName:    sheet.AllianceName,
		FactionID:       int(sheet.FactionID),
		URL:             sheet.URL,
		MemberCount:     int(sheet.MemberCount),
		Shares:          int(sheet.Shares),
	}
}

func (api *EveAPI) GetPublicCorporationSheet(corpID int) (*CorporationSheet, error) {
	sheet, err := api.client.EVEAPI.CorporationPublicSheetXML(int64(corpID))
	if err != nil {
		return nil, err
	}
	res := xmlToCorporationSheet(sheet)
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

func (api *EveAPI) GetCorporationSheet(ctx context.Context) (*CorporationSheetDetail, error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sheet, err := api.client.EVEAPI.CorporationSheetXML(tok)
	if err != nil {
		return nil, err
	}
	res := &CorporationSheetDetail{
		CorporationSheet: xmlToCorporationSheet(&sheet.CorporationSheetXML),
		Wallets:          make(Divisions),
		Hangars:          make(Divisions),
	}
	for _, set := range sheet.Divisions {
		switch set.Name {
		case "walletDivisions":
			for _, div := range set.Accounts {
				if div.Key == "1000" {
					res.Wallets["1000"] = "Master Wallet"
					continue
				}
				res.Wallets[div.Key] = div.Description
			}
		case "divisions":
			for _, div := range set.Accounts {
				res.Hangars[div.Key] = div.Description
			}

		default:
			// do nothing
		}
	}
	return res, nil
}
