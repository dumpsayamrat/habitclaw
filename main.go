package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dumpsayamrat/habitclaw/config"
	"github.com/dumpsayamrat/habitclaw/web"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	mux.Handle("/", web.NewHandler())

	fmt.Printf("HabitClaw starting on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
