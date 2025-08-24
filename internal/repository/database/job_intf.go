package database

import (
	"backend-assignment/internal/models"
	"context"
)

type Job struct {
	Service TblJob
}

type TblJob interface {
	CreateJob(ctx context.Context, job *models.Job) (*models.Job, error)
	GetJob(ctx context.Context, job *models.Job) (*models.Job, error)
	UpdateStatus(ctx context.Context, jobID, status string) error
	UpdateProgress(ctx context.Context, jobId string, progress float64, processedCount int, totalCount int) error
	UpdateResult(ctx context.Context, jobId, status, resultPath string) error
}
