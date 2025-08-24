package worker

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type Worker struct {
	DB      *sql.DB
	JobRepo database.TblJob
}

func NewWorker(db *sql.DB, jobRepo database.TblJob) *Worker {
	return &Worker{DB: db, JobRepo: jobRepo}
}

func (w *Worker) Start(jobQueue <-chan models.Job) {
	for job := range jobQueue {
		// tiap job diproses oleh goroutine sendiri
		go w.processJob(job)
	}
}

func (w *Worker) processJob(job models.Job) {
	ctx := context.Background()

	// set job jadi RUNNING
	_ = w.JobRepo.UpdateStatus(ctx, job.JobID, "RUNNING")

	// hitung total transaksi
	var total int
	_ = w.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM transactions 
		 WHERE paid_at BETWEEN $1 AND $2`,
		job.FromDate, job.ToDate).Scan(&total)

	job.TotalCount = total
	// update total_count di DB juga
	_ = w.JobRepo.UpdateProgress(ctx, job.JobID, job.Progress, job.ProcessedCount, job.TotalCount)

	batchSize := 10000
	offset := 0
	results := map[string]*models.Settlement{}

	for {
		rows, err := w.DB.QueryContext(ctx,
			`SELECT merchant_id, amount_cents, fee_cents, paid_at
			 FROM transactions
			 WHERE paid_at BETWEEN $1 AND $2
			 ORDER BY id
			 LIMIT $3 OFFSET $4`,
			job.FromDate, job.ToDate, batchSize, offset)
		if err != nil {
			_ = w.JobRepo.UpdateStatus(ctx, job.JobID, "FAILED")
			return
		}

		count := 0
		for rows.Next() {
			var mid string
			var amount, fee int64
			var paid time.Time
			_ = rows.Scan(&mid, &amount, &fee, &paid)

			key := fmt.Sprintf("%s_%s", mid, paid.Format("2006-01-02"))
			if _, ok := results[key]; !ok {
				results[key] = &models.Settlement{
					MerchantID:     mid,
					SettlementDate: paid,
				}
			}
			s := results[key]
			s.GrossCents += amount
			s.FeeCents += fee
			s.NetCents += (amount - fee)
			s.TxnCount++
			count++
		}
		rows.Close()

		if count == 0 {
			break
		}

		// update progress
		job.ProcessedCount += count
		if job.TotalCount > 0 {
			job.Progress = float64(job.ProcessedCount) / float64(job.TotalCount) * 100
		} else {
			job.Progress = 100
		}
		offset += batchSize

		_ = w.JobRepo.UpdateProgress(ctx, job.JobID, job.Progress, job.ProcessedCount, job.TotalCount)
	}

	// tulis CSV hasil settlement
	filename := fmt.Sprintf("/tmp/settlements/%s.csv", job.JobID)
	f, _ := os.Create(filename)
	writer := csv.NewWriter(f)
	defer f.Close()

	for _, s := range results {
		_ = writer.Write([]string{
			s.MerchantID,
			s.SettlementDate.Format("2006-01-02"),
			fmt.Sprint(s.GrossCents),
			fmt.Sprint(s.FeeCents),
			fmt.Sprint(s.NetCents),
			fmt.Sprint(s.TxnCount),
		})
	}
	writer.Flush()

	// update DONE
	downloadURL := fmt.Sprintf("/downloads/%s.csv", job.JobID)
	_ = w.JobRepo.UpdateResult(ctx, job.JobID, "COMPLETED", downloadURL)
}
