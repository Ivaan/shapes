package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	useBearings := false //false means direct on shaft
	tumblerFaceEdgeLength := 26.0
	tumblerRadius := tumblerFaceEdgeLength / 2 / math.Cos(sdf.Tau/12)
	tumblerCornerRound := tumblerRadius * 0.05
	tumblerSpacing := (tumblerRadius - (tumblerFaceEdgeLength / 2)) * 2
	tumblerShortRadius := math.Sqrt(tumblerRadius*tumblerRadius - (tumblerFaceEdgeLength/2)*(tumblerFaceEdgeLength/2))
	tumblerMinimumWallThickness := 2.0

	bearingOD := 22.0
	bearingThickness := 7.0
	bearingHolderStopConstriction := 1.0 //horizontle and vertical chamfer distance
	bearingHolderTolerance := 0.1

	shaftOD := 6.35
	tumblerShaftTollerance := 0.4
	tumblerShaftBearingSurfaceLength := 5.0 //length of shaft touched by tumbler and either end
	spacerThickness := 5.0
	spacerChamfer := 1.0
	spacerShaftTollerance := 0.3
	spacerTumblerTollerance := 0.5
	spacerEdgeInFromTumblerEdge := 1.0

	pusherNibSize := 2.0
	pusherLength := 2.2
	pusherTollerance := 1.5

	tumblerOutside := makeTumblerOutside(tumblerRadius, tumblerCornerRound, tumblerFaceEdgeLength)

	var insideHole sdf.SDF3
	if useBearings {
		insideHole = makeBearingHole(bearingOD, bearingThickness, bearingHolderStopConstriction, bearingHolderTolerance, tumblerFaceEdgeLength)
	} else {
		insideHole = makeShaftAndSpacerHole(shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeLength, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness)
	}

	tracks, nibs := makePusherTracksAndNibs(pusherNibSize, pusherLength, tumblerRadius, tumblerShortRadius-spacerEdgeInFromTumblerEdge, tumblerFaceEdgeLength, tumblerSpacing, pusherTollerance)

	holes := sdf.Union3D(insideHole, tracks)
	tumbler := sdf.Difference3D(tumblerOutside, holes)
	tumbler = sdf.Union3D(tumbler, nibs)
	spacerDisk := makeSpacerDisk(shaftOD, spacerThickness, spacerChamfer, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerShortRadius)

	sdf.RenderSTLSlow(tumbler, 400, "tumbler.stl")
	sdf.RenderSTLSlow(spacerDisk, 400, "spacerDisk.stl")
	// sdf.RenderSTL(tumbler, 200, "tumbler.stl")
	// sdf.RenderSTL(spacerDisk, 200, "spacerDisk.stl")

}

func makeTumblerOutside(tumblerRadius, tumblerCornerRound, tumblerFaceEdgeLength float64) sdf.SDF3 {
	//tumbler, triangle outside, bearing holder hole inside
	// triangle is an extruded 3 nagon
	triangle := sdf.Polygon2D(sdf.Nagon(3, tumblerRadius-tumblerCornerRound))
	triangle = sdf.Offset2D(triangle, tumblerCornerRound)
	tumblerOutside := sdf.Extrude3D(triangle, tumblerFaceEdgeLength)
	return tumblerOutside
}

func pusherNibProfile(width, height float64) sdf.SDF2 {
	return sdf.Polygon2D([]sdf.V2{
		{-width / 2, 0},
		{-width / 2, height},
		{0, height + width*0.5},
		{width / 2, height},
		{width / 2, 0},
		{-width / 2, 0},
	})
}

func makePusherTracksAndNibs(pusherNibSize, pusherLength, tumblerRadius, innerClearance, tumblerFaceEdgeLength, tumblerSpacing, pusherTollerance float64) (sdf.SDF3, sdf.SDF3) {
	pusher2D := pusherNibProfile(pusherNibSize, pusherNibSize+tumblerSpacing)
	track2D := pusherNibProfile(pusherNibSize+pusherTollerance/2, pusherNibSize+tumblerSpacing+pusherTollerance)

	pushers := make([]sdf.SDF3, 3)
	tracks := make([]sdf.SDF3, 3)

	for n := 0; n < 3; n++ {
		//distanceFromCenter := (tumblerRadius+innerClearance)/2 + (pusherNibSize-(pusherTollerance/2))*(float64(n-1)) // three tracks one halfway between tumbler radioius and the inner clearance, and one on each side
		distanceFromCenter := innerClearance + pusherTollerance + pusherNibSize*(float64(n)) // three tracks one halfway between tumbler radioius and the inner clearance, and one on each side
		startAngle := float64(n) * sdf.Tau / 3

		pusherArcAngle := pusherLength / distanceFromCenter

		pushers[n] = sdf.Transform3D(
			sdf.RevolveTheta3D(
				sdf.Transform2D(
					pusher2D,
					sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
				),
				pusherArcAngle,
			),
			sdf.RotateZ(startAngle-pusherArcAngle/2).Mul(
				sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeLength / 2}),
			),
		)

		tracks[n] = sdf.Transform3D(
			sdf.RevolveTheta3D(
				sdf.Transform2D(
					track2D,
					sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
				),
				sdf.Tau/3+pusherArcAngle,
			),
			sdf.RotateZ(startAngle-pusherArcAngle/2).Mul(
				sdf.Translate3d(sdf.V3{0, 0, -tumblerFaceEdgeLength / 2}),
			),
		)
	}
	return sdf.Union3D(tracks...), sdf.Union3D(pushers...)
}

//bearing holder hole is an extruded hex with a constriction in the middle
//this is accomplished by stacking 5 shapes
//top and bottom are bearing sized cylinders
//center is constricted cylinder
//between top and center (and bottom and center) are correctly oriented truncated between bearing and constricted sized cylinders
func makeBearingHole(bearingOD, bearingThickness, bearingHolderStopConstriction, bearingHolderTolerance, tumblerFaceEdgeLength float64) sdf.SDF3 {

	bearingHolder := sdf.Cylinder3D(bearingThickness, bearingOD/2+bearingHolderTolerance, 0)
	chamfer3d := sdf.Cone3D(
		bearingHolderStopConstriction,
		bearingOD/2+bearingHolderTolerance-bearingHolderStopConstriction,
		bearingOD/2+bearingHolderTolerance,
		0,
	)
	constricted := sdf.Cylinder3D(tumblerFaceEdgeLength-2*bearingThickness, bearingOD/2+bearingHolderTolerance-bearingHolderStopConstriction, 0)

	topBearingHolder := sdf.Transform3D(
		bearingHolder,
		sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeLength/2 - bearingThickness/2}),
	)
	topChamfer3d := sdf.Transform3D(
		chamfer3d,
		sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeLength/2 - bearingThickness - bearingHolderStopConstriction/2}),
	)
	flip := sdf.RotateX(sdf.Tau / 2)
	bottomBearingHolder := sdf.Transform3D(topBearingHolder, flip)
	bottomChamfer3d := sdf.Transform3D(topChamfer3d, flip)

	bearingHole := sdf.Union3D(topBearingHolder, topChamfer3d, constricted, bottomChamfer3d, bottomBearingHolder)
	return bearingHole
}

func makeShaftAndSpacerHole(shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeLength, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness float64) sdf.SDF3 {
	innerChamfer := (tumblerShortRadius - tumblerMinimumWallThickness) - (shaftOD/2 + tumblerShaftTollerance)
	//innerChamfer := 10.0
	fmt.Printf("shaftOD: %v, tumblerShaftTollerance: %v, tumblerShaftBearingSurfaceLength: %v, spacerThickness: %v, spacerShaftTollerance: %v, spacerEdgeInFromTumblerEdge: %v, spacerTumblerTollerance: %v, tumblerFaceEdgeLength: %v, tumblerShortRadius: %v, tumblerSpacing: %v, tumblerMinimumWallThickness: %v, innerChamfer: %v", shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeLength, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness, innerChamfer)

	return sdf.Revolve3D(
		sdf.Polygon2D([]sdf.V2{
			{0, tumblerFaceEdgeLength / 2},
			{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeLength / 2},
			{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeLength/2 - spacerThickness - spacerTumblerTollerance},
			{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeLength/2 - spacerThickness - spacerTumblerTollerance},
			{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeLength/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength},
			{shaftOD/2 + tumblerShaftTollerance + innerChamfer, tumblerFaceEdgeLength/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength - innerChamfer},
			{shaftOD/2 + tumblerShaftTollerance + innerChamfer, -tumblerFaceEdgeLength/2 + tumblerShaftBearingSurfaceLength + innerChamfer},
			{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeLength/2 + tumblerShaftBearingSurfaceLength},
			{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeLength / 2},
			{0, -tumblerFaceEdgeLength / 2},
			{0, tumblerFaceEdgeLength / 2},
		}),
	)
}
func makeSpacerDisk(shaftOD, spacerThickness, spacerChamfer, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerShortRadius float64) sdf.SDF3 {
	return sdf.Difference3D(
		sdf.Cylinder3D(spacerThickness, tumblerShortRadius-spacerEdgeInFromTumblerEdge-spacerTumblerTollerance, 0),
		sdf.Union3D(
			sdf.Cylinder3D(spacerThickness, shaftOD/2+spacerShaftTollerance, 0),
			sdf.Transform3D(
				sdf.Cone3D(spacerChamfer, shaftOD/2+spacerShaftTollerance, shaftOD/2+spacerShaftTollerance+spacerChamfer, 0),
				sdf.Translate3d(sdf.V3{0, 0, (spacerThickness - spacerChamfer) / 2}),
			),
		),
	)
}
