package helpers

import (
	"fmt"
	"strings"

	"github.com/howeyc/gopass"
)

func AskForConfirmation(q string) bool {
	fmt.Print(q)

	var response string
	read, err := fmt.Scanln(&response)
	if err != nil && read != 0 {
		fmt.Println("Error: " + PrintRed(err.Error()))
		return false
	} else if read == 0 {
		return false
	}

	response = strings.ToLower(response)

	if response == "y" || response == "yes" {
		return true
	} else if response == "n" || response == "no" || response == "" {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation(q)
	}
}

func GetUserInput(q string) string {
	fmt.Print(q)

	var response string
	read, err := fmt.Scanln(&response)
	if err != nil && read != 0 {
		fmt.Println("Error: " + PrintRed(err.Error()))
		return ""
	} else if read == 0 {
		return ""
	}

	return response
}

func GetUserInputPassword(q string) string {
	fmt.Print(q)

	pass, err := gopass.GetPasswd()
	if err != nil {
		return ""
	}

	return string(pass)
}
