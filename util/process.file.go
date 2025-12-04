package util

import (
	"bufio"
	"os"
)

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
	if err := pf.Scanner.Err(); err != nil {
		return err
	}
	return nil
}
