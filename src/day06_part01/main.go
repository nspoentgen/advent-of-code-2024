package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const STARTING_POSITION rune = '^'
const OBSTACLE rune = '#'

type DIRECTION int
const (
	North DIRECTION = iota
	East
	South
	West
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day06_part01\input.txt`
	guardMap := parseData(FILEPATH)
	startingPos := getStartingPosition(guardMap)
	numVisitedSpots := getNumVisitedSpots(guardMap, startingPos)

	fmt.Printf("Number of visited spots:%d\n", numVisitedSpots)
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

		for _, symbol := range scanner.Text() {
			line = append(line, symbol)
		}

		data = append(data, line)
	}

	return data
}

func getNumVisitedSpots(guardMap [][]rune, startingPos [2]int) int {
	const STARTING_DIRECTION DIRECTION = North
	NUM_ROWS := len(guardMap)
	NUM_COLS := len(guardMap[0])

	visitedMap := make(map[[2]int]bool)
	leftArea := false
	position := startingPos
	direction := STARTING_DIRECTION

	visitedMap[startingPos] = true
	for !leftArea {
		willLeaveArea := willLeaveArea(position, direction, NUM_ROWS, NUM_COLS)

		if !willLeaveArea {
			position, direction = getNextMove(position, direction, guardMap)
			visitedMap[position] = true
		}

		leftArea = willLeaveArea
	}

	visitedCount := 0
	for _, visited := range visitedMap {
		if visited {
			visitedCount++
		}
	}

	return visitedCount
}

func getStartingPosition(guardMap [][]rune) [2]int {
	for row := range len(guardMap) {
		for col := range len(guardMap[row]) {
			if guardMap[row][col] == STARTING_POSITION {
				return [2]int{row, col}
			}
		}
	}

	panic("No starting position found")
}

//Assumes next move can't leave the area. This indirectly means the current position is not
//on one of the boundaries
func getNextMove(position [2]int, direction DIRECTION, guardMap [][]rune) ([2]int, DIRECTION) {
	var nextPosition [2]int
	var nextDirection DIRECTION

	testPosition := getNextPosition(position, direction)
	turn := guardMap[testPosition[0]][testPosition[1]] == OBSTACLE

	if turn {
		nextDirection = getTurnedDirection(direction)
		nextPosition = getNextPosition(position, nextDirection)
	} else {
		nextPosition = testPosition
		nextDirection = direction
	}

	return nextPosition, nextDirection
}

func getNextPosition(position [2]int, direction DIRECTION) (nextPosition [2]int) {
	switch direction {
	case North:
		nextPosition = [2]int{position[0] - 1, position[1]}
	case East:
		nextPosition = [2]int{position[0], position[1] + 1}
	case South:
		nextPosition = [2]int{position[0] + 1, position[1]}
	case West:
		nextPosition = [2]int{position[0], position[1] - 1}
	}

	return
}

func getTurnedDirection(direction DIRECTION) (turnedDirection DIRECTION) {
	switch direction {
	case North:
		turnedDirection = East
	case East:
		turnedDirection = South
	case South:
		turnedDirection = West
	case West:
		turnedDirection = North
	}

	return
}

func willLeaveArea(position [2]int, direction DIRECTION, numRows int, numCols int) bool {
	switch direction {
	case North:
		return position[0] == 0
	case East:
		return position[1] == numCols - 1
	case South:
		return position[0] == numRows - 1
	case West:
		return position[1] == 0
	default:
		panic("Invalid direction")
	}
}
