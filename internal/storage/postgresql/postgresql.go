package postgresql

import (
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
)



const (
    host     = "localhost"
    port     = 5432
    user     = "todo"
    password = "123456"
    dbname   = "tododb"
)


type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {

	const op = "storage.postgresql.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50),
			note VARCHAR(1000),
			importance INTEGER DEFAULT 0
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt2, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(40),
			password VARCHAR(100)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) CheckUserInDb(username string, password string) (bool, error) {
	const op = "storage.postgresql.CheckUserInDb"
   
	query := `SELECT COUNT(*) FROM users WHERE username=$1 AND password=$2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
	 return false, fmt.Errorf("%s: %w", op, err)
	}
   
	var count int
	err = stmt.QueryRow(username, password).Scan(&count)
	if err != nil {
	 return false, fmt.Errorf("%s: %w", op, err)
	}
   
	return count > 0, nil
   }
   