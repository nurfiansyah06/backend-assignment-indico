package settlementservice

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SettlementServiceImpl struct {
	tblJob database.TblJob
}

func NewSettlementService(tblJob database.TblJob) service.SettlementService {
	return &SettlementServiceImpl{
		tblJob: tblJob,
	}
}

// CreateJob implements service.SettlementService.
func (s *SettlementServiceImpl) CreateJob(ctx context.Context, req models.RequestJob) (*models.Job, *models.Response, error) {
	fromDate, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		return nil, nil, err
	}
	toDate, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		return nil, nil, err
	}

	job := &models.Job{
		JobID:          uuid.New().String(),
		FromDate:       fromDate,
		ToDate:         toDate,
		Status:         "QUEUED",
		Progress:       0,
		ProcessedCount: 0,
		TotalCount:     0,
	}

	createdJob, err := s.tblJob.CreateJob(ctx, job)
	if err != nil {
		return nil, nil, err
	}

	response := models.Response{
		Status: http.StatusAccepted,
		Data: map[string]interface{}{
			"job_id": createdJob.JobID,
			"status": createdJob.Status,
		},
		Message: "Success Create Job",
	}

	return createdJob, &response, nil
}

// GetJobById implements service.SettlementService.
func (s *SettlementServiceImpl) GetJobById(ctx context.Context, jobID string) (*models.Response, error) {
	job := &models.Job{
		JobID: jobID,
	}

	foundJob, err := s.tblJob.GetJob(ctx, job)
	if err != nil {
		return nil, err
	}

	if foundJob.Progress < 100 {
		foundJob.Status = "RUNNING"
	}

	response := models.Response{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"job_id":       foundJob.JobID,
			"status":       foundJob.Status,
			"progress":     foundJob.Progress,
			"processed":    foundJob.ProcessedCount,
			"total":        foundJob.TotalCount,
			"download_url": foundJob.ResultPath,
		},
		Message: "Success Get Job",
	}

	return &response, nil
}
