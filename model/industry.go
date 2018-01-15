package model

import (
	"golang.org/x/net/context"

	"github.com/motki/core/eveapi"
)

func (m *Manager) GetCorporationIndustryJobs(ctx context.Context, corpID int) (jobs []*eveapi.IndustryJob, err error) {
	jobs, err = m.getCorporationIndustryJobsFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if jobs != nil {
		return jobs, nil
	}
	return m.getCorporationIndustryJobsFromAPI(ctx, corpID)
}

func (m *Manager) getCorporationIndustryJobsFromDB(corpID int) ([]*eveapi.IndustryJob, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.job_id
			, c.installer_id
			, c.installer_name
			, c.facility_id
			, c.solar_system_name
			, c.solar_system_id
			, c.station_id
			, c.activity_id
			, c.blueprint_id
			, c.blueprint_type_id
			, c.blueprint_type_name
			, c.blueprint_location_id
			, c.output_location_id
			, c.product_type_id
			, c.runs
			, c.cost
			, c.licensed_runs
			, c.probability
			, c.product_type_name
			, c.status
			, c.time_in_seconds
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
	var res []*eveapi.IndustryJob
	for rs.Next() {
		r := &eveapi.IndustryJob{}
		err := rs.Scan(
			&r.JobID,
			&r.InstallerID,
			&r.InstallerName,
			&r.FacilityID,
			&r.SolarSystemName,
			&r.SolarSystemID,
			&r.StationID,
			&r.ActivityID,
			&r.BlueprintID,
			&r.BlueprintTypeID,
			&r.BlueprintTypeName,
			&r.BlueprintLocationID,
			&r.OutputLocationID,
			&r.ProductTypeID,
			&r.Runs,
			&r.Cost,
			&r.LicensedRuns,
			&r.Probability,
			&r.ProductTypeName,
			&r.Status,
			&r.TimeInSeconds,
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

func (m *Manager) getCorporationIndustryJobsFromAPI(ctx context.Context, corpID int) ([]*eveapi.IndustryJob, error) {
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

func (m *Manager) apiCorporationIndustryJobsToDB(corpID int, jobs []*eveapi.IndustryJob) ([]*eveapi.IndustryJob, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, j := range jobs {
		_, err = db.Exec(
			`INSERT INTO app.industry_jobs
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, DEFAULT)
					ON CONFLICT ON CONSTRAINT "industry_jobs_pkey" DO
						UPDATE SET completed_date = EXCLUDED.completed_date,
							     completed_character_id = EXCLUDED.completed_character_id,
							     pause_date = EXCLUDED.pause_date,
							     successful_runs = EXCLUDED.successful_runs,
							     fetched_at = DEFAULT`,
			j.JobID,
			corpID,
			j.InstallerID,
			j.InstallerName,
			j.FacilityID,
			j.SolarSystemName,
			j.SolarSystemID,
			j.StationID,
			j.ActivityID,
			j.BlueprintID,
			j.BlueprintTypeID,
			j.BlueprintTypeName,
			j.BlueprintLocationID,
			j.OutputLocationID,
			j.ProductTypeID,
			j.Runs,
			j.Cost,
			j.LicensedRuns,
			j.Probability,
			j.ProductTypeName,
			j.Status,
			j.TimeInSeconds,
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
	}
	return jobs, nil
}
