package repositories

import "job-board/models"

type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) Create(user *models.User) error {

	query := `	
	INSERT INTO users (email, password)
	VALUES ($1, $2)
	RETURNING id
	`	

	return r.DB.QueryRow(
		query,
		user.Email,
		user.Password,
	).Scan(&user.ID)
}


func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {

	query := `
	SELECT id, email, password
	FROM users
	WHERE email=$1
	`

	user := &models.User{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}