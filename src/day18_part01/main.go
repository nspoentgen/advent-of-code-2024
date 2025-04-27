package main

import (
	"bufio"
	"container/heap"
	"log"
	"os"
	"slices"
	"strings"
	"strconv"
)

const (
	CORRUPTED  = '#'
	EMPTY = '.'
)

const MAX_COST int64 = 1<<63 - 1

type MazeState struct {
	Time int
	Position    [2]int
}

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day18_part01\input.txt`
	MAX_DIM_INDEX := 70
	INITIAL_STATE := MazeState{1024, [2]int{0, 0}}

	maze := parseData(INPUT_FILEPATH)
	minCost := solveMaze(&INITIAL_STATE, maze, MAX_DIM_INDEX)


	log.Printf("The min cost is %d\n", minCost)
}

func parseData(filepath string) map[[2]int]int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file")
	}

	defer file.Close()

	maze := make(map[[2]int]int)
	scanner := bufio.NewScanner(file)
	time := 1

	for scanner.Scan() {
		terms := strings.Split(scanner.Text(), ",")
		
		left, err := strconv.Atoi(terms[0])
		if err != nil {
			log.Fatal(err)
		}

		right, err := strconv.Atoi(terms[1])
		if err != nil {
			log.Fatal(err)
		}

		maze[[2]int{left, right}] = time
		time += 1
	}

	return maze
}

func solveMaze(initialState *MazeState, maze map[[2]int]int, maxDimIndex int) int64 {
	const MAX_ITERATIONS int = 10000
	const MOVE_COST int = 1
	
	goalPosition := [2]int{ maxDimIndex, maxDimIndex }
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

		if state.Position == goalPosition {
			return cost
		}

		if cost > minCost {
			continue
		}

		moves := getValidMoves(state, maze, maxDimIndex)

		for i := range len(moves) {
			updatedMincost := cost + int64(MOVE_COST)
			currentMinCost := getMinCost(minCosts, &moves[i])

			if updatedMincost < currentMinCost {
				item := Item{value: &moves[i], priority: updatedMincost}
				heap.Push(&searchQueue, &item)

				minCosts[moves[i]] = updatedMincost
				prev[moves[i]] = []*MazeState{state}
			} else if updatedMincost == currentMinCost {
				prev[moves[i]] = append(prev[moves[i]], state)
			}
		}

		iteration++
	}

	return MAX_COST
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

func getValidMoves(state *MazeState, maze map[[2]int]int, maxDimIndex int) []MazeState {
	const TIME_DELTA int = 0
	possibleMoves := [4][2]int {
		{ state.Position[0] + 1, state.Position[1] },
		{ state.Position[0] - 1, state.Position[1] },
		{ state.Position[0], state.Position[1] + 1 },
		{ state.Position[0], state.Position[1] - 1 }}

	validMoves := make([]MazeState, 0)

	for _, move := range possibleMoves {
		boundsCriterion := move[0] >= 0 && move[0] <= maxDimIndex && move[1] >= 0 && move[1] <= maxDimIndex
		
		corruptedTime, exists := maze[move]
		newTime := state.Time + TIME_DELTA
		
		openSpaceCriterion := !exists || newTime < corruptedTime

		if boundsCriterion && openSpaceCriterion {
			validMoves = append(validMoves, MazeState{Time: newTime, Position: move})
		}
	}

	return validMoves
}
