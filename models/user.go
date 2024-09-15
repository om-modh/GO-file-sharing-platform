// models/user.go
package models

import (
    "database/sql"
    "errors"
    "time"
)

type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

func CreateUser(db *sql.DB, email, passwordHash string) error {
    query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
    _, err := db.Exec(query, email, passwordHash)
    return err
}

func FindUserByEmail(db *sql.DB, email string) (*User, error) {
    var user User
    row := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", email)
    err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}
