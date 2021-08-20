package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

type treeLoftSDF3 struct {
	sdfs   []sdf.SDF2
	center sdf.V2
	height float64
	round  float64
	min    sdf.MinFunc
	twist  sdf.ExtrudeFunc
	bb     sdf.Box3
}

func treeLoft3D(sdfs []sdf.SDF2, center sdf.V2, height, twist, round float64) sdf.SDF3 {
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
	t.bb = sdf.Box3{Min: sdf.V3{X: bb.Min.X, Y: bb.Min.Y, Z: -t.height}.SubScalar(round), Max: sdf.V3{X: bb.Max.X, Y: bb.Max.Y, Z: t.height}.AddScalar(round)}
	return &t
}

func (t *treeLoftSDF3) Evaluate(p sdf.V3) float64 {
	// work out the mix value as a function of height
	k := sdf.Clamp((0.5*p.Z/t.height)+0.5, 0, 1)

	pt := t.twist(p.Sub(sdf.V3{X: t.center.X, Y: t.center.Y, Z: p.Z})).Add(sdf.V2{X: t.center.X, Y: t.center.Y})

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

func (t *treeLoftSDF3) shiftPoint(p sdf.V2, bb sdf.Box2, k float64) sdf.V2 {
	d := t.center.Sub(bb.Center())
	return sdf.V2{X: p.X, Y: p.Y}.Sub(d.MulScalar(1 - k))
}
