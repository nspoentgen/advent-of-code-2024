package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day03_part01\input.txt`

	lines := parseData(FILEPATH)
	operands := getMatches(lines)
	total := getTotal(operands)
	fmt.Printf("The total is %d\n", total)
}

func parseData(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func getMatches(lines []string) [][2]int {
	var matchExpression = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	var allArgs = make([][2]int, 0)

	for _, line := range lines {
		var lineArgs [][2]int

		for _, match := range matchExpression.FindAllStringSubmatch(line, -1) {
			if lineArgs == nil {
				lineArgs = make([][2]int, 0)
			}

			operand0, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}

			operand1, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			lineArgs = append(lineArgs, [2]int{operand0, operand1})
		}

		if lineArgs != nil {
			allArgs = slices.Concat(allArgs, lineArgs)
		}
	}

	return allArgs
}

func getTotal(operandPairs [][2]int) int64 {
	total := int64(0)

	for _, operandPair := range operandPairs {
		total += int64(operandPair[0] * operandPair[1])
	}

	return total
}
