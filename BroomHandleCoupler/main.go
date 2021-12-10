package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	//19, 17, 14
	tubeID := 20.1
	threadMajorID := 18.9
	//threadMinorID := 14.0
	threadLength := 25.0
	threadPitch := 5.08

	threadTaper := 1.1
	straightThreadLength := 8.0

	knerledDialThickness := 4.0
	knerledDialLength := 5.0

	sliceDepth := 18.0
	sliceWidth := 2.0
	sliceCount := 7

	cylinder, err := sdf.Cylinder3D(threadLength, tubeID/2.0, 1)
	if err != nil {
		panic(err)
	}
	cylinder = sdf.Transform3D(cylinder, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: threadLength / 2}))

	dial, err := obj.KnurledHead3D(tubeID/2.0+knerledDialThickness, knerledDialLength, knerledDialLength/2.0)
	if err != nil {
		panic(err)
	}
	dial = sdf.Transform3D(dial, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: knerledDialLength / 2}))

	acmeThread, _ := sdf.AcmeThread(threadMajorID/2.0, threadPitch)
	thread := Screw3DWithTaper(
		acmeThread,                        // 2D thread profile
		threadLength,                      // length of screw
		threadPitch,                       // thread to thread distance
		1,                                 // number of thread starts (< 0 for left hand threads)
		threadLength-straightThreadLength, // the amount of thread that is tapered
		threadTaper,                       // the maximum difference between thread and tapered thread
	)
	thread = sdf.Transform3D(thread, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: threadLength / 2}))

	triangle, err := sdf.Polygon2D(
		[]sdf.V2{
			{X: 0.0, Y: 0.0},
			{X: sliceDepth, Y: -sliceWidth / 2},
			{X: sliceDepth, Y: sliceWidth / 2},
		},
	)
	if err != nil {
		panic(err)
	}

	cut := sdf.Transform3D(
		sdf.Extrude3D(
			triangle,
			tubeID/2,
		),
		sdf.Translate3d(sdf.V3{X: sliceDepth / 2, Y: 0, Z: threadLength - sliceDepth}).Mul(
			sdf.RotateY(-sdf.Tau/4),
		),
	)
	cuts := make([]sdf.SDF3, sliceCount+1) //+1 = ew

	for i := 0; i < sliceCount; i++ {
		cuts[i] = sdf.Transform3D(
			cut,
			sdf.RotateZ(sdf.Tau/float64(sliceCount)*float64(i)),
		)
	}
	cuts = append(cuts, thread) //I'm sorry
	broomCoupler := sdf.Difference3D(
		sdf.Union3D(
			cylinder,
			dial,
		),
		sdf.Union3D(
			cuts...,
		),
	)

	render.RenderSTLSlow(broomCoupler, 300, "BroomCoupler.stl")

}
