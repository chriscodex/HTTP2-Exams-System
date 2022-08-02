package server

import (
	"context"

	"github.com/ChrisCodeX/gRPC/exampb"
	"github.com/ChrisCodeX/gRPC/models"
	"github.com/ChrisCodeX/gRPC/repository"
	"github.com/ChrisCodeX/gRPC/studentpb"
)

// Server struct for services
type Server struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
	exampb.UnimplementedExamServiceServer
}

/*Student Service*/
// Assign database to Student Server
func NewStudentServer(repo repository.Repository) *Server {
	return &Server{
		repo: repo,
	}
}

// Methods to Student Service (Defined in protobuf file)
// Request to Get a Student
func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	// Get student from database
	student, err := s.repo.GetStudent(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Map student struct into student protobuf
	return &studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

// Request to Set a Student
func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {
	// Map student protobuf into student struct
	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	// Send student to database
	err := s.repo.SetStudent(ctx, student)
	if err != nil {
		return nil, err
	}

	// Return SetStudentResponse protobuf
	return &studentpb.SetStudentResponse{
		Id: student.Id,
	}, nil
}

/*Exam Service*/
// Assign database to Exam Server
func NewExamServer(repo repository.Repository) *Server {
	return &Server{
		repo: repo,
	}
}

//Exams Service Methods
// Request to Get a Exam
func (s *Server) GetExam(ctx context.Context, req *exampb.GetExamRequest) (*exampb.Exam, error) {
	// Get exam from database
	exam, err := s.repo.GetExam(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Map student struct into student protobuf
	return &exampb.Exam{
		Id:   exam.Id,
		Name: exam.Name,
	}, nil
}

// Request to Set a Exam
func (s *Server) SetExam(ctx context.Context, req *exampb.Exam) (*exampb.SetExamResponse, error) {
	// Map student protobuf into student struct
	exam := &models.Exam{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	// Send student to database
	err := s.repo.SetExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	// Return SetStudentResponse protobuf
	return &exampb.SetExamResponse{
		Id: exam.Id,
	}, nil
}
