package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)

type Email struct {
	Email string `json:"email"`
}

// Get complete list of tasks
func GetList(w http.ResponseWriter, r *http.Request) {
	// Set header to json content, otherwise data appear as plain text
	w.Header().Set("Content-Type", "application/json")

	var getemail Email
	var userId string
	err := json.NewDecoder(r.Body).Decode(&getemail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	db := OpenConnection()

	val, err := client.Get(getemail.Email).Result()
	if err != nil {
		fmt.Println(err)
		fmt.Println("TESTING CACHE SERVICE")
		// Connect to database and get user_id
		sql := `Select user_id from users where email= $1 `
		// getUser := `SELECT user_id FROM users WHERE email = $1;`
		err = db.QueryRow(sql, getemail.Email).Scan(&userId)
		// Setting cache data
		defer db.Close()

		err = client.Set(getemail.Email, userId, 0).Err()
		if err != nil {
			fmt.Println(err)
			fmt.Println("TEsting cache setting service")
		}
	} else {
		userId = val
	}

	///#####
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
