package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day02_part01\input.txt`

	records := parseData(FILEPATH)
	numSafeRecords := getRecordSafeCount(records)
	fmt.Printf("The number of safe records is %d\n", numSafeRecords)
}

func parseData(filepath string) [][]int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	records := make([][]int, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := make([]int, 0)
		levels := strings.Split(scanner.Text(), " ")

		for _, level := range levels {
			value, err := strconv.Atoi(level)
			if err != nil {
				log.Fatal(err)
			}

			record = append(record, value)
		}

		records = append(records, record)
	}

	return records
}

func getRecordSafeCount(records [][]int) int {
	count := 0

	for _, originalRecord := range records {
		for i := -1; i < len(originalRecord); i++ {
			testRecord := getTestRecord(i, originalRecord)

			if isRecordSafe(testRecord) {
				count++
				break
			}
		}
	}

	return count
}

func getTestRecord(removalIndex int, originalRecord []int) []int {
	if removalIndex == -1 {
		return originalRecord
	} else if removalIndex == 0 {
		return originalRecord[1:]
	} else if removalIndex < len(originalRecord)-1 {
		testRecord := make([]int, removalIndex)
		copy(testRecord, originalRecord[0:removalIndex])
		testRecord = slices.Concat(testRecord, originalRecord[removalIndex+1:])

		return testRecord
	} else {
		return originalRecord[0:removalIndex]
	}
}

func isRecordSafe(record []int) bool {
	const MIN_DIFFERENCE int = 1
	const MAX_DIFFERENCE int = 3

	var derivativeSign int

	for i := 1; i < len(record); i++ {
		diff := record[i] - record[i-1]
		absDiff := int(math.Round(math.Abs(float64(diff))))

		if absDiff < MIN_DIFFERENCE || absDiff > MAX_DIFFERENCE {
			return false
		}

		if i == 1 {
			derivativeSign = signum(diff)
		} else if signum(diff) != derivativeSign {
			return false
		}
	}

	return true
}

func signum(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}
