package main

import (
	"fmt"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

//The idea is to rename the function you want to print to main
//later maybe I can make this a param but...

// dots
func main() {
	circleRadius := 15.0
	circlesRadius := 125.0
	numberOfCircles := 5
	thickness := 0.2

	c, err := sdf.Circle2D(circleRadius)
	if err != nil {
		panic(err)
	}

	c3D := sdf.Extrude3D(c, thickness)

	c3D = sdf.Transform3D(c3D, sdf.Translate3d(v3.Vec{X: circlesRadius}))

	dots := sdf.RotateCopy3D(c3D, numberOfCircles)

	render.ToSTL(dots, "dots.stl", render.NewMarchingCubesUniform(300))

}

// testtube - for vertical strength testing?
func main_testtube() {
	//this is going to print a vertical cylinder that hollow with a half sphere at the bottom
	//wall thickness is also the bottom of the cylinder
	//sphere is there only to keep it attached to the print plate incase that's also a problem

	cylinderHeight := 30.0
	cylinderRadius := 5.0
	thickness := 1.0
	sphereRadius := 7.5

	c, err := sdf.Cylinder3D(cylinderHeight, cylinderRadius, thickness)
	if err != nil {
		panic(err)
	}
	c = sdf.Transform3D(c, sdf.Translate3d(v3.Vec{Z: cylinderHeight / 2}))

	h, err := sdf.Cylinder3D(cylinderHeight-thickness, cylinderRadius-thickness, 0)
	if err != nil {
		panic(err)
	}
	h = sdf.Transform3D(h, sdf.Translate3d(v3.Vec{Z: (cylinderHeight-thickness)/2 + thickness}))

	s, err := sdf.Sphere3D(sphereRadius)
	if err != nil {
		panic(err)
	}

	testTube := sdf.Union3D(c, s)
	testTube = sdf.Difference3D(testTube, h)
	testTube = sdf.Cut3D(testTube, v3.Vec{}, v3.Vec{Z: 1})

	render.ToSTL(testTube, "testTube.stl", render.NewMarchingCubesUniform(300))

}

// testcube
func main_testcube() {

	f, err := sdf.LoadFont("cmr10.ttf")
	//f, err := LoadFont("Times_New_Roman.ttf")
	//f, err := LoadFont("wt064.ttf")

	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	xText := sdf.NewText("X")
	yText := sdf.NewText("Y")
	zText := sdf.NewText("Z")

	x2d, err := sdf.Text2D(f, xText, 15.0)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	y2d, err := sdf.Text2D(f, yText, 15.0)
	if err != nil {
		fmt.Printf("can't generate Y sdf2 %s\n", err)
		os.Exit(1)
	}

	z2d, err := sdf.Text2D(f, zText, 15.0)
	if err != nil {
		fmt.Printf("can't generate Z sdf2 %s\n", err)
		os.Exit(1)
	}

	x3d, err := sdf.Loft3D(sdf.Offset2D(x2d, -4), x2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate X sdf3 %s\n", err)
		os.Exit(1)
	}
	x3d = sdf.Transform3D(x3d, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -2}))

	y3d, err := sdf.Loft3D(sdf.Offset2D(y2d, -4), y2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate Y sdf3 %s\n", err)
		os.Exit(1)
	}
	y3d = sdf.Transform3D(y3d, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -2}))

	z3d, err := sdf.Loft3D(sdf.Offset2D(z2d, -4), z2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate Z sdf3 %s\n", err)
		os.Exit(1)
	}
	z3d = sdf.Transform3D(z3d, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -2}))

	// x3d := sdf.ExtrudeRounded3D(x2d, 1.0, 0.2)
	// y3d := sdf.ExtrudeRounded3D(y2d, 1.0, 0.2)
	// z3d := sdf.ExtrudeRounded3D(z2d, 1.0, 0.2)

	zTop := sdf.Transform3D(
		z3d,
		sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
	)
	zBottom := sdf.Transform3D(
		z3d,
		sdf.RotateY(sdf.DtoR(180)).Mul(
			sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
		),
	)
	xLeft := sdf.Transform3D(
		x3d, sdf.RotateZ(sdf.DtoR(-90)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
		),
	)
	xRight := sdf.Transform3D(
		x3d, sdf.RotateZ(sdf.DtoR(90)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
		),
	)
	yFront := sdf.Transform3D(
		y3d,
		sdf.RotateX(sdf.DtoR(90)).Mul(
			sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
		),
	)
	yBack := sdf.Transform3D(
		y3d, sdf.RotateZ(sdf.DtoR(180)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: 10}),
		),
	)
	box, err := sdf.Box3D(v3.Vec{X: 20.0, Y: 20.0, Z: 20.0}, 2)
	if err != nil {
		fmt.Printf("can't generate Box sdf3 %s\n", err)
		os.Exit(1)
	}
	box = sdf.Difference3D(
		box,
		sdf.Union3D(xLeft, xRight, yFront, yBack, zTop, zBottom),
	)
	//box := sdf.Union3D(sdf.Box3D(v3.Vec{X: 20.0,Y: 20.0,Z: 20.0}, 2), xLeft, xRight, yFront, yBack, zTop, zBottom)

	render.ToSTL(box, "textCube.stl", render.NewMarchingCubesUniform(300))
}
