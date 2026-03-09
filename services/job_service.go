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
	Update(id string, job *models.Job) error
	Delete(id string) error
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



func (s *JobService) UpdateJob(id string, job *models.Job) error {

	if job.Title == "" {
		return errors.New("title is required")
	}

	return s.Repo.Update(id, job)
}




func (s *JobService) DeleteJob(id string) error {

	return s.Repo.Delete(id)
}