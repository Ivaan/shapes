package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type waveyDistortSDF3 struct {
	sdf       sdf.SDF3
	frequency v2.Vec
	amplitude v2.Vec
	bb        sdf.Box3
}

func waveyDistort3d(shape sdf.SDF3, frequency, amplitude v2.Vec) sdf.SDF3 {
	wd := waveyDistortSDF3{
		sdf:       shape,
		frequency: frequency,
		amplitude: amplitude,
		bb:        shape.BoundingBox().Enlarge(v3.Vec{Z: amplitude.X*2.0 + amplitude.Y*2.0}),
	}

	return &wd
}

func (wd *waveyDistortSDF3) Evaluate(p v3.Vec) float64 {
	wp := v3.Vec{X: p.X, Y: p.Y, Z: p.Z + math.Sin(p.X*wd.frequency.X)*wd.amplitude.X + math.Sin(p.Y*wd.frequency.Y)*wd.amplitude.Y}
	return wd.sdf.Evaluate(wp)
}

func (wd *waveyDistortSDF3) BoundingBox() sdf.Box3 {
	return wd.bb
}
