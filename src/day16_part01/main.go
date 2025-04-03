package main

import (
	"bufio"
	"log"
	"os"
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

type MazeState struct {
	Orientation Orientation
	Position [2]int
}

type JobState struct {
	Position *[2]int
	Orientation Orientation
	Visited map[MazeState]bool
}

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day16_part01\test_input2.txt`
	const STARTING_ORIENTATION Orientation = EAST

	maze, startingLocation, goalLocation := parseData(INPUT_FILEPATH)
	minCost := solveMaze(startingLocation, STARTING_ORIENTATION, maze, goalLocation)
	log.Printf("The min cost is %d\n", minCost)
}

func parseData(filepath string) ([][]rune, *[2]int, *[2]int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file")
	}

	defer file.Close()

	maze := make([][]rune, 0)
	startingLocation := [2]int{-1, -1}
	endingLocation := [2]int{-1,-1}

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

func solveMaze(startingPosition *[2]int, startingOrientation Orientation, maze [][]rune, goalPosition *[2]int) int64 {
	const MAX_COST int64 = 1<<63 - 1
	
	cache := make(map[MazeState]int64)
	workStack := NewStack[*JobState]() 
	
	initialVisited := make(map[MazeState]bool)
	initialStateKey := MazeState{Position: *startingPosition, Orientation: startingOrientation}
	initialVisited[initialStateKey] = true

	initialState := JobState{
		Position: startingPosition,
		Orientation: startingOrientation,
		Visited: initialVisited}
	workStack.Push(&initialState)

	workLeft := true
	for workLeft {
		state, err := workStack.Pop()

		if err == nil {
			moves, moveCosts, reachedGoal := getValidMoves(state.Position, state.Orientation, maze, state.Visited, goalPosition)

			if reachedGoal {
				key := MazeState{Position: *state.Position, Orientation: state.Orientation}
				cache[key] = 0
			} else {
				uncompletedMoves := make([]*MazeState, 0)
				minPathCost := MAX_COST

				for i := range len(moves) {
					childPathCost, exists := cache[*moves[i]]

					if exists && childPathCost < MAX_COST {
						pathCost := moveCosts[i] + childPathCost

						if pathCost < minPathCost {
							minPathCost = pathCost
						}
					} else if !exists {
						uncompletedMoves = append(uncompletedMoves, moves[i])
					}
				}

				if len(uncompletedMoves) == 0 && minPathCost == MAX_COST {
					cacheKey := MazeState{Position: *state.Position, Orientation: state.Orientation}
					cache[cacheKey] = MAX_COST
				} else if len(uncompletedMoves) == 0 && minPathCost < MAX_COST {
					cacheKey := MazeState{Position: *state.Position, Orientation: state.Orientation}
					cache[cacheKey] = minPathCost
				} else {
					workStack.Push(state)
					visitedCopy := cloneMap(state.Visited)
					visitedKey := MazeState{Position: *state.Position, Orientation: state.Orientation}
					visitedCopy[visitedKey] = true

					for _, uncompletedMove := range uncompletedMoves {
						workStack.Push(&JobState{
							Position: &uncompletedMove.Position,
							Orientation: uncompletedMove.Orientation,
							Visited: visitedCopy})
					}
				}
			}
		} else {
			workLeft = false
		}
	}

	minCost, exists := cache[initialStateKey]
	
	if !exists {
		log.Fatal("Something is wrong. Could not calculate min cost")
	}

	return minCost
}

func getValidMoves(position *[2]int, orientation Orientation, maze [][]rune, visited map[MazeState]bool, goalPosition *[2]int) ([]*MazeState, []int64, bool) {
	const MOVE_COST int64 = 1
	const ROTATE_COST int64 = 1000
	
	moves := make([]*MazeState, 0)
	costs := make([]int64, 0)
	reachedGoal := false

	if *position == *goalPosition {
		reachedGoal = true
	} else {
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

		movePosition := [2]int {position[0] + rowOffset, position[1] + colOffset}
		if maze[movePosition[0]][movePosition[1]] == EMPTY {
			move1 := MazeState{Position: movePosition, Orientation: orientation}

			_, exists := visited[move1]
			if !exists {
				moves = append(moves, &move1)
				costs = append(costs, MOVE_COST)
			}
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
		
		_, exists := visited[move2]
		if !exists {
			moves = append(moves, &move2)
			costs = append(costs, ROTATE_COST)
		}

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
		
		_, exists = visited[move3]
		if !exists {
			moves = append(moves, &move3)
			costs = append(costs, ROTATE_COST)
		}
	}

	return moves, costs, reachedGoal
}
