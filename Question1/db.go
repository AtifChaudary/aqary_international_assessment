package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	CreateUser(ctx context.Context, name, phoneNumber string) error
	GenerateOTP(ctx context.Context, phoneNumber string) (string, error)
	VerifyOTP(ctx context.Context, phoneNumber, otp string) error
}

type database struct {
	conn *pgxpool.Pool
}

func NewDB(conn *pgxpool.Pool) DB {
	return &database{conn: conn}
}

func (d *database) CreateUser(ctx context.Context, name, phoneNumber string) error {
	query := "INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id"
	var userID int
	err := d.conn.QueryRow(ctx, query, name, phoneNumber).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (d *database) GenerateOTP(ctx context.Context, phoneNumber string) (string, error) {
	otp := generateRandomOTP()
	expirationTime := time.Now().Add(1 * time.Minute)

	query := "UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE phone_number = $3 RETURNING otp"
	var generatedOTP string
	err := d.conn.QueryRow(ctx, query, otp, expirationTime, phoneNumber).Scan(&generatedOTP)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %v", err)
	}

	return generatedOTP, nil
}

func (d *database) VerifyOTP(ctx context.Context, phoneNumber, otp string) error {
	var storedOTP string
	var expirationTime time.Time

	query := "SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1"
	err := d.conn.QueryRow(ctx, query, phoneNumber).Scan(&storedOTP, &expirationTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to verify OTP: %v", err)
	}

	if storedOTP != otp {
		return fmt.Errorf("incorrect OTP")
	}

	if time.Now().After(expirationTime) {
		return fmt.Errorf("expired OTP")
	}

	return nil
}

func generateRandomOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
