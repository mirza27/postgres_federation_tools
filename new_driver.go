package main_2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// User struct mewakili tabel users
type User struct {
	ID       string  `db:"user_id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
	UserType string `db:"user_type"`
	Status   string `db:"status"`
}

// Driver struct mewakili tabel drivers
type Driver struct {
	ID     string  `db:"driver_id"`
	UserID string  `db:"user_id"`
	Name   string `db:"name"`
}

func main() {
	// Konfigurasi database (ganti dengan konfigurasi Anda)
	targetDBConfig := struct {
		Host     string
		Port     int
		User     string
		Name     string
		Password string
		SSLMode  string
	}{
		Host:     "localhost",
		Port:     5433,
		User:     "new_user",
		Password: "dsav9das",
		Name:     "new_ojol_db",
		SSLMode:  "disable",
	}

	targetConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		targetDBConfig.Host, targetDBConfig.Port, targetDBConfig.User,
		targetDBConfig.Password, targetDBConfig.Name, targetDBConfig.SSLMode)

	targetDB, err := sqlx.Connect("postgres", targetConnStr)
	if err != nil {
		log.Fatal("Error connecting to target DB:", err)
	}
	defer targetDB.Close()

	new_user := "tedasd"

	// Data untuk dimasukkan
	user := User{
		Email:    new_user +"@gmail.com",
		Username: new_user,
		Password: "dasbsdaiu12",
		UserType: "driver",
		Status:   "active",
	}

	// Panggil fungsi untuk melakukan transaksi
	userID, driverID, err := insertUserAndDriver(targetDB, user, new_user)
	if err != nil {
		log.Fatal("Error inserting user and driver:", err)
	}

	fmt.Printf("Successfully inserted user with ID: %d and driver with ID: %d\n", userID, driverID)
}

// insertUserAndDriver melakukan transaksi untuk insert ke tabel users dan drivers
func insertUserAndDriver(db *sqlx.DB, user User, driverName string) (string, string, error) {
	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return "", "", fmt.Errorf("begin transaction: %w", err)
	}

	// Defer rollback jika terjadi error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("rollback error: %v", rbErr)
			}
		}
	}()

	// Insert ke tabel users
	var userID string
	query := `
		INSERT INTO users (email, username, password, user_type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id
	`
	
	err = tx.QueryRowxContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.UserType,
		user.Status,
		time.Now(),
		time.Now(),
	).Scan(&userID)
	
	if err != nil {
		return "", "", fmt.Errorf("insert user: %w", err)
	}

	// Insert ke tabel drivers menggunakan userID yang didapat
	var driverID string
	query = `
		INSERT INTO drivers (user_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING driver_id
	`
	
	err = tx.QueryRowxContext(ctx, query,
		userID,
		driverName,
		time.Now(),
		time.Now(),
	).Scan(&driverID)
	
	if err != nil {
		return "", "", fmt.Errorf("insert driver: %w", err)
	}

	// Commit transaksi
	if err = tx.Commit(); err != nil {
		return "", "", fmt.Errorf("commit transaction: %w", err)
	}

	return userID, driverID, nil
}