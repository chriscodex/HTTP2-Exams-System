package repository

import (
	"context"

	"github.com/ChrisCodeX/gRPC/models"
)

// Repository interface
type Repository interface {
	GetStudent(ctx context.Context, id string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error
}

// Assign Repository
var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

// Table Students Operations
func GetStudent(ctx context.Context, id string) (*models.Student, error) {
	return implementation.GetStudent(ctx, id)
}

func SetStudent(ctx context.Context, student *models.Student) error {
	return implementation.SetStudent(ctx, student)
}
