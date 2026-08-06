package main

import (
	"os"

	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/data/model"
)

func main() {
	student := &model.Student{
		StudentId: "9999999999",
		Name:      "Test Student",
		Gender:    "male",
		AcademicData: &model.AcademicData{
			Faculty:         "School of Computing",
			StudyProgram:    "Informatics",
			CurrentSemester: 1,
		},
		AcademicHistory: []*model.SemesterRecord{
			{
				Semester:    1,
				SemesterGpa: 3.5,
				Courses: []*model.Course{
					{Code: "C101", Name: "Test Course A", Credits: 3, Score: 85},
					{Code: "C102", Name: "Test Course B", Credits: 3, Score: 90},
					{Code: "C103", Name: "Test Course C", Credits: 4, Score: 88},
				},
			},
			{
				Semester:    2,
				SemesterGpa: 3.6,
				Courses: []*model.Course{
					{Code: "C201", Name: "Test Course D", Credits: 3, Score: 92},
					{Code: "C202", Name: "Test Course E", Credits: 3, Score: 87},
					{Code: "C203", Name: "Test Course F", Credits: 4, Score: 89},
				},
			},
		},
		CumulativeGpa: 3.55,
	}

	b, err := proto.Marshal(student)
	if err != nil {
		panic(err)
	}
	os.WriteFile("tmp/student.bin", b, 0644)
}
