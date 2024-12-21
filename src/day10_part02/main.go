package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type WorkItem struct {
	visited map[[2]int]byte
	position [2]int
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day10_part01\input.txt`
	topoMap := parseData(FILEPATH)

	sum := getTrailRatingsSum(topoMap)
	fmt.Printf("The sum is %d\n", sum)
}

func parseData(filepath string) [][]int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	topoMap := make([][]int, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineVals := make([]int, 0)
		line := scanner.Text()

		for _, char := range line {
			lineVals = append(lineVals, int(char - '0'))
		}

		topoMap = append(topoMap, lineVals)
	}

	return topoMap
}

func getTrailRatingsSum(topoMap [][]int) int {
	sum := 0

	for row := range len(topoMap) {
		for col := range len(topoMap[row]) {
			if topoMap[row][col] == 0 {
				sum += getTrailRating([2]int{row, col}, topoMap)
			}
		}
	}

	return sum
}


//Assumes starting at valid position
func getTrailRating(position [2]int, topoMap [][]int) int {
	workStack := NewStack[*WorkItem]()
	visited := make(map[[2]int]byte)
	visited[position] = 0
	rating := 0

	workStack.Push(&WorkItem{visited: visited, position: position})
	workLeft := true

	for workLeft {
		val, err := workStack.Pop()
		workLeft = err == nil

		if workLeft{
			for _, move := range getValidMoves(val.position, topoMap, val.visited) {
				height := topoMap[move[0]][move[1]]

				if height == 9 {
					rating++
				} else {
					visited = cloneMap(val.visited)
					visited[move] = 0
					workStack.Push(&WorkItem{visited: visited, position: move})
				}
			}
		}
	}

	return rating
}

func getValidMoves(position [2]int, topoMap [][]int, visited map[[2]int]byte) [][2]int {
	var DELTAS = [4][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1}}

	moves := make([][2]int, 0)
	height := topoMap[position[0]][position[1]]

	for _, delta := range DELTAS {
		testPosition := [2]int{position[0] + delta[0], position[1] + delta[1]}
		inBounds := testPosition[0] >= 0 && testPosition[0] < len(topoMap) &&
			testPosition[1] >= 0 && testPosition[1] < len(topoMap[0])
		_, haveVisited := visited[testPosition]

		if inBounds && !haveVisited && topoMap[testPosition[0]][testPosition[1]] == height + 1 {
			moves = append(moves, testPosition)
		}
	}

	return moves
}

func cloneMap(input map[[2]int]byte) map[[2]int]byte {
	clonedMap := make(map[[2]int]byte)

	for k, v := range input {
		clonedMap[k] = v
	}

	return clonedMap
}
