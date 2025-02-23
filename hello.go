package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

const (
	ENV_VARIABLE = "environment variables cannot be empty"
)

func Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	rd_endpoint := os.Getenv("REDIS_ENDPOINT")
	rd_password := os.Getenv("REDIS_PASSWORD")
	rd_PORT := os.Getenv("REDIS_PORT")
	if rd_endpoint == "" || rd_password == "" || rd_PORT == "" {
		log.Fatal(errors.New(ENV_VARIABLE))
	}
	db_role := os.Getenv("DB_ROLE")
	db_password := os.Getenv("DB_PASSWORD")
	db_hostname := os.Getenv("DB_HOSTNAME")
	db_name := os.Getenv("DB_DATABASE_NAME")
	if db_role == "" || db_password == "" || db_hostname == "" || db_name == "" {
		log.Fatal(errors.New(ENV_VARIABLE))
	}

	rd_url := fmt.Sprintf("rediss://default:%s@%s:%s", rd_endpoint, rd_password, rd_PORT)
	opt, _ := redis.ParseURL(rd_url)
	client := redis.NewClient(opt)
	val := client.Get("key1").Val()
	print(val)

	db_url := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", db_role, db_password, db_hostname, db_name)
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	rows, err := db.QueryContext(context.Background(), "SELECT ID, NAME FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %d - name: %s\n", id, name)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/liveness", Liveness)
	http.HandleFunc("/readiness", Readiness)

	fmt.Printf("Server will running on port: %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
