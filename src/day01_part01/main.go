package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day01_part01\input.txt`

	left, right := parseData(FILEPATH)
	sum := getDistanceSum(left, right)
	fmt.Printf("The sum is %d\n", sum)
}

func parseData(filepath string) ([]int, []int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	left := make([]int, 0)
	right := make([]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		value, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		left = append(left, value)

		value, err = strconv.Atoi(parts[3])
		if err != nil {
			log.Fatal(err)
		}
		right = append(right, value)
	}

	return left, right
}

func getDistanceSum(left []int, right []int) int {
	sort.Ints(left)
	sort.Ints(right)

	var sum int = 0
	for i := range len(left) {
		sum += int(math.Round(math.Abs(float64(right[i] - left[i]))))
	}

	return sum
}
