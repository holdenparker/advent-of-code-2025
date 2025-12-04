package main

import (
	"fmt"
	"math"
	"strconv"

	util "github.com/holdenparker/advent-of-code-2025/util"
)

func main() {
	batts := Batteries{}
	pf := util.ProcessFile{
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
		fmt.Printf("Largest 12 Battery Joltage Sum: %v\n", batts.Largest12JoltageSum)
	}
}

type Batteries struct {
	LargestJoltageSum   int64
	Largest12JoltageSum int64
}

func (b *Batteries) ProcessBank(bank string) error {
	maxJoltage, err := joltage_calculator(bank, 2)
	if err != nil {
		return err
	}
	b.LargestJoltageSum += maxJoltage
	maxJoltage, err = joltage_calculator(bank, 12)
	if err != nil {
		return err
	}
	b.Largest12JoltageSum += maxJoltage
	return nil
}

func joltage_calculator(bank string, batteries int) (int64, error) {
	result := int64(0)
	prevj := -1
	for i := 0; i < batteries; i++ {
		prevdig := -1
		pj := 0
		for j, c := range bank[prevj+1 : len(bank)-(batteries-i-1)] {
			n, err := strconv.Atoi(string(c))
			if err != nil {
				return 0, err
			}
			if n > prevdig {
				prevdig = n
				pj = j
			}
		}
		prevj += pj + 1
		result += int64(int(math.Pow10(batteries-i-1)) * prevdig)
	}
	return result, nil
}
