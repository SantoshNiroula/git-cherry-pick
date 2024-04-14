package helper

import (
	"log"
	"os"
	"strings"
)

func ReadPRNumberFromFile(fileName string) []string {
	entry, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Unable to read file", err)
	}
	entries := string(entry)
	list := strings.Split(entries, "\n")
	if len(list) == 1 {
		log.Fatal("File is empty")
	}
	return list
}
