package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type WorkItem struct {
	number uint64
	blinkCount int
}

const NUM_BLINKS int = 75

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day11_part01\test_input.txt`

	initialStoneLine := parseData(FILEPATH)
	finalCount := getFinalStoneLineCount(initialStoneLine)
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

func getFinalStoneLineCount(initialStoneLine []uint64) uint64 {
	maxBlink := atomic.Int32{}

	workQueue := make(chan WorkItem, 1000)
	for _, stoneNumber := range(initialStoneLine) {
		workQueue <- WorkItem{number: stoneNumber, blinkCount: 0}
	}

	cache := sync.Map{}
	stoneCount := atomic.Uint64{}
	writerCount := atomic.Int64{}
	done := false

	for !done {
		dataRead := false

		WorkQueueWaitLoop:
		for {
			select {
			case nextWorkItem := <- workQueue:
				writerCount.Add(1)
				go writer(nextWorkItem, workQueue, &writerCount, &stoneCount, &maxBlink, &cache)
				dataRead = true

				break WorkQueueWaitLoop
			default:
				if writerCount.Load() == 0 {
					break WorkQueueWaitLoop
				}
			}
		}

		done = !dataRead && writerCount.Load() == 0
	}

	return stoneCount.Load()
}

func writer(stoneInfo WorkItem, workQueue chan WorkItem, writerCount *atomic.Int64, stoneCount *atomic.Uint64, maxBlink *atomic.Int32, cache *sync.Map) {
	newWorkItems := blink(stoneInfo, stoneCount, maxBlink, cache)

	for _, item := range newWorkItems {
		workQueue <- item
	}

	writerCount.Add(-1)
}

func blink(srcStone WorkItem, stoneCount *atomic.Uint64, maxBlink *atomic.Int32, cache *sync.Map) []WorkItem {
	newWork := make([]WorkItem, 0)

	if srcStone.blinkCount == NUM_BLINKS {
		stoneCount.Add(1)
	} else {
		var newStones []uint64
		newStonesAny, contains := cache.Load(srcStone.number)

		if contains {
			newStones = newStonesAny.([]uint64)
		} else {
			newStones = updateStone(srcStone.number)
			cache.Store(srcStone.number, newStones)
		}

		for _, newStoneNumber := range newStones {
			newWork = append(newWork, WorkItem{
				number:     newStoneNumber,
				blinkCount: srcStone.blinkCount + 1})
		}

		if int32(srcStone.blinkCount + 1) > maxBlink.Load(){
			maxBlink.Store(int32(srcStone.blinkCount + 1))
			fmt.Printf("max blink count is %d\n", srcStone.blinkCount + 1)
		}
	}

	return newWork
}

func updateStone(stone uint64) []uint64 {
	newStones := make([]uint64, 0)

	if stone == 0 {
		newStones = append(newStones, 1)
	} else {
		stringifiedNum := strconv.FormatUint(stone, 10)
		numLength := len(stringifiedNum)

		if numLength % 2 == 0 {
			leftStone, _ := strconv.ParseUint(stringifiedNum[0 : numLength / 2], 10, 64)
			rightStone, _ := strconv.ParseUint(stringifiedNum[numLength / 2:], 10, 64)
			newStones = append(newStones, leftStone, rightStone)
		} else {
			newStones = append(newStones, stone * 2024)
		}
	}

	return newStones
}
