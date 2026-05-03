package services

import (
	"errors"
	"job-board/models"
)

type JobService struct {
	Repo JobRepository
}

type JobRepository interface {
	Create(job *models.Job) error
	GetAll() ([]models.Job, error)
	Update(job *models.Job, userID int) error
	Delete(id int, userID int) error
	GetByUserID(userID int) ([]models.Job, error)
	GetByID(id int) (*models.Job, error)
}

func (s *JobService) CreateJob(job *models.Job) error {

	if job.Title == "" {
		return errors.New("title is required")
	}

	return s.Repo.Create(job)
}

func (s *JobService) GetJobs() ([]models.Job, error) {

	return s.Repo.GetAll()
}

func (s *JobService) GetByID(id int) (*models.Job, error) {
	return s.Repo.GetByID(id)
}

func (s *JobService) GetJobsByUser(userID int) ([]models.Job, error) {

	return s.Repo.GetByUserID(userID)
}

func (s *JobService) UpdateJob(id int, userID int, job *models.Job) error {

	// attach ID from URL
	job.ID = id

	// delegate ownership enforcement to repo
	return s.Repo.Update(job, userID)
}

func (s *JobService) DeleteJob(id int, userID int) error {
	return s.Repo.Delete(id, userID)
}
