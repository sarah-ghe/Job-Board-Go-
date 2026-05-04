package handlers

import "job-board/models"

func toJobResponse(job models.Job) JobResponse {
	userID := 0
	if job.UserID.Valid {
		userID = int(job.UserID.Int64)
	}

	return JobResponse{
		ID:          job.ID,
		Title:       job.Title,
		Description: job.Description,
		UserID:      userID,
	}
}