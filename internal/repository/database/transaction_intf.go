package database

import "context"

type Transaction struct {
	Service TblTransaction
}

type TblTransaction interface {
	CreateTransaction(ctx context.Context)
}
