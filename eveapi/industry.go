package eveapi

import (
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type IndustryJob struct {
	JobID                int
	InstallerID          int
	InstallerName        string
	FacilityID           int
	SolarSystemName      string
	SolarSystemID        int
	StationID            int
	ActivityID           int
	BlueprintID          int
	BlueprintTypeID      int
	BlueprintTypeName    string
	BlueprintLocationID  int
	OutputLocationID     int
	ProductTypeID        int
	Runs                 int
	Cost                 decimal.Decimal
	LicensedRuns         int
	Probability          decimal.Decimal
	ProductTypeName      string
	Status               int
	TimeInSeconds        int
	StartDate            time.Time
	EndDate              time.Time
	PauseDate            time.Time
	CompletedDate        time.Time
	CompletedCharacterID int
	SuccessfulRuns       int
}

func (api *EveAPI) GetCorporationIndustryJobs(ctx context.Context, corpID int) (jobs []*IndustryJob, err error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationIndustryJobsXML(tok, int64(corpID))
	if err != nil {
		return nil, err
	}
	for _, j := range res.Entries {
		job := &IndustryJob{
			JobID:                int(j.JobID),
			InstallerID:          int(j.InstallerID),
			InstallerName:        j.InstallerName,
			FacilityID:           int(j.FacilityID),
			SolarSystemName:      j.SolarSystemName,
			SolarSystemID:        int(j.SolarSystemID),
			StationID:            int(j.StationID),
			ActivityID:           int(j.ActivityID),
			BlueprintID:          int(j.BlueprintID),
			BlueprintTypeID:      int(j.BlueprintTypeID),
			BlueprintTypeName:    j.BlueprintTypeName,
			BlueprintLocationID:  int(j.BlueprintLocationID),
			OutputLocationID:     int(j.OutputLocationID),
			ProductTypeID:        int(j.ProductTypeID),
			Runs:                 int(j.Runs),
			Cost:                 decimal.NewFromFloat(j.Cost),
			LicensedRuns:         int(j.LicensedRuns),
			Probability:          decimal.NewFromFloat(j.Probability),
			ProductTypeName:      j.ProductTypeName,
			Status:               int(j.Status),
			TimeInSeconds:        int(j.TimeInSeconds),
			StartDate:            j.StartDate.Time,
			EndDate:              j.EndDate.Time,
			PauseDate:            j.PauseDate.Time,
			CompletedDate:        j.CompletedDate.Time,
			CompletedCharacterID: int(j.CompletedCharacterID),
			SuccessfulRuns:       int(j.SuccessfulRuns),
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (api *EveAPI) GetCorporationIndustryJobHistory(ctx context.Context, corpID int) (jobs []*IndustryJob, err error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationIndustryJobsHistoryXML(tok, int64(corpID))
	if err != nil {
		return nil, err
	}
	for _, j := range res.Entries {
		job := &IndustryJob{
			JobID:                int(j.JobID),
			InstallerID:          int(j.InstallerID),
			InstallerName:        j.InstallerName,
			FacilityID:           int(j.FacilityID),
			SolarSystemName:      j.SolarSystemName,
			SolarSystemID:        int(j.SolarSystemID),
			StationID:            int(j.StationID),
			ActivityID:           int(j.ActivityID),
			BlueprintID:          int(j.BlueprintID),
			BlueprintTypeID:      int(j.BlueprintTypeID),
			BlueprintTypeName:    j.BlueprintTypeName,
			BlueprintLocationID:  int(j.BlueprintLocationID),
			OutputLocationID:     int(j.OutputLocationID),
			ProductTypeID:        int(j.ProductTypeID),
			Runs:                 int(j.Runs),
			Cost:                 decimal.NewFromFloat(j.Cost),
			LicensedRuns:         int(j.LicensedRuns),
			Probability:          decimal.NewFromFloat(j.Probability),
			ProductTypeName:      j.ProductTypeName,
			Status:               int(j.Status),
			TimeInSeconds:        int(j.TimeInSeconds),
			StartDate:            j.StartDate.Time,
			EndDate:              j.EndDate.Time,
			PauseDate:            j.PauseDate.Time,
			CompletedDate:        j.CompletedDate.Time,
			CompletedCharacterID: int(j.CompletedCharacterID),
			SuccessfulRuns:       int(j.SuccessfulRuns),
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
