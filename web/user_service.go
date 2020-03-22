package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// UserService ...
type UserService struct {
	Host string
	Port string
}

// CreateUser ...
func (us UserService) CreateUser(user UserInput) error {
	requestBody, err := json.Marshal(user)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:4432/users", us.Host)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Error reading request. ", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)

	return nil
}
