package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

func main() {
	smallSpiralRadius := 15.0
	largeSpiralRadius := 30.0
	numberOfPoints := 9
	twistNumerator := -5.0
	twistDenominator := 1.0
	twistHeight := 170.0

	topRadius := 3.0
	outerRadius := 40.0
	toleranceRadius := 0.5

	halfPointAngle := sdf.Tau / (float64(numberOfPoints) * 2.0)
	innerPoly := sdf.Polygon{}
	outerCutPoly := sdf.Polygon{}

	v := v2.Vec{Y: 1}
	rot := sdf.Rotate2d(halfPointAngle)

	//Make a star-polygon with the numberOfPoints
	for i := 0; i < 2*numberOfPoints; i++ {
		if i%2 == 0 {
			innerPoly.AddV2(v.MulScalar(smallSpiralRadius))
			outerCutPoly.AddV2(v.MulScalar(smallSpiralRadius + toleranceRadius))
		} else {
			innerPoly.AddV2(v.MulScalar(largeSpiralRadius))
			outerCutPoly.AddV2(v.MulScalar(largeSpiralRadius + toleranceRadius))
		}
		v = rot.MulPosition(v)
	}
	inner2D, err := innerPoly.Mesh2D()
	if err != nil {
		panic(err)
	}
	outer2D, err := outerCutPoly.Mesh2D()
	if err != nil {
		panic(err)
	}

	cone, err := sdf.Cone3D(twistHeight, outerRadius, topRadius, 2)
	if err != nil {
		panic(err)
	}

	inner3D := sdf.TwistExtrude3D(inner2D, twistHeight, twistNumerator/twistDenominator)
	outer3Dcut := sdf.TwistExtrude3D(outer2D, twistHeight, twistNumerator/twistDenominator)

	outer3D, err := sdf.Cylinder3D(twistHeight, outerRadius, 2)
	outer3D = sdf.Difference3D(outer3D, outer3Dcut)

	inner3D = sdf.Intersect3D(inner3D, cone)
	outer3D = sdf.Intersect3D(outer3D, cone)

	render.ToSTL(inner3D, "Inner.stl", render.NewMarchingCubesUniform(500))
	render.ToSTL(outer3D, "Outer.stl", render.NewMarchingCubesUniform(500))

}
