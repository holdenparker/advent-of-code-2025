package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/holdenparker/advent-of-code-2025/util"
)

type CephalopodOperator string

const (
	Addition       CephalopodOperator = "+"
	Multiplication CephalopodOperator = "*"
)

type PreprocessState int

const (
	Starting PreprocessState = iota
	CharsFound
	OperatorFound
	OperatorParsed
)

const EMPTY_TOKEN = "empty-token"

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
		Process:  ch.PreprocessToken,
		Split:    preprocessParsing,
	}

	// Preprocessing to get maxlen for each problem
	err := pf.Go()
	if err != nil {
		fmt.Printf("Error preprocessing!\n%v\n", err)
		return
	}

	// Initial load for part 1 of the task
	ch.ResetStrings()
	ch.ResetNumbers()
	pf.Process = ch.NextToken
	pf.Split = ch.tokenSplit

	compSum := 0
	err = pf.Go()
	if err == nil {
		compSum, err = ch.Complete()
	}
	if err != nil {
		fmt.Printf("Error in part 1!\n%v\n", err)
	} else {
		fmt.Printf("Completed homework sum: %v\n", compSum)
	}

	// Reset numbers and transpose for part 2 of the task
	ch.ResetNumbers()

	ch.Transpose()
	compSum, err = ch.Complete()
	if err != nil {
		fmt.Printf("Error in part 2!\n%v\n", err)
	} else {
		fmt.Printf("Completed transposed homework sum: %v\n", compSum)
	}
}

func preprocessParsing(data []byte, atEof bool) (advance int, token []byte, err error) {
	state := Starting
	startingPos := -1
	endingPos := -1
	for i, c := range data {
		switch c {
		case '*', '+':
			switch state {
			case Starting, CharsFound:
				startingPos = i
				state = OperatorFound
			case OperatorFound:
				endingPos = i - 1
				state = OperatorParsed
			}
		case ' ', '\n':
		default:
			state = CharsFound
		}
		if state == OperatorParsed {
			return endingPos + 1, data[startingPos:endingPos], nil
		}
	}
	if atEof && state == OperatorFound {
		return len(data), data[startingPos:], nil
	}
	if state == CharsFound {
		return len(data), []byte(""), nil
	}
	return 0, nil, nil
}

func (ch *CephalopodHomework) tokenSplit(data []byte, atEof bool) (advance int, token []byte, err error) {
	maxLen := 0
	if ch.currProb >= len(ch.problems) {
		maxLen = len(data) + 1
	} else {
		maxLen = ch.problems[ch.currProb].MaxLen
	}
	for i, c := range data {
		switch c {
		case '\n':
			if i > 0 {
				return i, data[:i], nil
			}
			return 1, []byte("\n"), nil
		default:
			if i == maxLen {
				return i + 1, data[:i], nil
			}
		}
	}
	if atEof && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

type CephalopodHomework struct {
	problems []CephalopodProblem
	currProb int
}

func (ch *CephalopodHomework) PreprocessToken(t string) error {
	if len(t) > 0 {
		for ch.currProb >= len(ch.problems) {
			ch.problems = append(ch.problems, CephalopodProblem{})
		}
		switch string(t[0]) {
		case string(Addition):
			ch.problems[ch.currProb].Operator = Addition
		case string(Multiplication):
			ch.problems[ch.currProb].Operator = Multiplication
		}
		ch.problems[ch.currProb].MaxLen = len(t)
		ch.currProb++
	}
	return nil
}

func (ch *CephalopodHomework) NextToken(t string) error {
	if t == "\n" {
		ch.currProb = 0
	} else if len(t) > 0 {
		// This switch has been reduced to filtering out the ops
		switch strings.Trim(t, " ") {
		// We've already gathered the op in preprocessing, do nothing here
		case string(Addition), string(Multiplication):
		default:
			ch.problems[ch.currProb].Strings = append(ch.problems[ch.currProb].Strings, t)
		}
		ch.currProb++
	}
	return nil
}

func (ch *CephalopodHomework) ResetStrings() {
	ch.currProb = 0
	for i := range ch.problems {
		ch.problems[i].Strings = []string{}
	}
}

func (ch *CephalopodHomework) ResetNumbers() {
	for i := range ch.problems {
		ch.problems[i].Numbers = []int{}
	}
}

func (ch *CephalopodHomework) Complete() (int, error) {
	result := 0
	for i := range ch.problems {
		err := ch.problems[i].Atoi()
		if err != nil {
			return -1, err
		}
		calc, err := ch.problems[i].Calculate()
		if err != nil {
			return -1, err
		}
		result += calc
	}
	return result, nil
}

func (ch *CephalopodHomework) Transpose() {
	for i := range ch.problems {
		ch.problems[i].Transpose()
	}
}

type CephalopodProblem struct {
	Strings  []string
	Numbers  []int
	Operator CephalopodOperator
	MaxLen   int
}

func (cp *CephalopodProblem) Atoi() error {
	for _, str := range cp.Strings {
		n, err := strconv.Atoi(strings.Trim(str, " "))
		if err != nil {
			return err
		}
		cp.Numbers = append(cp.Numbers, n)
	}
	return nil
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

func (cp *CephalopodProblem) Transpose() {
	transposedStrings := []string{}
	for i := cp.MaxLen - 1; i >= 0; i-- {
		str := ""
		for _, s := range cp.Strings {
			str += string(s[i])
		}
		transposedStrings = append(transposedStrings, str)
	}
	cp.Strings = transposedStrings
}
