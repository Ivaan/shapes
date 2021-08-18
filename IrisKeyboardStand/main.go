package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func main() {
	thickThickness := 50.0
	thinThickness := 10.0
	irisScrewCircleRadius := 5.0
	irisScrewCircleThickness := 5.0
	irisScrewHoleRadius := 2.5
	nutCircleRadius := 8.0
	nutThickness := 8.0
	boardWidth := 149.2438
	treeCenter := sdf.V2{X: 20, Y: 45}
	bendOffset := thickThickness*(boardWidth/(thickThickness-thinThickness)) - boardWidth
	tiltAngle := math.Atan2(thickThickness-thinThickness, boardWidth)
	irisScrewLocations := [...]sdf.V2{{X: 3.9067, Y: 3.8671}, {X: 69.6776, Y: -14.9993}, {X: 125.7031, Y: -29.7570}, {X: 149.2438, Y: 11.1359}, {X: 135.4798, Y: 20.5876}, {X: 123.2997, Y: 38.1379}, {X: 122.8589, Y: 86.1195}, {X: 102.3013, Y: 94.3053}, {X: 52.4209, Y: 94.9352}, {X: 3.7082, Y: 84.1462}}
	nutLocations := [...]sdf.V2{{X: 28, Y: 10}, {X: 25, Y: 74}, {X: 112, Y: 77}, {X: 116, Y: 0}}
	threadRadius := 4.0 //thead, as in bolt thread
	threadPitch := 3.0
	threadTolerance := 0.20
	tallBoltHeight := 50.0
	shortBoltHeight := 30.0
	//bunch of circles at Iris screw locations

	topLayerCircles := make([]sdf.SDF2, len(irisScrewLocations))
	topLayerScrewHoles := make([]sdf.SDF3, len(irisScrewLocations))
	for i, v := range irisScrewLocations {
		topLayerCircles[i] = sdf.Transform2D(
			sdf.Circle2D(irisScrewCircleRadius),
			sdf.Translate2d(v),
		)
		topLayerScrewHoles[i] = sdf.Transform3D(
			sdf.Sphere3D(irisScrewHoleRadius),
			sdf.Translate3d(sdf.V3{X: v.X, Y: v.Y, Z: 0}),
		)
	}

	tpi45 := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	screwHole := sdf.Screw3D(
		tpi45,        // 2D thread profile
		nutThickness, // length of screw
		threadPitch,  // thread to thread distance
		1,            // number of thread starts (< 0 for left hand threads)
	)
	backwardScrewHole := sdf.Transform3D(screwHole, sdf.MirrorYZ())
	bottomLayerCircles := make([]sdf.SDF2, len(nutLocations))
	bottomLayerScrewHoles := make([]sdf.SDF3, len(nutLocations))
	backwardBottomLayerScrewHoles := make([]sdf.SDF3, len(nutLocations))
	for i, v := range nutLocations {
		bottomLayerCircles[i] = sdf.Transform2D(
			sdf.Circle2D(nutCircleRadius),
			sdf.Translate2d(v),
		)
		bottomLayerScrewHoles[i] = sdf.Transform3D(
			screwHole,
			sdf.Translate3d(sdf.V3{X: v.X, Y: v.Y, Z: 0}),
		)
		backwardBottomLayerScrewHoles[i] = sdf.Transform3D(
			backwardScrewHole,
			sdf.Translate3d(sdf.V3{X: v.X, Y: v.Y, Z: 0}),
		)
	}

	stand := sdf.Transform3D(
		sdf.Union3D(
			sdf.Transform3D(
				treeLoft3D(topLayerCircles, treeCenter, thickThickness-irisScrewCircleThickness/2, 0, 0),
				sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -irisScrewCircleThickness / 2}),
			),
			sdf.Transform3D(
				treeLoft3D(bottomLayerCircles, treeCenter, thickThickness-nutThickness/2, 0, 0),
				sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: nutThickness / 2}).Mul(
					sdf.MirrorXY(),
				),
			),
		),
		sdf.RotateX(-sdf.Tau/4).Mul(
			sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: 0}),
		),
	)
	stand = bend3d(stand, bendOffset+boardWidth)
	stand = sdf.Transform3D(
		stand,
		sdf.RotateX(sdf.Tau/4),
	)

	screwPads := sdf.Transform3D(
		sdf.Extrude3D(sdf.Union2D(topLayerCircles...), irisScrewCircleThickness),
		sdf.RotateY(-tiltAngle/2.0).Mul(
			sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: -irisScrewCircleThickness / 2}),
		),
	)

	nuts := sdf.Transform3D(
		sdf.Extrude3D(sdf.Union2D(bottomLayerCircles...), nutThickness),
		sdf.RotateY(tiltAngle/2.0).Mul(
			sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: nutThickness / 2}),
		),
	)

	topLayerScrewHoles3D := sdf.Transform3D(
		sdf.Union3D(topLayerScrewHoles...),
		sdf.RotateY(-tiltAngle/2.0).Mul(
			sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: 0}),
		),
	)

	screwHoles := sdf.Union3D(
		topLayerScrewHoles3D,
		sdf.Transform3D(
			sdf.Union3D(bottomLayerScrewHoles...),
			sdf.RotateY(tiltAngle/2.0).Mul(
				sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: nutThickness / 2}),
			),
		),
	)

	backwardScrewHoles := sdf.Union3D(
		topLayerScrewHoles3D,
		sdf.Transform3D(
			sdf.Union3D(backwardBottomLayerScrewHoles...),
			sdf.RotateY(tiltAngle/2.0).Mul(
				sdf.Translate3d(sdf.V3{X: bendOffset, Y: 0, Z: nutThickness / 2}),
			),
		),
	)

	standAddScewPads := sdf.Union3D(stand, screwPads, nuts)

	leftStand := sdf.Transform3D(
		sdf.Difference3D(
			standAddScewPads,
			screwHoles,
		),
		sdf.RotateY(-tiltAngle/2.0),
	)
	rightStand := sdf.Transform3D(
		sdf.Difference3D(
			standAddScewPads,
			backwardScrewHoles,
		),
		sdf.MirrorYZ().Mul(sdf.RotateY(-tiltAngle/2.0)),
	)

	tpe45 := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	tallScrewBolt := sdf.Screw3D(
		tpe45,          // 2D thread profile
		tallBoltHeight, // length of screw
		threadPitch,    // thread to thread distance
		1,              // number of thread starts (< 0 for left hand threads)
	)

	shortScrewBolt := sdf.Screw3D(
		tpe45,           // 2D thread profile
		shortBoltHeight, // length of screw
		threadPitch,     // thread to thread distance
		1,               // number of thread starts (< 0 for left hand threads)
	)

	head := sdf.KnurledHead3D(nutCircleRadius, nutThickness, nutThickness/4)

	nut := sdf.Difference3D(
		head,
		screwHole,
	)
	_ = leftStand
	_ = rightStand
	_ = tallScrewBolt
	_ = shortScrewBolt
	sdf.RenderSTLSlow(leftStand, 800, "leftStandLowRes.stl")
	sdf.RenderSTLSlow(rightStand, 200, "rightStandLowRes.stl")
	sdf.RenderSTLSlow(tallScrewBolt, 800, "tallScrewBolt.stl")
	sdf.RenderSTLSlow(shortScrewBolt, 800, "shortScrewBolt.stl")
	sdf.RenderSTLSlow(nut, 800, "jamNut.stl")
	//sdf.RenderSTLSlow(stand, 400, "stand.stl")
	//bunch of circles at magnet locations
	//bunch of other circles at locations to make the rest work?

	//topplate use a 2D union function that will blend circles

	//4 circles at adjustment screw locations

	//bottom plate is 2D union blend of these circles

	//connect top and bottom plate through custom extrude function
	//that's a combination of Loft3D and RevolveTheta3D
	//and _maybe_ an animation of 2D union blend function

}

// func lineFromTo(v1, v2 sdf.V2, round float64) sdf.SDF2 {
// 	return sdf.Transform2D(
// 		sdf.Line2D()

// 	)
// }
