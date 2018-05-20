package eveapi

import (
	"strconv"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type IndustryJob struct {
	JobID                int             `json:"job_id"`
	InstallerID          int             `json:"installer_id"`
	FacilityID           int             `json:"facility_id"`
	LocationID           int             `json:"location_id"`
	ActivityID           int             `json:"activity_id"`
	BlueprintID          int             `json:"blueprint_id"`
	BlueprintTypeID      int             `json:"blueprint_type_id"`
	BlueprintLocationID  int             `json:"blueprint_location_id"`
	OutputLocationID     int             `json:"output_location_id"`
	ProductTypeID        int             `json:"product_type_id"`
	Runs                 int             `json:"runs"`
	Cost                 decimal.Decimal `json:"cost"`
	LicensedRuns         int             `json:"licensed_runs"`
	Probability          decimal.Decimal `json:"probability"`
	Status               string          `json:"status"`
	StartDate            time.Time       `json:"start_date"`
	EndDate              time.Time       `json:"end_date"`
	PauseDate            time.Time       `json:"pause_date"`
	CompletedDate        time.Time       `json:"completed_date"`
	CompletedCharacterID int             `json:"completed_character_id"`
	SuccessfulRuns       int             `json:"successful_runs"`
}

func (api *EveAPI) GetCorporationIndustryJobs(ctx context.Context, corpID int) (jobs []*IndustryJob, err error) {
	_, err = TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	for max, p := 1, 1; p <= max; p++ {
		res, resp, err := api.client.ESI.IndustryApi.GetCorporationsCorporationIdIndustryJobs(
			ctx,
			int32(corpID),
			&esi.GetCorporationsCorporationIdIndustryJobsOpts{
				IncludeCompleted: optional.NewBool(true),
				Page:             optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, j := range res {
			job := &IndustryJob{
				JobID:                int(j.JobId),
				InstallerID:          int(j.InstallerId),
				FacilityID:           int(j.FacilityId),
				LocationID:           int(j.LocationId),
				ActivityID:           int(j.ActivityId),
				BlueprintID:          int(j.BlueprintId),
				BlueprintTypeID:      int(j.BlueprintTypeId),
				BlueprintLocationID:  int(j.BlueprintLocationId),
				OutputLocationID:     int(j.OutputLocationId),
				ProductTypeID:        int(j.ProductTypeId),
				Runs:                 int(j.Runs),
				Cost:                 decimal.NewFromFloat(j.Cost),
				LicensedRuns:         int(j.LicensedRuns),
				Probability:          decimal.NewFromFloat(float64(j.Probability)),
				Status:               j.Status,
				StartDate:            j.StartDate,
				EndDate:              j.EndDate,
				PauseDate:            j.PauseDate,
				CompletedDate:        j.CompletedDate,
				CompletedCharacterID: int(j.CompletedCharacterId),
				SuccessfulRuns:       int(j.SuccessfulRuns),
			}
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (api *EveAPI) GetCorporationIndustryJobHistory(ctx context.Context, corpID int) (jobs []*IndustryJob, err error) {
	return api.GetCorporationIndustryJobs(ctx, corpID)
}
