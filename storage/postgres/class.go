package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
)

const insertQueryClass = `
		INSERT INTO class(
			class_name
		) VALUES (
			:class_name
		) RETURNING *;
	`

func (p PostgresStorage) CreateClass(s storage.Class) (*storage.Class, error) {

	stmt, err := p.DB.PrepareNamed(insertQueryClass)
	if err != nil {
		log.Fatalln(err)
	}

	if err := stmt.Get(&s, s); err != nil {
		log.Println(err)
		return nil, err
	}

	return &s, nil
}

const deleteClassByIdQuery = `UPDATE class SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) DeleteClassByID(id string) error {
	res, err := s.DB.Exec(deleteClassByIdQuery, id)
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

const getClassByIDQuery = `SELECT id,class_name FROM class WHERE deleted_at IS NULL;`

func (s PostgresStorage) GetClassByIDQuery() ([]storage.Class, error) {
	var u []storage.Class
	if err := s.DB.Select(&u, getClassByIDQuery); err != nil {
		log.Println(err)
		return nil, err
	}

	return u, nil
}

const updateClassQuery = `
	UPDATE class SET
		class_name = :class_name
	WHERE id = :id AND deleted_at IS NULL RETURNING *;
	`

func (s PostgresStorage) Updateclass(u storage.Class) (*storage.Class, error) {
	stmt, err := s.DB.PrepareNamed(updateClassQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}



const listclassQuery =  `SELECT * FROM class WHERE deleted_at IS NULL ORDER BY id DESC`

func (s PostgresStorage) ListClass() ([]storage.Class, error) {
	var listCl []storage.Class
	if err := s.DB.Select(&listCl, listclassQuery); err != nil {
		log.Println(err)
		return nil, err
	}

	return listCl, nil
}

const getclassIDByIDQuery = `SELECT * FROM class WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetclassIDByIDQuery(id string) (*storage.Class, error) {
	var u storage.Class
	if err := s.DB.Get(&u, getclassIDByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

// Class ALready Exists Check
func (p PostgresStorage) CheckClassExists(ClassName string) (bool, error) {
	var Alreadyexists bool
	err := p.DB.QueryRow(`SELECT exists(SELECT 1 FROM class WHERE class_name = $1 AND deleted_at IS NULL)`, ClassName).Scan(&Alreadyexists)
	if err != nil {
		return false, err
	}
	return Alreadyexists, nil
}