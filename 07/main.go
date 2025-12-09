package main

import (
	"fmt"

	"github.com/holdenparker/advent-of-code-2025/util"
)

func main() {
	tmb := TachyonManifoldBeams{}
	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  tmb.NextLine,
	}

	err := pf.Go()
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Number of splits: %v\n", tmb.Splits)
		fmt.Printf("Number of timelines: %v\n", tmb.Timelines())
	}
}

type TachyonManifoldBeams struct {
	beams  []int
	Splits int
}

func (tmb *TachyonManifoldBeams) NextLine(line string) error {
	for len(tmb.beams) < len(line) {
		tmb.beams = append(tmb.beams, 0)
	}
	for i, c := range line {
		switch c {
		case '.':
		case '^':
			if tmb.beams[i] > 0 {
				tmb.split(i)
			}
		case 'S':
			tmb.beams[i] = 1
		}
	}
	return nil
}

func (tmb *TachyonManifoldBeams) split(i int) {
	if i-1 >= 0 {
		tmb.beams[i-1] += tmb.beams[i]
	}
	if i+1 < len(tmb.beams) {
		tmb.beams[i+1] += tmb.beams[i]
	}
	tmb.beams[i] = 0
	tmb.Splits++
}

func (tmb *TachyonManifoldBeams) Timelines() int {
	result := 0
	for _, t := range tmb.beams {
		result += t
	}
	return result
}
