package services

import (
	"job-board/models"
	"testing"
)

type FakeJobRepository struct {
	Jobs []models.Job
}

func (f *FakeJobRepository) Create(job *models.Job) error {
	job.ID = len(f.Jobs) + 1
	f.Jobs = append(f.Jobs, *job)
	return nil
}

func (f *FakeJobRepository) GetAll() ([]models.Job, error) {
	return f.Jobs, nil
}

func (f *FakeJobRepository) Update(id string, job *models.Job) error {
	return nil
}

func (f *FakeJobRepository) Delete(id string) error {
	return nil
}


func TestCreateJob_TitleRequired(t *testing.T) {

	fakeRepo := &FakeJobRepository{}

	service := JobService{
		Repo: fakeRepo,
	}

	job := models.Job{
		Title: "",
	}

	err := service.CreateJob(&job)

	if err == nil {
		t.Errorf("expected error when title is empty")
	}
}


func TestCreateJob_Success(t *testing.T) {

	fakeRepo := &FakeJobRepository{}

	service := JobService{
		Repo: fakeRepo,
	}

	job := models.Job{
		Title: "Go Developer",
		Description: "Build APIs",
	}

	err := service.CreateJob(&job)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if job.ID == 0 {
		t.Errorf("expected job ID to be set")
	}

	if len(fakeRepo.Jobs) != 1 {
		t.Errorf("expected job to be stored in repository")
	}
}