package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeJSON(w, map[string]string{"error": "empty title"})
		return
	}

	now := toDateOnly(time.Now())

	t.Date = strings.TrimSpace(t.Date)
	if t.Date == "" {
		t.Date = now.Format(Layout)
	}

	parsed, err := time.Parse(Layout, t.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "bad date format"})
		return
	}
	parsed = toDateOnly(parsed)

	t.Repeat = strings.TrimSpace(t.Repeat)
	var next string
	if t.Repeat != "" {
		next, err = NextDate(now, t.Date, t.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
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
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": int64ToString(id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	t.ID = strings.TrimSpace(t.ID)
	if t.ID == "" {
		writeJSON(w, map[string]string{"error": "missing id"})
		return
	}

	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeJSON(w, map[string]string{"error": "empty title"})
		return
	}

	now := toDateOnly(time.Now())

	t.Date = strings.TrimSpace(t.Date)
	if t.Date == "" {
		t.Date = now.Format(Layout)
	}
	parsed, err := time.Parse(Layout, t.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "bad date format"})
		return
	}
	parsed = toDateOnly(parsed)

	t.Repeat = strings.TrimSpace(t.Repeat)
	var next string
	if t.Repeat != "" {
		next, err = NextDate(now, t.Date, t.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
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
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, struct{}{})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
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
		writeJSON(w, map[string]string{"error": "method not allowed"})
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, map[string]string{"error": "missing id"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, struct{}{})
}
