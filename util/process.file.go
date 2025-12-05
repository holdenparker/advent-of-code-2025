package util

import (
	"bufio"
	"os"
)

type ProcessFile struct {
	Filename string
	Scanner  *bufio.Scanner
	Process  func(string) error
	Split    bufio.SplitFunc
}

func (pf *ProcessFile) Go() error {
	file, err := os.Open(pf.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pf.Scanner = bufio.NewScanner(file)
	if pf.Split != nil {
		pf.Scanner.Split(pf.Split)
	}

	for pf.Scanner.Scan() {
		segment := pf.Scanner.Text()
		err := pf.Process(segment)
		if err != nil {
			return err
		}
	}
	if err := pf.Scanner.Err(); err != nil {
		return err
	}
	return nil
}
