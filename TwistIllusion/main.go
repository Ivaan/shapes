package main

import (
	"fmt"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	smallSpiralRadius := 9 / 2.0
	largeSpiralRadius := 39.0 / 2.0
	numberOfPoints := 7
	twistNumerator := -2.0
	twistDenominator := 3.0
	twistHeight := 73.0
	polyRound := 2.5

	topRadius := 3.0
	outerRadius := 22.0
	toleranceRadius := 2.2
	outerRound := 2.0
	spindlyCutoff := 1.0 //this is the thickness at which the outer fingers at the top get cut off (so they don't get too spindly)
	//(ssr + tr + sc - tr)/x = (or-tr)/twistHeight
	cutOffHeight := (twistHeight / 2.0) - twistHeight*(smallSpiralRadius+topRadius+spindlyCutoff-topRadius)/(outerRadius-topRadius)
	fmt.Println(cutOffHeight)

	halfPointAngle := sdf.Tau / (float64(numberOfPoints) * 2.0)
	innerPoly := sdf.Polygon{}
	outerCutPoly := sdf.Polygon{}

	v := v2.Vec{Y: 1}
	rot := sdf.Rotate2d(halfPointAngle)

	//Make a star-polygon with the numberOfPoints
	for i := 0; i < 2*numberOfPoints; i++ {
		if i%2 == 0 {
			innerPoly.AddV2(v.MulScalar(smallSpiralRadius - polyRound))
			outerCutPoly.AddV2(v.MulScalar(smallSpiralRadius + toleranceRadius - polyRound))
		} else {
			innerPoly.AddV2(v.MulScalar(largeSpiralRadius - polyRound))
			outerCutPoly.AddV2(v.MulScalar(largeSpiralRadius + toleranceRadius - polyRound))
		}
		v = rot.MulPosition(v)
	}
	inner2D, err := innerPoly.Mesh2D()
	if err != nil {
		panic(err)
	}
	inner2D = sdf.Offset2D(inner2D, polyRound/2.0)

	outer2D, err := outerCutPoly.Mesh2D()
	if err != nil {
		panic(err)
	}
	outer2D = sdf.Offset2D(outer2D, polyRound/2.0)

	cone, err := sdf.Cone3D(twistHeight, outerRadius, topRadius, outerRound)
	if err != nil {
		panic(err)
	}

	inner3D := sdf.TwistExtrude3D(inner2D, twistHeight, twistNumerator/twistDenominator*sdf.Tau)
	outer3Dcut := sdf.TwistExtrude3D(outer2D, twistHeight, twistNumerator/twistDenominator*sdf.Tau)

	outer3D, err := sdf.Cylinder3D(twistHeight, outerRadius, outerRound)
	outer3D = sdf.Difference3D(outer3D, outer3Dcut)
	outer3D = sdf.Cut3D(outer3D, v3.Vec{Z: cutOffHeight}, v3.Vec{Z: -1})

	inner3D = sdf.Intersect3D(inner3D, cone)
	outer3D = sdf.Intersect3D(outer3D, cone)

	render.ToSTL(inner3D, "Inner.stl", render.NewMarchingCubesUniform(500))
	render.ToSTL(outer3D, "Outer.stl", render.NewMarchingCubesUniform(500))

}
