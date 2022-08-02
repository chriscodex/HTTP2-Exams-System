package server

import (
	"context"

	"github.com/ChrisCodeX/gRPC/models"
	"github.com/ChrisCodeX/gRPC/repository"
	"github.com/ChrisCodeX/gRPC/studentpb"
)

// Server for student service
type StudentServer struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
}

/*Student Service*/
// Assign database to Student Server
func NewStudentServer(repo repository.Repository) *StudentServer {
	return &StudentServer{
		repo: repo,
	}
}

// Methods to Student Service (Defined in protobuf file)
// Request to Get a Student
func (s *StudentServer) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
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
func (s *StudentServer) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {
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
