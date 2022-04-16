package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	//for M2 x 12 bolts and M2 x 8 x 3.2 insert
	//test results:
	//largest size for bolt
	// - thread diameter 3.0
	// - head diameter := 4.8

	//inserts melted the plastic up into the threads from the bottom
	//even the largest size did this
	//largest size seemed a bit large even though it might be strong enough

	insertLength := 8.0
	insertDiameter := 2.9
	insertWallThickness := 2.0
	screwThreadDiameter := 2.0
	screwThreadLength := 12.0
	screwHeadDiameter := 3.8

	incrementDiameters := 0.1
	increments := 5

	wallThickness := 3.0
	shapeRadius := 10.0

	bottom, _ := sdf.Sphere3D(shapeRadius)
	bottomHollow, _ := sdf.Sphere3D(shapeRadius - wallThickness)
	bottomShell := sdf.Difference3D(bottom, bottomHollow)
	bottomShell = sdf.Cut3D(bottomShell, sdf.V3{}, sdf.V3{Z: -1})

	top, _ := sdf.Cylinder3D(shapeRadius, shapeRadius, wallThickness)
	topHollow, _ := sdf.Cylinder3D(shapeRadius-wallThickness, shapeRadius-wallThickness, 0)
	topShell := sdf.Difference3D(top, topHollow)
	topShell = sdf.Cut3D(topShell, sdf.V3{}, sdf.V3{Z: 1})

	insertHolders := make([]sdf.SDF3, increments)
	insertHoldersHoles := make([]sdf.SDF3, increments)

	for i := 0; i < increments; i++ {
		radiusAdjust := float64(i-increments/2) * incrementDiameters
		holder, _ := sdf.Cylinder3D(insertLength+insertWallThickness, insertDiameter/2+insertWallThickness+radiusAdjust, 0)
		holderHole, _ := sdf.Cylinder3D(insertLength, insertDiameter/2+radiusAdjust, 0)
		holder = sdf.Transform3D(holder, sdf.RotateZ(float64(i)*sdf.Tau/float64(increments)).Mul(sdf.Translate3d(sdf.V3{X: shapeRadius - (insertDiameter/2 + insertWallThickness + radiusAdjust), Z: (insertLength + insertWallThickness) / 2})))
		holderHole = sdf.Transform3D(holderHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(increments)).Mul(sdf.Translate3d(sdf.V3{X: shapeRadius - (insertDiameter/2 + insertWallThickness + radiusAdjust), Z: insertLength / 2})))
		insertHolders[i] = holder
		insertHoldersHoles[i] = holderHole
	}

	screwChannel := make([]sdf.SDF3, increments)
	screwHoles := make([]sdf.SDF3, increments)

	for i := 0; i < increments; i++ {
		radiusAdjust := float64(i-increments/2) * incrementDiameters
		channel, _ := sdf.Cylinder3D(shapeRadius, screwHeadDiameter/2+insertWallThickness+radiusAdjust, 0)
		screwThreadHole, _ := sdf.Cylinder3D(screwThreadLength-insertLength, screwThreadDiameter/2+radiusAdjust, 0)
		screwHeadHole, _ := sdf.Cylinder3D(shapeRadius, screwHeadDiameter/2+radiusAdjust, 0)

		channel = sdf.Transform3D(channel, sdf.RotateZ(float64(i)*sdf.Tau/float64(increments)).Mul(sdf.Translate3d(sdf.V3{X: shapeRadius - (insertDiameter/2 + insertWallThickness + radiusAdjust), Z: -shapeRadius / 2})))
		channel = sdf.Intersect3D(bottom, channel)
		screwThreadHole = sdf.Transform3D(screwThreadHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(increments)).Mul(sdf.Translate3d(sdf.V3{X: shapeRadius - (insertDiameter/2 + insertWallThickness + radiusAdjust), Z: -(screwThreadLength - insertLength) / 2})))
		screwHeadHole = sdf.Transform3D(screwHeadHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(increments)).Mul(sdf.Translate3d(sdf.V3{X: shapeRadius - (insertDiameter/2 + insertWallThickness + radiusAdjust), Z: -shapeRadius/2 - (screwThreadLength - insertLength)})))
		hole := sdf.Union3D(screwThreadHole, screwHeadHole)
		screwChannel[i] = channel
		screwHoles[i] = hole
	}

	fitmentTestTop := sdf.Difference3D(
		sdf.Union3D(topShell, sdf.Union3D(insertHolders...)),
		sdf.Union3D(insertHoldersHoles...),
	)

	fitmentTestBottom := sdf.Difference3D(
		sdf.Union3D(bottomShell, sdf.Union3D(screwChannel...)),
		sdf.Union3D(screwHoles...),
	)

	render.RenderSTL(fitmentTestTop, 300, "fitmentTestTop.stl")
	render.RenderSTL(fitmentTestBottom, 300, "fitmentTestBottom.stl")

}
