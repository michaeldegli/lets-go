package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte("pass"), 12)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(hashedPassword))

	err := bcrypt.CompareHashAndPassword([]byte("$2a$12$w1vFjnnOmKnP9y5ec.1eP.F3azTzrytGgqVvD8xpl3SQ8BJ1PwQem"), []byte("pass"))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Good")
}
