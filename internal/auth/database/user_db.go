package database

import (
	"context"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

type UserDTO struct {
	CreatedAt      time.Time  `json:"created_at"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Id             int        `json:"id"`
	LastConnexion  *time.Time `json:"last_connexion"`
	UpdatedAt      time.Time  `json:"updated_at"`
	HashedPassword string     `json:"hashed_password"`
}

func InsertUser(ctx context.Context, email, hasedPassword, fullName string) (id int, err error) {
	query := `
        INSERT INTO usr.users
            (full_name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err = globals.DB.QueryRow(query, fullName, email, hasedPassword).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	fmt.Printf("User inserted: %d\n", id)
	return
}

func GetUserByEmail(ctx context.Context, email string) (user UserDTO, err error) {
	query := `
       SELECT id, full_name, email, password, created_at, last_connexion FROM usr.users WHERE email=$1
    `

	err = globals.DB.QueryRow(query, email).Scan(&user.Id, &user.FullName, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.LastConnexion)
	if err != nil {
		return UserDTO{}, fmt.Errorf("failed to get user: %w", err)
	}

	return
}
