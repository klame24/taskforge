package repository

import (
	"database/sql"
	"taskforge/internal/models"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	// GetByID(id int) (*models.User, error)
	// GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	query := `
		INSERT INTO users (first_name, last_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at FROM users WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(
		query,
		email,
	).Scan(
		&user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at FROM users WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// CREATE TABLE users (
//     id SERIAL PRIMARY KEY,
//     first_name VARCHAR(255) NOT NULL,
//     last_name VARCHAR(255) NOT NULL,
//     email VARCHAR(255) UNIQUE NOT NULL,
//     password_hash TEXT NOT NULL,
//     created_at TIMESTAMP,
//     updated_at TIMESTAMP
// );
