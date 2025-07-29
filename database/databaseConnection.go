package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func DBInstance() *sql.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "restaurant_management"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	fmt.Println("Database connection established successfully")
	return db
}

var Client *sql.DB

func InitializeDB() {
	Client = DBInstance()
}

func CloseDBConnection() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Fatal("Error closing the database connection:", err)
		} else {
			fmt.Println("Database connection closed successfully")
		}
	}
}
