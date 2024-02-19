package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	

	_ "github.com/lib/pq"
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


type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Importance int `json:"importance"`
}

type TaskById struct{
	Note string `json:"note"`
	Importance int `json:"importance"`
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




func (s *Storage) CheckAllUserTasks(username string) ([]Task, error) {
	const op = "storage.postgresql.CheckAllUserTasks"
   
	query := `SELECT id, note, importance FROM notes WHERE username=$1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
	 return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
   
	rows, err := stmt.Query(username)
	if err != nil {
	 return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
   
	tasks := []Task{}
	for rows.Next() {
	 var id int
	 var note string
	 var importance int
	 if err := rows.Scan(&id, &note, &importance); err != nil {
	  return nil, fmt.Errorf("%s: %w", op, err)
	 }
	 task := Task{id, note, importance}
	 tasks = append(tasks, task)
	}
   
	if err := rows.Err(); err != nil {
	 return nil, fmt.Errorf("%s: %w", op, err)
	}
   
	return tasks, nil
   }


func (s *Storage) CheckTaskById(username string, id string) (*TaskById, error) {
    const op = "storage.postgresql.CheckTaskById"

    query := "SELECT note, importance FROM notes WHERE username=$1 AND id=$2"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer stmt.Close()

    var task TaskById
    err = stmt.QueryRowContext(context.Background(), username, id).Scan(&task.Note, &task.Importance)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    
    return &task, nil
}


func (s *Storage) DeleteTaskById(username string, id string) (bool, error){
	const op = "storage.postgresql.DeleteTaskById"

	query := `DELETE FROM notes WHERE username=$1 and id=$2`
	stmt, err := s.db.Exec(query, username, id)
	if err != nil{
		return false, fmt.Errorf("Error %v", op)
	}
	rowAffected, _ := stmt.RowsAffected()
	if rowAffected == 0{
		return false, fmt.Errorf("no rows to delete")
	}

	return true, nil
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
   