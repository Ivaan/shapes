package main

import (
	"flag"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func main() {
	zeroLength := flag.Float64("zeroLength", 20.0, "Specify the length of the zero bar in mm")
	barLength := flag.Float64("barLength", 8.0, "Specify the length of the horizontle bars in mm")
	barThickness := flag.Float64("barThickness", 2.0, "Specify the thickness of the bars in mm (thickness is also added to the lengths because reasons")
	numeral := flag.Int64("numeral", 1410, "The numeral to encode")

	flag.Parse()

	zeroBar := sdf.Box3D(sdf.V3{X: *barThickness, Y: *barThickness, Z: *zeroLength + *barThickness}, 0) // not part of array because we don't need it added for every zero

	numeralShapes := makeNumeralShapes(*zeroLength, *barLength, *barThickness)
	transformations := makeTransformations()

	figures := make([]sdf.SDF3, 0)
	figures = append(figures, zeroBar)
	for i, t := range transformations {
		digit := *numeral / int64(math.Pow10(i)) % 10
		if digit > 0 {
			figures = append(
				figures,
				sdf.Transform3D(
					numeralShapes[digit],
					t,
				),
			)
		}
	}

	cistercian := sdf.Union3D(figures...)
	sdf.RenderSTLSlow(cistercian, 400, "cistercian.stl")
}

//makeTransformations constructs an array of transformation, one for each power of ten
func makeTransformations() []sdf.M44 {
	transformations := make([]sdf.M44, 8)

	transformations[0] = sdf.Identity3d()
	transformations[1] = sdf.MirrorYZ()
	transformations[2] = sdf.MirrorXY()
	transformations[3] = sdf.MirrorYZ().Mul(sdf.MirrorXY())
	transformations[4] = sdf.RotateZ(sdf.Tau * 1 / 4)
	transformations[5] = sdf.RotateZ(sdf.Tau * 1 / 4).Mul(sdf.MirrorYZ())
	transformations[6] = sdf.RotateZ(sdf.Tau * 1 / 4).Mul(sdf.MirrorXY())
	transformations[7] = sdf.RotateZ(sdf.Tau * 1 / 4).Mul(sdf.MirrorYZ().Mul(sdf.MirrorXY()))

	return transformations
}

//makeNumeralShapes constructs an array of Cictercian numeral shapes, 1 to 9 (0 isn't populated)
func makeNumeralShapes(zeroLength, barLength, barThickness float64) []sdf.SDF3 {
	bar := sdf.Box3D(sdf.V3{X: barThickness, Y: barThickness, Z: barLength + barThickness}, 0)
	barDiag := sdf.Box3D(sdf.V3{X: barThickness, Y: barThickness, Z: barLength * math.Sqrt(2)}, 0)

	paperSix := sdf.Transform3D( //the paper six doesn't work in 3d but is used to construct 7-9
		bar,
		sdf.Translate3d(sdf.V3{X: barLength, Y: 0, Z: zeroLength/2 - barLength/2}),
	)

	numeralShapes := make([]sdf.SDF3, 10)

	numeralShapes[1] = sdf.Transform3D(
		bar,
		sdf.Translate3d(sdf.V3{X: barLength / 2, Y: 0, Z: zeroLength / 2}).Mul(
			sdf.RotateY(sdf.Tau*1/4),
		),
	)

	numeralShapes[2] = sdf.Transform3D(
		bar,
		sdf.Translate3d(sdf.V3{X: barLength / 2, Y: 0, Z: zeroLength/2 - barLength}).Mul(
			sdf.RotateY(sdf.Tau*1/4),
		),
	)

	numeralShapes[3] = sdf.Transform3D(
		barDiag,
		sdf.Translate3d(sdf.V3{X: barLength / 2, Y: 0, Z: zeroLength/2 - barLength/2}).Mul(
			sdf.RotateY(sdf.Tau*(-1.0/8.0)),
		),
	)

	numeralShapes[4] = sdf.Transform3D(
		barDiag,
		sdf.Translate3d(sdf.V3{X: barLength / 2, Y: 0, Z: zeroLength/2 - barLength/2}).Mul(
			sdf.RotateY(sdf.Tau*(1.0/8.0)),
		),
	)

	numeralShapes[5] = sdf.Union3D(
		numeralShapes[1],
		numeralShapes[4],
	)

	numeralShapes[6] = sdf.Union3D(
		numeralShapes[2],
		numeralShapes[3],
	)

	numeralShapes[7] = sdf.Union3D(
		paperSix,
		numeralShapes[1],
	)

	numeralShapes[8] = sdf.Union3D(
		paperSix,
		numeralShapes[2],
	)

	numeralShapes[9] = sdf.Union3D(
		paperSix,
		numeralShapes[1],
		numeralShapes[2],
	)

	return numeralShapes
}
