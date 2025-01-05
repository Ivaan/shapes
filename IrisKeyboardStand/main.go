package main

import (
	"math"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	thickThickness := 50.0
	thinThickness := 10.0
	irisScrewCircleRadius := 5.0
	irisScrewCircleThickness := 7.0
	irisScrewCircleProtrusionAbovePlane := 2.2
	irisScrewHoleRadius := 3.5
	nutCircleRadius := 8.0
	nutThickness := 8.0
	boardWidth := 149.2438
	treeCenter := v2.Vec{X: 20, Y: 45}
	bendOffset := thickThickness*(boardWidth/(thickThickness-thinThickness)) - boardWidth
	tiltAngle := math.Atan2(thickThickness-thinThickness, boardWidth)
	//original Iris screw locations
	//irisScrewLocations := [...]v2.Vec{{X: 3.9067, Y: 3.8671}, {X: 69.6776, Y: -14.9993}, {X: 125.7031, Y: -29.7570}, {X: 149.2438, Y: 11.1359}, {X: 135.4798, Y: 20.5876}, {X: 123.2997, Y: 38.1379}, {X: 122.8589, Y: 86.1195}, {X: 102.3013, Y: 94.3053}, {X: 52.4209, Y: 94.9352}, {X: 3.7082, Y: 84.1462}}
	//Some experimentation with a not an Iris keyboard
	//irisScrewLocations := [...]sdf.V2{{X: 163.7082, Y: -40.2331}, {X: 190.2947, Y: 16.7791}, {X: 80.4060, Y: -3.4064}, {X: 128.5439, Y: 103.9666}, {X: 148.5527, Y: 124.3147}, {X: 95.0694, Y: 126.1336}, {X: 18.5439, Y: 103.9648}, {X: 0.0000, Y: 124.3148}, {X: 0.0000, Y: 0.0000}, {X: 0.0000, Y: -40.2331}}
	//Waterfowl screw locations
	irisScrewLocations := [...]v2.Vec{{X: 34.604, Y: -0.586}, {X: 0, Y: 0}, {X: -0.185, Y: 71.38}, {X: 92.629, Y: 82.255}, {X: 94.969, Y: 13.2}, {X: 122.991, Y: -15.996}, {X: 109.661, Y: -27.669}, {X: 36.564, Y: 40.152}, {X: 57.802, Y: 57.413}, {X: 57.385, Y: -1.184}, {X: 35.838, Y: 79.942}}
	//nutLocations := [...]v2.Vec{{X: 28, Y: 10}, {X: 25, Y: 74}, {X: 112, Y: 77}, {X: 116, Y: 0}}
	//Waterfowl nut locations
	nutLocations := [...]v2.Vec{{X: 18, Y: 10}, {X: 15, Y: 74}, {X: 112, Y: 77}, {X: 116, Y: 0}}
	threadRadius := 4.0 //thead, as in bolt thread
	threadPitch := 3.0
	threadTolerance := 0.20
	tallBoltHeight := 50.0
	shortBoltHeight := 30.0
	//bunch of circles at Iris screw locations

	screwpadCircle, err := sdf.Circle2D(irisScrewCircleRadius)
	if err != nil {
		panic(err)
	}
	screwHoleSphere, err := sdf.Sphere3D(irisScrewHoleRadius)
	if err != nil {
		panic(err)
	}
	topLayerCircles := make([]sdf.SDF2, len(irisScrewLocations))
	topLayerScrewHoles := make([]sdf.SDF3, len(irisScrewLocations))
	for i, v := range irisScrewLocations {
		topLayerCircles[i] = sdf.Transform2D(
			screwpadCircle,
			sdf.Translate2d(v),
		)
		topLayerScrewHoles[i] = sdf.Transform3D(
			screwHoleSphere,
			sdf.Translate3d(v3.Vec{X: v.X, Y: v.Y, Z: 0}),
		)
	}

	tpi45, err := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	if err != nil {
		panic(err)
	}
	screwHole, err := sdf.Screw3D(
		tpi45,        // 2D thread profile
		nutThickness, // length of screw
		0,            // thread taper angle
		threadPitch,  // thread to thread distance
		1,            // number of thread starts (< 0 for left hand threads)
	)
	if err != nil {
		panic(err)
	}

	nutCircle, err := sdf.Circle2D(nutCircleRadius)
	if err != nil {
		panic(err)
	}
	backwardScrewHole := sdf.Transform3D(screwHole, sdf.MirrorYZ())
	bottomLayerCircles := make([]sdf.SDF2, len(nutLocations))
	bottomLayerScrewHoles := make([]sdf.SDF3, len(nutLocations))
	backwardBottomLayerScrewHoles := make([]sdf.SDF3, len(nutLocations))
	for i, v := range nutLocations {
		bottomLayerCircles[i] = sdf.Transform2D(
			nutCircle,
			sdf.Translate2d(v),
		)
		bottomLayerScrewHoles[i] = sdf.Transform3D(
			screwHole,
			sdf.Translate3d(v3.Vec{X: v.X, Y: v.Y, Z: 0}),
		)
		backwardBottomLayerScrewHoles[i] = sdf.Transform3D(
			backwardScrewHole,
			sdf.Translate3d(v3.Vec{X: v.X, Y: v.Y, Z: 0}),
		)
	}

	stand := sdf.Transform3D(
		sdf.Union3D(
			sdf.Transform3D(
				treeLoft3D(topLayerCircles, treeCenter, thickThickness-irisScrewCircleThickness/2, 0, 0),
				sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -irisScrewCircleThickness/2 + irisScrewCircleProtrusionAbovePlane}),
			),
			sdf.Transform3D(
				treeLoft3D(bottomLayerCircles, treeCenter, thickThickness-nutThickness/2, 0, 0),
				sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: nutThickness / 2}).Mul(
					sdf.MirrorXY(),
				),
			),
		),
		sdf.RotateX(-sdf.Tau/4).Mul(
			sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: 0}),
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
			sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: -irisScrewCircleThickness/2 + irisScrewCircleProtrusionAbovePlane}),
		),
	)

	nuts := sdf.Transform3D(
		sdf.Extrude3D(sdf.Union2D(bottomLayerCircles...), nutThickness),
		sdf.RotateY(tiltAngle/2.0).Mul(
			sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: nutThickness / 2}),
		),
	)

	topLayerScrewHoles3D := sdf.Transform3D(
		sdf.Union3D(topLayerScrewHoles...),
		sdf.RotateY(-tiltAngle/2.0).Mul(
			sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: 0}),
		),
	)

	screwHoles := sdf.Union3D(
		topLayerScrewHoles3D,
		sdf.Transform3D(
			sdf.Union3D(bottomLayerScrewHoles...),
			sdf.RotateY(tiltAngle/2.0).Mul(
				sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: nutThickness / 2}),
			),
		),
	)

	backwardScrewHoles := sdf.Union3D(
		topLayerScrewHoles3D,
		sdf.Transform3D(
			sdf.Union3D(backwardBottomLayerScrewHoles...),
			sdf.RotateY(tiltAngle/2.0).Mul(
				sdf.Translate3d(v3.Vec{X: bendOffset, Y: 0, Z: nutThickness / 2}),
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

	tpe45, err := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	if err != nil {
		panic(err)
	}
	tallScrewBolt, err := sdf.Screw3D(
		tpe45,          // 2D thread profile
		tallBoltHeight, // length of screw
		0,              // thread taper angle
		threadPitch,    // thread to thread distance
		1,              // number of thread starts (< 0 for left hand threads)
	)
	if err != nil {
		panic(err)
	}

	shortScrewBolt, err := sdf.Screw3D(
		tpe45,           // 2D thread profile
		shortBoltHeight, // length of screw
		0,               // thread taper angle
		threadPitch,     // thread to thread distance
		1,               // number of thread starts (< 0 for left hand threads)
	)
	if err != nil {
		panic(err)
	}

	head, err := obj.KnurledHead3D(nutCircleRadius, nutThickness, nutThickness/4)
	if err != nil {
		panic(err)
	}

	nut := sdf.Difference3D(
		head,
		screwHole,
	)

	_ = rightStand
	_ = tallScrewBolt
	_ = shortScrewBolt
	_ = nut

	//render.RenderSTLSlow(leftStand, 800, "ergodoneleftStand.stl")
	//render.ToSTL(leftStand, "WaterfowlLeftStandLow.stl", render.NewMarchingCubesUniform(200))
	render.ToSTL(rightStand, "WaterfowlRightStand.stl", render.NewMarchingCubesUniform(800))
	render.ToSTL(leftStand, "WaterfowlLeftStand.stl", render.NewMarchingCubesUniform(800))
	render.ToSTL(tallScrewBolt, "TallScrewBolt.stl", render.NewMarchingCubesUniform(300))
	render.ToSTL(shortScrewBolt, "ShortScrewBolt.stl", render.NewMarchingCubesUniform(300))
	render.ToSTL(nut, "JamNut.stl", render.NewMarchingCubesUniform(200))
	//render.RenderSTLSlow(rightStand, 200, "ergodonerightStandLow.stl")
	//render.RenderSTLSlow(tallScrewBolt, 300, "tallScrewBoltSlow.stl")
	//render.RenderSTLSlow(shortScrewBolt, 200, "shortScrewBoltSlow.stl")
	//render.RenderSTLSlow(nut, 150, "jamNutSlow.stl")

}
