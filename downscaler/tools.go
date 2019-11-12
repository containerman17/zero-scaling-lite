package main

import (
	"log"
	"os"
	"strings"

	"github.com/sethvargo/go-password/password"
)

func genPassword() string {
	res, err := password.Generate(10, 4, 0, true, true)
	if err != nil {
		log.Println("Error: password", err)
	}
	return res
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func starPassword(password string, symbolsShowed int) string {
	openPart := password[0:min(symbolsShowed, len(password))]
	starPart := strings.Repeat("*", max(0, len(password)-symbolsShowed))
	return openPart + starPart
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
