package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day19_part01\input.txt`

	patterns, designs := parseData(INPUT_FILEPATH)
	numPossibleDesigns := getNumPossibleDesigns(designs, patterns)

	log.Printf("The number of possible designs is %d\n", numPossibleDesigns)
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

func getNumPossibleDesigns(designs []string, patterns []string) int {
	count := 0
	
	for _, design := range designs {
		if isPossible(design, 0, patterns) {
			count++
		}
	}

	return count
}

func isPossible(design string, index int, patterns []string) bool {	
	//Base case	
	if index == len(design) {
		return true
	}
	
	//Normal case
	success := false

	for _, pattern := range patterns {
		endIndexExclusive := index + len(pattern)

		if endIndexExclusive <= len(design) &&
		   design[index : endIndexExclusive] == pattern &&
		   isPossible(design, endIndexExclusive, patterns) {
			success = true
			break
		}
	}

	return success
}

/*
func getMatchMap(design string, patterns []string) map[string]map[int]bool {
	matchMap := make(map[string]map[int]bool)

	for _, pattern := range patterns {
		matchIndices := make([int]bool, 0)
		
		for index := 0; index + len(pattern) - 1 < len(design); index++ {
			testSlice := pattern[index : index + len(pattern)]
			
			if testSlice == pattern {
				matchIndices[index] = true
			}
		}

		if len(matchIndices) > 0 {
			matchMap[pattern] = matchIndices
		}
	}

	return matchMap
}
	*/
