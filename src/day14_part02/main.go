package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

const OUTPUT_FOLDER = "pics";

func main() {
	const FILEPATH string = `D:\Users\Nicolas\Documents\GoLandProjects\advent-of-code-2024\src\day14_part01\input.txt`
	var FIELD_DIMENSIONS = [2]int64{101, 103} //width, height
	var KERNEL_DIMENSIONS = [2]int64{3, 20}

	robots := parseData(FILEPATH, &FIELD_DIMENSIONS)
	easterEggSearch(robots, &FIELD_DIMENSIONS, &KERNEL_DIMENSIONS)
}

func parseData(filepath string, fieldDimensions *[2]int64) []Robot {
	robotRegex := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file for parsing")
	}
	defer file.Close()

	robots := make([]Robot, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		match := robotRegex.FindStringSubmatch(scanner.Text())

		xPosition, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		yPosition, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		xVelocity, err := strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		yVelocity, err := strconv.ParseInt(match[4], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		robot := NewRobot(
			&[2]int64{xPosition, yPosition},
			&[2]int64{xVelocity, yVelocity},
			fieldDimensions)
		robots = append(robots, *robot)
	}

	return robots
}

func easterEggSearch(robots []Robot, fieldDimensions *[2]int64, kernelDimensions *[2]int64) {
	const MAX_TIME_SEC int64 = 100000;
	var wg sync.WaitGroup;

	os.RemoveAll(OUTPUT_FOLDER);
	os.MkdirAll(OUTPUT_FOLDER, os.ModeDir);

	for index := range MAX_TIME_SEC {
		wg.Add(1)

		var time_sec = index + 1;
		robotsCopy := make([]Robot, len(robots))
		copy(robotsCopy, robots)

		go func() {
			defer wg.Done();
			occupied := moveRobots(robotsCopy, time_sec, fieldDimensions);
			match := kernelScan(occupied, fieldDimensions, kernelDimensions);

			if match {
				plotResult(time_sec, occupied);
			}
		}()
	}

	wg.Wait();
}

func moveRobots(robots []Robot, time_sec int64, fieldDimensions *[2]int64) [][]bool {
	occupied := make([][]bool, 0)

	for i := range fieldDimensions[0] {
		occupied = append(occupied, make([]bool, fieldDimensions[1]))

		for j := range fieldDimensions[1] {
			occupied[i][j] = false;
		}
	}

	for _, robot := range robots {
		robot.Move(time_sec)
		occupied[robot.Position[0]][robot.Position[1]] = true;
	}

	return occupied;
}

func kernelScan(occupied [][]bool, fieldDimensions *[2]int64, kernelDimensions *[2]int64) bool {
	for mapX := range fieldDimensions[0] {
		for mapY := range fieldDimensions[1] {
			if int64(mapX) + kernelDimensions[0] - 1 > fieldDimensions[0] - 1 || int64(mapY) + kernelDimensions[1] - 1 > fieldDimensions[1] - 1 {
				continue;
			}

			valid := true;

			for kernelX := mapX; kernelX < mapX + kernelDimensions[0]; kernelX++ {
				for kernelY := mapY; kernelY < mapY + kernelDimensions[1]; kernelY++ {
					if !occupied[kernelX][kernelY] {
						valid = false;
						break;
					}
				}

				if !valid{
					break;
				}
			}

			if valid {
				return true;
			}
		}
	}

	return false;
}

func plotResult(time_sec int64, occupied [][]bool) {
	p := plot.New()

	pts := make(plotter.XYs, 0);

	for i, _ := range occupied {
		for j, _ := range occupied[i] {
			if occupied[i][j] {
				pts = append(pts, struct{ X, Y float64 } {float64(i), float64(-j)})
			}
		}
	}

	scatter, err := plotter.NewScatter(pts)
	scatter.GlyphStyle.Color = color.Black;
	scatter.GlyphStyle.Radius = 3;
	scatter.GlyphStyle.Shape = draw.CircleGlyph{};

	if err != nil {
		log.Panic(err)
	}

	p.Add(scatter)

	err = p.Save(1000, 1000, fmt.Sprintf("%s\\%d.png", OUTPUT_FOLDER, time_sec));
	if err != nil {
		log.Panic(err)
	}
}
