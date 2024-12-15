package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/message"
	"log"
	"os"
	"strconv"
	"strings"
)

type BridgeData struct {
	value uint64
	operands []uint64
}

type OPERATION uint64
const (
	Add OPERATION = iota
	Multiply
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day07_part01\input.txt`
	data := parseData(FILEPATH)
	operationsSum := getPossibleOperationsSum(&data)
	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("The sum is %d\n", operationsSum)

	/*
	const n int = 4
	permutations := generatePermutations(n)
	debugDict := make(map[[n]OPERATION]bool)

	for _, permutation := range permutations {
		var key [n]OPERATION
		for index, val := range permutation {
			key[index] = val
		}
		debugDict[key] = true

		fmt.Printf("%d\n", permutation)
	}

	fmt.Printf("Num permutations = %d\n", len(debugDict))

	 */
}

func parseData(filepath string) [] BridgeData {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make([]BridgeData, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		rowTerms := strings.Split(scanner.Text(), ":")

		value, err := strconv.ParseUint(rowTerms[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		rowValue := BridgeData{
			value:    value,
			operands: make([]uint64, 0),
		}

		for _, operandRaw := range strings.Split(rowTerms[1], " ") {
			if operandRaw == "" { continue }

			operand, err := strconv.ParseUint(operandRaw, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			rowValue.operands = append(rowValue.operands, operand)
		}

		data = append(data, rowValue)
	}

	return data
}

func getPossibleOperationsSum(data *[]BridgeData) uint64 {
	var sum uint64 = 0

	for index, row := range *data {
		result := false

		if operationPossible(row.value, &row.operands) {
			sum += row.value
			result = true
		}

		fmt.Printf("Row %d = %t\n", index, result)
	}

	return sum
}

func operationPossible(value uint64, operands *[]uint64) bool {
	permutations := generatePermutations(len(*operands) - 1)
	for _, permutation := range permutations {
		var testResult uint64

		for j := 0; j < len(*operands) - 1; j++ {
			if j == 0 {
				testResult = doOperation((*operands)[j], (*operands)[j+1], permutation[j])
			} else {
				testResult = doOperation(testResult, (*operands)[j+1], permutation[j])
			}
		}

		if testResult == value {
			return true
		}
	}

	return false
}

func generatePermutations(n int) [][]OPERATION {
	operationValues := []OPERATION{Add, Multiply}
	permutations := make([][]OPERATION, 0)

	if n == 1 {
		return [][]OPERATION{
			{Add},
			{Multiply}}
	}

	subpermutations := generatePermutations(n - 1)
	for _, operation := range operationValues {
		temp := make([][]OPERATION, 0)
		for i := range len(subpermutations) {
			row := make([]OPERATION, 0)

			for j := range subpermutations[i] {
				row = append(row, subpermutations[i][j])
			}

			temp = append(temp, row)
		}

		for _, row := range temp {
			row = append(row, operation)
			permutations = append(permutations, row)
		}
	}

	return permutations
}

func doOperation(a uint64, b uint64, operation OPERATION) uint64 {
	if operation == Add {
		return a + b
	} else {
		return a * b
	}
}
