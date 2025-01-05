package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/v2"
	"github.com/deadsy/sdfx/vec/v3"
)

func main() {
	//Thread catcher clamps onto the lid of a sewing machine and holds a thread such that it doesn't wear the lid
	//Three independent pieces, Main jaw, floating jaw, nut
	//main jaw has the forward facing jaw, bolt threads, square peg to hold floating jaw from twisting, thread holder and a bar to connect these
	//floating jaw has a similar jaw to main, a square hole to connet to main jaw to keep from twisting, a friction surface to mate with nut
	//nut has inner threads, a friction surface to mate with floating jaw, and a knerled outside

	lidThickness := 3.0
	lidLipThickness := 4.3
	lidLipWidth := 2.5
	jawWidth := 25.0
	barWidth := 20.0
	jawThickness := 5.0
	antiRotationCubeTollerance := 1.5

	threadHoleRadius := 2.5 //thread as in sewing thread (gosh, this is a poor choice of naming that I just _know_ is going to stop being clever)
	threadSlotWidth := 2.0

	threadRadius := 4.0 //thead, as in bolt thread
	threadPitch := 3.0
	threadTolerance := 0.20
	threadExpand := 0.20 //last part of the bolt (about the height of the nut) tapers out to decreas tollerance and make the nut less easy to turn
	//tightThreadTolerance := 0.15

	nutThickness := 4.0
	nutHeight := 10.0
	captureAngleSize := 2.0
	captureAngleTolerance := 0.35
	captureRingShrink := 0.2 //To keep capture ring from touching knurl points

	boltHeight := 3*nutHeight + lidThickness + captureAngleSize

	barLength := 2*barWidth + threadRadius
	extrudeRound := jawThickness / 3

	captureAngleInnerRadius := barWidth/2 - extrudeRound - captureAngleSize

	knurlHead, err := obj.KnurledHead3D(threadRadius+nutThickness, nutHeight, nutHeight/4)
	if err != nil {
		panic(err)
	}
	cone, err := sdf.Cone3D(captureAngleSize, captureAngleInnerRadius-captureAngleTolerance, captureAngleInnerRadius+captureAngleSize-captureAngleTolerance, 0)
	if err != nil {
		panic(err)
	}
	head := sdf.Union3D(
		knurlHead,
		sdf.Transform3D(
			cone,
			sdf.Translate3d(v3.Vec{0, 0, nutHeight/2 + captureAngleSize/2}),
		),
	)

	tpi45, err := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	if err != nil {
		panic(err)
	}
	screwHole, err := sdf.Screw3D(
		tpi45,                      // 2D thread profile
		nutHeight+captureAngleSize, // length of screw
		0,                          // thread taper angle (radians)
		threadPitch,                // thread to thread distance
		1,                          // number of thread starts (< 0 for left hand threads)
	)

	if err != nil {
		panic(err)
	}

	tpe45, err := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	if err != nil {
		panic(err)
	}
	screwBolt := sdf.Transform3D(
		Screw3DWithTaper(
			tpe45,       // 2D thread profile
			boltHeight,  // length of screw
			threadPitch, // thread to thread distance
			1,           // number of thread starts (< 0 for left hand threads)
			nutHeight,
			-threadExpand,
		),
		sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: boltHeight/2 - jawThickness/2}),
	)

	nut := sdf.Difference3D(
		head,
		sdf.Transform3D(
			screwHole,
			sdf.Translate3d(v3.Vec{0, 0, captureAngleSize / 2}),
		),
	)

	mainJawPlan2D := sdf.Union2D(
		sdf.Transform2D( // clamp part
			sdf.Box2D(v2.Vec{barWidth - extrudeRound, jawWidth - extrudeRound}, barWidth/3),
			sdf.Translate2d(v2.Vec{-threadRadius - barWidth/2, 0}),
		),
		sdf.Box2D(v2.Vec{barLength - extrudeRound, barWidth - extrudeRound}, barWidth/3), //horizontle part
	)

	mainJaw3D, err := sdf.ExtrudeRounded3D( //Main Jaw T shape
		mainJawPlan2D,
		jawThickness,
		extrudeRound,
	)
	if err != nil {
		panic(err)
	}
	antiRotationCube, err := sdf.Box3D(v3.Vec{threadRadius * 2, threadRadius * 2, jawThickness * 1.5}, 0)
	if err != nil {
		panic(err)
	}
	threadHole, err := sdf.Cylinder3D(jawThickness, threadHoleRadius, 0)
	if err != nil {
		panic(err)
	}
	threadSlot, err := sdf.RevolveTheta3D(
		sdf.Transform2D(
			sdf.Box2D(v2.Vec{threadSlotWidth, jawThickness}, 0),
			sdf.Translate2d(v2.Vec{barWidth, 0}),
		),
		sdf.Tau/8,
	)
	if err != nil {
		panic(err)
	}
	mainJaw := sdf.Difference3D(
		sdf.Union3D(
			mainJaw3D,
			sdf.Transform3D( // anti rotation cube
				antiRotationCube,
				sdf.Translate3d(v3.Vec{0, 0, jawThickness * 0.75}),
			),
			screwBolt, //screwBolt
		),
		sdf.Union3D(
			sdf.Transform3D( //thread hole
				threadHole,
				sdf.Translate3d(v3.Vec{barLength/2 - barWidth/3, 0, 0}),
			),
			sdf.Transform3D( //thread slot
				threadSlot,
				sdf.Translate3d(v3.Vec{barLength/2 - barWidth/3 - barWidth + threadHoleRadius - threadSlotWidth, 0, 0}).Mul(
					sdf.RotateZ(-sdf.Tau/8),
				),
			),
		),
	)

	lidLipSlot, err := sdf.Box3D(v3.Vec{lidLipWidth, jawWidth, lidLipThickness - lidThickness}, 0)
	if err != nil {
		panic(err)
	}
	mainJaw.(*sdf.DifferenceSDF3).SetMax(sdf.PolyMax(0.5)) // soften thread hole edges
	mainJaw = sdf.Difference3D(
		mainJaw,
		sdf.Transform3D( //lid lip slot
			lidLipSlot,
			sdf.Translate3d(v3.Vec{-threadRadius - lidLipWidth/2, 0, jawThickness/2 - (lidLipThickness-lidThickness)/2}),
		),
	)

	floatingJawPlan2D := sdf.Union2D(
		sdf.Transform2D( // clamp part
			sdf.Box2D(v2.Vec{barWidth - extrudeRound, jawWidth - extrudeRound}, barWidth/3),
			sdf.Translate2d(v2.Vec{-threadRadius - barWidth/2, 0}),
		),
		sdf.Box2D(v2.Vec{barWidth - extrudeRound, barWidth - extrudeRound}, barWidth/3), //horizontle part
	)
	captureRingCircle, err := sdf.Circle2D(captureAngleSize - captureRingShrink)
	if err != nil {
		panic(err)
	}
	captureRingCrossSection := sdf.Transform2D(
		sdf.Cut2D(
			sdf.Cut2D(
				captureRingCircle,
				v2.Vec{0, 0},
				v2.Vec{-1, 0},
			),
			v2.Vec{0, 0},
			v2.Vec{-1, 1},
		),
		sdf.Translate2d(v2.Vec{captureAngleInnerRadius + captureAngleSize, 0}),
	)
	floatingJawTShape, err := sdf.ExtrudeRounded3D( //Floating Jaw T shape
		floatingJawPlan2D,
		jawThickness,
		extrudeRound,
	)
	if err != nil {
		panic(err)
	}
	captureRing, err := sdf.RevolveTheta3D(
		captureRingCrossSection,
		sdf.Tau/2,
	)
	if err != nil {
		panic(err)
	}
	antiRotationCubeHole, err := sdf.Box3D(v3.Vec{threadRadius*2 + antiRotationCubeTollerance, threadRadius*2 + antiRotationCubeTollerance, jawThickness}, 0)
	if err != nil {
		panic(err)
	}
	floatingJaw := sdf.Difference3D(
		sdf.Union3D(
			floatingJawTShape,
			sdf.Transform3D( // capture ring
				captureRing,
				sdf.Translate3d(v3.Vec{0, 0, jawThickness / 2}).Mul(
					sdf.RotateZ(sdf.Tau/4),
				),
			),
		),
		antiRotationCubeHole,
	)
	floatingJaw.(*sdf.DifferenceSDF3).SetMax(sdf.PolyMax(0.5)) // soften anti rotation cube hole edges

	_ = nut
	_ = mainJaw
	// sdf.RenderSTL(nut, 200, "nut.stl")
	// sdf.RenderSTL(mainJaw, 200, "mainJaw.stl")
	// sdf.RenderSTL(floatingJaw, 200, "floatJaw.stl")

	render.ToSTL(nut, "nut.stl", render.NewMarchingCubesOctree(400))
	render.ToSTL(mainJaw, "mainJaw.stl", render.NewMarchingCubesOctree(400))
	render.ToSTL(floatingJaw, "floatingJaw.stl", render.NewMarchingCubesOctree(400))

	// assembled := sdf.Union3D(
	// 	mainJaw,
	// 	sdf.Transform3D(
	// 		floatingJaw,
	// 		sdf.Translate3d(v3.Vec{0, 0, jawThickness}),
	// 	),
	// 	sdf.Transform3D(
	// 		nut,
	// 		sdf.Translate3d(v3.Vec{0, 0, 1.5*jawThickness + (nutHeight)/2 + captureAngleSize}).Mul(
	// 			sdf.RotateX(sdf.Tau/2),
	// 		),
	// 	),
	// )
	// sdf.RenderSTL(assembled, 200, "threadCatcher.stl")
}
