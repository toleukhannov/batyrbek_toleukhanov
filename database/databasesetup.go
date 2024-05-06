package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "development"
	password = "testpassword"
	dbname   = "canteenmanagement"
)

func DBSet() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Connected to PostgreSQL")
	return db
}

var DB = DBSet()

func UserData(db *sql.DB, tableName string) *sql.DB {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, tableName))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func ProductData(db *sql.DB, tableName string) *sql.DB {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name TEXT,
			price NUMERIC,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, tableName))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
