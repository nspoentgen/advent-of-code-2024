package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type CommandMarker struct {
	indexStart int
	command    bool
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day03_part01\input.txt`

	lines := parseData(FILEPATH)
	filteredOperands := getFilteredMatches(lines)
	total := getTotal(filteredOperands)
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

func getFilteredMatches(lines []string) [][]int {
	var operandsExpression = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	var commandsExpression = regexp.MustCompile(`don't\(\)|do\(\)`)
	operandMatchIndices := make([][][]int, 0)
	commandsMatchIndices := make([][][]int, 0)

	for _, line := range lines {
		operandMatchIndices = append(operandMatchIndices, operandsExpression.FindAllStringSubmatchIndex(line, -1))
		commandsMatchIndices = append(commandsMatchIndices, commandsExpression.FindAllStringSubmatchIndex(line, -1))
	}

	return filterMatches(operandMatchIndices, commandsMatchIndices, lines)
}

func filterMatches(operandMatchIndices [][][]int, commandsMatchIndices [][][]int, lines []string) [][]int {
	const ENABLE_COMMAND string = "do()"

	filteredOperands := make([][]int, 0)
	previousLineCommand := true

	for i, lineCommandIndices := range commandsMatchIndices {
		commandHistory := make([]CommandMarker, 0)
		commandHistory = append(commandHistory, CommandMarker{indexStart: -1, command: previousLineCommand})

		for _, indexRange := range lineCommandIndices {
			command := lines[i][indexRange[0]:indexRange[1]]
			if command == ENABLE_COMMAND {
				commandHistory = append(commandHistory, CommandMarker{indexStart: indexRange[0], command: true})
				previousLineCommand = true
			} else {
				commandHistory = append(commandHistory, CommandMarker{indexStart: indexRange[0], command: false})
				previousLineCommand = false
			}
		}

		filteredLineOperands := make([]int, 0)
		for _, matchIndices := range operandMatchIndices[i] {
			if operationIsEnabled(commandHistory, matchIndices[0]) {
				operand0, err := strconv.Atoi(lines[i][matchIndices[2]:matchIndices[3]])
				if err != nil {
					log.Fatal(err)
				}

				operand1, err := strconv.Atoi(lines[i][matchIndices[4]:matchIndices[5]])
				if err != nil {
					log.Fatal(err)
				}

				filteredLineOperands = append(filteredLineOperands, operand0, operand1)
			}
		}

		if len(filteredLineOperands) > 0 {
			filteredOperands = append(filteredOperands, filteredLineOperands)
		}
	}

	return filteredOperands
}

func operationIsEnabled(commandHistory []CommandMarker, operationIndex int) bool {
	for i := 0; i < len(commandHistory)-1; i++ {
		if operationIndex > commandHistory[i].indexStart && operationIndex < commandHistory[i+1].indexStart {
			return commandHistory[i].command
		}
	}

	return commandHistory[len(commandHistory)-1].command
}

func getTotal(operandPairs [][]int) int64 {
	total := int64(0)

	for _, lineOperandPairs := range operandPairs {
		for i := 0; i < len(lineOperandPairs); i += 2 {
			total += int64(lineOperandPairs[i] * lineOperandPairs[i+1])
		}
	}

	return total
}
