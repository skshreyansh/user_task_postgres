package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	// w.Header().Set("Content-Type", "application/json")
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		http.ServeFile(w, r, "login.html")
	} else {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		var credentials User
		credentials.Email = r.FormValue("email")
		credentials.Password = r.FormValue("password")

		// err := json.NewDecoder(r.Body).Decode(&credentials)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	panic(err)
		// }

		db := OpenConnection()
		var password string
		var err error
		sql := `Select passwrd from users where email= $1 `
		// getUser := `SELECT user_id FROM users WHERE email = $1;`
		err = db.QueryRow(sql, credentials.Email).Scan(&password)

		if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.ServeFile(w, r, "login.html")
			fmt.Fprintf(w, "Incorrect Username or Password Please Retry")
			return
		}
		expirationTime := time.Now().Add(time.Minute * 5)

		claims := &Claims{
			Username: credentials.Email,
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

		fmt.Fprintf(w, "User Logged In Successfully")
	}
}
