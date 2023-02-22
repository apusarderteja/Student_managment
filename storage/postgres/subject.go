package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
)


// insert subject query
const insertQuerySubject = `
		INSERT INTO subjects(
			class_id,
			subject_name
		) VALUES (
			:class_id,
			:subject_name
		) RETURNING *;
	`

func (p PostgresStorage) AddSubject(s storage.Subject) (*storage.Subject, error) {

	stmt, err := p.DB.PrepareNamed(insertQuerySubject)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}


// list subject query

const listSubjectQuery =  `SELECT subjects.id, class.class_name, subjects.subject_name 
FROM subjects
INNER JOIN class ON class.id = subjects.class_id WHERE subjects.deleted_at IS NULL AND class.deleted_at IS NULL;`

func (s PostgresStorage) ListSubjectQuery() ([]storage.Subject, error) {
	var listsubject []storage.Subject
	if err := s.DB.Select(&listsubject, listSubjectQuery); err != nil {
		log.Println(err)
		return nil, err
	}

	return listsubject, nil
}


// for edit subject id query

const getsubjectIDByIDQuery = `SELECT * FROM subjects WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetsubjectIDByIDQuery(id string) (*storage.Subject, error) {
	var u storage.Subject
	if err := s.DB.Get(&u, getsubjectIDByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

//  subject update query

const updateSubjectQuery = `
	UPDATE subjects SET
	subject_name = :subject_name
	WHERE id = :id AND deleted_at IS NULL RETURNING *;
	`

func (s PostgresStorage) UpdateSubjectFUNC(u storage.Subject) (*storage.Subject, error) {
	stmt, err := s.DB.PrepareNamed(updateSubjectQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

//  subject delete query

const deleteSubjectByIdQuery = `UPDATE subjects SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) DeleteSubjectByIdQuery(id string) error {
	res, err := s.DB.Exec(deleteSubjectByIdQuery, id)
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