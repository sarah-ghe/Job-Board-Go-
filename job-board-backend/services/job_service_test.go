package services

import (
	"database/sql"
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

func (f *FakeJobRepository) Update(job *models.Job) error {
	f.UpdatedJob = job
	for i := range f.Jobs {
		if f.Jobs[i].ID == job.ID {
			f.Jobs[i] = *job
			return nil
		}
	}
	return nil
}

func (f *FakeJobRepository) Delete(id int) error {
	f.DeletedID = id
	for i := range f.Jobs {
		if f.Jobs[i].ID == id {
			f.Jobs = append(f.Jobs[:i], f.Jobs[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestCreateJob_TitleRequired(t *testing.T) {
	fakeRepo := &FakeJobRepository{}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: ""}
	if err := service.CreateJob(&job); err == nil {
		t.Fatalf("expected error when title is empty")
	}
}

func TestCreateJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: "Go Developer", Description: "Build APIs"}
	if err := service.CreateJob(&job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if job.ID == 0 {
		t.Fatalf("expected job ID to be set")
	}
	if len(fakeRepo.Jobs) != 1 {
		t.Fatalf("expected job to be stored in repository, got %d", len(fakeRepo.Jobs))
	}
}

func TestGetJobs_ReturnsAllJobs(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, Title: "Job1"}, {ID: 2, Title: "Job2"}}}
	service := JobService{Repo: fakeRepo}

	jobs, err := service.GetJobs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestGetJobsByUser_FiltersByUserID(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}, Title: "Job1"}, {ID: 2, UserID: sql.NullInt64{Int64: 2, Valid: true}, Title: "Job2"}}}
	service := JobService{Repo: fakeRepo}

	jobs, err := service.GetJobsByUser(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job for user 1, got %d", len(jobs))
	}
	if jobs[0].UserID.Int64 != 1 {
		t.Fatalf("expected user ID 1, got %d", jobs[0].UserID.Int64)
	}
}

func TestUpdateJob_OwnerMismatchReturnsForbidden(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, UserID: 2, Title: "Original"}}}
	service := JobService{Repo: fakeRepo}

	job := models.Job{Title: "Updated Title", Description: "Updated"}
	err := service.UpdateJob(1, 1, &job)
	if err == nil {
		t.Fatalf("expected forbidden error for non-owner")
	}
}

func TestUpdateJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}, Title: "Original", Description: "Old"}}}
	service := JobService{Repo: fakeRepo}

	updated := models.Job{Title: "Updated Title", Description: "Updated"}
	if err := service.UpdateJob(1, 1, &updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fakeRepo.UpdatedJob == nil {
		t.Fatalf("expected Update to be called")
	}
	if fakeRepo.UpdatedJob.ID != 1 || fakeRepo.UpdatedJob.UserID.Int64 != 1 {
		t.Fatalf("expected updated job to have ID 1 and UserID 1, got ID %d UserID %d", fakeRepo.UpdatedJob.ID, fakeRepo.UpdatedJob.UserID.Int64)
	}
}

func TestDeleteJob_OwnerMismatchReturnsForbidden(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, UserID: sql.NullInt64{Int64: 2, Valid: true}, Title: "Original"}}}
	service := JobService{Repo: fakeRepo}

	if err := service.DeleteJob(1, 1); err == nil {
		t.Fatalf("expected forbidden error for non-owner")
	}
}

func TestDeleteJob_Success(t *testing.T) {
	fakeRepo := &FakeJobRepository{Jobs: []models.Job{{ID: 1, UserID: sql.NullInt64{Int64: 1, Valid: true}, Title: "Original"}}}
	service := JobService{Repo: fakeRepo}

	if err := service.DeleteJob(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fakeRepo.DeletedID != 1 {
		t.Fatalf("expected Delete to be called with ID 1, got %d", fakeRepo.DeletedID)
	}
	if len(fakeRepo.Jobs) != 0 {
		t.Fatalf("expected job list to be empty after delete, got %d", len(fakeRepo.Jobs))
	}
}
