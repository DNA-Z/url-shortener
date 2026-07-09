package main

import (
	"encoding/json"
	"log"
)

type User struct {
	Name     string
	LastName string
}

func main() {
	jsonUser := `{"Name": "John", "LastName": "Smith"}`
	var u User

	if err := json.Unmarshal([]byte(jsonUser), &u); err != nil {
		log.Printf("ошибка: %v", err)
	} else {
		log.Printf("пользователь: %v", u)
	}
}
