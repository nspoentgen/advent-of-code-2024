package main

import (
	"bufio"
	"log"
	"os"
)

const (
	WALL = '#'
	BOX = 'O'
	EMPTY = '.'
	ROBOT = '@'
)

const (
	UP = '^'
	RIGHT = '>'
	DOWN = 'v'
	LEFT = '<'
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day15_part01\input.txt`

	warehouseMap, actions, startingPosition := parseData(FILEPATH)
	doActions(actions, warehouseMap, startingPosition)
	gpsSum := calcGpsSum(warehouseMap)

	log.Printf("The GPS sum is %d\n", gpsSum)
}

func parseData(filepath string) ([][]rune, []rune, [2]int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Cannot open file")
	}

	defer file.Close()

	warehouseMap := make([][]rune, 0)
	startingPosition := [2]int{-1, -1}

	scanner := bufio.NewScanner(file)
	inMapSection := true
	rowIndex := 0
	colIndex := 0
	scanner.Scan()

	for inMapSection {
		row := make([]rune, 0)
		colIndex = 0

		for _, symbol := range scanner.Text() {
			row = append(row, symbol)

			if symbol == ROBOT {
				startingPosition[0] = rowIndex
				startingPosition[1] = colIndex
			}

			colIndex++
		}

		warehouseMap = append(warehouseMap, row)
		
		rowIndex++
		scanner.Scan()
		inMapSection = scanner.Text() != ""
	}

	actions := make([]rune, 0)
	for scanner.Scan() {
		for _, direction := range scanner.Text() {
			actions = append(actions, direction)
		}
	}

	return warehouseMap, actions, startingPosition
}

func doActions(actions []rune, warehouseMap [][]rune, startingPosition [2]int) {
	robotPosition := startingPosition
	
	for _, action := range actions {
		_, robotPosition = move(robotPosition, warehouseMap, action)
	}
}

func move(objectPosition [2]int, warehouseMap [][]rune, direction rune) (bool, [2]int) {
	var rowOffset int
	var colOffset int

	if direction == UP {
		rowOffset = -1
		colOffset = 0
	} else if direction == RIGHT {
		rowOffset = 0
		colOffset = 1
	} else if direction == DOWN {
		rowOffset = 1
		colOffset = 0
	} else {
		rowOffset = 0
		colOffset = -1
	}

	row := objectPosition[0] + rowOffset
	col := objectPosition[1] + colOffset

	var moved bool
	var newPosition [2]int
	copy(newPosition[:], objectPosition[:])

	if warehouseMap[row][col] == EMPTY {
		warehouseMap[row][col] = warehouseMap[objectPosition[0]][objectPosition[1]]
		warehouseMap[objectPosition[0]][objectPosition[1]] = EMPTY
		newPosition = [2]int{row, col}
		moved = true
	} else if warehouseMap[row][col] == WALL {
		moved = false
	} else if warehouseMap[row][col] == BOX {
		moved, _ = move([2]int{row, col}, warehouseMap, direction)
		
		if moved {
			warehouseMap[row][col] = warehouseMap[objectPosition[0]][objectPosition[1]]
			warehouseMap[objectPosition[0]][objectPosition[1]] = EMPTY
			newPosition = [2]int{row, col}
		}
	} else {
		log.Panic("Invalid map symbol for pushing")
	}

	return moved, newPosition
}

func calcGpsSum(warehouseMap [][]rune) int64 {
	const ROW_COEFFICIENT int64 = 100
	var sum int64 = 0

	for row := range len(warehouseMap) {
		for col := range len(warehouseMap[row]) {
			if warehouseMap[row][col] == BOX {
				sum += ROW_COEFFICIENT * int64(row) + int64(col)
			}
		}
	}

	return sum
}
