package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	const q = `INSERT INTO scheduler(date, title, comment, repeat) VALUES(?,?,?,?)`
	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}
	const q = `
SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat
FROM scheduler
ORDER BY date ASC
LIMIT ?;
`
	rows, err := DB.Query(q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = make([]*Task, 0)
	}
	return out, nil
}

func GetTask(id string) (*Task, error) {
	const q = `
SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat
FROM scheduler
WHERE id = ?;
`
	var t Task
	err := DB.QueryRow(q, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	const q = `
UPDATE scheduler
SET date = ?, title = ?, comment = ?, repeat = ?
WHERE id = ?;
`
	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}
func DeleteTask(id string) error {
	const q = `DELETE FROM scheduler WHERE id = ?;`
	res, err := DB.Exec(q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	const q = `UPDATE scheduler SET date = ? WHERE id = ?;`
	res, err := DB.Exec(q, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}
	return nil
}
