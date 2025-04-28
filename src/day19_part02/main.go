package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type OutputData struct {
	design string
	count int64
}

type CacheLookup struct {
	index int
	pattern string
}

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day19_part01\input.txt`

	patterns, designs := parseData(INPUT_FILEPATH)
	numPossiblePermutations := getAllPossibleDesignPermutations(designs, patterns)

	log.Printf("The number of possible design permutations is %d\n", numPossiblePermutations)
}

func parseData(filepath string) ([]string, []string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open input file")
	}

	defer file.Close()

	patterns := make([]string, 0)
	designs := make([]string, 0)

	scanner := bufio.NewScanner(file)
	inPatternBlock := true

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			inPatternBlock = false
			continue
		}

		if inPatternBlock {
			for _, term := range strings.Split(line, ",") {
				
			patterns = append(patterns, strings.TrimLeft(term, " "))
			}
		} else {
			designs = append(designs, line)
		}
	}

	return patterns, designs
}

func getAllPossibleDesignPermutations(designs []string, patterns []string) int64 {
	outputChannel := make(chan OutputData)
	
	for _, design := range designs {
		go func () {
			cache := make(map[CacheLookup]int64)
			count := getNumPossiblePermutations(design, 0, patterns, cache)
			outputChannel <- OutputData{design, count}
		} ()
	}

	var sum int64 = 0
	for range len(designs) {
		output := <- outputChannel
		sum += output.count
	}

	return sum
}

func getNumPossiblePermutations(design string, index int, patterns []string, cache map[CacheLookup]int64) int64 {	
	//Base case	
	if index == len(design) {
		return 1
	}
	
	//Normal case
	var successCount int64 = 0

	for _, pattern := range patterns {
		if val, ok := cache[CacheLookup{index, pattern}]; ok {
			successCount += val
		} else {
			endIndexExclusive := index + len(pattern)

			if endIndexExclusive <= len(design) &&
			   design[index : endIndexExclusive] == pattern {
				childCount := getNumPossiblePermutations(design, endIndexExclusive, patterns, cache)
				
				successCount += childCount
				cache[CacheLookup{index, pattern}] = childCount
			} else {
				cache[CacheLookup{index, pattern}] = 0
			}
		}
	}

	return successCount
}
