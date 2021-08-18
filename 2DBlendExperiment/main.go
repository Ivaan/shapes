package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func main() {
	circleRadius := 5.0
	nutCircleRadius := 8.0
	//circleLocations := [...]sdf.V2{{X: 0, Y: 0}, {X: 5, Y: 0}}
	circleLocations := [...]sdf.V2{{X: 3.9067, Y: 3.8671}, {X: 69.6776, Y: -14.9993}, {X: 125.7031, Y: -29.7570}, {X: 149.2438, Y: 11.1359}, {X: 135.4798, Y: 20.5876}, {X: 123.2997, Y: 38.1379}, {X: 122.8589, Y: 86.1195}, {X: 102.3013, Y: 94.3053}, {X: 52.4209, Y: 94.9352}, {X: 3.7082, Y: 84.1462}}
	nutLocations := [...]sdf.V2{{X: 10, Y: 10}, {X: 40, Y: 74}, {X: 112, Y: 77}, {X: 116, Y: 0}}
	//nutLocations := [...]sdf.V2{{X: 10, Y: 10}, {X: 40, Y: 74}, {X: 112, Y: 77}, {X: 116, Y: 0}}

	//bunch of circles at Iris screw locations

	circles2D := make([]sdf.SDF2, len(circleLocations))
	for i, v := range circleLocations {
		circles2D[i] = sdf.Transform2D(
			sdf.Circle2D(circleRadius),
			sdf.Translate2d(v),
		)
	}

	bottomLayerCircles := make([]sdf.SDF2, len(nutLocations))
	for i, v := range nutLocations {
		bottomLayerCircles[i] = sdf.Transform2D(
			sdf.Circle2D(nutCircleRadius),
			sdf.Translate2d(v),
		)
	}

	// layerCircles2D := sdf.Union2D(circles2D...)
	// layerCircles2D.(*sdf.UnionSDF2).SetMin(stickyMin(5))
	// topLayer := sdf.Extrude3D(
	// 	layerCircles2D,
	// 	2.0,
	// )

	_ = treeLoft3D(circles2D, bottomLayerCircles, 40, 1.0/12*sdf.Tau, 3)
	cube := sdf.Transform3D(
		sdf.Box3D(sdf.V3{X: 4, Y: 4, Z: 4}, .5),
		sdf.Translate3d(sdf.V3{X: 8, Y: 12, Z: 0}),
	)
	cube = bend3d(cube, 8)
	//bentTree := bend3d(tree, 120)

	sdf.RenderSTLSlow(cube, 200, "blend.stl")

}

func stickyMin(k float64) sdf.MinFunc {
	return func(a, b float64) float64 {
		return -k + (math.Log(k-a) + math.Log(k-b))
	}
}

type expandBoundsSDF3 struct {
	sdf sdf.SDF3
	bb  sdf.Box3
}

func expandBounds3D(sdf sdf.SDF3, factor float64) sdf.SDF3 {
	e := expandBoundsSDF3{}
	e.sdf = sdf
	e.bb = sdf.BoundingBox().ScaleAboutCenter(factor)
	return &e
}

func (e *expandBoundsSDF3) Evaluate(p sdf.V3) float64 {
	return e.sdf.Evaluate(p)
}

func (e *expandBoundsSDF3) BoundingBox() sdf.Box3 {
	return e.bb
}
