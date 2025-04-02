package main

import (
	"log"
)

const (
	UP    = '^'
	RIGHT = '>'
	DOWN  = 'v'
	LEFT  = '<'
)

const (
	EMPTY = iota
	WALL
	ROBOT
)


//INTERFACES
type IObject interface {
	TryMove(direction rune, warehouseMap [][]int, objects map[int]IObject) (bool, [2]int)
	CanMove(direction rune, warehouseMap [][]int, objects map[int]IObject) bool
	Move(direction rune, warehouseMap [][]int, objects map[int]IObject)
	GetId() int
}


//ROBOT
type Robot struct {
	Position [2]int
	Id int
}

func (r *Robot) GetId() int {
	return r.Id
}

func (r *Robot) CanMove(direction rune, warehouseMap [][]int, objects map[int]IObject) bool {
	var rowOffset int
	var colOffset int

	if direction == UP {
		rowOffset = -1
		colOffset = 0
	} else if direction == RIGHT {
		rowOffset = 0
		colOffset = 1
	} else if direction == DOWN {
		rowOffset = 1
		colOffset = 0
	} else {
		rowOffset = 0
		colOffset = -1
	}

	row := r.Position[0] + rowOffset
	col := r.Position[1] + colOffset

	var canMove bool

	if warehouseMap[row][col] == EMPTY {
		canMove = true
	} else if warehouseMap[row][col] == WALL {
		canMove = false
	} else if warehouseMap[row][col] > ROBOT {
		box := objects[warehouseMap[row][col]]
		canMove = box.CanMove(direction, warehouseMap, objects)
	} else {
		log.Panic("Invalid map symbol for pushing")
	}

	return canMove
}

func (r *Robot) TryMove(direction rune, warehouseMap [][]int, objects map[int]IObject) (bool, [2]int) {
	canMove := r.CanMove(direction, warehouseMap, objects)
	
	if canMove {
		r.Move(direction, warehouseMap, objects)
	}

	return canMove, r.Position
}

func (r *Robot) Move(direction rune, warehouseMap [][]int, objects map[int]IObject) {
	newPosition := r.calcPosition(direction)

	objectId := warehouseMap[newPosition[0]][newPosition[1]]
	if objectId >= ROBOT {
		object := objects[objectId]
		object.Move(direction, warehouseMap, objects)
	}

	warehouseMap[newPosition[0]][newPosition[1]] = ROBOT
	warehouseMap[r.Position[0]][r.Position[1]] = EMPTY
	
	r.Position[0] = newPosition[0]
	r.Position[1] = newPosition[1]
}

func (r *Robot) calcPosition(direction rune) [2]int {
	var rowOffset int
	var colOffset int

	if direction == UP {
		rowOffset = -1
		colOffset = 0
	} else if direction == RIGHT {
		rowOffset = 0
		colOffset = 1
	} else if direction == DOWN {
		rowOffset = 1
		colOffset = 0
	} else {
		rowOffset = 0
		colOffset = -1
	}

	return [2]int {r.Position[0] + rowOffset, r.Position[1] + colOffset}
}



//BOX
type Box struct {
	Position [2]int //left half of box. Width is always 2 and height always 1
	Id int
}

func (b *Box) GetId() int {
	return b.Id
}

func (b *Box) CanMove(direction rune, warehouseMap [][]int, objects map[int]IObject) bool {
	return b.halfCanMove(true, direction, warehouseMap, objects) && b.halfCanMove(false, direction, warehouseMap, objects)
}

func (b *Box) halfCanMove(leftSide bool, direction rune, warehouseMap [][]int, objects map[int]IObject) bool {
	position := b.calcPosition(leftSide, direction)
	var canMove bool

	if warehouseMap[position[0]][position[1]] == EMPTY {
		canMove = true
	} else if warehouseMap[position[0]][position[1]] == WALL {
		canMove = false
	} else if warehouseMap[position[0]][position[1]] > ROBOT {
		object := objects[warehouseMap[position[0]][position[1]]]

		if object.GetId() == b.Id {
			canMove = true
		} else {
			canMove = object.CanMove(direction, warehouseMap, objects)
		}
	} else {
		log.Panic("Invalid map symbol for pushing")
	}

	return canMove
}

func (b *Box) calcPosition(leftSide bool, direction rune) [2]int {
	var rowOffset int
	var colOffset int

	if direction == UP {
		rowOffset = -1
		colOffset = 0
	} else if direction == RIGHT {
		rowOffset = 0
		colOffset = 1
	} else if direction == DOWN {
		rowOffset = 1
		colOffset = 0
	} else {
		rowOffset = 0
		colOffset = -1
	}

	var position [2]int

	if leftSide {
		position = [2]int {b.Position[0] + rowOffset, b.Position[1] + colOffset}
	} else {
		position = [2]int {b.Position[0] + rowOffset, b.Position[1] + colOffset + 1}
	}

	return position
}

func (b *Box) TryMove(direction rune, warehouseMap [][]int, objects map[int]IObject) (bool, [2]int) {
	canMove := b.CanMove(direction, warehouseMap, objects)
	
	if canMove {
		b.Move(direction, warehouseMap, objects)
	}

	return canMove, b.Position
}

func (b *Box) Move(direction rune, warehouseMap [][]int, objects map[int]IObject) {
	newPosition := b.calcPosition(true, direction)

	objectId := warehouseMap[newPosition[0]][newPosition[1]]
	if objectId != b.Id && objectId >= ROBOT{
		object := objects[objectId]
		object.Move(direction, warehouseMap, objects)
	}

	objectId2 := warehouseMap[newPosition[0]][newPosition[1] + 1]
	if objectId2 != b.Id && objectId2 >= ROBOT && (direction == UP || direction == DOWN) && objectId2 != objectId {
		object2 := objects[objectId2]
		object2.Move(direction, warehouseMap, objects)
	}
	
	warehouseMap[newPosition[0]][newPosition[1]] = b.Id
	warehouseMap[newPosition[0]][newPosition[1] + 1] = b.Id

	if direction == UP || direction == DOWN {
		warehouseMap[b.Position[0]][b.Position[1]] = EMPTY
		warehouseMap[b.Position[0]][b.Position[1] + 1] = EMPTY
	} else if direction == LEFT {
		warehouseMap[b.Position[0]][b.Position[1] + 1] = EMPTY
	} else {
		warehouseMap[b.Position[0]][b.Position[1]] = EMPTY
	}

	b.Position[0] = newPosition[0]
	b.Position[1] = newPosition[1]
}
