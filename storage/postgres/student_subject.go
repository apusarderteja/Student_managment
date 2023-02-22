package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
)

const insertstudentsubjectQuery = `
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

func (p PostgresStorage) CreateStudentSubject(s storage.StudentSubject) (*storage.StudentSubject, error) {

	stmt, err := p.DB.PrepareNamed(insertstudentsubjectQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}


const getStudentSubjectByStudentQuery = `SELECT * FROM student_subject WHERE student_id=$1 AND deleted_at IS NULL;`

func (s PostgresStorage) GetStudentSubjectByStudentID(id string) ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentSubjectByStudentQuery,id); err != nil {
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




const deleteStudentSubjectByIdQuery = `UPDATE student_subject SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) DeleteStudentSubjectByID(id string) error {
	res, err := s.DB.Exec(deleteStudentSubjectByIdQuery, id)
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


const getStudentIdBySubjectID = `SELECT  sts.* , sub.subject_name
FROM student_subject as sts
INNER JOIN subjects as sub ON sts.subject_id = sub.id
WHERE sts.student_id =$1;`

func (s PostgresStorage) GetStudentIdBySubjectID(student_id string) ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentIdBySubjectID,student_id); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil
}
