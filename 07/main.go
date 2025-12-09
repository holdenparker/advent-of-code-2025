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
	}
}

type TachyonManifoldBeams struct {
	beams  []bool
	Splits int
}

func (tmb *TachyonManifoldBeams) NextLine(line string) error {
	for len(tmb.beams) < len(line) {
		tmb.beams = append(tmb.beams, false)
	}
	for i, c := range line {
		switch c {
		case '.':
		case '^':
			if tmb.beams[i] {
				tmb.split(i)
			}
		case 'S':
			tmb.beams[i] = true
		}
	}
	return nil
}

func (tmb *TachyonManifoldBeams) split(i int) {
	tmb.beams[i] = false
	if i-1 >= 0 {
		tmb.beams[i-1] = true
	}
	if i+1 < len(tmb.beams) {
		tmb.beams[i+1] = true
	}
	tmb.Splits++
}
