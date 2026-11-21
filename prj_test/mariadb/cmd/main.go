package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
)

type Note struct {
	ID        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func main() {
	// _ := context.Background()
	dbDSN := readEnv()
	conn, _ := sql.Open("mysql", dbDSN)
	defer conn.Close()

	fake := faker.New()
	// sb := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Question).
		Columns("title", "body").
		Values(fake.Address().City(), fake.Address().StreetAddress()).
		Suffix("RETURNING id")

	query, arg, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("Failed to build query %v", err)
	}

	var noteId int

	err = conn.QueryRow(query, arg...).Scan(&noteId)
	if err != nil {
		log.Fatalf("filaed to insert row %v", err)
	}

	log.Printf("inserted row with id: %v", noteId)

	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		// PlaceholderFormat(sq.Question).
		OrderBy("id ASC").
		Limit(10)

	query, arg, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("Failed to build query %v", err)
	}
	log.Printf(query, arg)
	rows, err := conn.Query(query)
	if err != nil {
		log.Fatalf("failed select query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var created_at time.Time
		var updated_at sql.NullTime

		err = rows.Scan(&id, &title, &body, &created_at, &updated_at)
		if err != nil {
			log.Fatalf("Failed to scan %v", err)
		}

		log.Printf("id: %d, title: %s, body: %s, created: %v, updated: %v", id, title, body, created_at, updated_at)

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
