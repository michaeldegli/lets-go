package models

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("models: invalid user credentials")
)

// Database ...
type Database struct {
	*sql.DB
}

// InsertSnippet ...
func (db *Database) InsertSnippet(title, content string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES($1, $2, CURRENT_DATE, CURRENT_DATE + INTERVAL '100000 seconds') returning id`

	result, err := db.Query(stmt, title, content)

	if err != nil {
		return 0, err
	}

	if err := result.Next(); !err {
		return 0, errors.New("no id was generated")
	}
	var id int
	err = result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

// LatestSnippets ...
func (db *Database) LatestSnippets() (Snippets, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_DATE ORDER BY created DESC LIMIT 10`

	rows, err := db.Query(stmt)

	if err != nil {
		return nil, err
	}

	// This should come after we check for an error
	// or the rows object could be nil
	defer rows.Close()

	snippets := Snippets{}

	for rows.Next() {

		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}

// GetSnippet ...
func (db *Database) GetSnippet(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_DATE AND id = $1`

	row := db.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

// InsertUser ...
func (db *Database) InsertUser(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	fmt.Println("password is ", string(hashedPassword), "original is ", password)

	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, password, created) VALUES ($1, $2, $3, CURRENT_DATE)"

	_, err = db.Exec(stmt, name, email, string(hashedPassword))

	if err != nil {
		return err
	}

	return err

}

// VerifyUser ...
func (db *Database) VerifyUser(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	row := db.QueryRow("SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil

}
