package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PageRules struct {
	predecessorPages []int
	ancestorPages    []int
}

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day05_part01\input.txt`
	pageRules, pageRows := parseData(FILEPATH)

	middlePageRowsSum := 0

	for _, pageRow := range pageRows {
		if updateValid(pageRow, pageRules) {
			middleIndex := (len(pageRow) - 1) / 2
			middlePageRowsSum += pageRow[middleIndex]
		}
	}

	fmt.Printf("The middle page rows sum is %d\n", middlePageRowsSum)
}

func parseData(filepath string) (map[int]*PageRules, [][]int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Rules section
	var rulesExpression = regexp.MustCompile(`(\d+)\|(\d+)`)

	rules := make(map[int]*PageRules)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		match := rulesExpression.FindStringSubmatch(line)
		page, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}

		ancestor, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := rules[page]; !ok {
			rules[page] = &PageRules{}
		}

		if _, ok := rules[ancestor]; !ok {
			rules[ancestor] = &PageRules{}
		}

		rules[page].ancestorPages = append(rules[page].ancestorPages, ancestor)
		rules[ancestor].predecessorPages = append(rules[ancestor].predecessorPages, page)
	}

	//Pages section
	pages := make([][]int, 0)

	for scanner.Scan() {
		pagesRecord := make([]int, 0)

		for _, term := range strings.Split(scanner.Text(), ",") {
			page, err := strconv.Atoi(term)
			if err != nil {
				log.Fatal(err)
			}

			pagesRecord = append(pagesRecord, page)
		}

		pages = append(pages, pagesRecord)
	}

	return rules, pages
}

func updateValid(pages []int, pageRules map[int]*PageRules) bool {
	valid := true

	for i, page := range pages {
		predecessorSlice := pages[0:i]

		var ancestorSlice []int
		if i == len(pages)-1 {
			ancestorSlice = make([]int, 0)
		} else {
			ancestorSlice = pages[i+1 : len(pages)]
		}

		_, hasRule := pageRules[page]
		predecessorsValid := !hasRule || !any(pageRules[page].ancestorPages, predecessorSlice)
		ancestorsValid := !hasRule || !any(pageRules[page].predecessorPages, ancestorSlice)

		if !predecessorsValid || !ancestorsValid {
			valid = false
			break
		}
	}

	return valid
}

func any(arr []int, testSlice []int) bool {
	for _, testVal := range testSlice {
		if contains(arr, testVal) {
			return true
		}
	}

	return false
}

func contains(arr []int, testVal int) bool {
	for _, arrVal := range arr {
		if arrVal == testVal {
			return true
		}
	}

	return false
}
