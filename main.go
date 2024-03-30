package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Email    string
	Password string
}

type Ad struct {
	Id       int
	UserId   int
	Title    string
	Text     string
	ImageURL string
	Price    int
}

var db *sql.DB

var jwtKey = []byte("secret_key")

func main() {
	var err error
	connStr := "user=postgres password=1234 dbname=db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/create-ad", createAd).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isPasswordSafe(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email уже занят", http.StatusConflict)
		return
	}

	if isEmailValid(user.Email) && isPasswordSafe(user.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		query := "INSERT INTO users(email, password) VALUES($1, $2) RETURNING id"
		err = db.QueryRow(query, user.Email, hashedPassword).Scan(&user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var userId int
	err = db.QueryRow("SELECT Id FROM users WHERE email = $1", user.Email).Scan(&userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func extractUserIDFromToken(r *http.Request) (int, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка при парсинге токена: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, fmt.Errorf("ошибка при извлечении идентификатора пользователя из токена: %w", err)
		}
		return userID, nil
	}

	return 0, errors.New("невозможно извлечь идентификатор пользователя из токена")
}

func createAd(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	fmt.Println(err)

	var ad Ad
	err = json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	maxTitleLength := 80
	maxTextLength := 1000

	if len(ad.Title) > maxTitleLength {
		http.Error(w, "Длина заголовка превышает максимальное значение", http.StatusBadRequest)
		return
	}

	if len(ad.Text) > maxTextLength {
		http.Error(w, "Длина текста объявления превышает максимальное значение", http.StatusBadRequest)
		return
	}

	if ad.Price < 0 {
		http.Error(w, "Цена не может быть отрицательной", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO ads (UserId, Title, Text, ImageUrl, Price) VALUES ($1, $2, $3, $4, $5) RETURNING Id"
	err = db.QueryRow(query, userID, ad.Title, ad.Text, ad.ImageURL, ad.Price).Scan(&ad.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}


