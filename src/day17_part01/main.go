package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
)

func main() {
	const INPUT_FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day17_part01\test_input.txt`

	initialRegisterValues, program := parseData(INPUT_FILEPATH)
	output := executeProgram(initialRegisterValues, program)

	log.Print("The output is ")
	log.Print(output)
}

func parseData(filepath string) (*[3]int, []int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file")
	}

	defer file.Close()

	registerValues := [3]int{}
	scanner := bufio.NewScanner(file)

	//Register values
	scanner.Scan()
	line := scanner.Text()
	terms := strings.Split(line, ":")
	registerValue, err := strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register A value")
	}

	registerValues[0] = registerValue

	scanner.Scan()
	line = scanner.Text()
	terms = strings.Split(line, ":")
	registerValue, err = strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register B value")
	}

	registerValues[1] = registerValue

	scanner.Scan()
	line = scanner.Text()
	terms = strings.Split(line, ":")
	registerValue, err = strconv.Atoi(strings.Trim(terms[1], " "))

	if err != nil {
		log.Panic("Could not parse register C value")
	}

	registerValues[2] = registerValue
	scanner.Scan()

	//Program
	program := make([]int, 0)
	line = scanner.Text()
	terms = strings.Split(line, ":")
	programRaw := strings.Trim(terms[1], " ")

	for _, inputRaw := range strings.Split(programRaw, ",") {
		input, err := strconv.Atoi(inputRaw)
		
		if err != nil {
			log.Panic("Could not parse program input")
		}
		
		program = append(program, input)
	}

	return &registerValues, program
}

func executeProgram(initialRegisterValues *[3]int, program []int) []int {
	var registers [3]int = *initialRegisterValues
	output := make([]int, 0)
	instructionPointer := 0

	commandTable := [8]Opcode{
		adv{},
		bxl{},
		bst{},
		jnz{},
		bxc{},
		out{},
		bdv{},
		cdv{}}

	for instructionPointer <= len(program) - 1 {
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
