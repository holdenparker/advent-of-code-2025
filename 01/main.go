package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	part01 := LockDial{
		Pos:    50,
		Zeroes: 0,
	}
	// err := process_file("test.txt", part01.PartOne)
	err := process_file("data.txt", part01.PartOne)
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", part01.Zeroes)
	}
}

type ProcessLine func(string) error

func process_file(filename string, procl ProcessLine) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		err = procl(line)

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

func (ld *LockDial) PartOne(line string) error {
	zeroes := ld.Zeroes
	dir, mag, err := read_line(line)
	if err != nil {
		return err
	}

	ld.AdjustPos(dir, mag)

	if ld.Pos == 0 {
		zeroes += 1
	}
	ld.Zeroes = zeroes

	return nil
}

func (ld *LockDial) AdjustPos(dir string, mov int) {
	switch dir {
	case "R":
		ld.Pos += mov
	case "L":
		ld.Pos -= mov
	}

	for ld.Pos < 0 {
		ld.Zeroes += 1
		ld.Pos += 100
	}
	for ld.Pos > 99 {
		ld.Zeroes += 1
		ld.Pos -= 100
	}
}

func read_line(line string) (string, int, error) {
	num, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", 0, err
	}
	return string(line[0]), num, nil
}
