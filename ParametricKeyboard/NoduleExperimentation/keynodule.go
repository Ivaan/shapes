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
	sphereRadius             float64
	sphereCut                float64
	plateThickness           float64
	sphereThicknes           float64
	backCoverkcut            float64
	switchHoleLength         float64
	switchHoleWidth          float64
	switchLatchWidth         float64
	switchLatchGrabThickness float64
	switchFlatzoneLength     float64
	switchFlatzoneWidth      float64
	keycapLength             float64
	keycapWidth              float64
	keycapMinHeight          float64
	keycapMaxHeight          float64
	keycapBottomRestHeight   float64
	keycapClearanced         float64
	keycapRound              float64
	huggingCylinderRound     float64
	laneWidth                float64 //as in "Stay in your lane" this restricts the holes to a max width
	insertLength             float64
	insertDiameter           float64
	insertWallThickness      float64
	screwThreadDiameter      float64
	screwThreadLength        float64
	screwHeadDiameter        float64
}

func (knp BubbleKeyNoduleProperties) MakeBubbleKey(orientAndMove sdf.M44) KeyNodule {
	shell, err := sdf.Sphere3D(knp.sphereRadius)
	if err != nil {
		panic(err)
	}

	shell = sdf.Transform3D(shell, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut - knp.keycapBottomRestHeight}))

	radiusAtCut := math.Sqrt(knp.sphereRadius*knp.sphereRadius - knp.sphereCut*knp.sphereCut)
	huggingCylinder, err := sdf.Cylinder3D((knp.keycapMaxHeight+knp.keycapMinHeight)/2+knp.huggingCylinderRound*2, radiusAtCut, knp.huggingCylinderRound*2)
	if err != nil {
		panic(err)
	}

	huggingCylinder = sdf.Transform3D(huggingCylinder, sdf.Translate3d(sdf.V3{Z: ((knp.keycapMaxHeight+knp.keycapMinHeight)/2+knp.huggingCylinderRound*2)/2 - knp.huggingCylinderRound*2 - knp.keycapBottomRestHeight}))

	hollow, err := sdf.Sphere3D(knp.sphereRadius - knp.sphereThicknes)
	if err != nil {
		panic(err)
	}

	hollow = sdf.Transform3D(hollow, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut - knp.keycapBottomRestHeight}))
	hollow = sdf.Cut3D(hollow, sdf.V3{X: 0, Y: 0, Z: -knp.plateThickness - knp.keycapBottomRestHeight}, sdf.V3{X: 0, Y: 0, Z: -1})

	clearingCylinder, err := sdf.Cylinder3D(knp.sphereRadius*2, knp.sphereRadius-knp.sphereThicknes, 0)
	if err != nil {
		panic(err)
	}

	topClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: -knp.sphereRadius - knp.backCoverkcut - knp.keycapBottomRestHeight}))
	bottomClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: knp.sphereRadius - knp.backCoverkcut - knp.keycapBottomRestHeight}))

	switchHole, err := sdf.Box3D(sdf.V3{X: knp.switchHoleWidth, Y: knp.switchHoleLength, Z: knp.plateThickness}, 0)
	if err != nil {
		panic(err)
	}

	switchHole = sdf.Transform3D(switchHole, sdf.Translate3d(sdf.V3{Z: -knp.plateThickness/2 - knp.keycapBottomRestHeight}))
	//todo: add latch reliefs

	switchFlatzone, err := sdf.Box3D(sdf.V3{X: knp.switchFlatzoneWidth, Y: knp.switchFlatzoneLength, Z: knp.keycapMinHeight}, 0)
	if err != nil {
		panic(err)
	}

	switchFlatzone = sdf.Transform3D(switchFlatzone, sdf.Translate3d(sdf.V3{Z: knp.keycapMinHeight/2 - knp.keycapBottomRestHeight}))

	keyCapClearanceShadow := sdf.Box2D(sdf.V2{X: knp.keycapWidth + knp.keycapClearanced, Y: knp.keycapLength + knp.keycapClearanced}, knp.keycapRound+knp.keycapClearanced)
	keyCapClearance, err := sdf.ExtrudeRounded3D(keyCapClearanceShadow, knp.keycapMaxHeight*2, 0)
	if err != nil {
		panic(err)
	}
	keyCapClearance = sdf.Transform3D(keyCapClearance, sdf.Translate3d(sdf.V3{Z: knp.keycapMaxHeight + knp.keycapMinHeight - knp.keycapBottomRestHeight}))

	lane, err := sdf.Box3D(sdf.V3{X: knp.laneWidth, Y: knp.sphereRadius * 2, Z: knp.sphereRadius * 2}, 0)
	if err != nil {
		panic(err)
	}
	lane = sdf.Transform3D(lane, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut - knp.keycapBottomRestHeight}))

	coverCutA := sdf.V3{Z: -knp.backCoverkcut - knp.keycapBottomRestHeight}
	plateCut := sdf.V3{Z: -knp.plateThickness - knp.keycapBottomRestHeight}
	coverTopV := sdf.V3{Z: 1}
	coverBottomtV := sdf.V3{Z: -1}
	shellTop := sdf.Cut3D(shell, coverCutA, coverTopV)
	shellBottom := sdf.Cut3D(shell, coverCutA, coverBottomtV)
	plate := sdf.Cut3D(shell, plateCut, coverTopV)

	numberOfScrews := 4
	insertHolders := make([]sdf.SDF3, numberOfScrews)
	screwChannels := make([]sdf.SDF3, numberOfScrews)
	screwHoles := make([]sdf.SDF3, numberOfScrews)
	insertHoldersHoles := make([]sdf.SDF3, numberOfScrews)

	for i := 0; i < numberOfScrews; i++ {

		holder, err := sdf.Cylinder3D(knp.insertLength+knp.insertWallThickness, knp.insertDiameter/2+knp.insertWallThickness, 0)
		if err != nil {
			panic(err)
		}
		holderHole, err := sdf.Cylinder3D(knp.insertLength, knp.insertDiameter/2, 0)
		if err != nil {
			panic(err)
		}
		holder = sdf.Transform3D(holder, sdf.RotateZ(float64(i)*sdf.Tau/float64(numberOfScrews)).Mul(sdf.Translate3d(sdf.V3{X: radiusAtCut - (knp.insertDiameter/2 + knp.insertWallThickness), Z: (knp.insertLength+knp.insertWallThickness)/2 - knp.sphereCut - knp.keycapBottomRestHeight})))
		holderHole = sdf.Transform3D(holderHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(numberOfScrews)).Mul(sdf.Translate3d(sdf.V3{X: radiusAtCut - (knp.insertDiameter/2 + knp.insertWallThickness), Z: knp.insertLength/2 - knp.sphereCut - knp.keycapBottomRestHeight})))
		insertHolders[i] = holder
		insertHoldersHoles[i] = holderHole

		channel, err := sdf.Cylinder3D(knp.sphereRadius, knp.screwHeadDiameter/2+knp.insertWallThickness, 0)
		if err != nil {
			panic(err)
		}
		screwThreadHole, err := sdf.Cylinder3D(knp.screwThreadLength-knp.insertLength, knp.screwThreadDiameter/2, 0)
		if err != nil {
			panic(err)
		}
		screwHeadHole, err := sdf.Cylinder3D(radiusAtCut, knp.screwHeadDiameter/2, 0)
		if err != nil {
			panic(err)
		}

		channel = sdf.Transform3D(channel, sdf.RotateZ(float64(i)*sdf.Tau/float64(numberOfScrews)).Mul(sdf.Translate3d(sdf.V3{X: radiusAtCut - (knp.insertDiameter/2 + knp.insertWallThickness), Z: -knp.sphereRadius / 2})))
		channel = sdf.Intersect3D(shellBottom, channel)
		screwThreadHole = sdf.Transform3D(screwThreadHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(numberOfScrews)).Mul(sdf.Translate3d(sdf.V3{X: radiusAtCut - (knp.insertDiameter/2 + knp.insertWallThickness), Z: -(knp.screwThreadLength-knp.insertLength)/2 - knp.sphereCut - knp.keycapBottomRestHeight})))
		screwHeadHole = sdf.Transform3D(screwHeadHole, sdf.RotateZ(float64(i)*sdf.Tau/float64(numberOfScrews)).Mul(sdf.Translate3d(sdf.V3{X: radiusAtCut - (knp.insertDiameter/2 + knp.insertWallThickness), Z: -knp.sphereRadius/2 - (knp.screwThreadLength - knp.insertLength) - knp.sphereCut - knp.keycapBottomRestHeight})))
		hole := sdf.Union3D(screwThreadHole, screwHeadHole)
		screwChannels[i] = channel
		screwHoles[i] = hole
	}

	lane = sdf.Transform3D(lane, orientAndMove)
	shellTop = sdf.Transform3D(shellTop, orientAndMove)
	plate = sdf.Transform3D(plate, orientAndMove)
	huggingCylinder = sdf.Transform3D(huggingCylinder, orientAndMove)
	hollow = sdf.Transform3D(hollow, orientAndMove)
	switchHole = sdf.Transform3D(switchHole, orientAndMove)
	switchFlatzone = sdf.Transform3D(switchFlatzone, orientAndMove)
	keyCapClearance = sdf.Transform3D(keyCapClearance, orientAndMove)
	shellBottom = sdf.Transform3D(shellBottom, orientAndMove)
	topClearingCylinder = sdf.Transform3D(topClearingCylinder, orientAndMove)
	bottomClearingCylinder = sdf.Transform3D(bottomClearingCylinder, orientAndMove)
	allInsertHolders := sdf.Transform3D(sdf.Union3D(insertHolders...), orientAndMove)
	allInsertHoldersHoles := sdf.Transform3D(sdf.Union3D(insertHoldersHoles...), orientAndMove)
	allScrewChannels := sdf.Transform3D(sdf.Union3D(screwChannels...), orientAndMove)
	allScrewHoles := sdf.Transform3D(sdf.Union3D(screwHoles...), orientAndMove)

	return KeyNodule{
		Top: MakeNodule(
			[]sdf.SDF3{switchHole, switchFlatzone, keyCapClearance, sdf.Intersect3D(hollow, lane), sdf.Intersect3D(topClearingCylinder, lane), allInsertHoldersHoles}, //hole rank 0
			[]sdf.SDF3{sdf.Intersect3D(plate, lane), huggingCylinder, allScrewChannels, allInsertHolders},                                                             //solid rank 0
			[]sdf.SDF3{hollow, shellBottom}, //hole rank 1
			[]sdf.SDF3{shellTop},            //solid rank 1
		),
		Bottom: MakeNodule(
			[]sdf.SDF3{switchHole, switchFlatzone, bottomClearingCylinder, allScrewHoles}, //hole rank 0
			[]sdf.SDF3{allScrewChannels}, //solid rank 0
			[]sdf.SDF3{hollow},           //hole rank 0
			[]sdf.SDF3{shellBottom},      //solid rank 0
		),
		//keycapHitbox sdf.SDF3
		//switchHitbox sdf.SDF3
	}
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
