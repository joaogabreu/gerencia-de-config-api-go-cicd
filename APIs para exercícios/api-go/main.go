package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"log"
	"os"
	"api-go/user"

	_ "github.com/lib/pq"
)

func main() {
	repo, cleanup := buildRepository()
	defer cleanup()

	service := user.NewUserService(repo)
	controller := user.NewUserController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListUsers(w, r)
		case http.MethodPost:
			controller.CreateUser(w, r)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUser(w, r)
		case http.MethodPut:
			controller.UpdateUser(w, r)
		case http.MethodDelete:
			controller.DeleteUser(w, r)
		}
	})

	log.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func buildRepository() (user.UserRepository, func()) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		name := getEnv("DB_NAME", "api_go")
		username := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		sslMode := getEnv("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, name, sslMode)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("failed to open postgres connection, using in-memory repository: %v", err)
		return user.NewUserRepository(), func() {}
	}

	if err := db.Ping(); err != nil {
		log.Printf("failed to connect to postgres, using in-memory repository: %v", err)
		_ = db.Close()
		return user.NewUserRepository(), func() {}
	}

	if err := user.EnsureSchema(db); err != nil {
		log.Printf("failed to initialize schema, using in-memory repository: %v", err)
		_ = db.Close()
		return user.NewUserRepository(), func() {}
	}

	log.Println("using PostgreSQL repository")
	return user.NewPostgresUserRepository(db), func() {
		_ = db.Close()
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}