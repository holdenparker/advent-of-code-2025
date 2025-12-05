package main

import (
	"fmt"
	"strings"

	"github.com/holdenparker/advent-of-code-2025/util"
)

func main() {
	pr := PaperRolls{}
	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  pr.NextLine,
	}
	err := pf.Go()
	// One more fake line to count the last line in the grid
	pr.NextLine("")
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", pr.Accessible)
	}

	pr.MarkAllAccessibleRolls()
	fmt.Printf("Part 2: %v\n", pr.Accessible)
}

const (
	PREVLINE   = -1
	CURRLINE   = 0
	NEXTLINE   = 1
	ACCESSIBLE = "x"
)

type PaperRolls struct {
	prevline   string
	currline   string
	nextline   string
	Accessible int
	nextpass   []string
}

func (pr *PaperRolls) NextLine(line string) error {
	pr.prevline = pr.currline
	pr.currline = pr.nextline
	pr.nextline = line

	marked := pr.MarkAccessibleRolls()
	pr.Accessible += strings.Count(marked, ACCESSIBLE)

	pr.nextpass = append(pr.nextpass, strings.ReplaceAll(marked, ACCESSIBLE, "."))
	// fmt.Println(marked)

	return nil
}

func (pr *PaperRolls) MarkAccessibleRolls() string {
	result := ""
	for i := 0; i < len(pr.currline); i++ {
		if pr.GetVal(CURRLINE, i) == "@" && pr.IsAccessible(i) {
			result += ACCESSIBLE
		} else {
			result += pr.GetVal(CURRLINE, i)
		}
	}
	return result
}

func (pr *PaperRolls) MarkAllAccessibleRolls() {
	currAccessible := 0
	for currAccessible != pr.Accessible {
		currAccessible = pr.Accessible
		nextpass := pr.nextpass[1:]
		pr.nextpass = []string{}
		pr.prevline = ""
		pr.currline = ""
		pr.nextline = ""
		for _, line := range nextpass {
			pr.NextLine(line)
		}
		pr.NextLine("")
		// fmt.Println(currAccessible, pr.Accessible)
	}
}

func (pr *PaperRolls) IsAccessible(col int) bool {
	neighbors := 0
	for _, l := range [3]int{PREVLINE, CURRLINE, NEXTLINE} {
		var cols []int
		if l == CURRLINE {
			cols = []int{col - 1, col + 1}
		} else {
			cols = []int{col - 1, col, col + 1}
		}
		for _, c := range cols {
			if pr.GetVal(l, c) == "@" {
				neighbors++
			}
		}
	}
	return neighbors < 4
}

func (pr *PaperRolls) GetVal(line int, col int) string {
	pr_line := ""
	switch line {
	case PREVLINE:
		pr_line = pr.prevline
	case CURRLINE:
		pr_line = pr.currline
	case NEXTLINE:
		pr_line = pr.nextline
	}
	if col >= len(pr_line) || col < 0 {
		return ""
	}
	return string(pr_line[col])
}
