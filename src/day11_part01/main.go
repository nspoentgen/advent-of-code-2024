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
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day11_part01\input.txt`

	initialStoneLine := parseData(FILEPATH)
	finalStoneLine := getFinalStoneLine(initialStoneLine)
	fmt.Printf("The number of stones is %d\n", len(finalStoneLine))
}

func parseData(filepath string) []uint64 {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	lineVals := make([]uint64, 0)
	for _, term := range strings.Split(line, " ") {
		number, _ := strconv.ParseUint(term, 10, 64)
		lineVals = append(lineVals, number)
	}

	return lineVals
}

func getFinalStoneLine(stones []uint64) []uint64 {
	const NUM_BLINKS int = 25

	for i := 0; i < NUM_BLINKS; i++ {
		stones = blink(stones)
	}

	return stones
}

func blink(stonesBefore []uint64) []uint64 {
	stonesAfter := make([]uint64, 0, len(stonesBefore))

	for _, oldStone := range stonesBefore {
		for _, newStone := range updateStone(oldStone) {
			stonesAfter = append(stonesAfter, newStone)
		}
	}

	return stonesAfter
}

func updateStone(stone uint64) []uint64 {
	newStones := make([]uint64, 0)

	if stone == 0 {
		newStones = append(newStones, 1)
	} else {
		stringifiedNum := strconv.FormatUint(stone, 10)
		numLength := len(stringifiedNum)

		if numLength % 2 == 0 {
			leftStone, _ := strconv.ParseUint(stringifiedNum[0 : numLength / 2], 10, 64)
			rightStone, _ := strconv.ParseUint(stringifiedNum[numLength / 2:], 10, 64)
			newStones = append(newStones, leftStone, rightStone)
		} else {
			newStones = append(newStones, stone * 2024)
		}
	}

	return newStones
}
