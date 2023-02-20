package postgres

import (
	// "fmt"
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


const getStudentSubjectByStudentQuery = `SELECT * FROM student_subject WHERE student_id=$1;`

func (s PostgresStorage) GetStudentSubjectByStudentID(id string) ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentSubjectByStudentQuery,id); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil
}



const updateStudentSubjectQuery = `
	UPDATE student_subject SET
		student_id = :student_id,
		subject_id = :subject_id,
		marks = :marks
	WHERE id = :id
		 RETURNING *;
	`

func (s PostgresStorage) UpdateStudentSubject(u storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := s.DB.PrepareNamed(updateStudentSubjectQuery)
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




// const listclassQuery =  `SELECT * FROM class WHERE deleted_at IS NULL ORDER BY id DESC`

// func (s PostgresStorage) ListClass() ([]storage.Class, error) {
// 	var listCl []storage.Class
// 	if err := s.DB.Select(&listCl, listclassQuery); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return listCl, nil
// }

// const getclassIDByIDQuery = `SELECT * FROM class WHERE id=$1 AND deleted_at IS NULL`

// func (s PostgresStorage) GetclassIDByIDQuery(id string) (*storage.Class, error) {
// 	var u storage.Class
// 	if err := s.DB.Get(&u, getclassIDByIDQuery, id); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return &u, nil
// }



const getStuIdBySubQuery = `SELECT * FROM subjects WHERE class_id= $1;`

func (s PostgresStorage) GetSubIdBYID(classID int) ([]storage.Subject, error) {
	var u []storage.Subject
	if err := s.DB.Select(&u, getStuIdBySubQuery,classID); err != nil {
		log.Println(err)
		return nil, err
	}
	
	return u, nil
}

// func (s PostgresStorage) ListStudentQuery() ([]storage.Student, error) {
// 	var liststudent []storage.Student
// 	if err := s.DB.Select(&liststudent, listStudentQuery); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return liststudent, nil
// }
