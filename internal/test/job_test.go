package test

import (
	"backend-assignment/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockSettlementService struct {
	mock.Mock
}

func (m *MockSettlementService) CreateJob(ctx context.Context, req models.RequestJob) (*models.Job, *models.Response, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Job), args.Get(1).(*models.Response), args.Error(2)
}

func (m *MockSettlementService) GetJobById(ctx context.Context, jobId string) (*models.Response, error) {
	args := m.Called(ctx, jobId)
	return args.Get(0).(*models.Response), args.Error(1)
}

func TestCreateJobSuccess(t *testing.T) {
	mockService := new(MockSettlementService)
	req := models.RequestJob{
		From: "2023-01-01",
		To:   "2023-01-31",
	}
	expectedJob := &models.Job{
		JobID:    "test-job-id",
		Status:   "QUEUED",
		Progress: 0,
	}
	expectedResponse := &models.Response{
		Status:  202,
		Message: "Success Create Job",
		Data: map[string]interface{}{
			"job_id": expectedJob.JobID,
			"status": expectedJob.Status,
		},
	}

	mockService.On("CreateJob", mock.Anything, req).Return(expectedJob, expectedResponse, nil)

	job, response, err := mockService.CreateJob(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if job.JobID != expectedJob.JobID {
		t.Errorf("expected job ID %v, got %v", expectedJob.JobID, job.JobID)
	}
	if response.Status != expectedResponse.Status {
		t.Errorf("expected response status %v, got %v", expectedResponse.Status, response.Status)
	}

	mockService.AssertExpectations(t)
}

func TestCreateJobFailed(t *testing.T) {
	mockService := new(MockSettlementService)

	req := models.RequestJob{
		From: "invalid-date",
		To:   "2023-01-31",
	}
	expectedError := "invalid date format"

	mockService.On("CreateJob", mock.Anything, req).
		Return((*models.Job)(nil), (*models.Response)(nil), errors.New(expectedError))

	job, response, err := mockService.CreateJob(context.Background(), req)

	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if err.Error() != expectedError {
		t.Errorf("expected error message %v, got %v", expectedError, err.Error())
	}
	if job != nil {
		t.Errorf("expected no job, got %v", job)
	}
	if response != nil {
		t.Errorf("expected no response, got %v", response)
	}

	mockService.AssertExpectations(t)
}

func TestGetJobSuccess(t *testing.T) {
	mockService := new(MockSettlementService)
	jobID := "test-job-id"
	expectedJob := &models.Job{
		JobID:    jobID,
		Status:   "RUNNING",
		Progress: 50,
	}
	expectedResponse := &models.Response{
		Status:  200,
		Message: "Success Get Job",
		Data:    expectedJob,
	}

	mockService.On("GetJobById", mock.Anything, jobID).Return(expectedResponse, nil)

	response, err := mockService.GetJobById(context.Background(), jobID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.Status != expectedResponse.Status {
		t.Errorf("expected response status %v, got %v", expectedResponse.Status, response.Status)
	}
	if response.Data.(*models.Job).JobID != expectedJob.JobID {
		t.Errorf("expected job ID %v, got %v", expectedJob.JobID, response.Data.(*models.Job).JobID)
	}

	mockService.AssertExpectations(t)
}

func TestGetJobFailed(t *testing.T) {
	mockService := new(MockSettlementService)
	jobID := "non-existent-job-id"
	expectedError := "job not found"

	mockService.On("GetJobById", mock.Anything, jobID).
		Return((*models.Response)(nil), errors.New(expectedError))

	response, err := mockService.GetJobById(context.Background(), jobID)

	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if err.Error() != expectedError {
		t.Errorf("expected error message %v, got %v", expectedError, err.Error())
	}
	if response != nil {
		t.Errorf("expected no response, got %v", response)
	}

	mockService.AssertExpectations(t)

}
