package main

import (
	"fmt"
	"strconv"

	util "github.com/holdenparker/advent-of-code-2025/util"
)

func main() {
	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
	}

	part01 := LockDial{
		Pos:    50,
		Zeroes: 0,
	}
	pf.Process = part01.PartOne
	err := pf.Go()
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", part01.Zeroes)
	}

	part02 := LockDial{
		Pos:    50,
		Zeroes: 0,
	}
	pf.Process = part02.PartTwo
	err = pf.Go()
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 2: %v\n", part02.Zeroes)
	}
}

type LockDial struct {
	Pos    int
	Zeroes int
}

func (ld *LockDial) PartOne(line string) error {
	dir, mag, err := read_line(line)
	if err != nil {
		return err
	}

	zeroes := ld.Zeroes

	ld.AdjustPos(dir, mag)

	if ld.Pos == 0 {
		zeroes += 1
	}
	ld.Zeroes = zeroes

	return nil
}

func (ld *LockDial) PartTwo(line string) error {
	dir, mag, err := read_line(line)
	if err != nil {
		return err
	}
	ld.AdjustPos(dir, mag)
	return nil
}

func (ld *LockDial) AdjustPos(dir string, mov int) {
	switch dir {
	case "R":
		ld.Pos += mov
		for ld.Pos > 99 {
			ld.Zeroes += 1
			ld.Pos -= 100
		}
	case "L":
		if ld.Pos == 0 {
			ld.Zeroes -= 1
		}
		ld.Pos -= mov
		for ld.Pos < 0 {
			ld.Zeroes += 1
			ld.Pos += 100
		}
		if ld.Pos == 0 {
			ld.Zeroes += 1
		}
	}

}

func read_line(line string) (string, int, error) {
	num, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", 0, err
	}
	return string(line[0]), num, nil
}
