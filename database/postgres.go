package database

import (
	"context"
	"database/sql"

	"github.com/ChrisCodeX/gRPC/models"
	_ "github.com/lib/pq"
)

// PostgresRepository
type PostgresRepository struct {
	db *sql.DB
}

// Constructor of Postgres Repository
func NewPostgreRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

// Methods

/*Student Service*/
// Get student from database by id
func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var student = models.Student{}

	for rows.Next() {
		if err = rows.Scan(&student.Id, &student.Name, &student.Age); err != nil {
			return &student, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &student, err
	}

	return &student, nil
}

// Insert a student into the database
func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)",
		student.Id, student.Name, student.Age)
	return err
}

/*Exams Service*/
// Get exam by id
func (repo *PostgresRepository) GetExam(ctx context.Context, id string) (*models.Exam, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name FROM exams WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var exam = models.Exam{}

	for rows.Next() {
		if err = rows.Scan(&exam.Id, &exam.Name); err != nil {
			return &exam, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &exam, err
	}

	return &exam, nil
}

// Insert a exam into the database
func (repo *PostgresRepository) SetExam(ctx context.Context, exam *models.Exam) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO exams (id, name) VALUES ($1, $2)",
		exam.Id, exam.Name)
	return err
}
