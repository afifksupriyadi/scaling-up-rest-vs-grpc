package seeder

import (
	"fmt"

	"scaling-up-rest-vs-grpc/internal/data/model"
)

func (f fakeStudent) toProto() *model.Student {
	history := make([]*model.SemesterRecord, len(f.AcademicHistory))
	for i, sem := range f.AcademicHistory {
		courses := make([]*model.Course, len(sem.Courses))
		for j, c := range sem.Courses {
			courses[j] = &model.Course{
				Code:    fmt.Sprintf("C%d%02d", sem.Semester, j+1),
				Name:    c.Name,
				Credits: c.Credits,
				Score:   c.Score,
			}
		}
		history[i] = &model.SemesterRecord{
			Semester:    sem.Semester,
			SemesterGpa: sem.SemesterGpa,
			Courses:     courses,
		}
	}

	return &model.Student{
		StudentId: f.StudentID,
		Name:      f.Name,
		Gender:    f.Gender,
		AcademicData: &model.AcademicData{
			Faculty:         f.AcademicData.Faculty,
			StudyProgram:    f.AcademicData.StudyProgram,
			CurrentSemester: f.AcademicData.CurrentSemester,
		},
		AcademicHistory: history,
		CumulativeGpa:   f.CumulativeGpa,
	}
}

// ToStudentResponse generates n fake students and wraps them into a
// model.StudentResponse, the wire shape used for both the 1-entry and
// 100-entry scenarios.
func ToStudentResponse(n int) (*model.StudentResponse, error) {
	students, err := GenerateStudents(n)
	if err != nil {
		return nil, err
	}

	data := make([]*model.Student, len(students))
	for i, s := range students {
		data[i] = s.toProto()
	}
	return &model.StudentResponse{Data: data}, nil
}
