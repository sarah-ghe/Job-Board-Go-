package services

import (
	"database/sql"
	"errors"
	"job-board/models"
	"testing"
)

type FakeJobRepository struct {
	Jobs       []models.Job
	UpdatedJob *models.Job
	DeletedID  int
}

func (f *FakeJobRepository) Create(job *models.Job) error {
	job.ID = len(f.Jobs) + 1
	f.Jobs = append(f.Jobs, *job)
	return nil
}

func (f *FakeJobRepository) GetAll() ([]models.Job, error) {
	return f.Jobs, nil
}

func (f *FakeJobRepository) GetByID(id int) (*models.Job, error) {
	for i := range f.Jobs {
		if f.Jobs[i].ID == id {
			return &f.Jobs[i], nil
		}
	}
	return nil, nil
}

func (f *FakeJobRepository) GetByUserID(userID int) ([]models.Job, error) {
	var userJobs []models.Job
	for _, job := range f.Jobs {
		if job.UserID.Valid && job.UserID.Int64 == int64(userID) {
			userJobs = append(userJobs, job)
		}
	}
	return userJobs, nil
}

func (f *FakeJobRepository) Update(job *models.Job, userID int) error {
	for i := range f.Jobs {
		if f.Jobs[i].ID == job.ID {

			// simulate ownership check
			if !f.Jobs[i].UserID.Valid || f.Jobs[i].UserID.Int64 != int64(userID) {
				return errors.New("forbidden")
			}

			job.UserID = f.Jobs[i].UserID
			f.Jobs[i] = *job
			f.UpdatedJob = job
			return nil
		}
	}
	return errors.New("not found")
}

func (f *FakeJobRepository) Delete(id int, userID int) error {
	for i := range f.Jobs {
		if f.Jobs[i].ID == id {

			if !f.Jobs[i].UserID.Valid || f.Jobs[i].UserID.Int64 != int64(userID) {
				return errors.New("forbidden")
			}

			f.DeletedID = id
			f.Jobs = append(f.Jobs[:i], f.Jobs[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

//tests

func TestCreateJob_TitleRequired(t *testing.T) {
	service := JobService{Repo: &FakeJobRepository{}}

	job := models.Job{Title: ""}
	err := service.CreateJob(&job)

	if err == nil {
		t.Fatal("expected error when title is empty")
	}
}

func TestCreateJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: "Go Developer", Description: "Build APIs"}

	err := service.CreateJob(&job)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID == 0 {
		t.Fatal("expected job ID to be set")
	}
	if len(fakeRepo.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(fakeRepo.Jobs))
	}
}

func TestGetJobs_ReturnsAllJobs(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, Title: "Job1"},
			{ID: 2, Title: "Job2"},
		},
	}
	service := JobService{Repo: fakeRepo}

	jobs, err := service.GetJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestGetJobsByUser_FiltersCorrectly(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}, Title: "Job1"},
			{ID: 2, UserID: sql.NullInt64{Int64: 2, Valid: true}, Title: "Job2"},
		},
	}
	service := JobService{Repo: fakeRepo}

	jobs, _ := service.GetJobsByUser(1)

	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
}

func TestUpdateJob_Forbidden(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, UserID: sql.NullInt64{Int64: 2, Valid: true}, Title: "Original"},
		},
	}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: "Updated"}
	err := service.UpdateJob(1, 1, &job)

	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestUpdateJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}, Title: "Old"},
		},
	}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: "Updated"}
	err := service.UpdateJob(1, 1, &job)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fakeRepo.UpdatedJob == nil {
		t.Fatal("expected update to be called")
	}
}

func TestDeleteJob_Forbidden(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, UserID: sql.NullInt64{Int64: 2, Valid: true}},
		},
	}
	service := JobService{Repo: fakeRepo}

	err := service.DeleteJob(1, 1)

	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestDeleteJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{
		Jobs: []models.Job{
			{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}},
		},
	}
	service := JobService{Repo: fakeRepo}

	err := service.DeleteJob(1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fakeRepo.Jobs) != 0 {
		t.Fatal("expected job to be deleted")
	}
}
