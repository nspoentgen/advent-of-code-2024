package main

import (
	"log"
	"math"
)

const NULL_OUTPUT int = -1

func decodeOperand(operand int, registers *[3]int, literalExpected bool) int {
	if operand >= 0 && operand <= 3 {
		if !literalExpected{
			log.Fatal("Expected a literal operand but got a combo operand")
		}

		return operand
	} else if operand >= 4 && operand <= 6 {
		if literalExpected {
			log.Fatal("Expected a combo operand but got a literal operand")
		}

		return registers[operand%4]
	}

	panic("Invalid operand")
}

type Opcode interface {
	execute(operand int, registers *[3]int) (int, int)
}

type adv struct {}

func (a adv) execute(operand int, registers *[3]int) (int, int) {
	registers[0] = registers[0] / int(math.Round(math.Pow(2, float64(decodeOperand(operand, registers, false)))))
	return NULL_OUTPUT, NULL_OUTPUT
}

type bxl struct {}

func (b bxl) execute(operand int, registers *[3]int) (int, int) {
	registers[1] |= decodeOperand(operand, registers, true)
	return NULL_OUTPUT, NULL_OUTPUT
}

type bst struct {}

func (b bst) execute(operand int, registers *[3]int) (int, int) {
	registers[1] = decodeOperand(operand, registers, false) % 8
	return NULL_OUTPUT, NULL_OUTPUT
}

type jnz struct {}

func (j jnz) execute(operand int, registers *[3]int) (int, int) {
	if registers[0] == 0 {
		return NULL_OUTPUT, NULL_OUTPUT
	}

	return NULL_OUTPUT, decodeOperand(operand, registers, true)
}

type bxc struct {}

func (b bxc) execute(operand int, registers *[3]int) (int, int) {
	registers[1] |= registers[2]
	return NULL_OUTPUT, NULL_OUTPUT
}

type out struct {}

func (o out) execute(operand int, registers *[3]int) (int, int) {
	return decodeOperand(operand, registers, false) % 8, NULL_OUTPUT
}

type bdv struct {}

func (b bdv) execute(operand int, registers *[3]int) (int, int) {
	registers[1] = registers[0] / int(math.Round(math.Pow(2, float64(decodeOperand(operand, registers, false)))))
	return NULL_OUTPUT, NULL_OUTPUT
}

type cdv struct {}

func (c cdv) execute(operand int, registers *[3]int) (int, int) {
	registers[2] = registers[0] / int(math.Round(math.Pow(2, float64(decodeOperand(operand, registers, false)))))
	return NULL_OUTPUT, NULL_OUTPUT
}
