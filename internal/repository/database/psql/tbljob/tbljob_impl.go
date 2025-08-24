package tbljob

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/logs"
	"context"
	"database/sql"
)

type TblJobImpl struct {
	DB *sql.DB
}

func New(DB *sql.DB) database.TblJob {
	return &TblJobImpl{
		DB: DB,
	}
}

// CreateJob implements database.TblJob.
func (repository *TblJobImpl) CreateJob(ctx context.Context, job *models.Job) (*models.Job, error) {
	tx, err := repository.DB.Begin()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	query := `
		INSERT INTO jobs (job_id, from_date, to_date, status, progress, processed_count, total_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = repository.DB.ExecContext(ctx, query,
		job.JobID, job.FromDate, job.ToDate, job.Status, job.Progress,
		job.ProcessedCount, job.TotalCount)

	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return job, nil
}

// GetJob implements database.TblJob.
func (repository *TblJobImpl) GetJob(ctx context.Context, job *models.Job) (*models.Job, error) {
	query := `
		SELECT job_id, from_date, to_date, status, progress, processed_count, total_count, result_path
		FROM jobs
		WHERE job_id = $1`

	row := repository.DB.QueryRowContext(ctx, query, job.JobID)

	err := row.Scan(&job.JobID, &job.FromDate, &job.ToDate, &job.Status,
		&job.Progress, &job.ProcessedCount, &job.TotalCount, &job.ResultPath)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// UpdateStatus implements database.TblJob.
func (repository *TblJobImpl) UpdateStatus(ctx context.Context, jobID string, status string) error {
	query := `
		UPDATE jobs 
		SET status = $1 
		WHERE job_id = $2`
	_, err := repository.DB.ExecContext(ctx, query, status, jobID)
	return err
}

// UpdateProgress implements database.TblJob.
func (repository *TblJobImpl) UpdateProgress(ctx context.Context, jobId string, progress float64, processedCount int, totalCount int) error {
	query := `
		UPDATE jobs 
		SET progress = $1, processed_count = $2, total_count = $3
		WHERE job_id = $4`
	_, err := repository.DB.ExecContext(ctx, query, progress, processedCount, totalCount, jobId)
	return err
}

// UpdateResult implements database.TblJob.
func (repository *TblJobImpl) UpdateResult(ctx context.Context, jobId, status, resultPath string) error {
	query := `
		UPDATE jobs 
		SET status = $1, result_path = $2, progress = 100.00
		WHERE job_id = $3`
	_, err := repository.DB.ExecContext(ctx, query, status, resultPath, jobId)
	return err
}
