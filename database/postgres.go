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

// Get student from database sending the id
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
