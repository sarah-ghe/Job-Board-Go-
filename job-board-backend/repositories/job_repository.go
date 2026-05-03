package repositories

import (
	"database/sql"
	"fmt"
	"job-board/models"
)

type PostgresJobRepository struct {
	DB *sql.DB
}

func (r *PostgresJobRepository) Create(job *models.Job) error {

	query := `
	INSERT INTO jobs (title, description, user_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	return r.DB.QueryRow(
		query,
		job.Title,
		job.Description,
		job.UserID,
	).Scan(&job.ID)
}

func (r *PostgresJobRepository) GetAll() ([]models.Job, error) {

	query := `SELECT id, title, description FROM jobs`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job

		err := rows.Scan(&job.ID, &job.Title, &job.Description)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *PostgresJobRepository) GetByUserID(userID int) ([]models.Job, error) {

	query := "SELECT id, title, description, user_id FROM jobs WHERE user_id = $1"

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.UserID)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *PostgresJobRepository) GetByID(id int) (*models.Job, error) {

	query := "SELECT id, title, description, user_id FROM jobs WHERE id = $1"

	job := &models.Job{}

	err := r.DB.QueryRow(query, id).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.UserID,
	)

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *PostgresJobRepository) Update(job *models.Job, userID int) error {

	query := `
	UPDATE jobs
	SET title=$1, description=$2
	WHERE id=$3 AND user_id=$4
	`

	result, err := r.DB.Exec(query, job.Title, job.Description, job.ID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("not found or not authorized")
	}

	return nil
}

func (r *PostgresJobRepository) Delete(id int, userID int) error {
	query := `
		DELETE FROM jobs
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("not found or not authorized")
	}

	return nil
}