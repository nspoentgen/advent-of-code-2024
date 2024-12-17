package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type DATA_TYPE int16
const (
	FREE DATA_TYPE = iota -1
	BLOCK
)

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
			value = int16(FREE)
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
	charIndex := len(data) - 1

	for !isDefragmented(data) && charIndex >= 1 {
		for i := 0; i <= charIndex - 1; i++ {
			if data[i] == int16(FREE) {
				data[i] = data[charIndex]
				data[charIndex] = int16(FREE)
				charIndex--
				break
			}
		}
	}
}

func isDefragmented(data []int16) bool {
	var previousDataType DATA_TYPE
	if data[0] == int16(FREE) { previousDataType = FREE } else { previousDataType = BLOCK }

	switched := false
	for i := 1; i < len(data); i++ {
		var dataType DATA_TYPE
		if data[i] == int16(FREE) { dataType = FREE } else { dataType = BLOCK }

		if dataType != previousDataType {
			if switched {
				return false
			} else {
				switched = true
			}
		}

		previousDataType = dataType
	}

	return true
}

func calcChecksum(data []int16) int64 {
	productSum := int64(0)

	for index, value := range data {
		if value == int16(FREE) {
			break
		}

		productSum += int64(index) * int64(value)
	}

	return productSum
}
