package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// {
//     "username":"abctest",
//     "password":"hello"
// }
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	if r.URL.Path != "/addUser" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if UserLogged(w, r) {
		w.Write([]byte(fmt.Sprintf("Cannot Add User, User Already Logged In")))
		return
	}
	if r.Method == "GET" {
		http.ServeFile(w, r, "addUser.html")
	} else {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		email := r.FormValue("email")
		passwrd := r.FormValue("password")

		var credentials User
		credentials.Email = email
		credentials.Password = passwrd
		// err := json.NewDecoder(r.Body).Decode(&credentials)
		var err error
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

		db := OpenConnection()
		addEmailAndPassword := `INSERT INTO users (email,passwrd) VALUES ($1, $2) RETURNING email, passwrd;`

		var updatedUser User
		err = db.QueryRow(addEmailAndPassword, credentials.Email, password).Scan(&updatedUser.Email, &updatedUser.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}

		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "User Created Successfully")
	}

}
