package main

type Robot struct {
	Position        [2]int //+x = right, +y = down, origin upper-left corner
	Velocity        [2]int //+x = right, +y = down
	fieldDimensions [2]int //[width, height]
}

func NewRobot(initialPosition *[2]int, velocity *[2]int, fieldDimensions *[2]int) *Robot {
	return &Robot{
		Position:        *initialPosition,
		Velocity:        *velocity,
		fieldDimensions: *fieldDimensions}
}

func (r *Robot) Move() {
	r.Position[0] = (r.Position[0] + r.Velocity[0] + r.fieldDimensions[0]) % r.fieldDimensions[0]
	r.Position[1] = (r.Position[1] + r.Velocity[1] + r.fieldDimensions[1]) % r.fieldDimensions[1]
}
