package main

import (
	"fmt"
	"math"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	tolerance := 0.4

	ironDiameter := 37.23 + tolerance
	wireWidth := 4.88
	wireHeight := 2.06

	holderDepth := 4.0
	baseThickness := 4.0
	sideThickness := 6.0

	bevelDepthX := 2.5
	bevelDepthZ := 6.0
	bevelAngle := -math.Atan2(bevelDepthX, bevelDepthZ)
	bevelAngleLength := math.Sqrt(bevelDepthX*bevelDepthX + bevelDepthZ*bevelDepthZ)
	holeCercumferance := math.Pi * ironDiameter

	logoDepth := 0.6

	edgeSoftening := 0.4

	base, _ := sdf.Cylinder3D(baseThickness+holderDepth, ironDiameter/2+sideThickness, edgeSoftening)
	base = sdf.Transform3D(base, sdf.Translate3d(sdf.V3{Z: (baseThickness + holderDepth) / 2}))
	hole, _ := sdf.Cylinder3D(holderDepth*2, ironDiameter/2, edgeSoftening) //longer holderDepth so softened top doesn't form an odd lip, double so math is easier
	hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{Z: holderDepth + baseThickness}))
	wireSlot, _ := sdf.Box3D(sdf.V3{X: wireHeight * 2, Y: wireWidth, Z: holderDepth * 2}, edgeSoftening) //again, double holderDepth for same reasons as above
	wireSlot = sdf.Transform3D(wireSlot, sdf.Translate3d(sdf.V3{X: ironDiameter / 2, Z: holderDepth + baseThickness}))

	f, err := sdf.LoadFont("data-latin.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}
	logoText := sdf.NewText("NazFab")
	logo2d, err := sdf.TextSDF2(f, logoText, bevelAngleLength)
	if err != nil {
		fmt.Printf("error rendering text %s\n", err)
		os.Exit(1)
	}
	logo3d, _ := sdf.ExtrudeRounded3D(logo2d, logoDepth*2, 0.2) //hard coded 10, could be related to knerl somehow, never did work out what makes the knerl points tall vs not ...
	logo3d = sdf.Transform3D(
		logo3d,
		sdf.Translate3d(sdf.V3{X: -bevelDepthX / 2, Y: holeCercumferance / 4}).Mul(
			sdf.RotateZ(sdf.Tau/4).Mul(
				sdf.RotateX(sdf.Tau/4),
			),
		),
	)

	bevel, _ := sdf.Box3D(sdf.V3{X: bevelDepthX, Y: holeCercumferance, Z: bevelAngleLength}, 0)
	bevel = sdf.Union3D(bevel, logo3d)
	bevel = sdf.Transform3D(
		bevel,
		sdf.Translate3d(sdf.V3{X: ironDiameter/2 + sideThickness, Z: holderDepth + baseThickness - bevelDepthZ}).Mul(
			sdf.RotateY(bevelAngle).Mul(
				sdf.Translate3d(sdf.V3{X: bevelDepthX / 2, Z: bevelAngleLength / 2}),
			),
		),
	)
	bevel = bend3d(bevel, ironDiameter/2)

	holeWithSlot := sdf.Union3D(hole, wireSlot, bevel)

	baseWithHole := sdf.Difference3D(base, holeWithSlot)
	//baseWithHole = sdf.Union3D(baseWithHole, bevel)

	render.RenderSTLSlow(baseWithHole, 400, "solderingIronBase.stl")

}

/*
dia - 37.23mm
dia with wire - 38.67
wire width - 4.88
wire height computed - 1.44
wire height messured - 2.06
*/
