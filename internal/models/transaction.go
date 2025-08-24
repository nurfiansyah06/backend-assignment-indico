package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	MerchantId  string    `json:"merchant_id"`
	AmountCents int       `json:"amount_cents"`
	FeeCents    int       `json:"fee_cents"`
	Status      string    `json:"status"`
	PaidAt      time.Time `json:"paid_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RequestTransaction struct {
	MerchantId  string `json:"merchant_id"`
	AmountCents int    `json:"amount_cents"`
}
