package main

import (
	"bufio"
	"container/heap"
	"log"
	"os"
	"slices"
)

const (
	WALL  = '#'
	EMPTY = '.'
	START = 'S'
	END   = 'E'
)

type Orientation int

const (
	NORTH Orientation = iota
	EAST
	SOUTH
	WEST
)

const MAX_COST int64 = 1<<63 - 1

type MazeState struct {
	Orientation Orientation
	Position    [2]int
}

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day16_part01\input.txt`
	const STARTING_ORIENTATION Orientation = EAST

	maze, startingLocation, goalLocation := parseData(INPUT_FILEPATH)
	initalState := MazeState{Position: *startingLocation, Orientation: STARTING_ORIENTATION}
	minCost, paths := solveMaze(&initalState, maze, goalLocation)
	
	uniqueTiles := make(map[[2]int]bool)
	for _, path := range paths {
		for _, tile := range path {
			uniqueTiles[tile.Position] = true
		}
	}

	log.Printf("The min cost is %d\n", minCost)
	log.Printf("The number of unique tiles is %d\n", len(uniqueTiles))
}

func parseData(filepath string) ([][]rune, *[2]int, *[2]int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file")
	}

	defer file.Close()

	maze := make([][]rune, 0)
	startingLocation := [2]int{-1, -1}
	endingLocation := [2]int{-1, -1}

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		row := make([]rune, 0)

		for j, symbol := range scanner.Text() {
			if symbol == START {
				startingLocation = [2]int{i, j}
				symbol = EMPTY
			} else if symbol == END {
				endingLocation = [2]int{i, j}
				symbol = EMPTY
			}

			row = append(row, symbol)
		}

		maze = append(maze, row)
		i++
	}

	return maze, &startingLocation, &endingLocation
}

func solveMaze(initialState *MazeState, maze [][]rune, goalPosition *[2]int) (int64, [][]*MazeState) {
	const MAX_ITERATIONS int = 1000000
	
	minCosts := make(map[MazeState]int64)
	prev := make(map[MazeState][]*MazeState)
	searchQueue := make(PriorityQueue, 0)
	searchQueue.Push(&Item{value: initialState, priority: 0, index: 0})
	heap.Init(&searchQueue)

	minCosts[*initialState] = 0
	iteration := 0

	for iteration < MAX_ITERATIONS {
		searchItem := heap.Pop(&searchQueue).(*Item)
		state := searchItem.value
		cost := searchItem.priority
		minCost := getMinCost(minCosts, state)

		if state.Position == *goalPosition {
			return cost, generatePaths(prev, state, make([]*MazeState, 0))
		}

		if cost > minCost {
			continue
		}

		moves, moveCosts := getValidMoves(&state.Position, state.Orientation, maze)
		for i := range len(moves) {
			updatedMincost := cost + moveCosts[i]
			currentMinCost := getMinCost(minCosts, moves[i])

			if updatedMincost < currentMinCost {
				item := Item{value: moves[i], priority: updatedMincost}
				heap.Push(&searchQueue, &item)

				minCosts[*moves[i]] = updatedMincost
				prev[*moves[i]] = []*MazeState{state}
			} else if updatedMincost == currentMinCost {
				prev[*moves[i]] = append(prev[*moves[i]], state)
			}
		}

		iteration++
	}

	return MAX_COST, make([][]*MazeState, 0)
}

func generatePaths(prev map[MazeState][]*MazeState, initalState *MazeState, path []*MazeState) [][]*MazeState {
	processing := true
	nodes := []*MazeState{initalState}
	paths := make([][]*MazeState, 0)

	for processing {
		if len(nodes) == 1 {
			path = append(path, nodes[0])
			nodes, processing = prev[*nodes[0]]
		} else {
			for _, node := range nodes {
				pathCopy := make([]*MazeState, len(path))
				copy(pathCopy, path)
				paths = append(paths, generatePaths(prev, node, pathCopy)...)
			}

			processing = false
		}
	}

	if len(paths) == 0 {
		slices.Reverse(path)
		paths = append(paths, path)
	}

	return paths
}

func getMinCost(minCosts map[MazeState]int64, state *MazeState) int64 {
	minCost, exists := minCosts[*state]
	if !exists {
		minCost = MAX_COST
	}

	return minCost
}

func getValidMoves(position *[2]int, orientation Orientation, maze [][]rune) ([]*MazeState, []int64) {
	const MOVE_COST int64 = 1
	const ROTATE_COST int64 = 1000

	moves := make([]*MazeState, 0)
	costs := make([]int64, 0)

	//Move move
	var rowOffset int
	var colOffset int

	if orientation == NORTH {
		rowOffset = -1
		colOffset = 0
	} else if orientation == EAST {
		rowOffset = 0
		colOffset = 1
	} else if orientation == SOUTH {
		rowOffset = 1
		colOffset = 0
	} else {
		rowOffset = 0
		colOffset = -1
	}

	movePosition := [2]int{position[0] + rowOffset, position[1] + colOffset}
	if maze[movePosition[0]][movePosition[1]] == EMPTY {
		move1 := MazeState{Position: movePosition, Orientation: orientation}
		moves = append(moves, &move1)
		costs = append(costs, MOVE_COST)
	}

	//Rotate CW move
	var cwMoveOrientation Orientation

	if orientation == NORTH {
		cwMoveOrientation = EAST
	} else if orientation == EAST {
		cwMoveOrientation = SOUTH
	} else if orientation == SOUTH {
		cwMoveOrientation = WEST
	} else {
		cwMoveOrientation = NORTH
	}

	move2 := MazeState{Position: *position, Orientation: cwMoveOrientation}
	moves = append(moves, &move2)
	costs = append(costs, ROTATE_COST)

	//Rotate CCW move
	var ccwMoveOrientation Orientation

	if orientation == NORTH {
		ccwMoveOrientation = WEST
	} else if orientation == WEST {
		ccwMoveOrientation = SOUTH
	} else if orientation == SOUTH {
		ccwMoveOrientation = EAST
	} else {
		ccwMoveOrientation = NORTH
	}

	move3 := MazeState{Position: *position, Orientation: ccwMoveOrientation}
	moves = append(moves, &move3)
	costs = append(costs, ROTATE_COST)

	return moves, costs
}
