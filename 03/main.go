package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	batts := Batteries{}
	pf := ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  batts.ProcessBank,
	}

	err := pf.Init()
	if err == nil {
		err = pf.Run()
	}
	if err != nil {
		fmt.Printf("Error processing batteries!\n%v\n", err)
	} else {
		fmt.Printf("Largest Joltage Sum: %v\n", batts.LargestJoltageSum)
	}
}

type ProcessFile struct {
	Filename string
	Scanner  *bufio.Scanner
	Process  func(string) error
	file     *os.File
}

func (pf *ProcessFile) Init() error {
	file, err := os.Open(pf.Filename)
	if err != nil {
		return err
	}

	pf.file = file
	pf.Scanner = bufio.NewScanner(pf.file)
	return nil
}

func (pf *ProcessFile) Run() error {
	defer pf.file.Close()

	for pf.Scanner.Scan() {
		segment := pf.Scanner.Text()
		err := pf.Process(segment)
		if err != nil {
			return err
		}
	}
	return nil
}

type Batteries struct {
	LargestJoltageSum int64
}

func (b *Batteries) ProcessBank(bank string) error {
	maxJoltage, err := largest_joltage(bank)
	if err != nil {
		return err
	}
	b.LargestJoltageSum += int64(maxJoltage)
	return nil
}

func largest_joltage(bank string) (int, error) {
	tens := 0
	tensi := -1

	for i, c := range bank[:len(bank)-1] {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			return 0, err
		}
		if n > tens {
			tens = n
			tensi = i
		}
	}

	ones := 0
	for _, c := range bank[tensi+1:] {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			return 0, err
		}
		if n > ones {
			ones = n
		}
	}

	return (tens * 10) + ones, nil
}
