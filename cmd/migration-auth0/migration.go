package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type Password struct {
	UserId string
	Email  string
	Name   string
}

func main() {
	users, err := ReadAuth0Users("users.json")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
	fmt.Printf("Users: %v", users)

}

func ReadAuth0Users(filename string) ([]User, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result []User
	for fileScanner.Scan() {
		data := User{}
		err = json.Unmarshal(fileScanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	file.Close()
	return result, nil
}

func ReadFile(filename string) (*os.File, *bufio.Scanner, error) {
	readFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	return readFile, fileScanner, nil
}
