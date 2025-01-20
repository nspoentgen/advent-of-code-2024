package main

import (
	"bufio"
	"log"
	"os"
)

type EdgeType int

const (
	NORTH EdgeType = iota
	EAST
	SOUTH
	WEST
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day12_part01\input.txt`

	symbols := parseData(FILEPATH)
	regions := getRegions(symbols)
	totalPrice := calculateTotalPrice(symbols, regions)
	log.Printf("The total price is %d\n", totalPrice)
}

func parseData(filepath string) [][]rune {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file for parsing")
	}
	defer file.Close()

	symbols := make([][]rune, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineVals := make([]rune, 0)

		for _, char := range scanner.Text() {
			lineVals = append(lineVals, char)
		}

		symbols = append(symbols, lineVals)
	}

	return symbols
}

func getRegions(symbols [][]rune) map[int]map[[2]int]byte {
	id := -1
	visited := make(map[[2]int]byte)
	regions := make(map[int]map[[2]int]byte)

	for row := range len(symbols) {
		for col := range len(symbols[row]) {
			if _, contains := visited[[2]int{row, col}]; !contains {
				id++
				mapRegion(symbols, [2]int{row, col}, visited, regions, id)
			}
		}
	}

	return regions
}

func mapRegion(symbols [][]rune, initialPos [2]int, visited map[[2]int]byte, regions map[int]map[[2]int]byte, id int) {
	regionSymbol := symbols[initialPos[0]][initialPos[1]]
	regionSet := make(map[[2]int]byte)
	regions[id] = regionSet

	dfs(symbols, regionSymbol, initialPos, visited, regionSet)
}

func dfs(symbols [][]rune, regionSymbol rune, position [2]int, visited map[[2]int]byte, regionSet map[[2]int]byte) {
	regionSet[position] = 1
	visited[position] = 1

	for _, nextPosition := range getNextPositions(position, len(symbols)-1, len(symbols[0])-1) {
		_, contains := visited[nextPosition]
		newRegionPosition := symbols[nextPosition[0]][nextPosition[1]] == regionSymbol && !contains

		if newRegionPosition {
			dfs(symbols, regionSymbol, nextPosition, visited, regionSet)
		}
	}
}

func getNextPositions(position [2]int, rowMax int, colMax int) [][2]int {
	deltas := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	nextPositions := make([][2]int, 0)

	for _, delta := range deltas {
		testPosition := [2]int{position[0] + delta[0], position[1] + delta[1]}
		if testPosition[0] >= 0 && testPosition[0] <= rowMax &&
			testPosition[1] >= 0 && testPosition[1] <= colMax {
			nextPositions = append(nextPositions, testPosition)
		}
	}

	return nextPositions
}

func calculateTotalPrice(symbols [][]rune, regions map[int]map[[2]int]byte) uint32 {
	var totalPrice uint32 = 0
	numRegions := 0
	outputChannel := make(chan uint32)

	for _, regionSet := range regions {
		numRegions++
		go calculateRegionPrice(symbols, regionSet, outputChannel)
	}

	for range numRegions {
		totalPrice += <-outputChannel
	}

	return totalPrice
}

func calculateRegionPrice(symbols [][]rune, regionSet map[[2]int]byte, outputChannel chan<- uint32) {
	outputChannel <- getRegionPermiter(symbols, regionSet) * uint32(len(regionSet))
}

func getRegionPermiter(symbols [][]rune, regionSet map[[2]int]byte) uint32 {
	var perimeter uint32 = 0

	for position := range regionSet {
		perimeter += getSquarePerimeter(symbols, position)
	}

	return perimeter
}

func getSquarePerimeter(symbols [][]rune, position [2]int) uint32 {
	regionSymbol := symbols[position[0]][position[1]]
	var perimeter uint32 = 0

	for _, edgeType := range []EdgeType{NORTH, EAST, SOUTH, WEST} {
		if isPermiterSegment(symbols, position, regionSymbol, edgeType) {
			perimeter += 1
		}
	}

	return perimeter
}

func isPermiterSegment(symbols [][]rune, position [2]int, regionSymbol rune, edgeType EdgeType) bool {
	onGridEdge := false
	isRegionBoundary := false

	if edgeType == NORTH {
		onGridEdge = position[0] == 0
	} else if edgeType == EAST {
		onGridEdge = position[1] == len(symbols[0])-1
	} else if edgeType == SOUTH {
		onGridEdge = position[0] == len(symbols)-1
	} else if edgeType == WEST {
		onGridEdge = position[1] == 0
	} else {
		log.Fatal("Invalid edge type enum value")
	}

	if !onGridEdge {
		if edgeType == NORTH {
			isRegionBoundary = symbols[position[0]-1][position[1]] != regionSymbol
		} else if edgeType == EAST {
			isRegionBoundary = symbols[position[0]][position[1]+1] != regionSymbol
		} else if edgeType == SOUTH {
			isRegionBoundary = symbols[position[0]+1][position[1]] != regionSymbol
		} else if edgeType == WEST {
			isRegionBoundary = symbols[position[0]][position[1]-1] != regionSymbol
		} else {
			log.Fatal("Invalid edge type enum value")
		}
	}

	return onGridEdge || isRegionBoundary
}
