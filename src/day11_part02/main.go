package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type WorkItem struct {
	number     uint64
	blinkCount int
}

type CacheEntry struct {
	numBlinks  int
	finalCount uint64
}

const NUM_BLINKS int = 75

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day11_part01\input.txt`

	initialStoneLine := parseData(FILEPATH)
	finalCount := getFinalCount(initialStoneLine)
	fmt.Printf("The number of stones is %d\n", finalCount)
}

func parseData(filepath string) []uint64 {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	lineVals := make([]uint64, 0)
	for _, term := range strings.Split(line, " ") {
		number, _ := strconv.ParseUint(term, 10, 64)
		lineVals = append(lineVals, number)
	}

	return lineVals
}

func getFinalCount(initialStoneLine []uint64) uint64 {
	cache := make(map[[2]uint64]uint64)
	workStack := NewStack[WorkItem]()

	for i := len(initialStoneLine) - 1; i >= 0; i-- {
		workStack.Push(WorkItem{
			number:     initialStoneLine[i],
			blinkCount: 75})
	}

	for workStack.Len() > 0 {
		task, _ := workStack.Pop()

		if task.blinkCount == 1 {
			rootCase(&task, cache)
		} else {
			mainCase(&task, workStack, cache)
		}
	}

	var finalCount uint64 = 0
	for _, stoneNumber := range initialStoneLine {
		subcount, contains := cache[[2]uint64{stoneNumber, 75}]

		if !contains {
			log.Fatal("75th blink of stone not found")
		}

		finalCount += subcount
	}

	return finalCount
}

func blink(stone uint64) []uint64 {
	newStones := make([]uint64, 0)

	if stone == 0 {
		newStones = append(newStones, 1)
	} else {
		stringifiedNum := strconv.FormatUint(stone, 10)
		numLength := len(stringifiedNum)

		if numLength%2 == 0 {
			leftStone, _ := strconv.ParseUint(stringifiedNum[0:numLength/2], 10, 64)
			rightStone, _ := strconv.ParseUint(stringifiedNum[numLength/2:], 10, 64)
			newStones = append(newStones, leftStone, rightStone)
		} else {
			newStones = append(newStones, stone*2024)
		}
	}

	return newStones
}

func rootCase(task *WorkItem, cache map[[2]uint64]uint64) {
	cache[[2]uint64{task.number, 1}] = uint64(len(blink(task.number)))
}

func mainCase(task *WorkItem, workStack *Stack[WorkItem], cache map[[2]uint64]uint64) {
	subproblemsDone := true
	newStones := blink(task.number)

	for _, newStone := range newStones {
		if _, contains := cache[[2]uint64{newStone, uint64(task.blinkCount - 1)}]; !contains {
			//On first detection that subproblems are incomplete, add back the original problem
			//to queue it after subproblems (since we are LIFO)
			if subproblemsDone {
				workStack.Push(*task)
			}

			subproblemsDone = false
			workStack.Push(WorkItem{
				number:     newStone,
				blinkCount: task.blinkCount - 1})
		}
	}

	if subproblemsDone {
		var count uint64 = 0

		for _, newStone := range newStones {
			subcount, contains := cache[[2]uint64{newStone, uint64(task.blinkCount - 1)}]
			if !contains {
				log.Fatal("Something is wrong. Entry should be available but was not found")
			}

			count += subcount
		}

		cache[[2]uint64{task.number, uint64(task.blinkCount)}] = count
	}
}
