package handlers

import (
	"encoding/json"
	"net/http"
)

type AddUserTask struct {
	Email  string `json:"email"`
	Task   string `json:"task"`
	Status bool   `json:"status"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	// Set header to json content, otherwise data appear as plain text
	w.Header().Set("Content-Type", "application/json")

	// decode the requested data to 'newTask'
	var newTask AddUserTask

	// NewDecoder: Decode FROM WHERE
	// Decode: WHERE TO STORE the decoded data
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	db := OpenConnection()
	defer db.Close()
	var userId string
	sql := `Select user_id from users where email= $1 `
	// getUser := `SELECT user_id FROM users WHERE email = $1;`
	err = db.QueryRow(sql, newTask.Email).Scan(&userId)

	sqlStatement := `INSERT INTO tasks (task, status, user_uuid) VALUES ($1, $2, $3) RETURNING id, task, status;`

	// retrieve the task after creation from the database and store its details in 'updatedTask'
	var updatedTask Item
	err = db.QueryRow(sqlStatement, newTask.Task, newTask.Status, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	// gives the new task as the output
	_ = json.NewEncoder(w).Encode(updatedTask)
}
