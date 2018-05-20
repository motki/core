package model

import (
	"golang.org/x/net/context"

	"time"

	"github.com/motki/core/eveapi"
	"github.com/shopspring/decimal"
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

type IndustryManager struct {
	bootstrap

	corp *CorpManager
}

func newIndustryManager(m bootstrap, corp *CorpManager) *IndustryManager {
	return &IndustryManager{m, corp}
}

func (m *IndustryManager) GetCorporationIndustryJobs(ctx context.Context, corpID int) (jobs []*IndustryJob, err error) {
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	jobs, err = m.getCorporationIndustryJobsFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if jobs != nil {
		return jobs, nil
	}
	return m.getCorporationIndustryJobsFromAPI(ctx, corpID)
}

func (m *IndustryManager) getCorporationIndustryJobsFromDB(corpID int) ([]*IndustryJob, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.job_id
			, c.installer_id
			, c.facility_id
			, c.location_id
			, c.activity_id
			, c.blueprint_id
			, c.blueprint_type_id
			, c.blueprint_location_id
			, c.output_location_id
			, c.product_type_id
			, c.runs
			, c.cost
			, c.licensed_runs
			, c.probability
			, c.status
			, c.start_date
			, c.end_date
			, c.pause_date
			, c.completed_date
			, c.completed_character_id
			, c.successful_runs
			FROM app.industry_jobs c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '1 hour')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*IndustryJob
	for rs.Next() {
		r := &IndustryJob{}
		err := rs.Scan(
			&r.JobID,
			&r.InstallerID,
			&r.FacilityID,
			&r.LocationID,
			&r.ActivityID,
			&r.BlueprintID,
			&r.BlueprintTypeID,
			&r.BlueprintLocationID,
			&r.OutputLocationID,
			&r.ProductTypeID,
			&r.Runs,
			&r.Cost,
			&r.LicensedRuns,
			&r.Probability,
			&r.Status,
			&r.StartDate,
			&r.EndDate,
			&r.PauseDate,
			&r.CompletedDate,
			&r.CompletedCharacterID,
			&r.SuccessfulRuns,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *IndustryManager) getCorporationIndustryJobsFromAPI(ctx context.Context, corpID int) ([]*IndustryJob, error) {
	jobs, err := m.eveapi.GetCorporationIndustryJobs(ctx, corpID)
	if err != nil {
		return nil, err
	}
	hjobs, err := m.eveapi.GetCorporationIndustryJobHistory(ctx, corpID)
	if err != nil {
		return nil, err
	}
	jobs = append(jobs, hjobs...)
	return m.apiCorporationIndustryJobsToDB(corpID, jobs)
}

func (m *IndustryManager) apiCorporationIndustryJobsToDB(corpID int, jobs []*eveapi.IndustryJob) ([]*IndustryJob, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	res := make([]*IndustryJob, len(jobs))
	for i, j := range jobs {
		_, err = db.Exec(
			`INSERT INTO app.industry_jobs
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, DEFAULT)
					ON CONFLICT ON CONSTRAINT "industry_jobs_pkey" DO
						UPDATE SET completed_date = EXCLUDED.completed_date,
							     completed_character_id = EXCLUDED.completed_character_id,
							     pause_date = EXCLUDED.pause_date,
							     successful_runs = EXCLUDED.successful_runs,
							     fetched_at = DEFAULT`,
			j.JobID,
			corpID,
			j.InstallerID,
			j.FacilityID,
			j.LocationID,
			j.ActivityID,
			j.BlueprintID,
			j.BlueprintTypeID,
			j.BlueprintLocationID,
			j.OutputLocationID,
			j.ProductTypeID,
			j.Runs,
			j.Cost,
			j.LicensedRuns,
			j.Probability,
			j.Status,
			j.StartDate,
			j.EndDate,
			j.PauseDate,
			j.CompletedDate,
			j.CompletedCharacterID,
			j.SuccessfulRuns,
		)
		if err != nil {
			return nil, err
		}
		res[i] = &IndustryJob{
			JobID:                j.JobID,
			InstallerID:          j.InstallerID,
			FacilityID:           j.FacilityID,
			LocationID:           j.LocationID,
			ActivityID:           j.ActivityID,
			BlueprintID:          j.BlueprintID,
			BlueprintTypeID:      j.BlueprintTypeID,
			BlueprintLocationID:  j.BlueprintLocationID,
			OutputLocationID:     j.OutputLocationID,
			ProductTypeID:        j.ProductTypeID,
			Runs:                 j.Runs,
			Cost:                 j.Cost,
			LicensedRuns:         j.LicensedRuns,
			Probability:          j.Probability,
			Status:               j.Status,
			StartDate:            j.StartDate,
			EndDate:              j.EndDate,
			PauseDate:            j.PauseDate,
			CompletedDate:        j.CompletedDate,
			CompletedCharacterID: j.CompletedCharacterID,
			SuccessfulRuns:       j.SuccessfulRuns,
		}
	}
	return res, nil
}
