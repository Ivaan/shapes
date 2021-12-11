package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

type tumblerBuilder struct {
	setup  ClockSetup
	facesA [9]int
	facesB [9]int
	// tumblerOutsideA    sdf.SDF3
	// tumblerOutsideB    sdf.SDF3
	tracks sdf.SDF3
	nibs   sdf.SDF3
	// tumblerHole        sdf.SDF3
	// gearHole           sdf.SDF3
	// tumblerAndGearHold sdf.SDF3
}

func makeTumblerBuilder(setup ClockSetup) tumblerBuilder {
	facesA := [9]int{0, 0, 1, 1, 0, 1, 1, 1, 1} // A faces
	facesB := [9]int{0, 0, 1, 1, 0, 0, 1, 0, 1} // B faces
	tracks, nibs := makePusherTracksAndNibs(setup.Transmission, setup.Tumbler, setup.Bearing.OD/2+setup.BearingHolder.Thickness)

	return tumblerBuilder{
		setup:  setup,
		facesA: facesA,
		facesB: facesB,
		// tumblerOutsideA:    makeTumblerOutside(setup.Tumbler, facesA),
		// tumblerOutsideB:    makeTumblerOutside(setup.Tumbler, facesB),
		tracks: tracks,
		nibs:   nibs,
		// tumblerHole:        makeBearingHole(setup.Bearing, setup.BearingHolder, setup.Tumbler.FaceEdgeHeight),
		// gearHole:           makeBearingHole(setup.Bearing, setup.BearingHolder, setup.Gear.Thickness),
		// tumblerAndGearHold: makeBearingHole(setup.Bearing, setup.BearingHolder, setup.Gear.Thickness+setup.Gear.SpaceToTumbler+setup.Tumbler.FaceEdgeHeight),
	}

}

//a tumbler is either the tiangular tumbler, a gear or both
//the triangular part is the only part that has pusher nibs
//the triangular part will all way have it's bottom on the xy plane
//the gear will be below the xy plane if there are both and on the xy plane if it is only a gear
//the pusher tracks only go on gears or triangls alone (not when there are both)
//therefore pushers tracks, when needed, will cut up from the xy plane
func (tb *tumblerBuilder) makeTumbler(p tumblerPart) sdf.SDF3 {
	var positives []sdf.SDF3
	var negatives []sdf.SDF3
	bearingHoleHeight := 0.0

	if p.tumbler != none {
		faces := tb.facesA
		if p.tumbler == bFace {
			faces = tb.facesB
		}
		if p.tens {
			faces = flipFacesForTens(faces)
		}

		tumbler := makeTumblerOutside(tb.setup.Tumbler, faces)
		if p.tens { //tens column digit is printed upside down
			tumbler = sdf.Transform3D(
				tumbler,
				sdf.RotateX(sdf.Tau/2.0),
			)
		}
		tumbler = sdf.Transform3D(
			tumbler,
			sdf.Translate3d(sdf.V3{0, 0, tb.setup.Tumbler.FaceEdgeHeight / 2}),
		)

		positives = append(positives, tumbler)
		positives = append(positives, tb.nibs)
		bearingHoleHeight += tb.setup.Tumbler.FaceEdgeHeight
	}

	if p.gear {
		numberOfTeeth := tb.setup.Gear.CouplerGearNumberOfTeeth
		if p.tens != (p.tumbler != none) { //p.tens && !p.tumbler || !p.tens && p.tumbler
			numberOfTeeth = tb.setup.Gear.DrivenGearNumberOfTeeth
		}
		gear, gearTumblerJoiner, joinerHeight := makeGear(tb.setup.Tumbler, tb.setup.Gear, numberOfTeeth)
		if p.tens && (p.tumbler != none) {
			gear = sdf.Transform3D(
				gear,
				sdf.MirrorXY(),
			)
		}
		if p.tumbler != none {
			gear = sdf.Transform3D(
				gear,
				sdf.Translate3d(sdf.V3{0, 0, -tb.setup.Gear.Thickness/2 - joinerHeight}),
			)
			gearTumblerJoiner = sdf.Transform3D(
				gearTumblerJoiner,
				sdf.Translate3d(sdf.V3{0, 0, -joinerHeight / 2}),
			)
			positives = append(positives, gear)
			bearingHoleHeight += tb.setup.Gear.Thickness
			if joinerHeight > 0 {
				positives = append(positives, gearTumblerJoiner)
				bearingHoleHeight += joinerHeight
			}
		} else {
			gear = sdf.Transform3D(
				gear,
				sdf.Translate3d(sdf.V3{0, 0, tb.setup.Gear.Thickness / 2}),
			)
			positives = append(positives, gear)
			bearingHoleHeight += tb.setup.Gear.Thickness
		}
	}

	if p.gear != (p.tumbler != none) { //one or the other but not both
		tracks := tb.tracks
		if p.hours {
			tracks = sdf.Transform3D(
				tracks,
				sdf.RotateZ(2.0*sdf.Tau/3.0), //rotate tracks so nibs engage ccw when faces aligned
			)
		}
		negatives = append(negatives, tracks)
	}
	bearingHole := makeBearingHole(tb.setup.Bearing, tb.setup.BearingHolder, bearingHoleHeight)

	if p.gear && p.tumbler != none {
		bearingHole = sdf.Transform3D(
			bearingHole,
			sdf.Translate3d(sdf.V3{0, 0, tb.setup.Tumbler.FaceEdgeHeight - bearingHoleHeight/2}),
		)
	} else {
		bearingHole = sdf.Transform3D(
			bearingHole,
			sdf.Translate3d(sdf.V3{0, 0, bearingHoleHeight / 2}),
		)
	}

	negatives = append(negatives, bearingHole)

	return sdf.Difference3D(sdf.Union3D(positives...), sdf.Union3D(negatives...))
}

func makeTumblerOutside(tumbler Tumbler, faces [9]int) sdf.SDF3 {
	//tumbler, triangle outside, bearing holder hole inside
	// triangle is an extruded 3 nagon
	triangle, _ := sdf.Polygon2D(sdf.Nagon(3, tumbler.Radius-tumbler.CornerRound))
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

func makeGear(tumbler Tumbler, gear Gear, numberOfTeeth int) (gear3D, joiner sdf.SDF3, joinerHeight float64) {
	gearModule := tumbler.Radius * 2.0 / float64(gear.CouplerGearNumberOfTeeth) //comput module for coupling gear, all other gears are different sizes by number of teeth
	pa := sdf.DtoR(20.0)

	pitchRadius := float64(numberOfTeeth) * gearModule / 2.0
	dedendum := gearModule + gear.clearance
	rootRadius := pitchRadius - dedendum

	gp := obj.InvoluteGearParms{
		NumberTeeth:   numberOfTeeth,
		Module:        gearModule,
		PressureAngle: pa,
		Backlash:      gear.backlash,
		Clearance:     gear.clearance,
		Facets:        7,
	}

	gear2d, _ := obj.InvoluteGear(&gp)

	gear3D = sdf.TwistExtrude3D(gear2d, gear.Thickness, sdf.Tau/float64(numberOfTeeth))
	if gearModule+rootRadius < tumbler.Radius {
		joinerHeight = gearModule
	} else {
		joinerHeight = tumbler.Radius - rootRadius
	}
	if joinerHeight > 0 {
		joiner, _ = sdf.Cone3D(joinerHeight, rootRadius, rootRadius+joinerHeight, 0)
	} else {
		joinerHeight = 0
		joiner = nil
	}

	return gear3D, joiner, joinerHeight
}

//bearing holder hole is a cylinder with a constriction in the middle
//this is accomplished by stacking 5 shapes
//top and bottom are bearing sized cylinders
//center is constricted cylinder
//between top and center (and bottom and center) are correctly oriented truncated between bearing and constricted sized cylinders
func makeBearingHole(bearing Bearing, bearingHolder BearingHolder, height float64) sdf.SDF3 {

	bearingHole, _ := sdf.Cylinder3D(bearing.Thickness, bearing.OD/2+bearingHolder.Tolerance, 0)
	chamfer3d, _ := sdf.Cone3D(
		bearingHolder.StopConstriction,
		bearing.OD/2+bearingHolder.Tolerance-bearingHolder.StopConstriction,
		bearing.OD/2+bearingHolder.Tolerance,
		0,
	)
	constricted, _ := sdf.Cylinder3D(height-2*bearing.Thickness, bearing.OD/2+bearingHolder.Tolerance-bearingHolder.StopConstriction, 0)

	topBearingHolder := sdf.Transform3D(
		bearingHole,
		sdf.Translate3d(sdf.V3{0, 0, height/2 - bearing.Thickness/2}),
	)
	topChamfer3d := sdf.Transform3D(
		chamfer3d,
		sdf.Translate3d(sdf.V3{0, 0, height/2 - bearing.Thickness - bearingHolder.StopConstriction/2}),
	)
	flip := sdf.RotateX(sdf.Tau / 2)
	bottomBearingHolder := sdf.Transform3D(topBearingHolder, flip)
	bottomChamfer3d := sdf.Transform3D(topChamfer3d, flip)

	return sdf.Union3D(topBearingHolder, topChamfer3d, constricted, bottomChamfer3d, bottomBearingHolder)
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

		pusher3D, _ := sdf.RevolveTheta3D(
			sdf.Transform2D(
				pusher2D,
				sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
			),
			pusherArcAngle,
		)
		pushers[n] = sdf.Transform3D(
			pusher3D,
			sdf.RotateZ(startAngle-pusherArcAngle/2).Mul(
				sdf.Translate3d(sdf.V3{0, 0, tumbler.FaceEdgeHeight}),
			),
		)

		track3D, _ := sdf.RevolveTheta3D(
			sdf.Transform2D(
				track2D,
				sdf.Translate2d(sdf.V2{distanceFromCenter, 0}),
			),
			2.0*sdf.Tau/3.0+pusherArcAngle,
		)
		tracks[n] = sdf.Transform3D(
			track3D,
			sdf.RotateZ(startAngle-pusherArcAngle/2-2.0*sdf.Tau/3.0),
		)
	}
	return sdf.Union3D(tracks...), sdf.Union3D(pushers...)
}

func pusherNibProfile(width, height float64) sdf.SDF2 {
	nib, _ := sdf.Polygon2D([]sdf.V2{
		{-width / 2, 0},
		{-width / 2, height},
		{0, height + width*0.5},
		{width / 2, height},
		{width / 2, 0},
		{-width / 2, 0},
	})

	return nib
}

func flipFacesForTens(faces [9]int) [9]int {
	var facesOut [9]int
	for face := 0; face < 3; face++ {
		for possition := 0; possition < 3; possition++ {
			facesOut[face*3+possition] = faces[(2-face)*3+possition]
		}
	}
	return facesOut
}

func makeUnTexturedPlane(length, width, thickness float64) sdf.SDF3 {
	box, _ := sdf.Box3D(sdf.V3{length, width, thickness}, 0)
	return sdf.Transform3D(
		box,
		sdf.RotateX(math.Atan(thickness/width)),
	)
}

func makeHexTexturePlane(length, width, thickness, hexRadius float64) sdf.SDF3 {
	minorRadius := math.Cos(sdf.Tau/12) * hexRadius
	hex, _ := sdf.Polygon2D(sdf.Nagon(6, hexRadius))
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

	box, _ := sdf.Box3D(sdf.V3{length, width, thickness * 2}, 0)
	return sdf.Intersect3D(
		sdf.Union3D(cells...),
		box,
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
