package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const OPEN_SPOT rune = '.'

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day08_part01\input.txt`
	cityMap, numRows, numCols := parseData(FILEPATH)
	antinodePositions := getAntinodePositions(cityMap, numRows, numCols)
	fmt.Printf("The number of antinode positions = %d\n", len(antinodePositions))
}

func parseData(filepath string) (map[rune][][2]int, int, int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make(map[rune][][2]int)
	row := 0
	numCols := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if row == 0 {
			numCols = len(line)
		}

		for col, symbol := range line {
			if symbol != OPEN_SPOT {
				if _, exists := data[symbol]; !exists {
					data[symbol] = make([][2]int, 0)
				}

				data[symbol] = append(data[symbol], [2]int{row, col})
			}
		}

		row++
	}

	return data, row, numCols
}

func getAntinodePositions(cityMap map[rune][][2]int, numRows int, numCols int) map[[2]int]byte {
	antinodePositions := make(map[[2]int]byte)

	for _, antennaPositions := range cityMap {
		for i := 0; i < len(antennaPositions) - 1; i++ {
			for j := i + 1; j < len(antennaPositions); j++ {
				antennaPairAntinodePositions := getAntennaPairAntinodePositions(antennaPositions[i], antennaPositions[j])
				for _, position := range antennaPairAntinodePositions {
					if position[0] >= 0 && position[0] < numRows &&
					   position[1] >= 0 && position[1] < numCols {
						antinodePositions[position] = 0
					}
				}
			}
		}
	}

	return antinodePositions
}

func getAntennaPairAntinodePositions(antenna1Pos [2]int, antenna2Pos [2]int) [2][2]int {
	gradient := [2]int { antenna2Pos[0] - antenna1Pos[0], antenna2Pos[1] - antenna1Pos[1] }
	return [2][2]int {
		{antenna2Pos[0] + gradient[0], antenna2Pos[1] + gradient[1]},
		{antenna1Pos[0] - gradient[0], antenna1Pos[1] - gradient[1]},
	}
}
