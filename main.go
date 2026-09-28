package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Fatal("API_TOKEN is required")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/app?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", newMux(db, token)))
}

func newMux(db *sql.DB, token string) http.Handler {
	mux := http.NewServeMux()
	// spider の起点。JSON API だけだとクロール対象が見つからないのでリンク一覧を返す
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<a href="/health">health</a> <a href="/items">items</a> <a href="/me">me</a>`))
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), "SELECT id, name FROM items ORDER BY id")
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		items := []item{}
		for rows.Next() {
			var it item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				http.Error(w, "db error", http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}
		writeJSON(w, items)
	})
	mux.Handle("GET /me", auth(token, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"user": "admin"})
	}))
	mux.Handle("POST /items", auth(token, func(w http.ResponseWriter, r *http.Request) {
		var in item
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in); err != nil || in.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if err := db.QueryRowContext(r.Context(), "INSERT INTO items (name) VALUES ($1) RETURNING id", in.Name).Scan(&in.ID); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, in)
	}))
	return mux
}

type item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ponytail: 固定 Bearer トークン。ユーザー管理が必要になったら DB 上のユーザー/セッションに置き換える
func auth(token string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
