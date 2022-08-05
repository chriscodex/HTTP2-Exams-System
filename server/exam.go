package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/ChrisCodeX/gRPC/exampb"
	"github.com/ChrisCodeX/gRPC/models"
	"github.com/ChrisCodeX/gRPC/repository"
	"github.com/ChrisCodeX/gRPC/studentpb"
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
// Enroll a student to the exam
func (s *ExamServer) EnrollStudents(stream exampb.ExamService_EnrollStudentsServer) error {
	for {
		// Recieve a message from the client
		msg, err := stream.Recv()

		// Response from server when client stop sending messages
		if err == io.EOF {
			// Response from server and close the stream
			return stream.SendAndClose(&exampb.EnrollmentResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}

		// Map the message into enrollment struct
		enrollment := &models.Enrollment{
			StudentId: msg.GetFkStudentId(),
			ExamId:    msg.GetFkExamId(),
		}

		// Insert question in database
		err = s.repo.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&exampb.EnrollmentResponse{
				Ok: false,
			})
		}
	}
}

// Get student by exam id
func (s *ExamServer) GetStudentsPerExam(req *exampb.GetStudentsPerExamRequest, stream exampb.ExamService_GetStudentsPerExamServer) error {
	// Get array of students
	students, err := s.repo.GetStudentsPerExam(context.Background(), req.GetFkExamId())
	if err != nil {
		return err
	}

	// Map student struct into student protobuf to be sended by the stream
	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}

		// Send the student to the client
		err := stream.Send(student)

		// Unnecessary code(Stream delay, Only to see the stream)
		time.Sleep(2 * time.Second)

		if err != nil {
			return err
		}
	}
	return nil
}

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
			ExamId:   msg.FkExamId,
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

//
func (s *ExamServer) TakeExam(stream exampb.ExamService_TakeExamServer) error {
	questions, err := s.repo.GetQuestionPerExam(context.Background(), "e1")
	if err != nil {
		return err
	}

	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &exampb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)

			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Answer: ", answer.GetAnswer())
	}
}
