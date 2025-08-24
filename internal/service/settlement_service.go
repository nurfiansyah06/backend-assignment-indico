package service

import (
	"backend-assignment/internal/models"
	"context"
)

type SettlementService interface {
	CreateJob(ctx context.Context, req models.RequestJob) (*models.Job, *models.Response, error)
	GetJobById(ctx context.Context, jobID string) (*models.Response, error)
}
