package seeder

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

const facultyName = "School of Computing"

// fakeCourse is the gofakeit template for a single course entry.
// Code is not generated here; the mapper derives it from position instead.
type fakeCourse struct {
	Name    string  `fake:"{randomstring:[Introduction to Algorithms,Data Structures,Computer Networks,Database Systems,Operating Systems,Software Engineering,Linear Algebra,Discrete Mathematics,Artificial Intelligence,Computer Architecture]}"`
	Credits int32   `fake:"{number:2,4}"`
	Score   float32 `fake:"{float32range:60,100}"`
}

// fakeSemesterRecord is the gofakeit template for one semester's academic
// record. Semester is overwritten manually in GenerateStudents with the
// correct sequential number (1, 2, ...).
type fakeSemesterRecord struct {
	Semester    int32
	SemesterGpa float32      `fake:"{float32range:2.5,4.0}"`
	Courses     []fakeCourse `fakesize:"3"`
}

// fakeAcademicData is the gofakeit template for a student's program information.
// Faculty is overwritten manually in GenerateStudents since every student
// belongs to the same fixed faculty.
type fakeAcademicData struct {
	Faculty         string
	StudyProgram    string `fake:"{randomstring:[Informatics,Software Engineering,Information Technology,Data Science]}"`
	CurrentSemester int32  `fake:"{number:1,8}"`
}

// fakeStudent is the gofakeit template for a single student record.
// StudentID is overwritten manually in GenerateStudents with a
// sequential, deterministic ID.
type fakeStudent struct {
	StudentID       string
	Name            string `fake:"{name}"`
	Gender          string `fake:"{gender}"`
	AcademicData    fakeAcademicData
	AcademicHistory []fakeSemesterRecord `fakesize:"2"`
	CumulativeGpa   float32              `fake:"{float32range:2.5,4.0}"`
}

// GenerateStudents produces n fake student records. StudentID, Faculty,
// and Semester are assigned sequentially/fixed after gofakeit fills the
// rest of the struct, since those fields must not be random.
func GenerateStudents(n int) ([]fakeStudent, error) {
	students := make([]fakeStudent, n)
	for i := range students {
		if err := gofakeit.Struct(&students[i]); err != nil {
			return nil, fmt.Errorf("generate student %d: %w", i, err)
		}
		students[i].StudentID = fmt.Sprintf("13%08d", i+1)
		students[i].AcademicData.Faculty = facultyName
		for j := range students[i].AcademicHistory {
			students[i].AcademicHistory[j].Semester = int32(j + 1)
		}
	}
	return students, nil
}
