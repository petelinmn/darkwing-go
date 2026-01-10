package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type helloResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok", Time: time.Now().UTC().Format(time.RFC3339)})
	})

	mux.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		a := r.URL.Query().Get("a")
		b := r.URL.Query().Get("b")

		var numA, numB float64
		if a != "" {
			if _, err := fmt.Sscanf(a, "%f", &numA); err != nil {
				http.Error(w, `{"error":"invalid parameter a"}`, http.StatusBadRequest)
				return
			}
		}
		if b != "" {
			if _, err := fmt.Sscanf(b, "%f", &numB); err != nil {
				http.Error(w, `{"error":"invalid parameter b"}`, http.StatusBadRequest)
				return
			}
		}

		sum := numA + numB
		_ = json.NewEncoder(w).Encode(map[string]float64{"sum": sum})
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		version := os.Getenv("VERSION")
		if version == "" {
			version = "dev"
		}
		_ = json.NewEncoder(w).Encode(helloResponse{Message: "hello from api 7", Version: version})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"darkwing-go-api","endpoints":["/health","/hello","/version"]}`))
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		version := os.Getenv("VERSION")
		if version == "" {
			version = "dev"
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"version": version})
	})

	addr := ":8080"
	log.Printf("API listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
