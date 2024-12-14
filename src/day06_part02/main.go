package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

const STARTING_POSITION rune = '^'
const OBSTACLE rune = '#'
const OPEN_SPOT rune = '.'

type DIRECTION int
const (
	North DIRECTION = iota
	East
	South
	West
)

type VisitedKey struct {
	position [2]int
	direction DIRECTION
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day06_part01\input.txt`
	guardMap := parseData(FILEPATH)
	startingPos := getStartingPosition(guardMap)
	numGuardLoops := getNumGuardLoops(guardMap, startingPos)

	fmt.Printf("Number of guard loops: %d\n", numGuardLoops)
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

func getNumGuardLoops(guardMap [][]rune, startingPos [2]int) uint32 {
	var numLoops atomic.Uint32
	var wg sync.WaitGroup

	for i := range len(guardMap) {
		for j := range len(guardMap[i]) {
			if guardMap[i][j] == OPEN_SPOT {
				wg.Add(1)

				go func() {
					copiedMap := copyGuardMap(guardMap)
					copiedMap[i][j] = OBSTACLE

					if guardWillLoop(copiedMap, startingPos) {
						numLoops.Add(1)
					}

					wg.Done()
				}()
			}
		}
	}

	wg.Wait()
	return numLoops.Load()
}

func copyGuardMap(guardMap [][]rune) [][]rune {
	copiedMap := make([][]rune, 0)

	for i := range guardMap {
		row := make([]rune, 0)

		for j := range guardMap[i] {
			row = append(row, guardMap[i][j])
		}

		copiedMap = append(copiedMap, row)
	}

	return copiedMap
}

func guardWillLoop(guardMap [][]rune, startingPos [2]int) bool {
	const STARTING_DIRECTION DIRECTION = North

	visitedMap := make(map[VisitedKey]bool)
	loopDetected := false
	position := startingPos
	direction := STARTING_DIRECTION
	inBounds := true

	for inBounds && !loopDetected {
		visitedKey := VisitedKey{position: position, direction: direction}

		if _, exists := visitedMap[visitedKey]; !exists {
			visitedMap[visitedKey] = true
		} else {
			loopDetected = true
		}

		position, direction, inBounds = getNextMove(position, direction, guardMap)
	}

	return loopDetected
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

func getNextMove(position [2]int, direction DIRECTION, guardMap [][]rune) ([2]int, DIRECTION, bool) {
	var nextPosition [2]int
	var nextDirection DIRECTION

	testPosition := getNextPosition(position, direction, guardMap)
	testPositionInBounds := isInBounds(testPosition, guardMap)

	if testPositionInBounds && guardMap[testPosition[0]][testPosition[1]] == OBSTACLE {
		nextDirection = getTurnedDirection(direction)
		nextPosition = getNextPosition(position, nextDirection, guardMap)
		if guardMap[nextPosition[0]][nextPosition[1]] == OBSTACLE {
			nextPosition = position;
		}
	} else {
		nextPosition = testPosition
		nextDirection = direction
	}

	return nextPosition, nextDirection, isInBounds(nextPosition, guardMap)
}

func isInBounds(position [2]int, guardMap [][]rune) bool {
	return position[0] >= 0 && position[0] < len(guardMap) &&
		position[1] >= 0 && position[1]  < len(guardMap[0])
}

func getNextPosition(position [2]int, direction DIRECTION, guardMap [][]rune) (nextPosition [2]int) {
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
