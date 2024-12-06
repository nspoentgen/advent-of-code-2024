package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const SEARCH_PATTERN string = "XMAS"

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day04_part01\input.txt`
	data := parseData(FILEPATH)
	numMatches := getNumMatches(data)
	fmt.Printf("The number of matches is %d\n", numMatches)
}

func parseData(filepath string) [][]rune {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make([][]rune, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := make([]rune, 0)

		for _, character := range scanner.Text() {
			line = append(line, character)
		}

		data = append(data, line)
	}

	return data
}

func getNumMatches(data [][]rune) int {
	numMatches := 0
	searchFunctions := []func(int, int, [][]rune, int, int) bool{northSearch, northeastSearch, eastSearch, southeastSearch, southSearch, southwestSearch, westSearch, northwestSearch}
	names := []string{"north", "northeast", "east", "southeast", "south", "southwest", "west", "northwest"}

	matchCounts := [8]int{}

	for rowIndex := range len(data) {
		for columnIndex := range len(data[0]) {
			for i, searchFunction := range searchFunctions {
				if searchFunction(rowIndex, columnIndex, data, len(data), len(data[0])) {
					numMatches++
					matchCounts[i]++
				}
			}
		}
	}

	for i := range len(matchCounts) {
		fmt.Printf("The number of local matches for %s is %d\n", names[i], matchCounts[i])
	}

	return numMatches
}

func eastSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if col > numCols-len(SEARCH_PATTERN) {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		colIndex := col + offset

		if data[row][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func westSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if col-len(SEARCH_PATTERN)+1 < 0 {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		colIndex := col - offset

		if data[row][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func southSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row > numRows-len(SEARCH_PATTERN) {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row + offset

		if data[rowIndex][col] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func northSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row-len(SEARCH_PATTERN)+1 < 0 {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row - offset

		if data[rowIndex][col] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func southeastSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row > numRows-len(SEARCH_PATTERN) || col > numCols-len(SEARCH_PATTERN) {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row + offset
		colIndex := col + offset

		if data[rowIndex][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func northeastSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row-len(SEARCH_PATTERN)+1 < 0 || col > numCols-len(SEARCH_PATTERN) {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row - offset
		colIndex := col + offset

		if data[rowIndex][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func northwestSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row-len(SEARCH_PATTERN)+1 < 0 || col-len(SEARCH_PATTERN)+1 < 0 {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row - offset
		colIndex := col - offset

		if data[rowIndex][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}

func southwestSearch(row int, col int, data [][]rune, numRows int, numCols int) bool {
	if row > numRows-len(SEARCH_PATTERN) || col-len(SEARCH_PATTERN)+1 < 0 {
		return false
	}

	searchIndex := 0
	for offset := 0; offset < len(SEARCH_PATTERN); offset++ {
		rowIndex := row + offset
		colIndex := col - offset

		if data[rowIndex][colIndex] == rune(SEARCH_PATTERN[searchIndex]) {
			searchIndex++
		}
	}

	return searchIndex == len(SEARCH_PATTERN)
}
