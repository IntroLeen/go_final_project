package api

import (
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]string{"error": "missing id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	rep := strings.TrimSpace(task.Repeat)
	if rep == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	now := toDateOnly(time.Now())
	next, err := NextDate(now, task.Date, rep)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, struct{}{})
}
