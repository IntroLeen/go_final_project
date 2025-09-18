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
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	const limit = 50
	tasks, err := db.Tasks(limit)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
