package api

import "net/http"

const Layout = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextdateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", doneHandler)
}
