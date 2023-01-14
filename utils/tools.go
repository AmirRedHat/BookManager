package utils

import (
	"fmt"
)

func Replace(str string, word string) string {
	wordLength := len(word)
	for index, _ := range str {
		fetchedWord := word[index : index+wordLength-1]
		if fetchedWord == word {
			fmt.Println("found but i can't remove it")
		}
	}

	return word
}
