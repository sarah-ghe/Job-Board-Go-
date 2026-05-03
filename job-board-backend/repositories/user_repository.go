package repositories

import (
	"database/sql"
	"errors"
	"job-board/models"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) Create(user *models.User) error {

	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	return r.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID)
}

func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {

	query := `
	SELECT id, username, email, password
	FROM users
	WHERE email=$1
	`

	user := &models.User{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByID(id int) (*models.User, error) {

	query := `
	SELECT id, username, email, password
	FROM users
	WHERE id = $1
	`

	user := &models.User{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
