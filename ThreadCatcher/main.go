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
	// jawWidth := 20.0
	// barwidth := 12.0
	// jawThickness := 5.0

	threadRadius := 4.0
	threadPitch := 3.0
	threadTolerance := 0.3

	nutThickness := 4.0
	nutHeight := 10.0

	boltHeight := 30.0

	head := sdf.KnurledHead3D(threadRadius+nutThickness, nutHeight, nutHeight/10)

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
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: boltHeight/2 - nutHeight/2}),
	)

	nut := sdf.Difference3D(head, screwHole)
	bolt := sdf.Union3D(head, screwBolt)

	// tp30 := threadProfile(threadRadius, 3, 30, "internal")
	// sdf.RenderSVG(tp30, 200, "tp30.svg", "fill:none;stroke:black;stroke-width:0.1")
	// tp40 := threadProfile(threadRadius, 3, 40, "internal")
	// sdf.RenderSVG(tp40, 200, "tp40.svg", "fill:none;stroke:black;stroke-width:0.1")
	// sdf.RenderSVG(tp45, 200, "tp45.svg", "fill:none;stroke:black;stroke-width:0.1")
	sdf.RenderSTLSlow(nut, 400, "nut.stl")
	sdf.RenderSTLSlow(bolt, 400, "bolt.stl")
}
