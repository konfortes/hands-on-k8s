package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// UserInput ...
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func handleUsers(w http.ResponseWriter, req *http.Request) {
	userInput := UserInput{}
	if err := json.NewDecoder(req.Body).Decode(&userInput); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid input"))
		return
	}

	log.Printf("persisting user %+v", userInput)

	w.WriteHeader(http.StatusCreated)
}
