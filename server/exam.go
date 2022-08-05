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
			Id:        msg.GetId(),
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
			Answer:   msg.Answer,
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

// Get student by exam id
func (s *ExamServer) GetQuestionsPerExam(req *exampb.GetQuestionsPerExamRequest, stream exampb.ExamService_GetQuestionsPerExamServer) error {
	// Get array of students
	questions, err := s.repo.GetQuestionPerExam(context.Background(), req.GetFkExamId())
	if err != nil {
		return err
	}

	// Map student struct into student protobuf to be sended by the stream
	for _, question := range questions {
		question := &exampb.Question{
			Id:       question.Id,
			Question: question.Question,
			Answer:   question.Answer,
		}

		// Send the student to the client
		err := stream.Send(question)

		// Unnecessary code(Stream delay, Only to see the stream)
		time.Sleep(2 * time.Second)

		if err != nil {
			return err
		}
	}
	return nil
}

// Take exam
func (s *ExamServer) TakeExam(stream exampb.ExamService_TakeExamServer) error {
	// Recieve a message from the client
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	enrollment, err := s.repo.GetEnrollmentById(context.Background(), msg.GetEnrollmentId())
	if err != nil {
		return err
	}

	// Get array of questions from database by exam id
	questions, err := s.repo.GetQuestionPerExam(context.Background(), enrollment.ExamId)
	if err != nil {
		return err
	}
	qts := &exampb.Question{
		Id: enrollment.ExamId,
	}
	err = stream.Send(qts)
	if err != nil {
		return err
	}

	//
	i := 0
	var count uint16
	var currentQuestion = &models.Question{}
	for i < len(questions) {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			// Send Question from protobuf file
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
		if answer.GetAnswer() == currentQuestion.Answer {
			count++
		}
		log.Println("Answer: ", answer.GetAnswer())
	}
	countQuestions, err := s.repo.GetCountQuestionsByExamId(context.Background(), enrollment.ExamId)
	if err != nil {
		return err
	}
	countQuest := float32(*countQuestions)
	score := float32(float32(count) * 10 / countQuest)
	log.Printf("Score: %.2f", score)

	// scoreString := fmt.Sprintf("%f", score)
	// err = s.repo.SetScore(context.Background(), enrollment.ExamId, scoreString)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	return nil
}

// func (s *ExamServer) GetScore(ctx context.Context, req *exampb.GetScoreRequest) (*exampb.GetScoreResponse, error) {
// 	// Get exam from database
// 	enrollment, err := s.repo.GetEnrollmentById(ctx, "m1")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Map enrollment struct into GetScoreResponse message protobuf
// 	return &exampb.GetScoreResponse{
// 		EnrollmentId: enrollment.Id,
// 		FkStudentId:  enrollment.StudentId,
// 		FkExamId:     enrollment.ExamId,
// 		Score:        enrollment.Score,
// 	}, nil
// }
