package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Code:    200,
		Message: "Success",
		Payload: "0.1.1",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func myIPHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		http.Error(w, "Unable to fetch IP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	response := Response{
		Code:    200,
		Message: "Success",
		Payload: string(ip),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/my-ip", myIPHandler)
	http.ListenAndServe(":8000", nil)
}

