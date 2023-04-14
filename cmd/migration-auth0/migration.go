package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	//admin_pb "github.com/zitadel/zitadel"
)

type User struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type Passwords struct {
	ID Password `json:"_ID"`
}

type Password struct {
	Oid          string `json:"$oid"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func main() {
	users, err := ReadAuth0Users("users.json")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
	fmt.Printf("Users: %v", users)

	passwords, pwerr := ReadAuth0UPasswords("passwords.json")
	if err != nil {
		fmt.Printf("ERROR: %v", pwerr)
		return
	}
	fmt.Printf("Passwords: %v", passwords)
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

func ReadAuth0UPasswords(filename string) ([]Password, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result []Password
	for fileScanner.Scan() {
		data := Password{}
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

//
//func CreateZITADELMigration() {
//		keyboard := &admin_pb.ImportDataOrg{
//		Layout:  randomKeyboardLayout(),
//		Backlit: randomBool(),
//	}
//
//		return keyboard
//	}
//}
