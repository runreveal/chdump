package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	// Replace these with your actual database credentials and connection details
	const (
		address  = "tcp://localhost:9000"
		username = "default"
		password = ""
		database = "default"
	)

	clickhouseDSN := os.Args[1]

	// Establish a connection to the database
	// connStr := fmt.Sprintf("%s/%s?username=%s&password=%s", address, database, username, password)
	db, err := sql.Open("clickhouse", clickhouseDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	// Ping the database to ensure connection is established
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Query to retrieve table names in the specified database
	tableNamesQuery := "SHOW TABLES"
	rows, err := db.Query(tableNamesQuery)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		// Query to retrieve the DDL of each table
		showCreateTableQuery := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
		tableDDL, err := db.Query(showCreateTableQuery)
		if err != nil {
			log.Fatalf("Failed to get DDL for table %s: %v", tableName, err)
		}
		defer tableDDL.Close()

		if tableDDL.Next() {
			var ddl string
			if err := tableDDL.Scan(&ddl); err != nil {
				log.Fatalf("Failed to scan DDL of table %s: %v", tableName, err)
			}
			fmt.Printf("%s\n;\n----------------------------------------\n", ddl)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
	}
}
