package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	const limit = 50
	tasks, err := db.Tasks(limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
