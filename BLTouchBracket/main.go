package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	printerHoleSpacing := 10.0
	printerHoleDiameter := 3.5

	sensorHoleSpacing := 18.0
	sensorHoleDiameter := 3.15

	bracketWidth := 12.0
	bracketThickness := 4.0
	bracketPastHoles := 7.0
	sensorDrop := 7.0

	printerSideLength := sensorDrop + printerHoleSpacing + bracketPastHoles

	printerScrewHole := sdf.Cylinder3D(bracketThickness, printerHoleDiameter/2.0, 0)
	printerBracket := sdf.Transform3D(sdf.Box3D(sdf.V3{printerSideLength, bracketWidth, bracketThickness}, 0.5), sdf.Translate3d(sdf.V3{printerSideLength / 2.0, 0, 0}))
	printerBracket = sdf.Difference3D(
		printerBracket,
		sdf.Union3D(
			sdf.Transform3D(printerScrewHole, sdf.Translate3d(sdf.V3{sensorDrop, 0, 0})),
			sdf.Transform3D(printerScrewHole, sdf.Translate3d(sdf.V3{sensorDrop + printerHoleSpacing, 0, 0})),
		),
	)

	sensorSideLength := bracketThickness + sensorHoleSpacing + 2*bracketPastHoles

	sensorScrewHole := sdf.Cylinder3D(bracketThickness, sensorHoleDiameter/2.0, 0)
	sensorBracket := sdf.Transform3D(sdf.Box3D(sdf.V3{sensorSideLength, bracketWidth, bracketThickness}, 0.5), sdf.Translate3d(sdf.V3{sensorSideLength / 2.0, 0, 0}))
	sensorBracket = sdf.Difference3D(
		sensorBracket,
		sdf.Union3D(
			sdf.Transform3D(sensorScrewHole, sdf.Translate3d(sdf.V3{bracketThickness + bracketPastHoles, 0, 0})),
			sdf.Transform3D(sensorScrewHole, sdf.Translate3d(sdf.V3{bracketThickness + bracketPastHoles + sensorHoleSpacing, 0, 0})),
		),
	)

	bracket := sdf.Union3D(
		sdf.Transform3D(
			sensorBracket,
			sdf.Translate3d(sdf.V3{0, 0, bracketThickness / 2}),
		),
		sdf.Transform3D(
			printerBracket,
			sdf.Translate3d(sdf.V3{bracketThickness / 2, 0, bracketThickness}).Mul(sdf.RotateY(-sdf.Tau/4)),
		),
	)
	//bracket.(*sdf.UnionSDF3).SetMin(sdf.RoundMin(4))
	bracket.(*sdf.UnionSDF3).ScaleBoundingBox(1.3)

	// bracket = sdf.Intersect3D(
	// 	bracket,
	// 	sdf.NewBox3(sdf.V3{sensorSideLength/2, })
	// )

	sdf.RenderSTLSlow(bracket, 400, "bracket.stl")
	//RoundMin(k)
}
