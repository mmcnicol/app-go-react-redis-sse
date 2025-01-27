package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

var ctx = context.Background()

// Initialize Redis client
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
}

// Data types to represent the different types of remote data
type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LabResult struct {
	ID    string `json:"id"`
	Test  string `json:"test"`
	Value string `json:"value"`
}

type EmergencyCareSummary struct {
	ID          string `json:"id"`
	PatientName string `json:"patient_name"`
	Summary     string `json:"summary"`
}

// Function to fetch documents from a REST API
func fetchDocuments(url string, userID string) ([]Document, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var documents []Document
	if err := json.Unmarshal(body, &documents); err != nil {
		return nil, err
	}

	// Store documents in Redis
	documentsJSON, _ := json.Marshal(documents)
	err = rdb.Set(ctx, fmt.Sprintf("%s:documents", userID), documentsJSON, 10*time.Minute).Err() // TTL: 10 minutes
	if err != nil {
		return nil, err
	}
	
	return documents, nil
}

// Function to fetch lab results from a REST API
func fetchLabResults(url string, userID string) ([]LabResult, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var labResults []LabResult
	if err := json.Unmarshal(body, &labResults); err != nil {
		return nil, err
	}
	
	// Store lab results in Redis
	labResultsJSON, _ := json.Marshal(labResults)
	cacheKey := fmt.Sprintf("%s:labResults", userID)
	err = rdb.Set(ctx, cacheKey, labResultsJSON, 10*time.Minute).Err() // TTL: 10 minutes
	if err != nil {
		return nil, err
	}

	return labResults, nil
}

// Function to fetch emergency care summaries from a REST API
func fetchEmergencyCareSummary(url string, userID string) ([]EmergencyCareSummary, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var summaries []EmergencyCareSummary
	if err := json.Unmarshal(body, &summaries); err != nil {
		return nil, err
	}
	
	// Store emergency care summary in Redis
	summariesJSON, _ := json.Marshal(summaries)
	cacheKey := fmt.Sprintf("%s:emergencyCareSummaries", userID)
	err = rdb.Set(ctx, cacheKey, summariesJSON, 10*time.Minute).Err() // TTL: 10 minutes
	if err != nil {
		fmt.Println("in fetchEmergencyCareSummary(), cache put error: ", err.Error())
		return nil, err
	}

	return summaries, nil
}

// SSE handler to stream updates to the frontend
func sseHandler(w http.ResponseWriter, r *http.Request) {

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userID := r.URL.Query().Get("userID") // Assume userID is passed as a query parameter

	// URLs for the REST APIs
	documentsURL := "http://mockapi:8081/documents"
	labResultsURL := "http://mockapi:8081/lab-results"
	emergencyCareSummariesURL := "http://mockapi:8081/emergency-care-summaries"

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(3)

	// Start goroutines to fetch data concurrently
	go func() {
		defer wg.Done()
		_, err := fetchDocuments(documentsURL, userID)
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			w.(http.Flusher).Flush()
		} else {
			fmt.Fprintf(w, "event: documents\ndata: %s\n\n", "fetched")
			w.(http.Flusher).Flush()
		}
	}()

	go func() {
		defer wg.Done()
		_, err := fetchLabResults(labResultsURL, userID)
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			w.(http.Flusher).Flush()
		} else {
			fmt.Fprintf(w, "event: labResults\ndata: %s\n\n", "fetched")
			w.(http.Flusher).Flush()
		}
	}()

	go func() {
		defer wg.Done()
		_, err := fetchEmergencyCareSummary(emergencyCareSummariesURL, userID)
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			w.(http.Flusher).Flush()
		} else {
			fmt.Fprintf(w, "event: emergencyCareSummaries\ndata: %s\n\n", "fetched")
			w.(http.Flusher).Flush()
		}
	}()

	wg.Wait()
	
	fmt.Fprintf(w, "event: eagerLoading\ndata: %s\n\n", "finished")
	w.(http.Flusher).Flush()
}

// Handler to retrieve data from Redis
func getDataHandler(w http.ResponseWriter, r *http.Request) {

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	userID := r.URL.Query().Get("userID")
	dataType := r.URL.Query().Get("type")

	// Retrieve data from Redis
	data, err := rdb.Get(ctx, fmt.Sprintf("%s:%s", userID, dataType)).Result()
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! backend\n")
	})
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		var health = "UP"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	})

	// Serve the SSE endpoint
	http.HandleFunc("/updates", sseHandler)

	// Serve the endpoint to retrieve data
	http.HandleFunc("/data", getDataHandler)

	// Start the HTTP server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
