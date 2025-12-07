package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/holdenparker/advent-of-code-2025/util"
)

func main() {
	ii := IngrediantInventory{}

	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  ii.NextLine,
	}

	err := pf.Go()
	if err != nil {
		fmt.Printf("Error in part 1!\n%v\n", err)
	} else {
		fmt.Printf("Number of fresh ingrediants: %v\n", ii.NumberOfFresh)
		fmt.Printf("Number of possible fresh ingrediants: %v\n", ii.CountPossibleFreshIds())
	}
}

type InventoryParsingState int

const (
	Ranges InventoryParsingState = iota
	Ingrediants
)

type IngrediantInventory struct {
	state         InventoryParsingState
	freshRanges   [][2]int
	NumberOfFresh int
}

func (ii *IngrediantInventory) NextLine(line string) error {
	if len(line) == 0 {
		ii.state = Ingrediants
	} else {
		switch ii.state {
		case Ranges:
			return ii.ParseIdRanges(line)
		case Ingrediants:
			return ii.ProcessIngrediantId(line)
		}
	}
	return nil
}

func (ii *IngrediantInventory) ParseIdRanges(line string) error {
	ids := strings.Split(line, "-")
	if len(ids) != 2 {
		return errors.New(fmt.Sprintf("Unexpected ids len! %v", len(ids)))
	}
	startId, err := strconv.Atoi(ids[0])
	if err != nil {
		return err
	}
	endId, err := strconv.Atoi(ids[1])
	if err != nil {
		return err
	}
	ii.AddFreshIdRange(startId, endId)
	return nil
}

func (ii *IngrediantInventory) ProcessIngrediantId(line string) error {
	id, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	if ii.IdIsFresh(id) {
		ii.NumberOfFresh++
	}
	return nil
}

func (ii *IngrediantInventory) AddFreshIdRange(startId int, endId int) {
	if startId > endId {
		tmp := startId
		startId = endId
		endId = tmp
	}
	for i, ir := range ii.freshRanges {
		if ir[0] <= startId && endId <= ir[1] {
			// inside
			return
		} else if startId <= ir[0] && ir[0] <= endId && endId <= ir[1] {
			// overlap lower
			ii.freshRanges[i] = [2]int{startId, ir[1]}
			ii.Rebalance(i)
			return
		} else if ir[0] <= startId && startId <= ir[1] && ir[1] <= endId {
			// overlap upper
			ii.freshRanges[i] = [2]int{ir[0], endId}
			ii.Rebalance(i)
			return
		} else if startId <= ir[0] && ir[1] <= endId {
			// contains
			ii.freshRanges[i] = [2]int{startId, endId}
			ii.Rebalance(i)
			return
		}
	}
	// no match
	ii.freshRanges = append(ii.freshRanges, [2]int{startId, endId})
}

func (ii *IngrediantInventory) Rebalance(i int) {
	newRange := ii.freshRanges[i]
	ii.freshRanges = append(ii.freshRanges[:i], ii.freshRanges[i+1:]...)
	ii.AddFreshIdRange(newRange[0], newRange[1])
}

func (ii *IngrediantInventory) IdIsFresh(id int) bool {
	for _, ir := range ii.freshRanges {
		if ir[0] <= id && id <= ir[1] {
			return true
		}
	}
	return false
}

func (ii *IngrediantInventory) CountPossibleFreshIds() int {
	result := 0
	for _, ir := range ii.freshRanges {
		result += (ir[1] - ir[0]) + 1
	}
	return result
}
