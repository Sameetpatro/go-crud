package storage

import "github.com/Sameetpatro/go-crud/internal/types"
type Storage interface{
	CreateStudent(name string, email string, age int) (int, error)
	GetStudentById(id int64)  (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateTheStudent(id int64, name string, email string, age int) (types.Student, error) 
}
