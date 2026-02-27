package main

import (
	"job-board/config"
	"job-board/handlers"
	"net/http"
)

func main() {

	config.ConnectDB()

	http.HandleFunc("/jobs", handlers.JobsHandler)

	http.ListenAndServe(":8080", nil)
}
