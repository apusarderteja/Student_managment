package postgres

import (
	"fmt"
	"log"

	"Student_managment/Project/storage"
	"golang.org/x/crypto/bcrypt"
)

const listQuery = `
	WITH tot AS (select count(*) as total FROM users
	WHERE
		deleted_at IS NULL
		AND (first_name ILIKE '%%' || $1 || '%%' OR last_name ILIKE '%%' || $1 || '%%' OR username ILIKE '%%' || $1 || '%%' OR email ILIKE '%%' || $1 || '%%'))
	SELECT *, tot.total as total FROM users
	LEFT JOIN tot ON TRUE
	WHERE
		deleted_at IS NULL
		AND (first_name ILIKE '%%' || $1 || '%%' OR last_name ILIKE '%%' || $1 || '%%' OR username ILIKE '%%' || $1 || '%%' OR email ILIKE '%%' || $1 || '%%')
		ORDER BY id DESC
		OFFSET $2
		LIMIT $3`

func (s PostgresStorage) ListUser(uf storage.UserFilter) ([]storage.User, error) {
	var listUser []storage.User
	if err := s.DB.Select(&listUser, listQuery, uf.SearchTerm, uf.Offset, uf.Limit); err != nil {
		log.Println(err)
		return nil, err
	}

	return listUser, nil
}



const insertQuery = `
		INSERT INTO users(
			first_name,
			last_name,
			username,
			email,
			password
		) VALUES (
			:first_name,
			:last_name,
			:username,
			:email,
			:password
		) RETURNING *;
	`

func (s PostgresStorage) CreateUser(u storage.User) (*storage.User, error) {
	stmt, err := s.DB.PrepareNamed(insertQuery)
	if err != nil {
		log.Fatalln(err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(hashPass)
	if err := stmt.Get(&u, u); err != nil {
		return nil, err
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("unable to insert user into db")
	}
	return &u, nil
}




const updateUserQuery = `
	UPDATE users SET
		first_name = :first_name,
		last_name = :last_name,
		status = :status
	WHERE id = :id AND deleted_at IS NULL RETURNING *;
	`

func (s PostgresStorage) UpdateUser(u storage.User) (*storage.User, error) {
	stmt, err := s.DB.PrepareNamed(updateUserQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

const getUserByIDQuery = `SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetUserByID(id string) (*storage.User, error) {
	var u storage.User
	if err := s.DB.Get(&u, getUserByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}


const getUserByUsernameQuery = `SELECT * FROM users WHERE username=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetUserByUsername(username string) (*storage.User, error) {
	var u storage.User
	if err := s.DB.Get(&u, getUserByUsernameQuery, username); err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

const deleteUserByIdQuery = `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) DeleteUserByID(id string) error {
	res, err := s.DB.Exec(deleteUserByIdQuery, id)
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

