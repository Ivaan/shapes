package main

import "github.com/deadsy/sdfx/sdf"

type KeyNodule struct {
	tops         []sdf.SDF3
	topHoles     []sdf.SDF3
	backs        []sdf.SDF3
	backHoles    []sdf.SDF3
	keycapHitbox sdf.SDF3
	switchHitbox sdf.SDF3
}

func (kn KeyNodule) GetTops() []sdf.SDF3 {
	return kn.tops
}

func (kn KeyNodule) GetTopHoles() []sdf.SDF3 {
	return kn.topHoles
}

func (kn KeyNodule) GetBacks() []sdf.SDF3 {
	return kn.backs
}

func (kn KeyNodule) GetBackHoles() []sdf.SDF3 {
	return kn.backHoles
}

func (kn KeyNodule) GetHitBoxes() []sdf.SDF3 {
	return []sdf.SDF3{kn.keycapHitbox, kn.switchHitbox}
}

type KeyNoduleProperties struct {
	sphereRadius             float64
	sphereCut                float64
	plateThickness           float64
	sphereThicknes           float64
	backCoverkcut            float64
	switchHoleLength         float64
	switchHoleWidth          float64
	switchLatchWidth         float64
	switchLatchGrabThickness float64
	switchFlatzoneWidth      float64
	switchFlatzoneHeight     float64
	keycapWidth              float64
	keycapHeight             float64
	keycapRound              float64
	keycapOffset             float64
}

func (knp KeyNoduleProperties) MakeKey(orientAndMove sdf.M44) (*KeyNodule, error) {
	shell, err := sdf.Sphere3D(knp.sphereRadius)
	if err != nil {
		return nil, err
	}

	shell = sdf.Transform3D(shell, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut}))
	shell = sdf.Cut3D(shell, sdf.V3{X: 0, Y: 0, Z: 0}, sdf.V3{X: 0, Y: 0, Z: -1})

	hollow, err := sdf.Sphere3D(knp.sphereRadius - knp.sphereThicknes)
	if err != nil {
		return nil, err
	}

	hollow = sdf.Transform3D(hollow, sdf.Translate3d(sdf.V3{Z: -knp.sphereCut}))
	hollow = sdf.Cut3D(hollow, sdf.V3{X: 0, Y: 0, Z: -knp.plateThickness}, sdf.V3{X: 0, Y: 0, Z: -1})

	clearingCylinder, err := sdf.Cylinder3D(knp.sphereRadius*2, knp.sphereRadius-knp.sphereThicknes, 0)
	if err != nil {
		return nil, err
	}

	topClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: -knp.sphereRadius - knp.backCoverkcut}))
	bottomClearingCylinder := sdf.Transform3D(clearingCylinder, sdf.Translate3d(sdf.V3{Z: knp.sphereRadius - knp.backCoverkcut}))

	switchHole, err := sdf.Box3D(sdf.V3{X: knp.switchHoleWidth, Y: knp.switchHoleLength, Z: knp.plateThickness}, 0)
	if err != nil {
		return nil, err
	}

	switchHole = sdf.Transform3D(switchHole, sdf.Translate3d(sdf.V3{Z: -knp.plateThickness / 2}))
	//todo: add latch reliefs

	switchFlatzone, err := sdf.Box3D(sdf.V3{X: knp.switchFlatzoneWidth, Y: knp.switchFlatzoneHeight, Z: knp.keycapHeight + knp.keycapOffset}, 0)
	if err != nil {
		return nil, err
	}

	switchFlatzone = sdf.Transform3D(switchFlatzone, sdf.Translate3d(sdf.V3{Z: (knp.keycapHeight + knp.keycapOffset) / 2}))

	coverCutA := sdf.V3{Z: -knp.backCoverkcut}
	coverTopV := sdf.V3{Z: 1}
	coverBottomtV := sdf.V3{Z: -1}
	shellTop := sdf.Cut3D(shell, coverCutA, coverTopV)
	shellBottom := sdf.Cut3D(shell, coverCutA, coverBottomtV)

	shellTop = sdf.Transform3D(shellTop, orientAndMove)
	hollow = sdf.Transform3D(hollow, orientAndMove)
	switchHole = sdf.Transform3D(switchHole, orientAndMove)
	switchFlatzone = sdf.Transform3D(switchFlatzone, orientAndMove)
	shellBottom = sdf.Transform3D(shellBottom, orientAndMove)
	topClearingCylinder = sdf.Transform3D(topClearingCylinder, orientAndMove)
	bottomClearingCylinder = sdf.Transform3D(bottomClearingCylinder, orientAndMove)

	return &KeyNodule{
			tops:      []sdf.SDF3{shellTop},
			topHoles:  []sdf.SDF3{hollow, switchHole, switchFlatzone, topClearingCylinder},
			backs:     []sdf.SDF3{shellBottom},
			backHoles: []sdf.SDF3{hollow, switchHole, switchFlatzone, bottomClearingCylinder},
			//keycapHitbox sdf.SDF3
			//switchHitbox sdf.SDF3
		},
		nil
}
