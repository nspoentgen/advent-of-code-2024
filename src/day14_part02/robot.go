package main

type Robot struct {
	Position        [2]int64 //+x = right, +y = down, origin upper-left corner
	Velocity        [2]int64 //+x = right, +y = down
	fieldDimensions [2]int64 //[width, height]
}

func NewRobot(initialPosition *[2]int64, velocity *[2]int64, fieldDimensions *[2]int64) *Robot {
	return &Robot{
		Position:        *initialPosition,
		Velocity:        *velocity,
		fieldDimensions: *fieldDimensions}
}

func (r *Robot) Move(time_sec int64) {
	r.Position[0] = (r.Position[0] + time_sec*(r.Velocity[0]+r.fieldDimensions[0])) % r.fieldDimensions[0]
	r.Position[1] = (r.Position[1] + time_sec*(r.Velocity[1]+r.fieldDimensions[1])) % r.fieldDimensions[1]
}
