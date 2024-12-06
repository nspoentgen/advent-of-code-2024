package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

	for rowIndex := range len(data) {
		for columnIndex := range len(data[0]) {
			if isMatch(rowIndex, columnIndex, data) {
				numMatches++
			}
		}
	}

	return numMatches
}

func isMatch(row int, col int, data [][]rune) bool {
	const FORWARD_PATTERN string = "MAS"
	const BACKWARD_PATTERN string = "SAM"

	if row+len(FORWARD_PATTERN) > len(data) || col+len(FORWARD_PATTERN) > len(data[0]) {
		return false
	}

	firstLeg := make([]rune, 0)
	for offset := range len(FORWARD_PATTERN) {
		firstLeg = append(firstLeg, data[row+offset][col+offset])
	}

	secondLeg := make([]rune, 0)
	startingOffset := len(FORWARD_PATTERN) - 1
	for offset := range len(FORWARD_PATTERN) {
		secondLeg = append(secondLeg, data[row+startingOffset-offset][col+offset])
	}

	firstString := string(firstLeg)
	secondString := string(secondLeg)
	return (firstString == FORWARD_PATTERN || firstString == BACKWARD_PATTERN) &&
		(secondString == FORWARD_PATTERN || secondString == BACKWARD_PATTERN)
}
