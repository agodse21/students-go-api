package storage

import "github.com/agodse21/students-go-api/internal/types"

type Storage interface {
	CreateStudent(student types.Student) (int64, error)
	GetStudentById(id int64) (types.Student, error)

	GetStudents() ([]types.Student, error)

	DeleteStudentById(id int64) error

	UpdateStudentById(id int64, student types.Student) error
}
