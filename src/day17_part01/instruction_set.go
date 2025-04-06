package main

import (
	"math"
)

const NULL_OUTPUT int64 = -1

func decodeComboOperand(operand int64, registers *[3]int64) int64 {
	if operand >= 0 && operand <= 3 {
		return operand
	} else if operand >= 4 && operand <= 6 {
		return registers[operand%4]
	}

	panic("Invalid operand")
}

type Opcode interface {
	execute(operand int64, registers *[3]int64) (int64, int64)
}

type adv struct {}

func (a adv) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[0] = registers[0] / int64(math.Round(math.Pow(2, float64(decodeComboOperand(operand, registers)))))
	return NULL_OUTPUT, NULL_OUTPUT
}

type bxl struct {}

func (b bxl) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[1] ^= operand
	return NULL_OUTPUT, NULL_OUTPUT
}

type bst struct {}

func (b bst) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[1] = decodeComboOperand(operand, registers) % 8
	return NULL_OUTPUT, NULL_OUTPUT
}

type jnz struct {}

func (j jnz) execute(operand int64, registers *[3]int64) (int64, int64) {
	if registers[0] == 0 {
		return NULL_OUTPUT, NULL_OUTPUT
	}

	return NULL_OUTPUT, operand
}

type bxc struct {}

func (b bxc) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[1] ^= registers[2]
	return NULL_OUTPUT, NULL_OUTPUT
}

type out struct {}

func (o out) execute(operand int64, registers *[3]int64) (int64, int64) {
	return decodeComboOperand(operand, registers) % 8, NULL_OUTPUT
}

type bdv struct {}

func (b bdv) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[1] = registers[0] / int64(math.Round(math.Pow(2, float64(decodeComboOperand(operand, registers)))))
	return NULL_OUTPUT, NULL_OUTPUT
}

type cdv struct {}

func (c cdv) execute(operand int64, registers *[3]int64) (int64, int64) {
	registers[2] = registers[0] / int64(math.Round(math.Pow(2, float64(decodeComboOperand(operand, registers)))))
	return NULL_OUTPUT, NULL_OUTPUT
}
