package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// UserService ...
type UserService struct {
	Host string
}

// CreateUser ...
func (us UserService) CreateUser(user UserInput) error {
	requestBody, err := json.Marshal(us)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/users", us.Host)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO: remove after debug
	log.Println(string(body))

	return nil
}
