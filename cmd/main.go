package main

import (
	"backend-assignment/config"
	"backend-assignment/internal/handler/orderhandler"
	"backend-assignment/internal/handler/settlementhandler"
	"backend-assignment/internal/handler/userhandler"
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database/psql/tbljob"
	"backend-assignment/internal/repository/database/psql/tblorder"
	"backend-assignment/internal/repository/database/psql/tbluser"
	"backend-assignment/internal/service/orderservice"
	"backend-assignment/internal/service/settlementservice"
	"backend-assignment/internal/service/userservice"
	"backend-assignment/internal/worker"
	"backend-assignment/logs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	file, err := logs.InitLog()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	validate := validator.New()
	userRepo := tbluser.New(db)
	userService := userservice.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService, validate)

	// order
	orderRepo := tblorder.New(db)
	orderService := orderservice.NewOrderService(orderRepo, userRepo)
	orderHandler := orderhandler.NewOrderHandler(orderService, validate)

	// jobs
	jobRepo := tbljob.New(db)

	jobQueue := make(chan models.Job, 100)

	w := worker.NewWorker(db, jobRepo)
	go w.Start(jobQueue)

	jobService := settlementservice.NewSettlementService(jobRepo)
	jobHandler := settlementhandler.NewSettlementHandler(jobService, validate, jobQueue)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	// users API
	users := v1.Group("/users")
	users.POST("/register", userHandler.Register)

	// orders API
	orders := v1.Group("/orders")
	orders.POST("", orderHandler.CreateOrder)
	orders.GET("/:id", orderHandler.GetOrderDetail)

	// jobs API
	jobs := v1.Group("/jobs")
	jobs.POST("/settlement", jobHandler.CreateJob)
	jobs.GET("/:id", jobHandler.GetJobById)

	r.Static("/downloads", "/tmp/settlements")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
