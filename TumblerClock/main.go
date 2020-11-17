package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	tumblerFaceEdgeWidth := 75.0
	tumblerFaceEdgeHeight := tumblerFaceEdgeWidth / 3
	tumblerRadius := tumblerFaceEdgeWidth / 2 / math.Cos(sdf.Tau/12)
	tumblerCornerRound := tumblerRadius * 0.05
	//tumblerShortRadius := math.Sqrt(tumblerRadius*tumblerRadius - (tumblerFaceEdgeWidth/2)*(tumblerFaceEdgeWidth/2))
	tumblerSpacing := 1.0

	bearingOD := 22.0
	bearingID := 8.0
	bearingThickness := 7.0
	bearingHolderStopConstriction := 1.0 //horizontle and vertical chamfer distance
	bearingHolderTolerance := 0.1
	bearingHolderThickness := 4.0

	shaftOD := 6.35
	spacerShaftTollerance := 0.3
	spacerBearingTollerance := 0.0
	spacerGapAngle := 3.0 / 360.0 * sdf.Tau
	spacerBearingPenetrationDepth := 3.0

	pusherNibSize := 3.0
	pusherLength := 8.5
	pusherTollerance := 1.5

	tumblerOutside := makeTumblerOutside(tumblerRadius, tumblerCornerRound, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight)

	insideHole := makeBearingHole(bearingOD, bearingThickness, bearingHolderStopConstriction, bearingHolderTolerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight)

	tracks, nibs := makePusherTracksAndNibs(pusherNibSize, pusherLength, tumblerRadius, bearingOD/2+bearingHolderThickness, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerSpacing, pusherTollerance)

	holes := sdf.Union3D(insideHole, tracks)
	tumbler := sdf.Difference3D(tumblerOutside, holes)
	tumbler = sdf.Union3D(tumbler, nibs)

	spacerDisk := makeSpacerDisk(shaftOD, spacerShaftTollerance, bearingID, spacerBearingTollerance, spacerBearingPenetrationDepth, tumblerSpacing, spacerGapAngle)

	//sdf.RenderSTLSlow(tumbler, 400, "tumbler.stl")
	//sdf.RenderSTLSlow(spacerDisk, 400, "spacerDisk.stl")
	//sdf.RenderSTL(tumbler, 200, "tumbler.stl")
	sdf.RenderSTL(spacerDisk, 200, "spacerDisk.stl")
	sdf.RenderSTL(makeHexTexturePlane(25, 25, 2, 3), 200, "hexPlane.stl")

}

func makeTumblerOutside(tumblerRadius, tumblerCornerRound, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight float64) sdf.SDF3 {
	//tumbler, triangle outside, bearing holder hole inside
	// triangle is an extruded 3 nagon
	triangle := sdf.Polygon2D(sdf.Nagon(3, tumblerRadius-tumblerCornerRound))
	triangle = sdf.Offset2D(triangle, tumblerCornerRound)
	tumblerOutside := sdf.Extrude3D(triangle, tumblerFaceEdgeHeight)
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

func makePusherTracksAndNibs(pusherNibSize, pusherLength, tumblerRadius, innerClearance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerSpacing, pusherTollerance float64) (sdf.SDF3, sdf.SDF3) {
	pusher2D := pusherNibProfile(pusherNibSize, pusherNibSize+tumblerSpacing)
	track2D := pusherNibProfile(pusherNibSize+pusherTollerance/2, pusherNibSize+tumblerSpacing+pusherTollerance)

	pushers := make([]sdf.SDF3, 3)
	tracks := make([]sdf.SDF3, 3)

	for n := 0; n < 2; n++ {
		//distanceFromCenter := (tumblerRadius+innerClearance)/2 + (pusherNibSize-(pusherTollerance/2))*(float64(n-1)) // three tracks one halfway between tumbler radioius and the inner clearance, and one on each side
		//distanceFromCenter := innerClearance + pusherTollerance + pusherNibSize*(float64(n)) // three tracks
		distanceFromCenter := innerClearance + pusherTollerance + pusherNibSize*(float64(n)) // two tracks
		startAngle := float64(n+1) * sdf.Tau / 2

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
				sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeHeight / 2}),
			),
		)

		tracks[n] = sdf.Transform3D(
			sdf.RevolveTheta3D(
				sdf.Transform2D(
					track2D,
					sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
				),
				2.0*sdf.Tau/3.0+pusherArcAngle,
			),
			sdf.RotateZ(startAngle-pusherArcAngle/2).Mul(
				sdf.Translate3d(sdf.V3{0, 0, -tumblerFaceEdgeHeight / 2}),
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
func makeBearingHole(bearingOD, bearingThickness, bearingHolderStopConstriction, bearingHolderTolerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight float64) sdf.SDF3 {

	bearingHolder := sdf.Cylinder3D(bearingThickness, bearingOD/2+bearingHolderTolerance, 0)
	chamfer3d := sdf.Cone3D(
		bearingHolderStopConstriction,
		bearingOD/2+bearingHolderTolerance-bearingHolderStopConstriction,
		bearingOD/2+bearingHolderTolerance,
		0,
	)
	constricted := sdf.Cylinder3D(tumblerFaceEdgeHeight-2*bearingThickness, bearingOD/2+bearingHolderTolerance-bearingHolderStopConstriction, 0)

	topBearingHolder := sdf.Transform3D(
		bearingHolder,
		sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeHeight/2 - bearingThickness/2}),
	)
	topChamfer3d := sdf.Transform3D(
		chamfer3d,
		sdf.Translate3d(sdf.V3{0, 0, tumblerFaceEdgeHeight/2 - bearingThickness - bearingHolderStopConstriction/2}),
	)
	flip := sdf.RotateX(sdf.Tau / 2)
	bottomBearingHolder := sdf.Transform3D(topBearingHolder, flip)
	bottomChamfer3d := sdf.Transform3D(topChamfer3d, flip)

	bearingHole := sdf.Union3D(topBearingHolder, topChamfer3d, constricted, bottomChamfer3d, bottomBearingHolder)
	return bearingHole
}

func makeShaftAndSpacerHole(shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness float64) sdf.SDF3 {
	innerChamfer := (tumblerShortRadius - tumblerMinimumWallThickness) - (shaftOD/2 + tumblerShaftTollerance)
	//innerChamfer := 10.0
	fmt.Printf("shaftOD: %v, tumblerShaftTollerance: %v, tumblerShaftBearingSurfaceLength: %v, spacerThickness: %v, spacerShaftTollerance: %v, spacerEdgeInFromTumblerEdge: %v, spacerTumblerTollerance: %v, tumblerFaceEdgeWidth: %v, tumblerFaceEdgeHeight: %v, tumblerShortRadius: %v, tumblerSpacing: %v, tumblerMinimumWallThickness: %v, innerChamfer: %v", shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness, innerChamfer)

	return sdf.Revolve3D(
		sdf.Polygon2D([]sdf.V2{
			{0, tumblerFaceEdgeHeight / 2},
			{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeHeight / 2},
			{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance},
			{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance},
			{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength},
			{shaftOD/2 + tumblerShaftTollerance + innerChamfer, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength - innerChamfer},
			{shaftOD/2 + tumblerShaftTollerance + innerChamfer, -tumblerFaceEdgeHeight/2 + tumblerShaftBearingSurfaceLength + innerChamfer},
			{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeHeight/2 + tumblerShaftBearingSurfaceLength},
			{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeHeight / 2},
			{0, -tumblerFaceEdgeHeight / 2},
			{0, tumblerFaceEdgeHeight / 2},
		}),
	)
}

func makeSpacerDisk(shaftOD, spacerShaftTollerance, bearingID, spacerBearingTollerance, spacerBearingPenetrationDepth, tumblerSpacing, spacerGapAngle float64) sdf.SDF3 {
	spacerHeight := spacerBearingPenetrationDepth*2 + tumblerSpacing
	return sdf.RevolveTheta3D(
		sdf.Polygon2D([]sdf.V2{
			{shaftOD/2 + spacerShaftTollerance, -spacerHeight / 2},
			{bearingID/2 - spacerBearingTollerance, -spacerHeight / 2},
			{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth},
			{bearingID/2 - spacerBearingTollerance + tumblerSpacing/2, -spacerHeight/2 + spacerBearingPenetrationDepth + tumblerSpacing/2},
			{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth + tumblerSpacing},
			{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth*2 + tumblerSpacing},
			{shaftOD/2 + spacerShaftTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth*2 + tumblerSpacing},
			{shaftOD/2 + spacerShaftTollerance, -spacerHeight / 2},
		}),
		sdf.Tau-spacerGapAngle,
	)
}

func makeHexTexturePlane(length, width, thickness, hexRadius float64) sdf.SDF3 {
	minorRadius := math.Cos(sdf.Tau/12) * hexRadius
	hex := sdf.Polygon2D(sdf.Nagon(6, hexRadius))
	cell := sdf.Extrude3D(hex, thickness)
	cell = sdf.Transform3D(
		cell,
		sdf.RotateX(math.Atan(thickness/2/minorRadius)),
	)

	startX := -length / 2 //we can center later
	startY := -width / 2  //we can center later
	countX := int(math.Ceil(length/(hexRadius*1.5)) + 1)
	countY := int(math.Ceil(width/(minorRadius*2)) + 1)
	incrementX := hexRadius * 1.5
	incrementY := minorRadius * 2 // math.Sqrt(minorRadius*minorRadius+thickness*thickness) * 2

	cells := make([]sdf.SDF3, countX*countY)

	for x := 0; x < countX; x++ {
		for y := 0; y < countY; y++ {
			cells[x*countY+y] = sdf.Transform3D(
				cell,
				sdf.Translate3d(sdf.V3{startX + float64(x)*incrementX, startY + float64(y)*incrementY - (float64(x%2) * incrementY / 2), 0}),
			)
		}
	}

	return sdf.Intersect3D(
		sdf.Union3D(cells...),
		sdf.Box3D(sdf.V3{length, width, thickness * 2}, 0),
	)
	//return sdf.Union3D(cells...)

}

/*
face patterns
A|111|001|111|111|101|111|111|111|111|111
B|101|001|001|001|101|100|100|001|101|101
C|101|001|111|111|111|111|111|001|111|111
D|101|001|100|001|001|001|101|001|101|001
E|111|001|111|111|001|111|111|001|111|001


A
001
101
111

B
001
100
101
*/
