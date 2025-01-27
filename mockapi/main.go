package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! mockapi\n")
	})

	http.HandleFunc("/documents", func(w http.ResponseWriter, r *http.Request) {
	
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
		
		time.Sleep(1000 * time.Millisecond)

		documents := []Document{
			{ID: "1", Name: "Document 1"},
			{ID: "2", Name: "Document 2"},
			{ID: "3", Name: "Document 3"},
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(documents)
	})

	http.HandleFunc("/lab-results", func(w http.ResponseWriter, r *http.Request) {
		
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
		
		time.Sleep(100 * time.Millisecond)
		
		labResults := []LabResult{
			{ID: "1", Test: "Blood Test", Value: "Normal"},
			{ID: "2", Test: "Urine Test", Value: "Abnormal"},
			{ID: "3", Test: "X-Ray", Value: "Clear"},
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(labResults)
	})

	http.HandleFunc("/emergency-care-summaries", func(w http.ResponseWriter, r *http.Request) {
		
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

		time.Sleep(3000 * time.Millisecond)

		summaries := []EmergencyCareSummary{
			{ ID: "3", PatientName: "John Doe", Summary: "Patient was treated for a mild fever." },
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(summaries)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
