package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	//for i, arg := range os.Args {
	//	fmt.Println(i, arg)
	//}
	//switch os.Args[1] {
	//case "hash":
	//	hash(os.Args[2])
	//case "compare":
	//	compare(os.Args[2], os.Args[3])
	//}

	passwordHash, err := hash(fmt.Sprintf("Super secret password"))
	if err != nil {
		panic(err)
	}

	err = compare(passwordHash, "Super secret password")
	if err != nil {
		panic(err)
	}

	fmt.Printf("They're a match hash was %q", passwordHash)

}

func hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Errorf("hash() error %v", err)
		return "", err
	}
	hash := string(hashedBytes)
	return hash, nil
}

func compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		fmt.Errorf("compare() error: %v", err)
		return err
	}
	return nil
}
