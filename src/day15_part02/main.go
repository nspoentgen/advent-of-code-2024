package main

import (
	"bufio"
	"log"
	"os"
)

const (
	WALL_INPUT = '#'
	SMALL_BOX_INPUT = 'O'
	BIG_BOX_LEFT_INPUT = '['
	BIG_BOX_RIGHT_INPUT = ']'
	EMPTY_INPUT = '.'
	ROBOT_INPUT = '@'
)



func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day15_part01\input.txt`

	warehouseMap, actions, startingPosition := parseData(FILEPATH)
	warehouseIdMap, robot, objectIdMap := generateSolverData(warehouseMap, startingPosition)
	doActions(actions, warehouseIdMap, robot, objectIdMap)
	gpsSum := calcGpsSum(objectIdMap)

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
			if symbol == WALL_INPUT {
				row = append(row, WALL_INPUT, WALL_INPUT)
			} else if symbol == SMALL_BOX_INPUT {
				row = append(row, BIG_BOX_LEFT_INPUT, BIG_BOX_RIGHT_INPUT)
			} else if symbol == EMPTY_INPUT {
				row = append(row, EMPTY_INPUT, EMPTY_INPUT)
			} else if symbol == ROBOT_INPUT {
				row = append(row, ROBOT_INPUT, EMPTY_INPUT)

				startingPosition[0] = rowIndex
				startingPosition[1] = colIndex
			} else {
				log.Panic("Invalid map symobl to parse")
			}

			colIndex += 2
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

func generateSolverData(warehouseMap [][]rune, startingPosition [2]int) ([][]int, *Robot, map[int]IObject) {
	warehouseIdMap := make([][]int, 0)
	var robot Robot
	objectIdMap := make(map[int]IObject)
	boxId := ROBOT + 1

	for rowIndex, row := range warehouseMap {
		mappedRow := make([]int, 0)

		for colIndex, element := range row {
			var value int

			if element == EMPTY_INPUT {
				value = EMPTY
			} else if element == WALL_INPUT {
				value = WALL
			} else if element == ROBOT_INPUT {
				value = ROBOT
				robot = Robot{Position: startingPosition, Id: ROBOT}
				objectIdMap[ROBOT] = &robot
			} else if element == BIG_BOX_LEFT_INPUT {
				value = boxId
				box := Box{Position: [2]int{rowIndex, colIndex}, Id: boxId}
				objectIdMap[boxId] = &box
				boxId++
			} else if element == BIG_BOX_RIGHT_INPUT {
				value = boxId - 1
			} else {
				log.Panic("Invalid input element to map")
			}

			mappedRow = append(mappedRow, value)
		}

		warehouseIdMap = append(warehouseIdMap, mappedRow)
	}

	return warehouseIdMap, &robot, objectIdMap
}

func doActions(actions []rune, warehouseIdMap [][]int, robot *Robot, objectIdMap map[int]IObject) {	
	for _, action := range actions {
		robot.TryMove(action, warehouseIdMap, objectIdMap)
	}
}

func calcGpsSum(objects map[int]IObject) int64 {
	const ROW_COEFFICIENT int64 = 100
	var sum int64 = 0
	
	for _, object := range objects {
		box, ok := object.(*Box)
		
		if ok {
			gps := ROW_COEFFICIENT * int64(box.Position[0]) + int64(box.Position[1])
			sum += gps
		}
	}
	

	return sum
}
