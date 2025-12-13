package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/holdenparker/advent-of-code-2025/util"
)

var (
	MULTIPLE_CIRCUITS_ERR = errors.New("Found multiple circuits for a junction box!")
	CIRCUIT_LOOP_ERR      = errors.New("Adding this edge creates a circuit loop!")
	LINESTRING_LOOP_ERR   = errors.New("Linestring contains the same junction box twice!")
	ALREADY_CONNECTED_ERR = errors.New("JunctionBoxes are already connected!")
)

func main() {
	pg := Playground{}
	pf := util.ProcessFile{
		// Filename: "test.txt",
		Filename: "data.txt",
		Process:  pg.NextLine,
	}

	err := pf.Go()
	if err != nil {
		fmt.Printf("Error in part 1!\n%v\n", err)
	} else {
		maxIter := 0
		if pf.Filename == "test.txt" {
			maxIter = 10
		} else {
			maxIter = 1000
		}
		err = pg.BuildCircuits(maxIter)
		if err != nil {
			fmt.Printf("Error in iterating part 1!\n%v\n", err)
		} else {
			fmt.Printf("The three largest circuit product is: %v\n", pg.ThreeLargestProduct())
			err = pg.BuildSingleCircuit()
			if err != nil {
				fmt.Printf("Error in iterating part 2!\n%v\n", err)
			} else {
				fmt.Printf("The product of the last two X values is: %f\n", pg.LastAddedXProduct())
			}
		}
	}
}

type Playground struct {
	junctionBoxes  []*JunctionBox
	shortestLights []*LightString
	circuits       [][]*JunctionBox
	lastConnected  *LightString
}

func (pg *Playground) NextLine(line string) error {
	coords := strings.Split(line, ",")
	if len(coords) != 3 {
		return errors.New(fmt.Sprintf("Coords should have exactly three values!  Found %v\n[%v]\n", len(coords), line))
	}
	x, err := strconv.Atoi(coords[0])
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(coords[1])
	if err != nil {
		return err
	}
	z, err := strconv.Atoi(coords[2])
	if err != nil {
		return err
	}
	jb := &JunctionBox{
		X: float64(x),
		Y: float64(y),
		Z: float64(z),
	}
	pg.junctionBoxes = append(pg.junctionBoxes, jb)
	pg.shortestLights = append(pg.shortestLights, &LightString{
		Junctions: [2]*JunctionBox{jb, jb},
		Length:    math.Inf(1),
	})
	return nil
}

func (pg *Playground) ThreeLargestProduct() int {
	lens := []float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for i := range pg.circuits {
		l := float64(len(pg.circuits[i]))
		if l > lens[2] {
			lens[2] = l
			sort.Slice(lens, func(i, j int) bool {
				return lens[i] > lens[j]
			})
		}
	}
	return int(lens[0]) * int(lens[1]) * int(lens[2])
}

func (pg *Playground) LastAddedXProduct() float64 {
	return pg.lastConnected.Junctions[0].X * pg.lastConnected.Junctions[1].X
}

func (pg *Playground) BuildCircuits(iterations int) error {
	for i := 0; i < iterations; i++ {
		err := pg.Iterate()
		if err != nil && err != CIRCUIT_LOOP_ERR {
			return err
		}
	}
	return nil
}

func (pg *Playground) BuildSingleCircuit() error {
	sort.Slice(pg.circuits, func(i, j int) bool {
		return len(pg.circuits[i]) > len(pg.circuits[j])
	})
	maxLen := len(pg.junctionBoxes)
	for len(pg.circuits[0]) < maxLen {
		err := pg.Iterate()
		if err != nil && err != CIRCUIT_LOOP_ERR {
			return err
		}
		sort.Slice(pg.circuits, func(i, j int) bool {
			return len(pg.circuits[i]) > len(pg.circuits[j])
		})
	}
	return nil
}

func (pg *Playground) Iterate() error {
	pg.RecheckShortestLights()
	shortest := pg.shortestLights[0]
	err := pg.ConnectJunctionBoxes(shortest)
	for i := 0; i < len(pg.shortestLights); i++ {
	}
	if err == CIRCUIT_LOOP_ERR {
		shortest.Junctions[0].Connections = append(shortest.Junctions[0].Connections, shortest.Junctions[1])
		shortest.Junctions[1].Connections = append(shortest.Junctions[1].Connections, shortest.Junctions[0])
		return CIRCUIT_LOOP_ERR
	} else if err != nil {
		return err
	}
	return nil
}

func (pg *Playground) ConnectJunctionBoxes(ls *LightString) error {
	jba := ls.Junctions[0]
	jbb := ls.Junctions[1]
	if jba.HasConnection(jbb) || jbb.HasConnection(jba) {
		return ALREADY_CONNECTED_ERR
	}
	err := pg.AddToCircuit(ls)
	if err != nil {
		return err
	}
	jba.Connections = append(jba.Connections, jbb)
	jbb.Connections = append(jbb.Connections, jba)
	return nil
}

func (pg *Playground) AddToCircuit(ls *LightString) error {
	if ls.Junctions[0].Equals(*ls.Junctions[1]) {
		return LINESTRING_LOOP_ERR
	}
	aCircuit := -1
	bCircuit := -1
	for i := range pg.circuits {
		for _, box := range pg.circuits[i] {
			if box.Equals(*ls.Junctions[0]) {
				if aCircuit != -1 {
					return MULTIPLE_CIRCUITS_ERR
				}
				aCircuit = i
			}
			if box.Equals(*ls.Junctions[1]) {
				if bCircuit != -1 {
					return MULTIPLE_CIRCUITS_ERR
				}
				bCircuit = i
			}
		}
	}
	if aCircuit == -1 && bCircuit == -1 {
		pg.circuits = append(pg.circuits, ls.Junctions[:])
	} else if aCircuit == bCircuit {
		return CIRCUIT_LOOP_ERR
	} else if aCircuit != -1 && bCircuit != -1 {
		if aCircuit > bCircuit {
			tmp := aCircuit
			aCircuit = bCircuit
			bCircuit = tmp
		}
		pg.circuits[aCircuit] = append(pg.circuits[aCircuit], pg.circuits[bCircuit]...)
		if bCircuit == len(pg.circuits)-1 {
			pg.circuits = pg.circuits[:bCircuit]
		} else {
			pg.circuits = append(pg.circuits[:bCircuit], pg.circuits[bCircuit+1:]...)
		}
	} else if aCircuit != -1 {
		pg.circuits[aCircuit] = append(pg.circuits[aCircuit], ls.Junctions[1])
	} else if bCircuit != -1 {
		pg.circuits[bCircuit] = append(pg.circuits[bCircuit], ls.Junctions[0])
	}
	pg.lastConnected = ls
	return nil
}

func (pg *Playground) RecheckShortestLights() {
	for i := range pg.shortestLights {
		jba := pg.shortestLights[i].Junctions[0]
		jbb := pg.shortestLights[i].Junctions[1]
		if jba.HasConnection(jbb) || jbb.HasConnection(jba) {
			pg.shortestLights[i] = pg.smallestEdge(jba)
		}
	}
	pg.shortestLights = slices.DeleteFunc(pg.shortestLights, func(light *LightString) bool {
		return light.Junctions[0].Equals(*light.Junctions[1])
	})
	sort.Slice(pg.shortestLights, func(i, j int) bool {
		return pg.shortestLights[i].Length < pg.shortestLights[j].Length
	})
}

/**
* This function focuses solely on the shortest unused edge from this junction.
*	This does not focus on whether the shortest unused edge from this junction is
*	already in a related circuit.  That is for a different function to determine.
 */
func (pg *Playground) smallestEdge(jb *JunctionBox) *LightString {
	result := LightString{
		Junctions: [2]*JunctionBox{jb, jb},
		Length:    math.Inf(1),
	}
	for i := range pg.junctionBoxes {
		if !jb.HasConnection(pg.junctionBoxes[i]) {
			lng := jb.DistanceTo(pg.junctionBoxes[i])
			if lng < result.Length {
				result.Length = lng
				result.Junctions[1] = pg.junctionBoxes[i]
			}
		}
	}
	return &result
}

type JunctionBox struct {
	X           float64
	Y           float64
	Z           float64
	Connections []*JunctionBox
}

func (jb *JunctionBox) DistanceTo(ob *JunctionBox) float64 {
	return math.Sqrt(math.Pow(jb.X-ob.X, 2) + math.Pow(jb.Y-ob.Y, 2) + math.Pow(jb.Z-ob.Z, 2))
}

func (jb *JunctionBox) Equals(ob JunctionBox) bool {
	if jb.X == ob.X && jb.Y == ob.Y && jb.Z == ob.Z {
		return true
	}
	return false
}

func (jb *JunctionBox) HasConnection(ob *JunctionBox) bool {
	if jb.Equals(*ob) {
		return true
	}
	for _, cb := range jb.Connections {
		if ob.X == cb.X && ob.Y == cb.Y && ob.Z == cb.Z {
			return true
		}
	}
	return false
}

type LightString struct {
	Junctions [2]*JunctionBox
	Length    float64
}
