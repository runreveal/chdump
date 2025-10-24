package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	database := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9999"},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	// Ping the database to ensure connection is established
	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Query to retrieve table names in the specified database
	tableNamesQuery := "SHOW TABLES"
	rows, err := conn.Query(tableNamesQuery)
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
		tableDDL, err := conn.Query(showCreateTableQuery)
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
