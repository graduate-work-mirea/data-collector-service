package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	db *sql.DB
}

// RawData represents the raw data collected from external sources
type RawData struct {
	ID           int64
	DataSourceID int64
	Data         json.RawMessage
	CollectedAt  time.Time
	CreatedAt    time.Time
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// SaveRawData saves raw data to the database
func (p *PostgresDB) SaveRawData(ctx context.Context, data json.RawMessage) (int64, error) {
	var id int64
	query := `
		INSERT INTO raw_data (data_source_id, data)
		VALUES (1, $1)
		RETURNING id
	`
	err := p.db.QueryRowContext(ctx, query, data).Scan(&id)
	return id, err
}

// GetRawData retrieves raw data from the database by ID
func (p *PostgresDB) GetRawData(ctx context.Context, id int64) (*RawData, error) {
	var rawData RawData
	query := `
		SELECT id, data_source_id, data, collected_at, created_at
		FROM raw_data
		WHERE id = $1
	`
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&rawData.ID,
		&rawData.DataSourceID,
		&rawData.Data,
		&rawData.CollectedAt,
		&rawData.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rawData, nil
}
