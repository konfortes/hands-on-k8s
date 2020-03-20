package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// UserInput ...
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
	userInput := UserInput{}
	if err := json.NewDecoder(req.Body).Decode(&userInput); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid input"))
		return
	}

	log.Printf("handling user %+v", userInput)

	processUser(&userInput)
	if err := persistUser(userInput); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func processUser(user *UserInput) {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
}

func persistUser(user UserInput) error {
	return userService.CreateUser(user)
}
