package models

// Model of student
type Student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

// Model of exam
type Exam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Model of enrollment
type Enrollment struct {
	Id        string `json:"id"`
	StudentId string `json:"fk_student_id"`
	ExamId    string `json:"fk_exam_id"`
}

// Model of question
type Question struct {
	Id       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	ExamId   string `json:"fk_exam_id"`
}
