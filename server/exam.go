package server

import (
	"context"
	"io"

	"github.com/ChrisCodeX/gRPC/exampb"
	"github.com/ChrisCodeX/gRPC/models"
	"github.com/ChrisCodeX/gRPC/repository"
)

type ExamServer struct {
	repo repository.Repository
	exampb.UnimplementedExamServiceServer
}

/*Exam Service*/
// Assign database to Exam Server
func NewExamServer(repo repository.Repository) *ExamServer {
	return &ExamServer{
		repo: repo,
	}
}

//Exams Service Unary Methods
// Request to Get a Exam
func (s *ExamServer) GetExam(ctx context.Context, req *exampb.GetExamRequest) (*exampb.Exam, error) {
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
func (s *ExamServer) SetExam(ctx context.Context, req *exampb.Exam) (*exampb.SetExamResponse, error) {
	// Map student protobuf into student struct
	exam := &models.Exam{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	// Send exam to database
	err := s.repo.SetExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	// Return SetExamResponse protobuf
	return &exampb.SetExamResponse{
		Id: exam.Id,
	}, nil
}

// Exam Service Stream Methods
// Methods SetQuestions
func (s *ExamServer) SetQuestions(stream exampb.ExamService_SetQuestionsServer) error {
	for {
		// Recieve a message from the client
		msg, err := stream.Recv()

		// Response from server when client stop sending messages
		if err == io.EOF {
			// Response from server and close the stream
			return stream.SendAndClose(&exampb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		// Map the message into question struct
		question := &models.Question{
			Id:       msg.Id,
			Question: msg.Question,
			Answer:   msg.Answer,
			ExamId:   msg.ExamId,
		}

		// Insert question in database
		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(&exampb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}
