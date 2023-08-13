package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	numberOfHoles := 12
	holeDiameter := 1.60
	chamferDiameter := 0.3

	barWidth := 2.5
	spacerHeight := 1.3
	holeSpacing := 2.54

	barLength := float64(numberOfHoles) * holeSpacing
	c, err := obj.ChamferedHole3D(spacerHeight, holeDiameter/2.0, chamferDiameter/2.0)
	if err != nil {
		panic(err)
	}

	b, err := sdf.Box3D(v3.Vec{X: barLength, Y: barWidth, Z: spacerHeight}, 0.1)
	if err != nil {
		panic(err)
	}

	holes := make([]sdf.SDF3, numberOfHoles)
	for i := 0; i < numberOfHoles; i++ {
		holes[i] = sdf.Transform3D(
			c,
			sdf.Translate3d(v3.Vec{X: float64(i)*holeSpacing - barLength/2.0 + holeSpacing/2.0}),
		)
	}

	holes3D := sdf.Union3D(holes...)

	spacer := sdf.Difference3D(b, holes3D)
	render.ToSTL(spacer, "spacerBar.stl", render.NewMarchingCubesOctree(600))
}
