package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"airline-voucher/internal/models"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS vouchers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crew_name TEXT NOT NULL,
	crew_id TEXT NOT NULL,
	flight_number TEXT NOT NULL,
	flight_date TEXT NOT NULL,
	aircraft_type TEXT NOT NULL,
	seat1 TEXT NOT NULL,
	seat2 TEXT NOT NULL,
	seat3 TEXT NOT NULL,
	created_at TEXT NOT NULL,
	UNIQUE (flight_number, flight_date)
);`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("migrate vouchers: %w", err)
	}
	return nil
}

func (s *Store) ExistsForFlight(flightNumber, date string) (bool, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(1) FROM vouchers WHERE flight_number = ? AND flight_date = ?`,
		flightNumber,
		date,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) Create(voucher *models.Voucher) error {
	if voucher.CreatedAt == "" {
		voucher.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	result, err := s.db.Exec(
		`INSERT INTO vouchers (
			crew_name, crew_id, flight_number, flight_date, aircraft_type,
			seat1, seat2, seat3, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		voucher.CrewName,
		voucher.CrewID,
		voucher.FlightNumber,
		voucher.FlightDate,
		voucher.AircraftType,
		voucher.Seat1,
		voucher.Seat2,
		voucher.Seat3,
		voucher.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	voucher.ID = id
	return nil
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint failed")
}
