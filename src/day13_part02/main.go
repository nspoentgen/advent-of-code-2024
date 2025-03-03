package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type MachineParams struct {
	aDeltaX   int64
	aDeltaY   int64
	bDeltaX   int64
	bDeltaY   int64
	xDistance int64
	yDistance int64
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
	const PRIZE_DISTANCE_OFFSET int64 = 10000000000000

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

			deltaX, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.aDeltaX = deltaX

			deltaY, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.aDeltaY = deltaY

		} else if relativeIndex == 1 {
			match := buttonBRegex.FindStringSubmatch(line)

			deltaX, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.bDeltaX = deltaX

			deltaY, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.bDeltaY = deltaY
		} else if relativeIndex == 2 {
			match := prizeRegex.FindStringSubmatch(line)

			xDistance, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.xDistance = PRIZE_DISTANCE_OFFSET + xDistance

			yDistance, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			machineParams.yDistance = PRIZE_DISTANCE_OFFSET + yDistance
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
	output := make(chan int64)

	for _, machineParams := range allMachineParams {
		go minimizeMachine(machineParams, output)
		numTasks++
	}

	for range numTasks {
		minNumTokens += int64(<-output)
	}

	return minNumTokens
}

func minimizeMachine(machineParams MachineParams, output chan<- int64) {
	const COST_A int64 = 3
	const COST_B int64 = 1

	constraint := func(a float64, b float64) bool {
		var integerA = int64(a)
		var integerB = int64(b)

		if math.Abs(a-float64(integerA)) > 1e-15 ||
			math.Abs(b-float64(integerB)) > 1e-15 {
			return false
		}

		xCorrect := integerA*machineParams.aDeltaX+integerB*machineParams.bDeltaX-machineParams.xDistance == 0
		yCorrect := integerA*machineParams.aDeltaY+integerB*machineParams.bDeltaY-machineParams.yDistance == 0
		return xCorrect && yCorrect
	}

	const MAX_COST int64 = 1<<63 - 1
	minCost := MAX_COST

	//Using Lagrange multipliers we can find potential optimal solution. Just need to check if solution is valid
	var a float64 = float64(machineParams.aDeltaX)
	var b float64 = float64(machineParams.aDeltaY)
	var c float64 = float64(machineParams.bDeltaX)
	var d float64 = float64(machineParams.bDeltaY)
	var e float64 = float64(machineParams.xDistance)
	var f float64 = float64(machineParams.yDistance)

	var numAPushes = (d*e - c*f) / (a*d - b*c)
	var numBPushes = (b*e - a*f) / (b*c - a*d)

	if constraint(numAPushes, numBPushes) {
		var integerNumAPushes = int64(numAPushes)
		var integerNumBPushes = int64(numBPushes)
		minCost = COST_A*integerNumAPushes + COST_B*integerNumBPushes
	}

	if minCost == MAX_COST {
		minCost = 0
	}

	output <- minCost
}
