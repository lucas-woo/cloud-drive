package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySql() *sql.DB {

  user := os.Getenv("DB_USER")
  pass := os.Getenv("DB_PASSWORD")
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  dbName := os.Getenv("DB_NAME")

  dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
    user, pass, host, port, dbName,
  ) 

  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal("failed to open db:", err)
  } 

  fmt.Println("connected to mysql")
  return db
}