package sqlite

import (
	"database/sql"
	"github.com/Sameetpatro/go-crud/internal/config"
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