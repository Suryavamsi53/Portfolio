package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Configurations
var (
	jwtKey       = []byte("your_super_secret_key_change_me") // In production, use environment variables
	adminUser    = "admin"
	adminHash, _ = bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost) // Default password: admin123
	messagesFile = "messages.json"
	mu           sync.Mutex
)

type Message struct {
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	Company string    `json:"company"`
	Region  string    `json:"region"`
	Country string    `json:"country"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Middleware: Authenticate JWT
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		tokenStr := cookie.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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

		next.ServeHTTP(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Simple check against hardcoded admin
	if creds.Username != adminUser || bcrypt.CompareHashAndPassword(adminHash, []byte(creds.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(messagesFile)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var msg Message
	json.NewDecoder(r.Body).Decode(&msg)
	msg.Time = time.Now()

	mu.Lock()
	defer mu.Unlock()

	var messages []Message
	data, _ := os.ReadFile(messagesFile)
	json.Unmarshal(data, &messages)
	messages = append(messages, msg)

	finalData, _ := json.MarshalIndent(messages, "", "  ")
	os.WriteFile(messagesFile, finalData, 0644)
	
	w.WriteHeader(http.StatusOK)
}

func handleUpdateMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		return
	}

	var req struct {
		Index int `json:"index"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	var messages []Message
	data, _ := os.ReadFile(messagesFile)
	json.Unmarshal(data, &messages)

	if req.Index >= 0 && req.Index < len(messages) {
		// New slice without the deleted item
		messages = append(messages[:req.Index], messages[req.Index+1:]...)
		finalData, _ := json.MarshalIndent(messages, "", "  ")
		os.WriteFile(messagesFile, finalData, 0644)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/contact", handleContact)
	http.HandleFunc("/api/messages", authenticate(handleGetMessages))
	http.HandleFunc("/api/messages/update", authenticate(handleUpdateMessages))

	fmt.Println("🚀 Secure Mail Backend running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
