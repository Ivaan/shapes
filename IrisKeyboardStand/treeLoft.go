package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type treeLoftSDF3 struct {
	sdfs   []sdf.SDF2
	center v2.Vec
	height float64
	round  float64
	min    sdf.MinFunc
	twist  sdf.ExtrudeFunc
	bb     sdf.Box3
}

func treeLoft3D(sdfs []sdf.SDF2, center v2.Vec, height, twist, round float64) sdf.SDF3 {
	t := treeLoftSDF3{
		sdfs:   sdfs,
		center: center,
		height: (height / 2) - round,
		round:  round,
		min:    math.Min,
		twist:  sdf.TwistExtrude(height, twist),
	}
	if t.height < 0 {
		panic("height < 2 * round")
	}
	// work out the bounding box
	bb := t.sdfs[0].BoundingBox()
	for _, x := range t.sdfs {
		bb = bb.Extend(x.BoundingBox()) // for now we presume center is withing the collective BB
	}
	t.bb = sdf.Box3{Min: v3.Vec{X: bb.Min.X, Y: bb.Min.Y, Z: -t.height}.SubScalar(round), Max: v3.Vec{X: bb.Max.X, Y: bb.Max.Y, Z: t.height}.AddScalar(round)}
	return &t
}

func (t *treeLoftSDF3) Evaluate(p v3.Vec) float64 {
	// work out the mix value as a function of height
	k := sdf.Clamp((0.5*p.Z/t.height)+0.5, 0, 1)

	pt := t.twist(p.Sub(v3.Vec{X: t.center.X, Y: t.center.Y, Z: p.Z})).Add(v2.Vec{X: t.center.X, Y: t.center.Y})

	var a float64
	for i, sh := range t.sdfs {
		ps := t.shiftPoint(pt, sh.BoundingBox(), k)
		if i == 0 {
			a = sh.Evaluate(ps)
		} else {
			a = t.min(a, sh.Evaluate(ps))
		}
	}

	//this code borrowed from sdf.Loft
	//it does the right thing for the various z ranges
	b := math.Abs(p.Z) - t.height
	var d float64
	if b > 0 {
		// outside the object Z extent
		if a < 0 {
			// inside the boundary
			d = b
		} else {
			// outside the boundary
			d = math.Sqrt((a * a) + (b * b))
		}
	} else {
		// within the object Z extent
		if a < 0 {
			// inside the boundary
			d = math.Max(a, b)
		} else {
			// outside the boundary
			d = a
		}
	}
	return d - t.round
}

func (t *treeLoftSDF3) BoundingBox() sdf.Box3 {
	return t.bb
}

func (t *treeLoftSDF3) shiftPoint(p v2.Vec, bb sdf.Box2, k float64) v2.Vec {
	d := t.center.Sub(bb.Center())
	return v2.Vec{X: p.X, Y: p.Y}.Sub(d.MulScalar(1 - k))
}
