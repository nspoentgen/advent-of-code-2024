package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
)

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day17_part01\input.txt`

	initialRegisterValues, program := parseData(INPUT_FILEPATH)
	output := executeProgram(initialRegisterValues, program)

	stringOutput := make([]string, 0)
	for _, element := range output {
		stringOutput = append(stringOutput, strconv.Itoa(int(element)))
	}

	log.Printf("Output: %s\n", strings.Join(stringOutput, ","))
}

func parseData(filepath string) (*[3]int64, []int64) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file")
	}

	defer file.Close()

	registerValues := [3]int64{}
	scanner := bufio.NewScanner(file)

	//Register values
	scanner.Scan()
	line := scanner.Text()
	terms := strings.Split(line, ":")
	registerValue, err := strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register A value")
	}

	registerValues[0] = int64(registerValue)

	scanner.Scan()
	line = scanner.Text()
	terms = strings.Split(line, ":")
	registerValue, err = strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register B value")
	}

	registerValues[1] = int64(registerValue)

	scanner.Scan()
	line = scanner.Text()
	terms = strings.Split(line, ":")
	registerValue, err = strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register C value")
	}

	registerValues[2] = int64(registerValue)
	scanner.Scan()
	scanner.Scan()

	//Program
	program := make([]int64, 0)
	line = scanner.Text()
	terms = strings.Split(line, ":")
	programRaw := strings.Trim(terms[1], " ")

	for _, inputRaw := range strings.Split(programRaw, ",") {
		input, err := strconv.Atoi(inputRaw)
		
		if err != nil {
			log.Panic("Could not parse program input")
		}
		
		program = append(program, int64(input))
	}

	return &registerValues, program
}

func executeProgram(initialRegisterValues *[3]int64, program []int64) []int64 {
	var registers [3]int64 = *initialRegisterValues
	output := make([]int64, 0)
	var instructionPointer int64 = 0

	commandTable := [8]Opcode{
		adv{},
		bxl{},
		bst{},
		jnz{},
		bxc{},
		out{},
		bdv{},
		cdv{}}

	for instructionPointer <= int64(len(program) - 1){
		instruction := program[instructionPointer]
		operand := program[instructionPointer + 1]

		commandOutput, instructionPointerJump := commandTable[instruction].execute(operand, &registers)
		
		if commandOutput != NULL_OUTPUT {
			output = append(output, commandOutput)
		}

		if instructionPointerJump != NULL_OUTPUT {
			instructionPointer = instructionPointerJump
		} else {
			instructionPointer += 2
		}
	}

	return output
}
