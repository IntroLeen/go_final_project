package api

import (
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}

	rep := strings.TrimSpace(task.Repeat)
	if rep == "" {
		if err := db.DeleteTask(id); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	now := toDateOnly(time.Now())
	next, err := NextDate(now, task.Date, rep)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, struct{}{})
}
