package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	// filename := "test.txt"
	filename := "data.txt"

	partOne := ProductValidation{}
	err := process_file(filename, partOne.PartOne)
	if err != nil {
		fmt.Printf("Error with part 1!\n%v\n", err)
	} else {
		fmt.Printf("Part 1: %v\n", partOne.InvalidSum)
	}

}

type ProcessSegment func(string) error

func process_file(filename string, procs ProcessSegment) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(commaSplit)

	for scanner.Scan() {
		segment := scanner.Text()
		err = procs(segment)
		if err != nil {
			return err
		}
	}

	return nil
}

func commaSplit(data []byte, atEof bool) (advance int, token []byte, err error) {
	if i := bytes.IndexRune(data, ','); i >= 0 {
		// return the chunk
		return i + 1, data[0:i], nil
	}
	if atEof {
		if len(data) == 0 {
			return 0, nil, nil
		}
		// return the last data
		return len(data), data, nil
	}
	// request more data
	return 0, nil, nil
}

type ProductValidation struct {
	InvalidSum int64
}

func (pv *ProductValidation) PartOne(seg string) error {
	ids := strings.Split(seg, "-")
	if len(ids) != 2 {
		return errors.New("There should be exactly 2 ids!")
	}
	idStart, err := strconv.Atoi(ids[0])
	if err != nil {
		return err
	}
	idEnd, err := strconv.Atoi(ids[1])
	if err != nil {
		return err
	}
	for id := idStart; id <= idEnd; id++ {
		if isInvalid(id) {
			pv.InvalidSum += int64(id)
		}
	}
	return nil
}

func isInvalid(id int) bool {
	idLen := intLen(id)
	halfLen := int(idLen / 2)

	firstHalf := id
	for i := 0; i < halfLen; i++ {
		firstHalf /= 10
	}

	return firstHalf == (id - (firstHalf * int(math.Pow10(halfLen))))
}

func intLen(num int) int {
	if num == 0 {
		return 1
	}
	result := 0
	for num > 0 {
		num /= 10
		result++
	}
	return result
}
