package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day01_part01\input.txt`

	left, right := parseData(FILEPATH)
	similarityScore := getSimilarityScore(left, right)
	fmt.Printf("The similarity score is %d\n", similarityScore)
}

func parseData(filepath string) ([]int, []int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	left := make([]int, 0)
	right := make([]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		value, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		left = append(left, value)

		value, err = strconv.Atoi(parts[3])
		if err != nil {
			log.Fatal(err)
		}
		right = append(right, value)
	}

	return left, right
}

func getSimilarityScore(left []int, right []int) int {
	rightSet := make(map[int]int)
	for _, key := range right {
		if _, ok := rightSet[key]; ok {
			rightSet[key]++
		} else {
			rightSet[key] = 1
		}
	}

	similarityScore := 0
	for _, value := range left {
		if count, ok := rightSet[value]; ok {
			similarityScore += value * count
		}
	}

	return similarityScore
}
