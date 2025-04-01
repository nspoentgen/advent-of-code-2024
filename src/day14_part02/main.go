package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day14_part01\input.txt`
	var FIELD_DIMENSIONS = [2]int{101, 103} //width, height

	robots := parseData(FILEPATH, &FIELD_DIMENSIONS)
	moveRobots(robots)
	robotMapping := mapPositions(robots)
	quadrantSums := getQuadrantSums(robotMapping, &FIELD_DIMENSIONS)

	product := 1
	for _, quadrantSum := range quadrantSums {
		product *= quadrantSum
	}

	log.Printf("The product is %d\n", product)
}

func parseData(filepath string, fieldDimension *[2]int64) []Robot {
	robotRegex := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file for parsing")
	}
	defer file.Close()

	robots := make([]Robot, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		match := robotRegex.FindStringSubmatch(scanner.Text())

		xPosition, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}

		yPosition, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}

		xVelocity, err := strconv.Atoi(match[3])
		if err != nil {
			log.Fatal(err)
		}

		yVelocity, err := strconv.Atoi(match[4])
		if err != nil {
			log.Fatal(err)
		}

		robot := NewRobot(
			&[2]int{xPosition, yPosition},
			&[2]int{xVelocity, yVelocity},
			fieldDimension)
		robots = append(robots, *robot)
	}

	return robots
}

func moveRobots(robots []Robot) {
	const NUM_MOVES int = 100
	var wg sync.WaitGroup

	for i := range robots {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range NUM_MOVES {
				robots[i].Move()
			}
		}()
	}

	wg.Wait()
}

func mapPositions(robots []Robot) map[[2]int]int {
	positionMap := make(map[[2]int]int)

	for i := range robots {
		_, ok := positionMap[robots[i].Position]

		if ok {
			positionMap[robots[i].Position]++
		} else {
			positionMap[robots[i].Position] = 1
		}
	}

	return positionMap
}

func getQuadrantSums(positions map[[2]int]int, fieldDimension *[2]int) *[4]int {
	var xMidpointIndex int = (fieldDimension[0] - 1) / 2
	var yMidpointIndex int = (fieldDimension[1] - 1) / 2
	counts := new([4]int)

	for position, count := range positions {
		if position[0] == xMidpointIndex || position[1] == yMidpointIndex {
			continue
		}

		xPositive := position[0] > xMidpointIndex
		yPositive := position[1] > yMidpointIndex

		if xPositive && yPositive {
			counts[0] += count
		} else if !xPositive && yPositive {
			counts[1] += count
		} else if !xPositive && !yPositive {
			counts[2] += count
		} else {
			counts[3] += count
		}
	}

	return counts
}
