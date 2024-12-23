package main

import (
    "database/sql"
    "log"
    "fmt"
    "net/http"
    "os"
    "strings" // Add this import

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "github.com/rs/cors"

    "STRIVEBackend/internal/api/http/server"
    "STRIVEBackend/internal/config"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file:", err)
    }

    // Log GitHub-related environment variables
    log.Printf("GITHUB_CLIENT_ID: %s", os.Getenv("GITHUB_CLIENT_ID"))
    log.Printf("GITHUB_CLIENT_SECRET: %s", maskSecret(os.Getenv("GITHUB_CLIENT_SECRET")))
    log.Printf("GITHUB_REDIRECT_URI: %s", os.Getenv("GITHUB_REDIRECT_URI"))

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
    log.Printf("Database connection string: %s", maskPassword(connStr))
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    router := server.NewRouter(db)

    corsOptions := cors.Options{
        AllowedOrigins:   []string{"https://habittrackerfordevs.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    }

    log.Printf("CORS AllowedOrigins: %v", corsOptions.AllowedOrigins)

    handler := cors.New(corsOptions).Handler(router)

    log.Println("Starting server on 0.0.0.0:8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

// maskSecret masks most of the characters in a secret string
func maskSecret(secret string) string {
    if len(secret) <= 4 {
        return "****"
    }
    return secret[:2] + "****" + secret[len(secret)-2:]
}

// maskPassword masks the password in a connection string
func maskPassword(connStr string) string {
    parts := strings.Split(connStr, ":")
    if len(parts) >= 3 {
        parts[2] = "****" + parts[2][strings.Index(parts[2], "@"):]
    }
    return strings.Join(parts, ":")
}
