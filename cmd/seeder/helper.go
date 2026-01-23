package main

import (
	"fmt"
	"math/rand"
)

func getRandomCountry() Country {
	idx := rand.Intn(len(CountryList))
	return CountryList[idx]
}

func getRandomResearcher(researcher_list []Researcher) Researcher {
	idx := rand.Intn(len(researcher_list))
	return researcher_list[idx]
}

func randomString(l int, isOnlyCapital bool, isWithNumber bool, isNumberOny bool) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	if isOnlyCapital {
		letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}

	if isWithNumber {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}

	if isNumberOny {
		letters = []rune("0123456789")
	}

	s := make([]rune, l)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)

}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func randomDate() string {
	year := rand.Intn(22) + 2000
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func randomInt(max int) int {

	return rand.Intn(max)
}
