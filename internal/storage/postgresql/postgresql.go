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

	const op = "storage.sqlite.New"

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

	return &Storage{db: db}, nil

}

// func Connect(){
//     psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
//         "password=%s dbname=%s sslmode=disable",
//         host, port, user, password, dbname)
//     db, err := sql.Open("postgres", psqlInfo)
//     if err != nil {
//         panic(err)
//     }
//     defer db.Close()
  
//     err = db.Ping()
//     if err != nil {
//         panic(err)
//     }
  
//     fmt.Println("Successfully connected!")
// }