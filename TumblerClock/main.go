package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	tumblerFaceEdgeLength := 50.0
	tumblerRadius := tumblerFaceEdgeLength / 2 / math.Cos(sdf.Tau/12)
	tumblerCornerRound := tumblerRadius * 0.05

	bearingOD := 22.0
	bearingThickness := 7.0
	bearingHolderStopConstriction := 1.0 //horizontle and vertical chamfer distance
	bearingHolderTolerance := 0.1

	pusherNibSize := 3.0
	pusherTollerance := 0.2

	//tumbler, triangle outside, bearing holder hole inside
	// triangle is an extruded 3 nagon
	triangle := sdf.Polygon2D(sdf.Nagon(3, tumblerRadius-tumblerCornerRound))
	triangle = sdf.Offset2D(triangle, tumblerCornerRound)
	tumblerOutside := sdf.Extrude3D(triangle, tumblerFaceEdgeLength)

	//bearing holder hole is an extruded hex with a constriction in the middle
	//this is accomplished by stacking 5 shapes
	//top and bottom are bearing sized cylinders
	//center is constricted cylinder
	//between top and center (and bottom and center) are correctly oriented truncated between bearing and constricted sized cylinders

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

	pusherNib2D := pusherNibProfile(pusherNibSize)
	pusherNib3D := sdf.Revolve3D(pusherNib2D)

	pusherTracks := make([]sdf.SDF3, 3)
	pusherNibs := make([]sdf.SDF3, 3)

	for n := 0; n < 3; n++ {
		distanceFromCenter := (tumblerRadius+bearingOD/2)/2 + pusherNibSize*(float64(n)-1.0)
		startAngle := float64(n) * sdf.Tau / 3
		pusherNibs[n] = sdf.Transform3D(
			pusherNib3D,
			sdf.RotateZ(startAngle).Mul(
				sdf.Translate3d(sdf.V3{distanceFromCenter, 0, tumblerFaceEdgeLength / 2}),
			),
		)

		pusherNibHole2D := sdf.Offset2D(pusherNib2D, pusherTollerance)
		track := sdf.RevolveTheta3D(
			sdf.Transform2D(
				pusherNibHole2D,
				sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
			),
			-sdf.Tau/3,
		)

		track = sdf.Union3D(
			track,
			sdf.Transform3D(
				sdf.Revolve3D(pusherNibHole2D),
				sdf.Translate3d(sdf.V3{distanceFromCenter, 0, 0}),
			),
			sdf.Transform3D(
				sdf.Revolve3D(pusherNibHole2D),
				sdf.RotateZ(sdf.Tau/3).Mul(
					sdf.Translate3d(sdf.V3{distanceFromCenter, 0, 0}),
				),
			),
		)

		pusherTracks[n] = sdf.Transform3D(
			track,
			sdf.RotateZ(startAngle).Mul(
				sdf.Translate3d(sdf.V3{0, 0, -tumblerFaceEdgeLength / 2}),
			),
		)
	}

	tracks := sdf.Union3D(pusherTracks...)
	nibs := sdf.Union3D(pusherNibs...)
	holes := sdf.Union3D(bearingHole, tracks)
	tumbler := sdf.Difference3D(tumblerOutside, holes)
	tumbler = sdf.Union3D(tumbler, nibs)

	sdf.RenderSTL(tumbler, 200, "tumbler.stl")

}

func pusherNibProfile(size float64) sdf.SDF2 {
	return sdf.Polygon2D([]sdf.V2{
		{-size / 2, 0},
		{-size / 2, size},
		{0, size * 1.5},
		{size / 2, size},
		{size / 2, 0},
		{-size / 2, 0},
	})
}
