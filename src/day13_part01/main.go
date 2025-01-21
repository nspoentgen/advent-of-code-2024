package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type MachineParams struct {
	aDeltaX   int
	aDeltaY   int
	bDeltaX   int
	bDeltaY   int
	xDistance int
	yDistance int
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day13_part01\input.txt`

	allMachineParams := parseData(FILEPATH)
	minNumTokens := getMinNumTokens(allMachineParams)
	fmt.Printf("Min tokens = %d\n", minNumTokens)
}

func parseData(filepath string) []MachineParams {
	buttonARegex := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	buttonBRegex := regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	prizeRegex := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	allMachineParams := make([]MachineParams, 0, 128)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Cannot open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCounter := -1
	machineParams := MachineParams{}

	for scanner.Scan() {
		lineCounter++
		relativeIndex := lineCounter % 4
		line := scanner.Text()

		if relativeIndex == 0 {
			match := buttonARegex.FindStringSubmatch(line)

			deltaX, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.aDeltaX = deltaX

			deltaY, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.aDeltaY = deltaY

		} else if relativeIndex == 1 {
			match := buttonBRegex.FindStringSubmatch(line)

			deltaX, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.bDeltaX = deltaX

			deltaY, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.bDeltaY = deltaY
		} else if relativeIndex == 2 {
			match := prizeRegex.FindStringSubmatch(line)

			xDistance, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.xDistance = xDistance

			yDistance, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			machineParams.yDistance = yDistance
		} else {
			allMachineParams = append(allMachineParams, machineParams)
			machineParams = MachineParams{}
		}
	}

	allMachineParams = append(allMachineParams, machineParams)
	machineParams = MachineParams{}

	return allMachineParams
}

func getMinNumTokens(allMachineParams []MachineParams) int64 {
	var minNumTokens int64 = 0
	numTasks := 0
	output := make(chan int)

	for _, machineParams := range allMachineParams {
		go minimizeMachine(machineParams, output)
		numTasks++
	}

	for range numTasks {
		minNumTokens += int64(<-output)
	}

	return minNumTokens
}

func minimizeMachine(machineParams MachineParams, output chan<- int) {
	const COST_A = 3
	const COST_B = 1
	const MAX_PUSHES = 100

	positionConstraint := func(a int, b int) bool {
		xCorrect := a*machineParams.aDeltaX+b*machineParams.bDeltaX-machineParams.xDistance == 0
		yCorrect := a*machineParams.aDeltaY+b*machineParams.bDeltaY-machineParams.yDistance == 0
		return xCorrect && yCorrect
	}

	const MAX_COST int = 1<<31 - 1
	minCost := MAX_COST

	for aPushes := 0; aPushes <= MAX_PUSHES; aPushes++ {
		for bPushes := 0; bPushes <= MAX_PUSHES; bPushes++ {
			if positionConstraint(aPushes, bPushes) {
				cost := aPushes*COST_A + bPushes*COST_B

				if cost < minCost {
					minCost = cost
				}
			}
		}
	}

	if minCost == MAX_COST {
		minCost = 0
	}

	output <- int(minCost)
}
