package server

import (
	"context"

	"github.com/ChrisCodeX/gRPC/models"
	"github.com/ChrisCodeX/gRPC/repository"
	"github.com/ChrisCodeX/gRPC/studentpb"
)

// Server struct for student service
type Server struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
}

// Server Constructor
func NewStudentServer(repo repository.Repository) *Server {
	return &Server{
		repo: repo,
	}
}

// Methods to student service (Defined in protobuf file)
// Protobuffer
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

//
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
