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

	// lidThickness := 2.0
	// lidLipThickness := 4.0
	// lidLipInset := 0.0
	// lidLipWidth := 1.0
	jawWidth := 20.0
	barWidth := 12.0
	jawThickness := 5.0

	threadHoleRadius := 2.5 //thread as in sewing thread (gosh, this is a poor choice of naming that I just _know_ is going to stop being clever)

	threadRadius := 5.0 //thead, as in bold thread
	threadPitch := 3.0
	threadTolerance := 0.25
	//tightThreadTolerance := 1.5

	nutThickness := 4.0
	nutHeight := 10.0

	boltHeight := 30.0

	head := sdf.KnurledHead3D(threadRadius+nutThickness, nutHeight, nutHeight/4)

	tpi45 := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	screwHole := sdf.Screw3D(
		tpi45,       // 2D thread profile
		nutHeight,   // length of screw
		threadPitch, // thread to thread distance
		1,           // number of thread starts (< 0 for left hand threads)
	)

	tpe45 := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	screwBolt := sdf.Transform3D(
		sdf.Screw3D(
			tpe45,       // 2D thread profile
			boltHeight,  // length of screw
			threadPitch, // thread to thread distance
			1,           // number of thread starts (< 0 for left hand threads)
		),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: boltHeight/2 - jawThickness/2}),
	)

	//cubeHead := sdf.Box3D(sdf.V3{10, 10, 10}, 2)

	nut := sdf.Difference3D(head, screwHole)

	mainJawPlan2D := sdf.Union2D(
		sdf.Transform2D(
			sdf.Box2D(sdf.V2{barWidth, jawWidth}, barWidth/3),
			sdf.Translate2d(sdf.V2{-threadRadius - barWidth/2, 0}),
		),
		sdf.Box2D(sdf.V2{2*barWidth + threadRadius, barWidth}, barWidth/3),
	)

	mainJaw := sdf.Difference3D(
		sdf.Union3D(
			sdf.ExtrudeRounded3D(
				mainJawPlan2D,
				jawThickness,
				jawThickness/3,
			),
			sdf.Transform3D(
				sdf.Box3D(sdf.V3{threadRadius * 2, threadRadius * 2, jawThickness}, 0),
				sdf.Translate3d(sdf.V3{0, 0, jawThickness}),
			),
			screwBolt,
		),
		sdf.Transform3D(
			sdf.Cylinder3D(jawThickness, threadHoleRadius, 0),
			sdf.Translate3d(sdf.V3{threadRadius + barWidth/3, 0, 0}),
		),
	)
	mainJaw.(*sdf.DifferenceSDF3).SetMax(sdf.PolyMax(0.5))

	_ = nut
	//sdf.RenderSTL(nut, 400, "nut.stl")
	sdf.RenderSTL(mainJaw, 400, "mainJaw.stl")
}
