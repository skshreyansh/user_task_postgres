package handlers

import (
	// Includes all packages to be used in this file
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	// An HTTP router
	// For getting the env variables

	_ "github.com/lib/pq" // Postgres driver for database/sql, _ indicates it won't be referenced directly in code
)

// Constants for database, can be set in .env file
const (
	// host = "localhost"
	// port = 5432 // default port for PostgreSQL
	host     = "ec2-44-193-178-122.compute-1.amazonaws.com"
	port     = 5432
	user     = "uhieudfnmtzmgx"
	password = "ba2a59f35c97b854b67e590af27714461b2dd3be9fa620940889667d060771aa"
	dbname   = "dc8l72r4a7e68n"
)

// The struct for a task, excluding the user_uuid which is added separately.
// Tasks in JSON will use the JSON tags like "id" instead of "TaskNum".
type Item struct {
	TaskNum int    `json:"id"`
	Task    string `json:"task"`
	Status  bool   `json:"status"`
}

// Connect to PostgreSQL database and also retrieve user_id from users table
func OpenConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// 2. validates the arguments provided, doesn't create connection to database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// 3. actually opens connection to database
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GetEmail() string {
	// To be explained, related to authentication
	return "hello@gmail.com"
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func UserLogged(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte("secretkey"), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))

}

// delete task
// var DeleteTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// Set header to json content, otherwise data appear as plain text
// 	w.Header().Set("Content-Type", "application/json")

// 	// getting the task id from the request URL
// 	vars := mux.Vars(r) // vars includes all variables in the request URL route.
// 	// For example, in "/list/delete/{id}", "id" is a variable (of type string)

// 	number, err := strconv.Atoi(vars["id"]) // convert the string id to integer and assign it to variable "number"
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	db, userId := OpenConnection()
// 	sqlStatement := `DELETE FROM tasks WHERE id = $1 AND user_uuid = $2;`

// 	// Note that unlike before, we assign a variable instead of _ to the first returned value by db.Exec,
// 	// as we need it to confirm that the row was deleted
// 	res, err := db.Exec(sqlStatement, number, userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	// verifying if row was deleted
// 	_, err = res.RowsAffected()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	// to get the remaining tasks, same as the GET function
// 	rows, err := db.Query("SELECT id, task, status FROM tasks JOIN users ON tasks.user_uuid = users.user_id WHERE user_id = $1;", userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}
// 	defer rows.Close()
// 	defer db.Close()

// 	// var items []Item
// 	items := make([]Item, 0)
// 	for rows.Next() {
// 		var item Item
// 		err := rows.Scan(&item.TaskNum, &item.Task, &item.Status)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			panic(err)
// 		}
// 		items = append(items, item)
// 	}

// 	// output with indentation
// 	// convert items into byte stream
// 	itemBytes, _ := json.MarshalIndent(items, "", "\t")

// 	// write to w
// 	_, err = w.Write(itemBytes)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		panic(err)
// 	}

// 	w.WriteHeader(http.StatusOK)
// })

// edit task
// var EditTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// Set header to json content, otherwise data appear as plain text
// 	w.Header().Set("Content-Type", "application/json")

// 	// get the task id from the request url
// 	vars := mux.Vars(r)
// 	number, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	sqlStatement := `UPDATE tasks SET task = $2 WHERE id = $1 AND user_uuid = $3 RETURNING id, task, status;`

// 	// decode the requested data to 'newTask'
// 	var newTask Item

// 	// NewDecoder: Decode FROM WHERE
// 	// Decode: WHERE TO STORE the decoded data
// 	err = json.NewDecoder(r.Body).Decode(&newTask)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	db, userId := OpenConnection()
// 	defer db.Close()

// 	// retrieve the task after creation from the database and store its details in 'updatedTask'
// 	var updatedTask Item
// 	err = db.QueryRow(sqlStatement, number, newTask.Task, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	w.WriteHeader(http.StatusOK)

// 	// gives the new task as the output
// 	_ = json.NewEncoder(w).Encode(updatedTask)
// })

// // change task status
// var DoneTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// Set header to json content, otherwise data appear as plain text
// 	w.Header().Set("Content-Type", "application/json")

// 	// get the task id from the request url
// 	vars := mux.Vars(r)
// 	number, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	// store current status of the task from database
// 	var currStatus bool

// 	// store updated task
// 	var updatedTask Item

// 	sqlStatement1 := `SELECT status FROM tasks WHERE id = $1 AND user_uuid = $2;`
// 	sqlStatement2 := `UPDATE tasks SET status = $2 WHERE id = $1 AND user_uuid = $3 RETURNING id, task, status;`

// 	db, userId := OpenConnection()
// 	defer db.Close()

// 	// getting current status of the task
// 	err = db.QueryRow(sqlStatement1, number, userId).Scan(&currStatus)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}

// 	// changing the status of the task
// 	err = db.QueryRow(sqlStatement2, number, !currStatus, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		panic(err)
// 	}
// 	w.WriteHeader(http.StatusOK)

// 	// gives the new task as the output
// 	_ = json.NewEncoder(w).Encode(updatedTask)
// })
