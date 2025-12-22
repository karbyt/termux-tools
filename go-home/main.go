package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gen2brain/beeep"
)

type NotificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload NotificationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Title == "" || payload.Message == "" {
		http.Error(w, "title and message are required", http.StatusBadRequest)
		return
	}

	// Trigger system notification
	err = beeep.Notify(payload.Title, payload.Message, "")
	if err != nil {
		log.Println("Notification error:", err)
		http.Error(w, "Failed to show notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("notification sent"))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Notification server is running. Send POST requests to /notify with payload { \"title\": \"Your Title\", \"message\": \"Your Message\" }"))
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/notify", notifyHandler)

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
