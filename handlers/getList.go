package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Email struct {
	Email string `json:"email"`
}

// Get complete list of tasks
func GetList(w http.ResponseWriter, r *http.Request) {
	// Set header to json content, otherwise data appear as plain text
	// w.Header().Set("Content-Type", "application/json")
	if r.URL.Path != "/list" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		http.ServeFile(w, r, "listUser.html")
	} else {
		if !UserLogged(w, r) {
			w.Write([]byte(fmt.Sprintf("Cannot List Users Task, User Not Logged In")))
			return
		}
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		email := r.FormValue("email")

		var getemail Email
		getemail.Email = email
		var userId string
		// err := json.NewDecoder(r.Body).Decode(&getemail)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	panic(err)
		// }

		// client := redis.NewClient(&redis.Options{
		// 	Addr:     "localhost:6379",
		// 	Password: "",
		// 	DB:       0,
		// })
		db := OpenConnection()

		// val, err := client.Get(getemail.Email).Result()
		// if err != nil {
		var err error
		fmt.Println(err)
		fmt.Println("TESTING CACHE SERVICE")
		// Connect to database and get user_id
		sql := `Select user_id from users where email= $1 `
		// getUser := `SELECT user_id FROM users WHERE email = $1;`
		err = db.QueryRow(sql, getemail.Email).Scan(&userId)
		// Setting cache data
		defer db.Close()

		// err = client.Set(getemail.Email, userId, 0).Err()
		// if err != nil {
		// 	fmt.Println(err)
		// 	fmt.Println("TEsting cache setting service")
		// }
		// } else {
		// 	userId = val
		// }
		fmt.Println("TEST", userId)
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
	}
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

type Profile struct {
	Name    string
	Hobbies []string
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
