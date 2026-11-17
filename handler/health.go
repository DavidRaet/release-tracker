package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func hello(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string {
		"status": "ok",
		"message": "Server is healthy",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/healthz", hello)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


