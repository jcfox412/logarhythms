package main

import (
	"fmt"
	"os"

	"github.com/jcfox412/logarhythms/internal/input"
)

func main() {

	userInput := input.UserInput{
		Reader: os.Stdin,
	}

	if err := userInput.PrintMainMenu(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
