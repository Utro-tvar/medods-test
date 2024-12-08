package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"
	host, exists := os.LookupEnv("DB_HOST")
	if !exists {
		return nil, fmt.Errorf("%s: DB_HOST env var not found", op)
	}
	port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		return nil, fmt.Errorf("%s: DB_PORT env var not found", op)
	}
	user, exists := os.LookupEnv("DB_USER")
	if !exists {
		return nil, fmt.Errorf("%s: DB_USER env var not found", op)
	}
	password, exists := os.LookupEnv("DB_PASS")
	if !exists {
		return nil, fmt.Errorf("%s: DB_PASS env var not found", op)
	}
	dbname, exists := os.LookupEnv("DB_NAME")
	if !exists {
		return nil, fmt.Errorf("%s: DB_NAME env var not found", op)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Store(guid, token string) error {
	const op = "storage.postgres.Store"
	stmt, err := s.db.Prepare("INSERT INTO refreshes VALUES ($1, $2);")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(fmt.Sprintf("'%s'", guid), fmt.Sprintf("'%s'", token))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetHash(guid string) (hash string, err error) {
	const op = "storage.postgres.CheckAndRemove"
	stmt, err := s.db.Prepare("DELETE FROM refreshes WHERE guid = $1 RETURNING hash;")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	row := ""
	res := stmt.QueryRow(fmt.Sprintf("'%s'", guid))

	err = res.Scan(&row)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if row == "" {
		return "", nil
	}
	return row, nil
}

func (s *Storage) Close() {
	s.db.Close()
}
