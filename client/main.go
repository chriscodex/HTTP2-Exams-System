package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/ChrisCodeX/gRPC/exampb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect with the server
	// Using insecure.NewCredentials because server isn't encrypted
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	// Close the connection with server
	defer cc.Close()

	// Create the client
	c := exampb.NewExamServiceClient(cc)

	// Invoke Unary Method
	//DoUnary(c, "e1")

	// Invoke Client Stream
	//DoClientStreaming(c, "e1")

	// Invoke Server Stream
	//DoServerStreaming(c, "e1")

	// Invoke Bidirection Stream
	DoBidirectionalStreaming(c, "m2")
}

// Unary Connections
func DoUnary(c exampb.ExamServiceClient, examId string) {
	// Create the request
	req := &exampb.GetExamRequest{
		Id: examId,
	}

	// Invoke server function in the client by gRPC
	res, err := c.GetExam(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetExam: %v", err)
	}

	// Response from server
	log.Printf("response from server: %v", res)
}

// Client Streaming
func DoClientStreaming(c exampb.ExamServiceClient, examId string) {
	// Create questions
	questions := []*exampb.Question{
		{
			Id:       "q5",
			Question: "6x6",
			Answer:   "36",
			FkExamId: examId,
		},
		{
			Id:       "q6",
			Question: "7x7",
			Answer:   "49",
			FkExamId: examId,
		},
		{
			Id:       "q7",
			Question: "8x8",
			Answer:   "64",
			FkExamId: examId,
		},
	}

	// Get streaming
	stream, err := c.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("error while calling SetQuestions: %v", err)
	}

	// Client: Send questions
	for _, question := range questions {
		stream.Send(question)
		log.Println("Sending question: ", question.Id)
		time.Sleep(1 * time.Second)
	}

	// Tell the server to close stream
	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	// Response from server
	log.Printf("response from server: %v", msg)
}

// Server Streaming
// rpc GetStudentsPerExam(GetStudentsPerExamRequest) returns (stream student.Student);
func DoServerStreaming(c exampb.ExamServiceClient, examId string) {
	// Request from client
	req := &exampb.GetStudentsPerExamRequest{
		FkExamId: examId,
	}

	// Get stream
	stream, err := c.GetStudentsPerExam(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetStudentsPerExam: %v", err)
	}

	// Receive students from stream
	for {
		msg, err := stream.Recv()

		// End streaming
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		// Response from server
		log.Printf("response from server: %v", msg)
	}
}

// Bidirectional Streaming
// rpc TakeExam(stream TakeExamRequest) returns (stream Question)
func DoBidirectionalStreaming(c exampb.ExamServiceClient, enrollmentId string) {
	// Create the answer request
	answerRequest := exampb.TakeExamRequest{
		EnrollmentId: enrollmentId,
		Answer:       "42",
	}

	numberOfQuestion := 6

	// Control channel
	waitChannel := make(chan struct{})

	// Get stream
	stream, err := c.TakeExam(context.Background())
	if err != nil {
		log.Fatalf("error while calling TakeExam: %v", err)
	}

	// Write stream
	// Send answers
	go func() {
		for i := 0; i < numberOfQuestion; i++ {
			// Send answer request
			stream.Send(&answerRequest)
			time.Sleep(5 * time.Second)
		}
	}()

	// Read stream
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}

			// Response from server
			log.Printf("response from server: %v", res)
		}
		close(waitChannel)
	}()
	<-waitChannel
}
