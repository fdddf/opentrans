package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
}