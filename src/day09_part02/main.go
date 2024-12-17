package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const FREE_ID int16 = -1

type DISK_CHUNK struct {
	startIndex int
	length int
	id int16
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day09_part01\input.txt`
	data := parseData(FILEPATH)
	defragment(data)
	checksum := calcChecksum(data)
	fmt.Printf("The checksum is %d\n", checksum)
}

func parseData(filepath string) []int16 {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	diskMap := make([]int16, 0, 10000)
	var id int16 = 0
	isBlockData := true

	//We are told problem is one line
	scanner := bufio.NewReader(file)
	char, _, err := scanner.ReadRune()

	for err == nil {
		size, _ := strconv.Atoi(string(char))

		var value int16
		if isBlockData {
			value = id
			id++
		} else {
			value = FREE_ID
		}

		for range size {
			diskMap = append(diskMap, value)
		}

		isBlockData = !isBlockData
		char, _, err = scanner.ReadRune()
	}

	return diskMap
}

func defragment(data []int16) {
	spaceMap := calcSpaceMap(data)

	for i := len(spaceMap) - 1; i >= 1; i-- {
		if spaceMap[i].id != FREE_ID {
			for j := 0; j < i; j++ {
				if spaceMap[j].id == FREE_ID && spaceMap[j].length >= spaceMap[i].length {
					srcStartIndex := spaceMap[i].startIndex
					srcEndIndex := spaceMap[i].startIndex + spaceMap[i].length
					dstStartIndex := spaceMap[j].startIndex
					dstEndIndex := dstStartIndex + spaceMap[i].length

					copy(data[dstStartIndex:dstEndIndex], data[srcStartIndex:srcEndIndex])
					for k := srcStartIndex; k < srcEndIndex; k++ {
						data[k] = FREE_ID
					}

					spaceMap[j].startIndex += spaceMap[i].length
					spaceMap[j].length -= spaceMap[i].length
					break
				}
			}
		}

		foo := 1
		foo += 1
	}
}

func calcSpaceMap(data []int16) []DISK_CHUNK {
	lastId := data[0]
	lastBlockIndex := 0
	diskMap := make([]DISK_CHUNK, 0, 1000)

	for index := 0; index < len(data); index++ {
		if data[index] != lastId {
			diskMap = append(diskMap, DISK_CHUNK{lastBlockIndex, index - lastBlockIndex, lastId})
			lastBlockIndex = index
			lastId = data[lastBlockIndex]
		}
	}

	diskMap = append(diskMap, DISK_CHUNK{lastBlockIndex, len(data) - lastBlockIndex, lastId})
	return diskMap
}

func calcChecksum(data []int16) int64 {
	productSum := int64(0)

	for index, value := range data {
		if value == FREE_ID { continue }

		productSum += int64(index) * int64(value)
	}

	return productSum
}
