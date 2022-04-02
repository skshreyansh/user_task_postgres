package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// {
//     "username":"abctest",
//     "password":"hello"
// }
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
