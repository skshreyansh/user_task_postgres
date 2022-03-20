package handlers

import (
	// Includes all packages to be used in this file
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	// An HTTP router
	// For getting the env variables

	_ "github.com/lib/pq" // Postgres driver for database/sql, _ indicates it won't be referenced directly in code
	"golang.org/x/crypto/bcrypt"
)

// Constants for database, can be set in .env file
const (
	// host = "localhost"
	// port = 5432 // default port for PostgreSQL
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
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
	// Getting constants from .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// user, ok := os.LookupEnv("USER")
	// if !ok {
	// 	log.Fatal("Error loading env variables")
	// }
	// password, ok := os.LookupEnv("PASSWORD")
	// if !ok {
	// 	log.Fatal("Error loading env variables")
	// }
	// dbname, ok := os.LookupEnv("DB_NAME")
	// if !ok {
	// 	log.Fatal("Error loading env variables")
	// }

	// connecting to database
	// 1. creating the connection string
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

	// add email to users table if not present
	// email := GetEmail()
	// addEmail := `INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO NOTHING;`
	// _, err = db.Exec(addEmail, email)
	// if err != nil {
	// 	panic(err)
	// }

	// // get user_id
	// var userId string
	// getUser := `SELECT user_id FROM users WHERE email = $1;`
	// err = db.QueryRow(getUser, email).Scan(&userId)
	// if err != nil {
	// 	panic(err)
	// }

	return db
}

func GetEmail() string {
	// To be explained, related to authentication
	return "hello@gmail.com"
}

type Email struct {
	Email string `json:"email"`
}

// Get complete list of tasks
func GetList(w http.ResponseWriter, r *http.Request) {
	// Set header to json content, otherwise data appear as plain text
	w.Header().Set("Content-Type", "application/json")

	// Connect to database and get user_id
	db := OpenConnection()
	var getemail Email
	err := json.NewDecoder(r.Body).Decode(&getemail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	var userId string
	sql := `Select user_id from users where email= $1 `
	// getUser := `SELECT user_id FROM users WHERE email = $1;`
	err = db.QueryRow(sql, getemail.Email).Scan(&userId)

	// Return all tasks (rows) as id, task, status where the user_uuid of the task is the same as user_id we have obtained in the previous step
	rows, err := db.Query("SELECT id, task, status FROM tasks JOIN users ON tasks.user_uuid = users.user_id WHERE user_id = $1;", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	defer rows.Close()
	defer db.Close()

	// Initializing slice like this and not "var items []Item" because aforementioned method returns null when empty thus leading to errors,
	// while used method returns empty slice
	items := make([]Item, 0)
	// Add each task to array of Items
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.TaskNum, &item.Task, &item.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
	// Output with indentation
	// convert items into byte stream
	// itemBytes, _ := json.MarshalIndent(items, "", "\t")
	// // write to w
	// _, err = w.Write(itemBytes)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	panic(err)
	// }

	// w.WriteHeader(http.StatusOK)

	// Alternatively, output without indentation
	// NewEncoder: WHERE should the encoder write to
	// Encode: encode WHAT
	// _ = json.NewEncoder(w).Encode(items)
}

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	db := OpenConnection()
	addEmailAndPassword := `INSERT INTO users (email,passwrd) VALUES ($1, $2) RETURNING email, passwrd;`

	var updatedUser User
	err = db.QueryRow(addEmailAndPassword, credentials.Username, password).Scan(&updatedUser.Username, &updatedUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(updatedUser)

}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	db := OpenConnection()
	var password string
	sql := `Select passwrd from users where email= $1 `
	// getUser := `SELECT user_id FROM users WHERE email = $1;`
	err = db.QueryRow(sql, credentials.Username).Scan(&password)

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secretkey"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
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

func Logout(w http.ResponseWriter, r *http.Request) {
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
	expirationTime := time.Now().Add(-time.Minute * 5)

	claims = &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secretkey"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
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
