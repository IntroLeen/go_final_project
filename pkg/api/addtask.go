package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeErr(w, http.StatusBadRequest, "empty title")
		return
	}

	now := toDateOnly(time.Now())

	t.Date = strings.TrimSpace(t.Date)
	if t.Date == "" {
		t.Date = now.Format(Layout)
	}

	parsed, err := time.Parse(Layout, t.Date)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad date format")
		return
	}
	parsed = toDateOnly(parsed)

	t.Repeat = strings.TrimSpace(t.Repeat)
	var next string
	if t.Repeat != "" {
		next, err = NextDate(now, t.Date, t.Repeat)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if parsed.Before(now) {
		if t.Repeat == "" {
			t.Date = now.Format(Layout)
		} else {
			t.Date = next
		}
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]string{"id": int64ToString(id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	t.ID = strings.TrimSpace(t.ID)
	if t.ID == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}

	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeErr(w, http.StatusBadRequest, "empty title")
		return
	}

	now := toDateOnly(time.Now())

	t.Date = strings.TrimSpace(t.Date)
	if t.Date == "" {
		t.Date = now.Format(Layout)
	}
	parsed, err := time.Parse(Layout, t.Date)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad date format")
		return
	}
	parsed = toDateOnly(parsed)

	t.Repeat = strings.TrimSpace(t.Repeat)
	var next string
	if t.Repeat != "" {
		next, err = NextDate(now, t.Date, t.Repeat)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if parsed.Before(now) {
		if t.Repeat == "" {
			t.Date = now.Format(Layout)
		} else {
			t.Date = next
		}
	}

	if err := db.UpdateTask(&t); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, struct{}{})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode error: %v", err)
	}

}

func int64ToString(v int64) string {
	return strconvFormatInt(v)
}

func strconvFormatInt(v int64) string {
	const digits = "0123456789"
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = digits[v%10]
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func checkDate(task *db.Task) error {
	if task == nil {
		return errors.New("nil task")
	}
	return nil
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, struct{}{})
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		log.Printf("writeErr encode error: %v", err)
	}
}
