package models

import "time"

type Settlement struct {
	ID             int64     `db:"id"`
	MerchantID     string    `db:"merchant_id"`
	SettlementDate time.Time `db:"settlement_date"`
	GrossCents     int64     `db:"gross_cents"`
	FeeCents       int64     `db:"fee_cents"`
	NetCents       int64     `db:"net_cents"`
	TxnCount       int64     `db:"txn_count"`
	GeneratedAt    time.Time `db:"generated_at"`
	UniqueJobID    string    `db:"unique_job_id"`
}

type RequestSettlement struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}
