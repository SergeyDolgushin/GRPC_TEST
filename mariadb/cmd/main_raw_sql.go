package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
)

func main() {
	dbDSN := readEnv()
	ctx := context.Background()

	fake := faker.New()

	conn, err := sql.Open("mysql", dbDSN)
	defer conn.Close()

	res, err := conn.ExecContext(ctx,
		"INSERT INTO note (title, body) VALUES (?, ?)",
		fake.Address().City(), //
		fake.Lorem().Sentence(30),
	)
	if err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}

	// Получение ID вставленной записи
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last ID: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	fmt.Printf("Inserted ID: %d, Rows affected: %d\n", lastID, rowsAffected)

	rows, err := conn.QueryContext(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Fatalf("Error %+v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var created_at time.Time
		var updated_at sql.NullTime

		err = rows.Scan(&id, &title, &body, &created_at, &updated_at)
		if err != nil {
			log.Fatalf("Failed to scan", err)
		}

		log.Printf("id: %d, title: %s, bode: %s, created: %v, updated: %v", id, title, body, created_at, updated_at)

	}

}

func readEnv() string {
	envFile := flag.String("env", ".env", "path to .env file")
	flag.Parse()
	flag.CommandLine.Parse(os.Args[1:])

	err := godotenv.Load(*envFile)
	if err != nil {
		log.Printf("Warning: .env file not found at %s: %v", *envFile, err)
		// Продолжаем работу — используем системные переменные
	}

	dbDSN := os.Getenv("MYSQL_DSN")
	fmt.Println("Using DSN from:", *envFile)
	fmt.Println("DB DSN:", dbDSN)

	return dbDSN
}
