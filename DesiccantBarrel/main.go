package main

import (
	// "fmt"
	// "math"
	// "os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	barrelRadius := 25.0
	barrelHeight := 64.0 - 7.0
	barrelThickness := 1.0
	cuffHeight := 5.0

	plane, err := sdf.Box3D(v3.Vec{Y: barrelRadius*sdf.Tau + 1.0, X: barrelHeight - barrelThickness, Z: 1}, barrelThickness)
	if err != nil {
		panic(err)
	}
	membrain := waveyDistort3d(plane, v2.Vec{X: 1, Y: 1}, v2.Vec{X: 1, Y: 1})
	membrain = sdf.Cut3D(membrain, v3.Vec{Z: -barrelThickness / 2}, v3.Vec{Z: 1.0})
	membrain = sdf.Cut3D(membrain, v3.Vec{Z: barrelThickness / 2}, v3.Vec{Z: -1.0})

	membrain = sdf.Transform3D(
		membrain,
		sdf.Translate3d(v3.Vec{X: barrelRadius}).Mul(
			sdf.RotateY(sdf.Tau/4.0),
		),
	)

	holeyBarrel := bend3d(membrain, barrelRadius)

	holeyBarrel.(*bendSDF3).SetBoundingBox(
		sdf.Box3{
			Min: v3.Vec{X: -barrelRadius - barrelThickness, Y: -barrelRadius - barrelThickness, Z: -barrelHeight / 2},
			Max: v3.Vec{X: barrelRadius + barrelThickness, Y: barrelRadius + barrelThickness, Z: barrelHeight / 2},
		},
	)

	cuff, err := sdf.Cylinder3D(cuffHeight, barrelRadius+barrelThickness/2.0, barrelThickness)
	if err != nil {
		panic(err)
	}

	cuffHole, err := sdf.Cylinder3D(cuffHeight, barrelRadius-barrelThickness/2.0, 0)
	cuff = sdf.Difference3D(cuff, cuffHole)

	holeyBarrel = sdf.Union3D(
		holeyBarrel,
		sdf.Transform3D(cuff, sdf.Translate3d(v3.Vec{Z: barrelHeight/2.0 - cuffHeight/2.0})),
		sdf.Transform3D(cuff, sdf.Translate3d(v3.Vec{Z: -barrelHeight/2.0 + cuffHeight/2.0})),
	)

	waveyDisc, err := sdf.Circle2D(barrelRadius - cuffHeight)
	if err != nil {
		panic(err)
	}
	waveyDisc3D := sdf.Extrude3D(waveyDisc, barrelThickness)
	waveyDisc3D = waveyDistort3d(waveyDisc3D, v2.Vec{X: 1, Y: 1}, v2.Vec{X: 1, Y: 1})
	waveyDisc3D = sdf.Cut3D(waveyDisc3D, v3.Vec{Z: -barrelThickness / 2}, v3.Vec{Z: 1.0})
	waveyDisc3D = sdf.Cut3D(waveyDisc3D, v3.Vec{Z: barrelThickness / 2}, v3.Vec{Z: -1.0})

	endCap, err := sdf.Cylinder3D(barrelThickness, barrelRadius, 0)
	if err != nil {
		panic(err)
	}

	endCap = sdf.Difference3D(endCap, sdf.Extrude3D(waveyDisc, barrelThickness))
	endCap = sdf.Union3D(endCap, waveyDisc3D)
	endCap = sdf.Transform3D(endCap, sdf.Translate3d(v3.Vec{Z: -barrelHeight/2.0 + barrelThickness/2.0}))

	holeyBarrel = sdf.Union3D(holeyBarrel, endCap)

	render.ToSTL(holeyBarrel, "holeyBarrel.stl", render.NewMarchingCubesUniform(400))
}
