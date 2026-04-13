package repositories

import (
	"database/sql"
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

	query := `SELECT id, title, description, user_id FROM jobs`

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



func (r *PostgresJobRepository) Update(id string, job *models.Job) error {

	query := `
	UPDATE jobs
	SET title=$1, description=$2
	WHERE id=$3
	`

	result, err := r.DB.Exec(query, job.Title, job.Description, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}


func (r *PostgresJobRepository) Delete(id string) error {

	query := `DELETE FROM jobs WHERE id=$1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}