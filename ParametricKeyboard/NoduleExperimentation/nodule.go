package main

import (
	"github.com/deadsy/sdfx/sdf"
)

type Nodule struct {
	holesThenSolidPairs [][]sdf.SDF3
}

func MakeNodule(holesAndPairs ...[]sdf.SDF3) Nodule {
	return Nodule{holesThenSolidPairs: holesAndPairs}
}

func (n Nodule) OrientAndMove(transform sdf.M44) Nodule {
	movedHolesThenSolidPairs := make([][]sdf.SDF3, len(n.holesThenSolidPairs))
	for i, sdf3s := range n.holesThenSolidPairs {
		movedHolesThenSolidPairs[i] = make([]sdf.SDF3, len(n.holesThenSolidPairs[i]))
		for j, sdf3 := range sdf3s {
			if sdf3 != nil {
				movedHolesThenSolidPairs[i][j] = sdf.Transform3D(sdf3, transform)
			}
		}
	}
	return Nodule{holesThenSolidPairs: movedHolesThenSolidPairs}
}

type NoduleCollection []Nodule

func (nc NoduleCollection) Combine() sdf.SDF3 {
	startingPair := 0
	for _, n := range nc {

		numPairs := len(n.holesThenSolidPairs) / 2
		if numPairs > startingPair {
			startingPair = numPairs
		}
	}

	getSDFsAtRank := func(r int) []sdf.SDF3 {
		totalLength := 0
		for _, n := range nc {
			if len(n.holesThenSolidPairs) > r {
				totalLength += len(n.holesThenSolidPairs[r])
			}
		}
		sdfsAtRank := make([]sdf.SDF3, totalLength)
		var i int
		for _, n := range nc {
			if len(n.holesThenSolidPairs) > r {
				i += copy(sdfsAtRank[i:], n.holesThenSolidPairs[r])
			}
		}
		return sdfsAtRank
	}

	var currentSDF3 sdf.SDF3

	for pair := startingPair; pair >= 0; pair-- {
		currentSDF3 = sdf.Difference3D(
			sdf.Union3D(append(getSDFsAtRank(pair*2+1), currentSDF3)...),
			sdf.Union3D(getSDFsAtRank(pair*2)...),
		)
	}

	return currentSDF3
}
