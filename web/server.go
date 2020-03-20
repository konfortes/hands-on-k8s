package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var userService UserService

func main() {
	initialize()
	http.HandleFunc("/health", health)
	http.HandleFunc("/users", usersHandler)

	port := getEnvOr("PORT", "4431")
	log.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func initialize() {
	// TODO: service discovery
	userService = UserService{
		Host: getEnvOr("USERS_HOST", "users"),
	}
}

func getEnvOr(env, ifNotFound string) string {
	foundEnv, found := os.LookupEnv(env)

	if found {
		return foundEnv
	}

	return ifNotFound
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}
