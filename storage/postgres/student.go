package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
)

const insertStudentQuery = `
		INSERT INTO students(
			class_id,
			first_name,
			last_name,
			roll
		) VALUES (
			:class_id,
			:first_name,
			:last_name,
			:roll
		) RETURNING *;
	`

func (p PostgresStorage) CreateStudent(s storage.Student) (*storage.Student, error) {

	stmt, err := p.DB.PrepareNamed(insertStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const getStudentByIdQuery = `SELECT * FROM students WHERE id=$1;`

func (s PostgresStorage) GetStudentByID(id string) (*storage.StudentSubject, error) {
	var u storage.StudentSubject
	if err := s.DB.Get(&u, getStudentByIdQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil

}

const listStudentQuery = `SELECT students.id, class.class_name, students.first_name ,students.last_name, students.roll
FROM students
INNER JOIN class ON class.id = students.class_id
WHERE students.deleted_at IS NULL AND class.deleted_at IS NULL;`

func (s PostgresStorage) ListStudentQuery() ([]storage.Student, error) {
	var liststudent []storage.Student
	if err := s.DB.Select(&liststudent, listStudentQuery); err != nil {
		log.Println(err)
		return nil, err
	}

	return liststudent, nil
}

const deleteStudentByIdQuery = `UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL;`

func (s PostgresStorage) DeleteStudentByIdQuery(id string) error {
	res, err := s.DB.Exec(deleteStudentByIdQuery, id)
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

const getstudentIDByIDQuery = `SELECT * FROM students WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetstudentIDByIDQuery(id string) (*storage.Student, error) {
	var u storage.Student
	if err := s.DB.Get(&u, getstudentIDByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

const updateStudentQuery = `
	UPDATE students SET
	first_name = :first_name ,
	last_name = :last_name,
	roll = :roll
	WHERE id = :id AND deleted_at IS NULL RETURNING *;
	`

func (s PostgresStorage) UpdateStudent(u storage.Student) (*storage.Student, error) {
	stmt, err := s.DB.PrepareNamed(updateStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}
