package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	start, err := time.Parse(Layout, dstart)
	if err != nil {
		return "", fmt.Errorf("bad start date: %w", err)
	}

	now = toDateOnly(now)

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("bad repeat format: d <days>")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", errors.New("bad repeat days: must be 1..400")
		}
		date := start.AddDate(0, 0, n)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, n)
		}
		return date.Format(Layout), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("bad repeat format: y")
		}
		date := addOneYearWithLeapRule(start)
		for !afterNow(date, now) {
			date = addOneYearWithLeapRule(date)
		}
		return date.Format(Layout), nil

	case "w", "m":
		return "", errors.New("unsupported repeat format")

	default:
		return "", errors.New("invalid repeat prefix")
	}
}

func toDateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()
	if dy != ny {
		return dy > ny
	}
	if dm != nm {
		return dm > nm
	}
	return dd > nd
}

func addOneYearWithLeapRule(t time.Time) time.Time {
	_, origM, origD := t.Date()
	next := t.AddDate(1, 0, 0)
	if origM == time.February && origD == 29 {
		if next.Month() == time.February && next.Day() == 28 {
			return time.Date(next.Year(), time.March, 1, 0, 0, 0, 0, next.Location())
		}
	}
	return next
}

func nextdateHandler(w http.ResponseWriter, r *http.Request) {
	qNow := strings.TrimSpace(r.FormValue("now"))
	qDate := strings.TrimSpace(r.FormValue("date"))
	qRep := strings.TrimSpace(r.FormValue("repeat"))

	var now time.Time
	var err error
	if qNow == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(Layout, qNow)
		if err != nil {
			http.Error(w, "bad now date", http.StatusBadRequest)
			return
		}
	}

	if qDate == "" {
		http.Error(w, "missing date", http.StatusBadRequest)
		return
	}
	if qRep == "" {
		http.Error(w, "missing repeat", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, qDate, qRep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}
