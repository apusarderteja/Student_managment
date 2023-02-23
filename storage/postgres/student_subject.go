package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
)

const insertstudentMarkQuery = `
		INSERT INTO student_subject(
			student_id,
			subject_id,
			marks
		) VALUES (
			:student_id,
			:subject_id,
			:marks
		) RETURNING *;
	`

func (p PostgresStorage) InsertstudentMarkQuery(s storage.StudentSubject) (*storage.StudentSubject, error) {

	stmt, err := p.DB.PrepareNamed(insertstudentMarkQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}


// const getStudentSubjectByStudentQuery = `SELECT * FROM student_subject WHERE student_id=$1 AND deleted_at IS NULL;`
const getStudentSubjectByStudentQuery = `SELECT student_subject.student_id,student_subject.marks,  students.first_name ,students.last_name, students.roll
FROM student_subject
INNER JOIN students ON student_subject.student_id = students.id
WHERE student_subject.deleted_at IS NULL AND students.deleted_at IS NULL;`

func (s PostgresStorage) GetStudentSubjectByStudentID() ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentSubjectByStudentQuery); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil
}



const updateStudentMarkQuery = `
	UPDATE student_subject SET
		student_id = :student_id,
		subject_id = :subject_id,
		marks = :marks
	WHERE student_id = :student_id AND subject_id = :subject_id
		 RETURNING *;
	`

func (s PostgresStorage) UpdateStudentMark(u storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := s.DB.PrepareNamed(updateStudentMarkQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}


const getStudentProfileQuery = `SELECT  student_subject.* , subjects.subject_name ,students.first_name ,students.last_name ,students.roll 
FROM student_subject
INNER JOIN students ON student_subject.student_id = students.id
INNER JOIN subjects ON student_subject.subject_id = subjects.id
WHERE student_subject.student_id =$1;`


func (s PostgresStorage) GetStudentProfileQuery(StudentId string) ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentProfileQuery , StudentId); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil

}
	




const deleteStudentSubjectByIdQuery = `UPDATE student_subject SET deleted_at = CURRENT_TIMESTAMP WHERE student_id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) DeleteStudentSubjectByID(StudentId string) error {
	res, err := s.DB.Exec(deleteStudentSubjectByIdQuery, StudentId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if rowCount <= 0 {
		return fmt.Errorf("unable to delete user")
	}

	return nil
}






const getStuIdBySubQuery = `SELECT * FROM subjects WHERE class_id= $1 AND deleted_at IS NULL;`

func (s PostgresStorage) GetSubIdBYID(classID int) ([]storage.Subject, error) {
	var u []storage.Subject
	if err := s.DB.Select(&u, getStuIdBySubQuery,classID); err != nil {
		log.Println(err)
		return nil, err
	}
	
	return u, nil
}


const getStudentIdBySubjectID = `SELECT  student_subject.* , subjects.subject_name
FROM student_subject 
INNER JOIN subjects 
ON student_subject.subject_id = subjects.id
WHERE student_subject.student_id =$1;`

func (s PostgresStorage) GetStudentIdBySubjectID(student_id string) ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentIdBySubjectID,student_id); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil
}
