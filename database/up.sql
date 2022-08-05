DROP TABLE IF EXISTS students;

CREATE TABLE students (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL
);

DROP TABLE IF EXISTS exams;

CREATE TABLE exams (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS enrollments;

CREATE TABLE enrollments (
    id VARCHAR(32) PRIMARY KEY,
    score VARCHAR(32),
    fk_student_id VARCHAR(32) NOT NULL,
    fk_exam_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (fk_student_id) REFERENCES students(id),
    FOREIGN KEY (fk_exam_id) REFERENCES exams(id)
);

DROP TABLE IF EXISTS questions;

CREATE TABLE questions (
    id VARCHAR(32) PRIMARY KEY,
    fk_exam_id VARCHAR(32) NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    FOREIGN KEY (fk_exam_id) REFERENCES exams(id)
);