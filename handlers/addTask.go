package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AddUserTask struct {
	Email  string `json:"email"`
	Task   string `json:"task"`
	Status bool   `json:"status"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	// Set header to json content, otherwise data appear as plain text
	// w.Header().Set("Content-Type", "application/json")
	if !UserLogged(w, r) {
		w.Write([]byte(fmt.Sprintf("Cannot Add Task, User Not Logged In")))
		return
	}
	if r.Method == "GET" {
		http.ServeFile(w, r, "addTask.html")
	} else {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		email := r.FormValue("email")
		task := r.FormValue("task")
		status := r.FormValue("status")
		// decode the requested data to 'newTask'
		var newTask AddUserTask
		newTask.Email = email
		newTask.Task = task
		newTask.Status = status == "true"
		// NewDecoder: Decode FROM WHERE
		// Decode: WHERE TO STORE the decoded data
		// err := json.NewDecoder(r.Body).Decode(&newTask)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	panic(err)
		// }
		var err error
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
}
