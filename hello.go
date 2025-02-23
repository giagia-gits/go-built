package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
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
	opt, _ := redis.ParseURL("rediss://default:AUidAAIjcDE0NTg5NTY3MzJiOTg0YTcxOWE2YzI2MDBkYTU5NDQxY3AxMA@stirring-gar-18589.upstash.io:6379")
	client := redis.NewClient(opt)

	val := client.Get("key1").Val()
	print(val)

	connStr := "user='koyeb-adm' password=npg_0XkYaq9BCONs host=ep-icy-river-a2npoafz.eu-central-1.pg.koyeb.app dbname='koyebdb'"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(context.Background(), "INSERT INTO users (ID, NAME) VALUES ($1, $2)", 1, "user1"); err != nil {
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
