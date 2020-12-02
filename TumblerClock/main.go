package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/deadsy/sdfx/sdf"
	"gopkg.in/yaml.v3"
)

//-----------------------------------------------------------------------------

func main() {
	//setup := makeDefaultClockSetup()
	//bytes, _ := yaml.Marshal(setup)
	//ioutil.WriteFile("setup.yaml", bytes, 0644)

	filename := flag.String("SetupFile", "setup.yaml", "File name with the setup for the Tumbler Clock")
	//partsList := flag.String("Parts", "all", "Parts list to print") //(h|m)(t|u)(a|b|ag|g)

	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	var setup ClockSetup
	err = yaml.Unmarshal(yamlFile, &setup)
	if err != nil {
		panic(err)
	}
	setup = setup.computeSynthetics()

	facesA := [9]int{0, 0, 1, 1, 0, 1, 1, 1, 1} // A faces
	facesB := [9]int{0, 0, 1, 1, 0, 0, 1, 0, 1} // B faces
	tumblerOutsideA := makeTumblerOutside(setup.Tumbler, facesA)
	tumblerOutsideB := makeTumblerOutside(setup.Tumbler, facesB)

	insideHole := makeBearingHole(setup.Bearing, setup.BearingHolder, setup.Tumbler)

	tracks, nibs := makePusherTracksAndNibs(setup.Transmission, setup.Tumbler, setup.Bearing.OD/2+setup.BearingHolder.Thickness)

	holes := sdf.Union3D(insideHole, tracks)
	tumblerA := sdf.Difference3D(tumblerOutsideA, holes)
	tumblerA = sdf.Union3D(tumblerA, nibs)

	tumblerB := sdf.Difference3D(tumblerOutsideB, holes)
	tumblerB = sdf.Union3D(tumblerB, nibs)

	//spacerDisk := makeSpacerDisk(shaftOD, spacerShaftTollerance, bearingID, spacerBearingTollerance, spacerBearingPenetrationDepth, tumblerSpacing, spacerGapAngle)
	spacerDisk := makeSimpleSpacerDisk(setup.Shaft, setup.Spacer, setup.Tumbler)
	// sdf.RenderSTLSlow(tumblerA, 400, "tumblerA.stl")
	// sdf.RenderSTLSlow(tumblerB, 400, "tumblerB.stl")
	//sdf.RenderSTLSlow(spacerDisk, 100, "spacerDisk.stl")
	sdf.RenderSTL(tumblerA, 200, "tumblerA.stl")
	sdf.RenderSTL(tumblerB, 200, "tumblerB.stl")
	sdf.RenderSTL(spacerDisk, 400, "spacerDiskWood.stl")

}

func makeTumblerOutside(tumbler Tumbler, faces [9]int) sdf.SDF3 {
	//tumbler, triangle outside, bearing holder hole inside
	// triangle is an extruded 3 nagon
	triangle := sdf.Polygon2D(sdf.Nagon(3, tumbler.Radius-tumbler.CornerRound))
	triangle = sdf.Offset2D(triangle, tumbler.CornerRound)
	tumblerOutside := sdf.Extrude3D(triangle, tumbler.FaceEdgeHeight)

	textureWidth := (tumbler.FaceEdgeWidth - 2.0*tumbler.CornerRound) / 3.0

	textureOff := sdf.Transform3D(
		makeUnTexturedPlane(textureWidth, tumbler.FaceEdgeHeight, 2),
		sdf.Translate3d(sdf.V3{textureWidth / 2.0, -tumbler.CornerRound, 0}).Mul( //shift to possition 0 on face
			sdf.RotateX(sdf.Tau/4),
		),
	)
	textureOn := sdf.Transform3D(
		makeHexTexturePlane(textureWidth, tumbler.FaceEdgeHeight, 2, 3),
		sdf.Translate3d(sdf.V3{textureWidth / 2.0, -tumbler.CornerRound, 0}).Mul( //shift to possition 0 on face
			sdf.RotateY(sdf.Tau/2.0).Mul( //flip for 1
				sdf.RotateX(sdf.Tau/4),
			),
		),
	)

	//		sdf.Translate3d(sdf.V3{textureWidth/2.0 + tumblerCornerRound, 0, 0}).Mul(

	for face := 0; face < 3; face++ {
		for possition := 0; possition < 3; possition++ {
			textureToAdd := textureOff
			if faces[face*3+possition] == 1 {
				textureToAdd = textureOn
			}

			faceToAdd := sdf.Transform3D(
				textureToAdd,
				sdf.RotateZ(sdf.Tau*float64(face)/3.0).Mul( //rotate to face
					sdf.Translate3d(sdf.V3{tumbler.Radius - tumbler.CornerRound, 0, 0}).Mul( //translate to first face
						sdf.RotateZ(sdf.Tau*5.0/12.0).Mul( //allign to first face
							sdf.Translate3d(sdf.V3{float64(possition) * textureWidth, 0, 0}), //shift to possition on face
						),
					),
				),
			)
			if possition == 1 {
				tumblerOutside = sdf.Union3D(tumblerOutside, faceToAdd)
			} else {
				tumblerOutside = sdf.Difference3D(tumblerOutside, faceToAdd)
			}

		}
	}

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

func makePusherTracksAndNibs(transmission Transmission, tumbler Tumbler, innerClearance float64) (sdf.SDF3, sdf.SDF3) {
	pusher2D := pusherNibProfile(transmission.NibSize, transmission.NibSize+tumbler.Spacing)
	track2D := pusherNibProfile(transmission.NibSize+transmission.TrackTollerance/2, transmission.NibSize+tumbler.Spacing+transmission.TrackTollerance)

	pushers := make([]sdf.SDF3, 3)
	tracks := make([]sdf.SDF3, 3)

	for n := 0; n < 2; n++ {
		distanceFromCenter := innerClearance + transmission.TrackTollerance + transmission.NibSize*(float64(n)) // two tracks
		startAngle := float64(n+1) * sdf.Tau / 2

		pusherArcAngle := transmission.NibLength / distanceFromCenter

		pushers[n] = sdf.Transform3D(
			sdf.RevolveTheta3D(
				sdf.Transform2D(
					pusher2D,
					sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
				),
				pusherArcAngle,
			),
			sdf.RotateZ(startAngle-pusherArcAngle/2).Mul(
				sdf.Translate3d(sdf.V3{0, 0, tumbler.FaceEdgeHeight / 2}),
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
				sdf.Translate3d(sdf.V3{0, 0, -tumbler.FaceEdgeHeight / 2}),
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
func makeBearingHole(bearing Bearing, bearingHolder BearingHolder, tumbler Tumbler) sdf.SDF3 {

	bearingHole := sdf.Cylinder3D(bearing.Thickness, bearing.OD/2+bearingHolder.Tolerance, 0)
	chamfer3d := sdf.Cone3D(
		bearingHolder.StopConstriction,
		bearing.OD/2+bearingHolder.Tolerance-bearingHolder.StopConstriction,
		bearing.OD/2+bearingHolder.Tolerance,
		0,
	)
	constricted := sdf.Cylinder3D(tumbler.FaceEdgeHeight-2*bearing.Thickness, bearing.OD/2+bearingHolder.Tolerance-bearingHolder.StopConstriction, 0)

	topBearingHolder := sdf.Transform3D(
		bearingHole,
		sdf.Translate3d(sdf.V3{0, 0, tumbler.FaceEdgeHeight/2 - bearing.Thickness/2}),
	)
	topChamfer3d := sdf.Transform3D(
		chamfer3d,
		sdf.Translate3d(sdf.V3{0, 0, tumbler.FaceEdgeHeight/2 - bearing.Thickness - bearingHolder.StopConstriction/2}),
	)
	flip := sdf.RotateX(sdf.Tau / 2)
	bottomBearingHolder := sdf.Transform3D(topBearingHolder, flip)
	bottomChamfer3d := sdf.Transform3D(topChamfer3d, flip)

	return sdf.Union3D(topBearingHolder, topChamfer3d, constricted, bottomChamfer3d, bottomBearingHolder)
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

func makeSimpleSpacerDisk(shaft Shaft, spacer Spacer, tumbler Tumbler) sdf.SDF3 {
	return sdf.RevolveTheta3D(
		sdf.Polygon2D([]sdf.V2{
			{shaft.OD/2 + spacer.ShaftTollerance, 0},
			{shaft.OD/2 + spacer.ShaftTollerance + spacer.DiskWidth, 0},
			{shaft.OD/2 + spacer.ShaftTollerance + spacer.DiskWidth, tumbler.Spacing},
			{shaft.OD/2 + spacer.ShaftTollerance, tumbler.Spacing},
			{shaft.OD/2 + spacer.ShaftTollerance, 0},
		}),
		sdf.Tau-spacer.GapAngle,
	)
}

func makeUnTexturedPlane(length, width, thickness float64) sdf.SDF3 {
	return sdf.Transform3D(
		sdf.Box3D(sdf.V3{length, width, thickness}, 0),
		sdf.RotateX(math.Atan(thickness/width)),
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
