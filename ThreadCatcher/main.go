package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	//Thread catcher clamps onto the lid of a sewing machine and holds a thread such that it doesn't wear the lid
	//Three independent pieces, Main jaw, floating jaw, nut
	//main jaw has the forward facing jaw, bolt threads, square peg to hold floating jaw from twisting, thread holder and a bar to connect these
	//floating jaw has a similar jaw to main, a square hole to connet to main jaw to keep from twisting, a friction surface to mate with nut
	//nut has inner threads, a friction surface to mate with floating jaw, and a knerled outside

	lidThickness := 2.0
	lidLipThickness := 4.0
	lidLipWidth := 1.0
	jawWidth := 25.0
	barWidth := 20.0
	jawThickness := 5.0
	antiRotationCubeTollerance := 0.35
	threadHoleRadius := 2.5 //thread as in sewing thread (gosh, this is a poor choice of naming that I just _know_ is going to stop being clever)

	threadRadius := 4.0 //thead, as in bold thread
	threadPitch := 3.0
	threadTolerance := 0.20
	threadExpand := 0.20
	//tightThreadTolerance := 0.15

	nutThickness := 3.0
	nutHeight := 10.0
	captureAngleSize := 1.5
	captureAngleTolerance := 0.35
	captureRingShrink := 0.2 //To keep capture ring from touching knurl points

	boltHeight := 30.0

	barLength := 2*barWidth + threadRadius
	extrudeRound := jawThickness / 3

	head := sdf.Union3D(
		sdf.KnurledHead3D(threadRadius+nutThickness, nutHeight, nutHeight/4),
		sdf.Transform3D(
			sdf.Cone3D(captureAngleSize, threadRadius+nutThickness-captureAngleTolerance, threadRadius+nutThickness+captureAngleSize-captureAngleTolerance, 0),
			sdf.Translate3d(sdf.V3{0, 0, nutHeight/2 - captureAngleSize/2}),
		),
	)

	tpi45 := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	screwHole := sdf.Screw3D(
		tpi45,       // 2D thread profile
		nutHeight,   // length of screw
		threadPitch, // thread to thread distance
		1,           // number of thread starts (< 0 for left hand threads)
	)

	tpe45 := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	screwBolt := sdf.Transform3D(
		Screw3DWithTaper(
			tpe45,       // 2D thread profile
			boltHeight,  // length of screw
			threadPitch, // thread to thread distance
			1,           // number of thread starts (< 0 for left hand threads)
			nutHeight,
			-threadExpand,
		),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: boltHeight/2 - jawThickness/2}),
	)

	nut := sdf.Difference3D(head, screwHole)

	mainJawPlan2D := sdf.Union2D(
		sdf.Transform2D( // clamp part
			sdf.Box2D(sdf.V2{barWidth - extrudeRound, jawWidth - extrudeRound}, barWidth/3),
			sdf.Translate2d(sdf.V2{-threadRadius - barWidth/2, 0}),
		),
		sdf.Box2D(sdf.V2{barLength - extrudeRound, barWidth - extrudeRound}, barWidth/3), //horizontle part
	)

	mainJaw := sdf.Difference3D(
		sdf.Union3D(
			sdf.ExtrudeRounded3D( //Main Jaw T shape
				mainJawPlan2D,
				jawThickness,
				extrudeRound,
			),
			sdf.Transform3D( // anti rotation cube
				sdf.Box3D(sdf.V3{threadRadius * 2, threadRadius * 2, jawThickness * 1.5}, 0),
				sdf.Translate3d(sdf.V3{0, 0, jawThickness * 0.75}),
			),
			screwBolt, //screwBolt
		),
		sdf.Union3D(
			sdf.Transform3D( //thread hole
				sdf.Cylinder3D(jawThickness, threadHoleRadius, 0),
				sdf.Translate3d(sdf.V3{barLength/2 - barWidth/3, 0, 0}),
			),
		),
	)
	mainJaw.(*sdf.DifferenceSDF3).SetMax(sdf.PolyMax(0.5)) // soften thread hole edges
	mainJaw = sdf.Difference3D(
		mainJaw,
		sdf.Transform3D( //lid lip slot
			sdf.Box3D(sdf.V3{lidLipWidth, jawWidth, lidLipThickness - lidThickness}, 0),
			sdf.Translate3d(sdf.V3{-threadRadius - lidLipWidth/2, 0, jawThickness/2 - (lidLipThickness-lidThickness)/2}),
		),
	)

	floatingJawPlan2D := sdf.Union2D(
		sdf.Transform2D( // clamp part
			sdf.Box2D(sdf.V2{barWidth - extrudeRound, jawWidth - extrudeRound}, barWidth/3),
			sdf.Translate2d(sdf.V2{-threadRadius - barWidth/2, 0}),
		),
		sdf.Box2D(sdf.V2{barWidth - extrudeRound, barWidth - extrudeRound}, barWidth/3), //horizontle part
	)
	captureRingCrossSection := sdf.Transform2D(
		sdf.Cut2D(
			sdf.Cut2D(
				sdf.Circle2D(captureAngleSize-captureRingShrink),
				sdf.V2{0, 0},
				sdf.V2{-1, 0},
			),
			sdf.V2{0, 0},
			sdf.V2{-1, 1},
		),
		sdf.Translate2d(sdf.V2{threadRadius + nutThickness + captureAngleSize, 0}),
	)
	floatingJaw := sdf.Difference3D(
		sdf.Union3D(
			sdf.ExtrudeRounded3D( //Floating Jaw T shape
				floatingJawPlan2D,
				jawThickness,
				extrudeRound,
			),
			sdf.Transform3D( // capture ring
				sdf.RevolveTheta3D(
					captureRingCrossSection,
					sdf.Tau/2,
				),
				sdf.Translate3d(sdf.V3{0, 0, jawThickness / 2}).Mul(
					sdf.RotateZ(sdf.Tau/4),
				),
			),
		),
		sdf.Box3D(sdf.V3{threadRadius*2 + antiRotationCubeTollerance, threadRadius*2 + antiRotationCubeTollerance, jawThickness}, 0),
	)
	floatingJaw.(*sdf.DifferenceSDF3).SetMax(sdf.PolyMax(0.5)) // soften anti rotation cube hole edges

	_ = nut
	_ = mainJaw
	//sdf.RenderSTL(nut, 200, "nut.stl")
	//sdf.RenderSTL(mainJaw, 200, "mainJaw.stl")
	//sdf.RenderSTL(floatingJaw, 200, "floatJaw.stl")

	sdf.RenderSTLSlow(nut, 400, "nut.stl")
	sdf.RenderSTLSlow(mainJaw, 400, "mainJaw.stl")
	sdf.RenderSTLSlow(floatingJaw, 400, "floatJaw.stl")

	// assembled := sdf.Union3D(
	// 	mainJaw,
	// 	sdf.Transform3D(
	// 		floatingJaw,
	// 		sdf.Translate3d(sdf.V3{0, 0, jawThickness}),
	// 	),
	// 	sdf.Transform3D(
	// 		nut,
	// 		sdf.Translate3d(sdf.V3{0, 0, 1.5*jawThickness + (nutHeight)/2}).Mul(
	// 			sdf.RotateX(sdf.Tau/2),
	// 		),
	// 	),
	// )
	// sdf.RenderSTL(assembled, 200, "threadCatcher.stl")
}
