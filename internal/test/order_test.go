package test

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database/psql/tblorder"
	"context"
	"database/sql"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5433/backendtest?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`TRUNCATE orders RESTART IDENTITY; UPDATE products SET stock = 500 WHERE id = 1;`)
	assert.NoError(t, err)

	return db
}

func TestConcurrentOrders(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute)

	orderRepo := tblorder.New(db)

	var success int64
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(uid int) {
			defer wg.Done()

			_, err := orderRepo.CreateOrder(ctx, &models.RequestOrder{
				ProductId: 1,
				UserId:    1,
				Quantity:  1,
			})

			if err == nil {
				atomic.AddInt64(&success, 1)
			} else if !strings.Contains(err.Error(), "insufficient stock") {
				t.Errorf("unexpected error: %v", err)
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int64(500), success, "should only allow 500 successful orders")

	var stock int
	err := db.QueryRow(`SELECT stock FROM products WHERE id=1`).Scan(&stock)
	assert.NoError(t, err)
	assert.Equal(t, 0, stock, "stock should be 0 after 500 successful orders")

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 500, count, "orders table should have 100 rows")
}
