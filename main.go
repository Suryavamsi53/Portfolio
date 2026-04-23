package main

import (
	"encoding/json"
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
	jwtKey       []byte
	adminUser    string
	adminHash    []byte
	messagesFile = "messages.json"
	mu           sync.Mutex
)

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "svv-portfolio-secure-random-key-2026"
		log.Println("WARNING: JWT_SECRET not set, using default.")
	}
	jwtKey = []byte(secret)

	adminUser = os.Getenv("ADMIN_USER")
	if adminUser == "" {
		adminUser = "Surya"
	}

	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		pass = "suryavamsi"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to generate admin hash: %v", err)
	}
	adminHash = hash
}

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

	// Simple check against admin credentials
	if creds.Username != adminUser || bcrypt.CompareHashAndPassword(adminHash, []byte(creds.Password)) != nil {
		log.Printf("Failed login attempt for user: %s", creds.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("Successful login for user: %s", creds.Username)

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
	log.Printf("Received transmission from %s (%s)", msg.Name, msg.Email)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("SVV Portfolio Backend starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
