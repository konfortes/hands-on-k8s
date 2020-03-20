package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/users", usersHandler)

	port := getEnvOr("PORT", "4432")
	log.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
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
