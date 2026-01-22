package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type URLEntry struct {
	OriginalURL string
	CreatedAt   time.Time
}

var urlMap = make(map[string]URLEntry)
var counter = 0
var mu sync.Mutex // Protect concurrent access

// Background cleanup that runs every hour
func startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	go func() {
		for range ticker.C {
			cleanupExpiredURLs()
		}
	}()
}

func cleanupExpiredURLs() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	count := 0
	for shortCode, entry := range urlMap {
		if now.Sub(entry.CreatedAt) > 24*time.Hour {
			delete(urlMap, shortCode)
			count++
		}
	}
	if count > 0 {
		log.Printf("[Background] Cleaned up %d expired URL(s)\n", count)
	}
}

func shortenURL(url string) string {
	mu.Lock()
	defer mu.Unlock()

	counter++
	shortCode := fmt.Sprintf("%d", counter)

	urlMap[shortCode] = URLEntry{
		OriginalURL: url,
		CreatedAt:   time.Now(),
	}

	return shortCode
}

func getOriginalURL(shortCode string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	entry, exists := urlMap[shortCode]

	if !exists {
		return "", false
	}

	if time.Since(entry.CreatedAt) > 24*time.Hour {
		delete(urlMap, shortCode)
		return "", false
	}

	return entry.OriginalURL, true
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	htmlFile, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		log.Printf("Error reading index.html: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlFile)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	url := strings.TrimSpace(req.URL)
	if url == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL cannot be empty"})
		return
	}

	shortCode := shortenURL(url)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_code": shortCode,
		"expires_in": "24 hours",
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		homeHandler(w, r)
		return
	}

	originalURL, exists := getOriginalURL(shortCode)
	if !exists {
		http.Error(w, "URL not found or expired", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	startCleanupRoutine()

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/api/shorten", shortenHandler)

	port := ":8080"

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
