package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connStr = "user=myuser password=mypassword dbname=animetavern_db sslmode=disable"

var DB *sql.DB

func InitDB() error {
	var err error
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	user := os.Getenv("myuser")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	sslmode := os.Getenv("sslmode")
	conn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	DB, err = sql.Open("postgres", conn)
	return err

}
