package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// filename := "test.txt"
	filename := "data.txt"

	part01 := LockDial{
		Pos:    50,
		Zeroes: 0,
	}
	err := process_file(filename, part01.PartOne)
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", part01.Zeroes)
	}

	part02 := LockDial{
		Pos:    50,
		Zeroes: 0,
	}
	err = process_file(filename, part02.PartTwo)
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 2: %v\n", part02.Zeroes)
	}
}

type ProcessLine func(string, int) error

func process_file(filename string, procl ProcessLine) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		dir, mag, err := read_line(line)
		if err != nil {
			return err
		}
		err = procl(dir, mag)

		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type LockDial struct {
	Pos    int
	Zeroes int
}

func (ld *LockDial) PartOne(dir string, mag int) error {
	zeroes := ld.Zeroes

	ld.AdjustPos(dir, mag)

	if ld.Pos == 0 {
		zeroes += 1
	}
	ld.Zeroes = zeroes

	return nil
}

func (ld *LockDial) PartTwo(dir string, mag int) error {
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
