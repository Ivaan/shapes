package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

type KeyNodule struct {
	Top          Nodule
	Bottom       Nodule
	keycapHitbox sdf.SDF3
	switchHitbox sdf.SDF3
}

func (kn KeyNodule) GetHitBoxes() []sdf.SDF3 {
	return []sdf.SDF3{kn.keycapHitbox, kn.switchHitbox}
}

type BubbleKeyNoduleProperties struct {
	sphereRadius                     float64
	plateTopAtRadius                 float64
	plateThickness                   float64
	sphereThicknes                   float64
	backCoverCutAtRadius             float64
	switchHoleLength                 float64
	switchHoleWidth                  float64
	switchLatchWidth                 float64
	switchLatchGrabThickness         float64
	switchFlatzoneLength             float64
	switchFlatzoneWidth              float64
	keycapLength                     float64
	keycapWidth                      float64
	keycapBottomHeightAbovePlateDown float64
	keycapHeight                     float64
	keycapBottomHeightAbovePlateUp   float64
	keycapClearanced                 float64
	keycapRound                      float64
	huggingCylinderRound             float64
	laneWidth                        float64 //as in "Stay in your lane" this restricts the holes to a max width
	insertLength                     float64
	insertDiameter                   float64
	insertWallThickness              float64
	screwThreadDiameter              float64
	screwThreadLength                float64
	screwHeadDiameter                float64
}

func (knp BubbleKeyNoduleProperties) MakeBubbleKey(screwPossitionsBits int64) KeyNodule {
	sphereCenterZ := -knp.plateTopAtRadius - knp.keycapBottomHeightAbovePlateUp
	topOfPlateZ := -knp.keycapBottomHeightAbovePlateUp
	bottomOfPlateZ := topOfPlateZ - knp.plateThickness
	radiusAtTopOfPlate := math.Sqrt(knp.sphereRadius*knp.sphereRadius - knp.plateTopAtRadius*knp.plateTopAtRadius)
	backCoverCutZ := sphereCenterZ + knp.backCoverCutAtRadius
	keycapBottomWhenDownZ := topOfPlateZ + knp.keycapBottomHeightAbovePlateDown
	screwRadiusFromCenter := radiusAtTopOfPlate - (knp.insertDiameter/2 + knp.insertWallThickness)

	//solidSphere is the main outer shell
	solidSphere, err := Sphere3DAtHeight(knp.sphereRadius, sphereCenterZ)
	if err != nil {
		panic(err)
	}

	//huggingCylinder sits on top of the plate, forms the case around the keycaps
	huggingCylinder, err := sdf.Cylinder3D((knp.keycapHeight+knp.keycapBottomHeightAbovePlateDown)/2+knp.huggingCylinderRound, radiusAtTopOfPlate, knp.huggingCylinderRound)
	if err != nil {
		panic(err)
	}
	huggingCylinder = sdf.Transform3D(huggingCylinder, sdf.Translate3d(sdf.V3{Z: ((knp.keycapHeight+knp.keycapBottomHeightAbovePlateDown)/2+knp.huggingCylinderRound)/2 - knp.huggingCylinderRound - knp.keycapBottomHeightAbovePlateUp}))

	//hollow subtracts from the solidSphere up to the bottom of the plate
	hollow, err := Sphere3DAtHeight(knp.sphereRadius-knp.sphereThicknes, sphereCenterZ)
	if err != nil {
		panic(err)
	}
	hollow = sdf.Cut3D(hollow, sdf.V3{X: 0, Y: 0, Z: bottomOfPlateZ}, sdf.V3{X: 0, Y: 0, Z: -1})

	//Clearing cylinders are to remove artifacts from only nodules partially subtracted
	topClearingCylinder, err := Cylinder3DBelow(knp.sphereRadius*2, knp.sphereRadius-knp.sphereThicknes, 0, backCoverCutZ)
	if err != nil {
		panic(err)
	}
	bottomClearingCylinder, err := Cylinder3DAbove(knp.sphereRadius*2, knp.sphereRadius-knp.sphereThicknes, 0, backCoverCutZ)
	if err != nil {
		panic(err)
	}

	//Hole through the plate for the switch
	switchHole, err := Box3DBelow(sdf.V3{X: knp.switchHoleWidth, Y: knp.switchHoleLength, Z: knp.plateThickness}, 0, topOfPlateZ)
	if err != nil {
		panic(err)
	}

	//todo: add latch reliefs

	//switchFlatzone is the area on the top of the plate reserved for the switch foot print
	switchFlatzone, err := Box3DAbove(sdf.V3{X: knp.switchFlatzoneWidth, Y: knp.switchFlatzoneLength, Z: knp.keycapBottomHeightAbovePlateDown}, 0, topOfPlateZ)
	if err != nil {
		panic(err)
	}

	keyCapClearanceShadow := sdf.Box2D(sdf.V2{X: knp.keycapWidth + knp.keycapClearanced, Y: knp.keycapLength + knp.keycapClearanced}, knp.keycapRound+knp.keycapClearanced)
	keyCapClearance, err := ExtrudeRounded3DAbove(keyCapClearanceShadow, knp.keycapHeight*2, 0, keycapBottomWhenDownZ)
	if err != nil {
		panic(err)
	}

	lane, err := Box3DAndTranslate(sdf.V3{X: knp.laneWidth, Y: knp.sphereRadius * 2, Z: knp.sphereRadius * 2}, 0, sdf.V3{Z: sphereCenterZ})
	if err != nil {
		panic(err)
	}

	coverCutA := sdf.V3{Z: backCoverCutZ}
	plateCut := sdf.V3{Z: bottomOfPlateZ}
	coverTopV := sdf.V3{Z: 1}
	coverBottomtV := sdf.V3{Z: -1}
	shellTop := sdf.Cut3D(solidSphere, coverCutA, coverTopV)
	shellBottom := sdf.Cut3D(solidSphere, coverCutA, coverBottomtV)
	plate := sdf.Cut3D(solidSphere, plateCut, coverTopV)

	numberOfScrews := 4 //max 64 cuz screwPossitionsBits is int64
	insertHolders := make([]sdf.SDF3, numberOfScrews)
	screwChannels := make([]sdf.SDF3, numberOfScrews)
	screwHoles := make([]sdf.SDF3, numberOfScrews)
	insertHoldersHoles := make([]sdf.SDF3, numberOfScrews)

	for i := 0; i < numberOfScrews; i++ {
		if screwPossitionsBits&(1<<i) == 0 {
			continue
		}

		angle := float64(i) * sdf.Tau / float64(numberOfScrews)
		rotateIntoPlace := sdf.RotateZ(angle).Mul(sdf.Translate3d(sdf.V3{X: screwRadiusFromCenter}))

		holder, err := Cylinder3DAbove(knp.insertLength+knp.insertWallThickness, knp.insertDiameter/2+knp.insertWallThickness, 0, backCoverCutZ)
		if err != nil {
			panic(err)
		}
		holderHole, err := Cylinder3DAbove(knp.insertLength, knp.insertDiameter/2, 0, backCoverCutZ)
		if err != nil {
			panic(err)
		}
		holder = sdf.Transform3D(holder, rotateIntoPlace)
		holderHole = sdf.Transform3D(holderHole, rotateIntoPlace)
		insertHolders[i] = holder
		insertHoldersHoles[i] = holderHole

		screwChannel, err := Cylinder3DBelow(knp.sphereRadius, knp.screwHeadDiameter/2+knp.insertWallThickness, 0, backCoverCutZ)
		if err != nil {
			panic(err)
		}
		screwThreadHole, err := Cylinder3DBelow(knp.screwThreadLength-knp.insertLength, knp.screwThreadDiameter/2, 0, backCoverCutZ)
		if err != nil {
			panic(err)
		}
		screwHeadHole, err := Cylinder3DBelow(knp.sphereRadius, knp.screwHeadDiameter/2, 0, backCoverCutZ-(knp.screwThreadLength-knp.insertLength))
		if err != nil {
			panic(err)
		}

		screwChannel = sdf.Transform3D(screwChannel, rotateIntoPlace)
		screwChannel = sdf.Intersect3D(shellBottom, screwChannel)
		screwThreadHole = sdf.Transform3D(screwThreadHole, rotateIntoPlace)
		screwHeadHole = sdf.Transform3D(screwHeadHole, rotateIntoPlace)
		screwHole := sdf.Union3D(screwThreadHole, screwHeadHole)
		screwChannels[i] = screwChannel
		screwHoles[i] = screwHole
	}

	var allInsertHolders, allInsertHoldersHoles, allScrewChannels, allScrewHoles sdf.SDF3
	if screwPossitionsBits > 0 {
		allInsertHolders = sdf.Union3D(insertHolders...)
		allInsertHoldersHoles = sdf.Union3D(insertHoldersHoles...)
		allScrewChannels = sdf.Union3D(screwChannels...)
		allScrewHoles = sdf.Union3D(screwHoles...)
	}

	return KeyNodule{
		Top: MakeNodule(
			[]sdf.SDF3{},
			[]sdf.SDF3{},
			[]sdf.SDF3{switchHole, switchFlatzone, keyCapClearance, sdf.Intersect3D(topClearingCylinder, lane), allInsertHoldersHoles},
			[]sdf.SDF3{allInsertHolders},
			[]sdf.SDF3{sdf.Intersect3D(hollow, lane)},                 //hole rank 0
			[]sdf.SDF3{sdf.Intersect3D(plate, lane), huggingCylinder}, //solid rank 0
			[]sdf.SDF3{hollow, shellBottom},                           //hole rank 1
			[]sdf.SDF3{shellTop},                                      //solid rank 1
		),
		Bottom: MakeNodule(
			[]sdf.SDF3{},
			[]sdf.SDF3{},
			[]sdf.SDF3{switchHole, switchFlatzone, bottomClearingCylinder, allScrewHoles}, //hole rank 0
			[]sdf.SDF3{allScrewChannels}, //solid rank 0
			[]sdf.SDF3{hollow},           //hole rank 0
			[]sdf.SDF3{shellBottom},      //solid rank 0
		),
		//keycapHitbox sdf.SDF3
		//switchHitbox sdf.SDF3
	}
}

func Cylinder3DAbove(height, radius, round, bottomZ float64) (sdf.SDF3, error) {
	return Cylinder3DAndTranslate(height, radius, round, sdf.V3{Z: height/2 + bottomZ})
}

func Cylinder3DBelow(height, radius, round, topZ float64) (sdf.SDF3, error) {
	return Cylinder3DAndTranslate(height, radius, round, sdf.V3{Z: -height/2 + topZ})
}

func Cylinder3DAndTranslate(height, radius, round float64, move sdf.V3) (sdf.SDF3, error) {
	c, err := sdf.Cylinder3D(height, radius, round)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(c, sdf.Translate3d(move)), nil
}

func Box3DAbove(size sdf.V3, round, bottomZ float64) (sdf.SDF3, error) {
	return Box3DAndTranslate(size, round, sdf.V3{Z: size.Z/2 + bottomZ})
}

func Box3DBelow(size sdf.V3, round, topZ float64) (sdf.SDF3, error) {
	return Box3DAndTranslate(size, round, sdf.V3{Z: -size.Z/2 + topZ})
}

func Box3DAndTranslate(size sdf.V3, round float64, move sdf.V3) (sdf.SDF3, error) {
	b, err := sdf.Box3D(size, round)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(b, sdf.Translate3d(move)), nil
}

func ExtrudeRounded3DAbove(sdf2 sdf.SDF2, height, round, bottomZ float64) (sdf.SDF3, error) {
	return ExtrudeRounded3DAndTranslate(sdf2, height, round, sdf.V3{Z: height/2 + bottomZ})
}

func ExtrudeRounded3DBelow(sdf2 sdf.SDF2, height, round, topZ float64) (sdf.SDF3, error) {
	return ExtrudeRounded3DAndTranslate(sdf2, height, round, sdf.V3{Z: -height/2 + topZ})
}

func ExtrudeRounded3DAndTranslate(sdf2 sdf.SDF2, height, round float64, move sdf.V3) (sdf.SDF3, error) {
	e, err := sdf.ExtrudeRounded3D(sdf2, height, round)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(e, sdf.Translate3d(move)), nil
}

func Sphere3DAtHeight(radius, height float64) (sdf.SDF3, error) {
	s, err := sdf.Sphere3D(radius)
	if err != nil {
		return nil, err
	}
	return sdf.Transform3D(s, sdf.Translate3d(sdf.V3{Z: height})), nil
}

type FlatterKeyNoduleProperties struct {
	sphereRadius             float64
	sphereCut                float64
	plateThickness           float64
	sphereThicknes           float64
	backCoverLipCut          float64
	switchHoleLength         float64
	switchHoleWidth          float64
	switchHoleDepth          float64
	switchLatchWidth         float64
	switchLatchGrabThickness float64
	switchFlatzoneWidth      float64
	switchFlatzoneLength     float64
	pcbLength                float64
	pcbWidth                 float64
	keycapWidth              float64
	keycapHeight             float64
	keycapRound              float64
	keycapOffset             float64
}

func (knp FlatterKeyNoduleProperties) MakeFlatterKey(orientAndMove sdf.M44) (*KeyNodule, error) {
	shell, err := sdf.Sphere3D(knp.sphereRadius)
	if err != nil {
		panic(err)
	}

	shell = sdf.Transform3D(shell, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut}))

	top := sdf.Cut3D(shell, sdf.V3{X: 0, Y: 0, Z: 0}, sdf.V3{X: 0, Y: 0, Z: -1})
	top = sdf.Cut3D(top, sdf.V3{X: 0, Y: 0, Z: -knp.plateThickness}, sdf.V3{X: 0, Y: 0, Z: 1})

	hollow, err := sdf.Sphere3D(knp.sphereRadius - knp.sphereThicknes)
	if err != nil {
		panic(err)
	}

	hollow = sdf.Transform3D(hollow, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut}))
	hollow = sdf.Cut3D(hollow, sdf.V3{X: 0, Y: 0, Z: -knp.plateThickness}, sdf.V3{X: 0, Y: 0, Z: -1})

	clearingCylinder, err := sdf.Cylinder3D(knp.sphereRadius*2, knp.sphereRadius-knp.sphereThicknes, 0)
	if err != nil {
		panic(err)
	}

	//topClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: -knp.sphereRadius - knp.backCoverkcut}))
	bottomClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: knp.sphereRadius - knp.plateThickness + knp.backCoverLipCut}))

	switchHole, err := sdf.Box3D(sdf.V3{X: knp.switchHoleWidth, Y: knp.switchHoleLength, Z: knp.plateThickness}, 0)
	if err != nil {
		panic(err)
	}

	switchHole = sdf.Transform3D(switchHole, sdf.Translate3d(sdf.V3{Z: -knp.plateThickness / 2}))
	//todo: add latch reliefs

	switchFlatzone, err := sdf.Box3D(sdf.V3{X: knp.switchFlatzoneWidth, Y: knp.switchFlatzoneLength, Z: knp.keycapHeight + knp.keycapOffset}, 0)
	if err != nil {
		panic(err)
	}

	switchFlatzone = sdf.Transform3D(switchFlatzone, sdf.Translate3d(sdf.V3{Z: (knp.keycapHeight + knp.keycapOffset) / 2}))

	pcbCutAway, err := sdf.Box3D(sdf.V3{X: knp.pcbWidth, Y: knp.pcbLength, Z: knp.plateThickness - knp.switchHoleDepth}, 0)
	if err != nil {
		panic(err)
	}
	pcbCutAway = sdf.Transform3D(pcbCutAway, sdf.Translate3d(sdf.V3{Z: -(knp.plateThickness-knp.switchHoleDepth)/2 - knp.switchHoleDepth}))

	shellBottom := sdf.Cut3D(shell, sdf.V3{X: 0, Y: 0, Z: -knp.plateThickness + knp.backCoverLipCut}, sdf.V3{X: 0, Y: 0, Z: -1})
	top = sdf.Difference3D(top, shellBottom)

	top = sdf.Transform3D(top, orientAndMove)
	hollow = sdf.Transform3D(hollow, orientAndMove)
	switchHole = sdf.Transform3D(switchHole, orientAndMove)
	switchFlatzone = sdf.Transform3D(switchFlatzone, orientAndMove)
	pcbCutAway = sdf.Transform3D(pcbCutAway, orientAndMove)
	shellBottom = sdf.Transform3D(shellBottom, orientAndMove)
	//topClearingCylinder = sdf.Transform3D(topClearingCylinder, orientAndMove)
	bottomClearingCylinder = sdf.Transform3D(bottomClearingCylinder, orientAndMove)

	return &KeyNodule{
			// tops:      []sdf.SDF3{top},
			// topHoles:  []sdf.SDF3{switchHole, switchFlatzone, pcbCutAway},
			// backs:     []sdf.SDF3{shellBottom},
			// backHoles: []sdf.SDF3{hollow, switchHole, switchFlatzone, bottomClearingCylinder},
			//keycapHitbox sdf.SDF3
			//switchHitbox sdf.SDF3
		},
		nil
}
