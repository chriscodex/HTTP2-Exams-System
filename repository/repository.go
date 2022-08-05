package repository

import (
	"context"

	"github.com/ChrisCodeX/gRPC/models"
)

// Repository interface
type Repository interface {
	/* Student Service */
	GetStudent(ctx context.Context, id string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error

	/* Exam Service */
	//Unary Methods
	GetExam(ctx context.Context, id string) (*models.Exam, error)
	SetExam(ctx context.Context, exam *models.Exam) error

	// Stream from client
	// Enrollment
	SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error

	GetStudentsPerExam(ctx context.Context, examId string) ([]*models.Student, error)

	// Questions
	SetQuestion(ctx context.Context, question *models.Question) error

	// Questions per Exam
	GetQuestionPerExam(ctx context.Context, examId string) ([]*models.Question, error)

	// Questions Count per Exam
	GetCountQuestionsByExamId(ctx context.Context, examId string) (*uint16, error)

	// Get Enrollment by ID
	GetEnrollmentById(ctx context.Context, enrollmentId string) (*models.Enrollment, error)

	// Set Score
	//SetScore(ctx context.Context, enrollmentId string, score string) error
}

// Assign Repository
var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

// Abstract Implementations

/* Service Students Operations */
// Get a student
func GetStudent(ctx context.Context, id string) (*models.Student, error) {
	return implementation.GetStudent(ctx, id)
}

// Insert a student
func SetStudent(ctx context.Context, student *models.Student) error {
	return implementation.SetStudent(ctx, student)
}

/* Service Exams Operations */
// Get a Exam
func GetExam(ctx context.Context, id string) (*models.Exam, error) {
	return implementation.GetExam(ctx, id)
}

// Insert a Exam
func SetExam(ctx context.Context, exam *models.Exam) error {
	return implementation.SetExam(ctx, exam)
}

// Insert Questions
func SetQuestion(ctx context.Context, question *models.Question) error {
	return implementation.SetQuestion(ctx, question)
}

// Enroll a student to a exam
func SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return implementation.SetEnrollment(ctx, enrollment)
}

// Get students by exam id
func GetStudentsPerExam(ctx context.Context, examId string) ([]*models.Student, error) {
	return implementation.GetStudentsPerExam(ctx, examId)
}

// Get Questions per Exam
func GetQuestionPerExam(ctx context.Context, examId string) ([]*models.Question, error) {
	return implementation.GetQuestionPerExam(ctx, examId)
}

// Get Count of questions by exam id
func GetCountQuestionsByExamId(ctx context.Context, examId string) (*uint16, error) {
	return implementation.GetCountQuestionsByExamId(ctx, examId)
}

// Get Enrollment by id
func GetEnrollmentById(ctx context.Context, enrollmentId string) (*models.Enrollment, error) {
	return implementation.GetEnrollmentById(ctx, enrollmentId)
}

// Set Score
// func SetScore(ctx context.Context, enrollmentId string, score string) error {
// 	return implementation.SetScore(ctx, enrollmentId, score)
// }
