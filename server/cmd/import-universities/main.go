// Import an extracted university CSV into the existing universities table.
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"modmapper/server/internal/platform/config"
	"modmapper/server/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// cli args
	filePath := flag.String("csv", "../data/universities.csv", "CSV with name,country columns")
	envPath := flag.String("env", "../.env", "Required dotenv file")
	flag.Parse()

	// open & read csv
	file, err := os.Open(*filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// validate headers
	header, err := reader.Read()
	if err != nil {
		return err
	}
	if len(header) != 2 || strings.TrimPrefix(header[0], "\uFEFF") != "name" || header[1] != "country" {
		return errors.New("CSV header must be name,country")
	}

	// read and validate every row
	var records [][]string
	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		row[0], row[1] = strings.TrimSpace(row[0]), strings.TrimSpace(row[1])
		if row[0] == "" || row[1] == "" {
			line, _ := reader.FieldPos(0)
			return fmt.Errorf("CSV line %d: empty name or country", line)
		}
		records = append(records, row)
	}
	if len(records) == 0 {
		return errors.New("CSV contains no universities")
	}

	// load DB config
	cfg, err := config.Load(*envPath)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	var inserted int64
	for _, row := range records {
		// skip insert if it's a duplicate
		tag, err := tx.Exec(ctx, `INSERT INTO universities (name, country)
			VALUES ($1, $2) ON CONFLICT (name, country) DO NOTHING`, row[0], row[1])
		if err != nil {
			return fmt.Errorf("failed to insert university: %w", err)
		}
		inserted += tag.RowsAffected()
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	fmt.Printf("Inserted %d universities; skipped %d existing/duplicate pairs.\n", inserted, int64(len(records))-inserted)
	return nil
}
