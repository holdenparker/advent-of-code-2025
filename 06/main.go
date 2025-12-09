package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/holdenparker/advent-of-code-2025/util"
)

type CephalopodOperator string

const (
	Addition       CephalopodOperator = "+"
	Multiplication CephalopodOperator = "*"
)

type ParsingState int

const (
	FindingChars ParsingState = iota
	CharVals
	SpaceVals
	TokenParsed
)

func main() {
	ch := CephalopodHomework{}
	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  ch.NextToken,
		Split:    tokenSplit,
	}

	compSum := 0
	err := pf.Go()
	if err == nil {
		compSum, err = ch.Complete()
	}
	if err != nil {
		fmt.Printf("Error in part 1!\n%v\n", err)
	} else {
		fmt.Printf("Completed homework sum: %v\n", compSum)
	}
}

func tokenSplit(data []byte, atEof bool) (advance int, token []byte, err error) {
	state := FindingChars
	startingPos := -1
	endingPos := -1
	nextStart := -1
	for i, c := range data {
		switch c {
		case ' ':
			switch state {
			case CharVals:
				endingPos = i
				state = SpaceVals
			}
		case '\n':
			switch state {
			case FindingChars:
				startingPos = i
				endingPos = i + 1
				nextStart = i + 1
				state = TokenParsed
			case CharVals:
				endingPos = i
				nextStart = i
				state = TokenParsed
			case SpaceVals:
				nextStart = i
				state = TokenParsed
			}
		default:
			switch state {
			case FindingChars:
				startingPos = i
				state = CharVals
			case SpaceVals:
				nextStart = i
				state = TokenParsed
			}
		}
		if state == TokenParsed {
			return nextStart, data[startingPos:endingPos], nil
		}
	}
	if atEof {
		switch state {
		case CharVals:
			return len(data), data[startingPos:], nil
		case SpaceVals:
			return len(data), data[startingPos:endingPos], nil
		}
	}
	return 0, nil, nil
}

type CephalopodHomework struct {
	Problems []CephalopodProblem
	currProb int
}

func (ch *CephalopodHomework) NextToken(t string) error {
	if t == "\n" {
		ch.currProb = 0
	} else {
		for ch.currProb >= len(ch.Problems) {
			ch.Problems = append(ch.Problems, CephalopodProblem{})
		}
		switch t {
		case string(Addition):
			ch.Problems[ch.currProb].Operator = Addition
		case string(Multiplication):
			ch.Problems[ch.currProb].Operator = Multiplication
		default:
			num, err := strconv.Atoi(t)
			if err != nil {
				return err
			}
			ch.Problems[ch.currProb].Numbers = append(ch.Problems[ch.currProb].Numbers, num)
		}
		ch.currProb++
	}
	return nil
}

func (ch *CephalopodHomework) Complete() (int, error) {
	result := 0
	for _, p := range ch.Problems {
		calc, err := p.Calculate()
		if err != nil {
			return -1, err
		}
		result += calc
	}
	return result, nil
}

type CephalopodProblem struct {
	Numbers  []int
	Operator CephalopodOperator
}

func (cp *CephalopodProblem) Calculate() (int, error) {
	var result int
	switch cp.Operator {
	case Addition:
		result = 0
		for _, n := range cp.Numbers {
			result += n
		}
	case Multiplication:
		result = 1
		for _, n := range cp.Numbers {
			result *= n
		}
	default:
		return -1, errors.New(fmt.Sprintf("Unexpected operator!%v\n", cp.Operator))
	}
	return result, nil
}
