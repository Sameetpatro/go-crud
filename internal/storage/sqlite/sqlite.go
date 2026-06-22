package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/Sameetpatro/go-crud/internal/config"
	"github.com/Sameetpatro/go-crud/internal/types"
	_ "github.com/mattn/go-sqlite3"
) 

type Sqlite struct{
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error){
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, errr := db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL
	);`)
	if errr != nil {
		return nil, errr
	}
	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int, error){
	statement, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)") //we use '?' as placeholder to prevent sql injection
	if err != nil{
		return 0, err
	}
	defer statement.Close()
	res, errr := statement.Exec(name, email, age)
	if errr != nil{
		return 0, errr
	}
	lastid, errrr := res.LastInsertId()
	if errrr != nil{
		return 0, errrr
	}
	return int(lastid), nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error){
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil{
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	errr := stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if errr != nil {
		if errr == sql.ErrNoRows{
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", errr)
	}
	return student, nil
}