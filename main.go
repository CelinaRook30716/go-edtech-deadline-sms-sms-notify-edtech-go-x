package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	c := NewSMSClient()
	http.HandleFunc("/deadlines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var d Deadline
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		result, err := SendDeadlineAlert(c, d, time.Now())
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
