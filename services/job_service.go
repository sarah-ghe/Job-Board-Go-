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
	Update(job *models.Job) error
	Delete(id int) error
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

func (s *JobService) GetJobsByUser(userID int) ([]models.Job, error) {

	return s.Repo.GetByUserID(userID)
}


func (s *JobService) UpdateJob(id int, userID int, job *models.Job) error {

	existingJob, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	// Ownership check
	if existingJob.UserID != userID {
		return errors.New("forbidden")
	}

	job.ID = id
	job.UserID = userID

	return s.Repo.Update(job)
}

func (s *JobService) DeleteJob(id int, userID int) error {

	existingJob, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	// 🔥 Ownership check
	if existingJob.UserID != userID {
		return errors.New("forbidden")
	}

	return s.Repo.Delete(id)
}
