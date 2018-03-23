package model

import (
	"golang.org/x/net/context"

	"time"

	"github.com/motki/core/eveapi"
	"github.com/shopspring/decimal"
)

type IndustryJob struct {
	JobID                int
	InstallerID          int
	FacilityID           int
	LocationID           int
	ActivityID           int
	BlueprintID          int
	BlueprintTypeID      int
	BlueprintLocationID  int
	OutputLocationID     int
	ProductTypeID        int
	Runs                 int
	Cost                 decimal.Decimal
	LicensedRuns         int
	Probability          decimal.Decimal
	Status               string
	StartDate            time.Time
	EndDate              time.Time
	PauseDate            time.Time
	CompletedDate        time.Time
	CompletedCharacterID int
	SuccessfulRuns       int
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
		res[i] = (*IndustryJob)(j)
	}
	return res, nil
}
