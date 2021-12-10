package main

import (
	"fmt"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

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

	x2d, err := sdf.TextSDF2(f, xText, 15.0)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	y2d, err := sdf.TextSDF2(f, yText, 15.0)
	if err != nil {
		fmt.Printf("can't generate Y sdf2 %s\n", err)
		os.Exit(1)
	}

	z2d, err := sdf.TextSDF2(f, zText, 15.0)
	if err != nil {
		fmt.Printf("can't generate Z sdf2 %s\n", err)
		os.Exit(1)
	}

	x3d, err := sdf.Loft3D(sdf.Offset2D(x2d, -4), x2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate X sdf3 %s\n", err)
		os.Exit(1)
	}
	x3d = sdf.Transform3D(x3d, sdf.Translate3d(sdf.V3{0, 0, -2}))

	y3d, err := sdf.Loft3D(sdf.Offset2D(y2d, -4), y2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate Y sdf3 %s\n", err)
		os.Exit(1)
	}
	y3d = sdf.Transform3D(y3d, sdf.Translate3d(sdf.V3{0, 0, -2}))

	z3d, err := sdf.Loft3D(sdf.Offset2D(z2d, -4), z2d, 4.0, 0)
	if err != nil {
		fmt.Printf("can't generate Z sdf3 %s\n", err)
		os.Exit(1)
	}
	z3d = sdf.Transform3D(z3d, sdf.Translate3d(sdf.V3{0, 0, -2}))

	// x3d := sdf.ExtrudeRounded3D(x2d, 1.0, 0.2)
	// y3d := sdf.ExtrudeRounded3D(y2d, 1.0, 0.2)
	// z3d := sdf.ExtrudeRounded3D(z2d, 1.0, 0.2)

	zTop := sdf.Transform3D(
		z3d,
		sdf.Translate3d(sdf.V3{0, 0, 10}),
	)
	zBottom := sdf.Transform3D(
		z3d,
		sdf.RotateY(sdf.DtoR(180)).Mul(
			sdf.Translate3d(sdf.V3{0, 0, 10}),
		),
	)
	xLeft := sdf.Transform3D(
		x3d, sdf.RotateZ(sdf.DtoR(-90)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(sdf.V3{0, 0, 10}),
		),
	)
	xRight := sdf.Transform3D(
		x3d, sdf.RotateZ(sdf.DtoR(90)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(sdf.V3{0, 0, 10}),
		),
	)
	yFront := sdf.Transform3D(
		y3d,
		sdf.RotateX(sdf.DtoR(90)).Mul(
			sdf.Translate3d(sdf.V3{0, 0, 10}),
		),
	)
	yBack := sdf.Transform3D(
		y3d, sdf.RotateZ(sdf.DtoR(180)).Mul(
			sdf.RotateX(sdf.DtoR(90)),
		).Mul(
			sdf.Translate3d(sdf.V3{0, 0, 10}),
		),
	)
	box, err := sdf.Box3D(sdf.V3{20.0, 20.0, 20.0}, 2)
	if err != nil {
		fmt.Printf("can't generate Box sdf3 %s\n", err)
		os.Exit(1)
	}
	box = sdf.Difference3D(
		box,
		sdf.Union3D(xLeft, xRight, yFront, yBack, zTop, zBottom),
	)
	//box := sdf.Union3D(sdf.Box3D(sdf.V3{20.0, 20.0, 20.0}, 2), xLeft, xRight, yFront, yBack, zTop, zBottom)

	render.RenderSTLSlow(box, 300, "textCube.stl")
}
