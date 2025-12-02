package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// result1, err := part_1("test.txt")
	result1, err := part_1("data.txt")
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", result1)
	}
}

func part_1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	pos := 50
	zeroes := 0

	for scanner.Scan() {
		line := scanner.Text()
		dir, mag, err := read_line(line)
		if err != nil {
			return 0, err
		}
		pos = adjust_pos(pos, dir, mag)
		if pos == 0 {
			zeroes += 1
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return zeroes, nil
}

func read_line(line string) (string, int, error) {
	num, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", 0, err
	}
	return string(line[0]), num, nil
}

func adjust_pos(curr int, dir string, mov int) int {
	switch dir {
	case "R":
		curr += mov
	case "L":
		curr -= mov
	}

	for curr < 0 {
		curr += 100
	}
	for curr > 99 {
		curr -= 100
	}

	return curr
}
