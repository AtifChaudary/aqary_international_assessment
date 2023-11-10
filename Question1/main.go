package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db DB

func main() {
	initDB()
	defer db.(*database).conn.Close()

	router := gin.Default()

	router.POST("/api/users", createUser)
	router.POST("/api/users/generateotp", generateOTP)
	router.POST("/api/users/verifyotp", verifyOTP)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting Gin server:", err)
	}
}

func initDB() {
	config, err := pgxpool.ParseConfig("postgres://atif:atif@localhost:5432/aqaryint")
	if err != nil {
		log.Fatal("Error parsing db URL:", err)
	}

	pool, err := pgxpool.ConnectConfig(nil, config)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	db = NewDB(pool)
}

func createUser(c *gin.Context) {
	var user struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.CreateUser(context.Background(), user.Name, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func generateOTP(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := db.GenerateOTP(context.Background(), request.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp})
}

func verifyOTP(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
		OTP         string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.VerifyOTP(context.Background(), request.PhoneNumber, request.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
