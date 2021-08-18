package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

type treeLoftSDF3 struct {
	topSdfs    []sdf.SDF2
	bottomSdfs []sdf.SDF2
	height     float64
	round      float64
	min        sdf.MinFunc
	twist      sdf.ExtrudeFunc
	bb         sdf.Box3
}

func treeLoft3D(topSdfs, bottomSdfs []sdf.SDF2, height, twist, round float64) sdf.SDF3 {
	t := treeLoftSDF3{
		topSdfs:    topSdfs,
		bottomSdfs: bottomSdfs,
		height:     (height / 2) - round,
		round:      round,
		min:        sdf.Min,
		twist:      sdf.TwistExtrude(height, twist),
	}
	if t.height < 0 {
		panic("height < 2 * round")
	}
	// work out the bounding box
	bb := t.topSdfs[0].BoundingBox()
	for _, x := range t.topSdfs {
		bb = bb.Extend(x.BoundingBox())
	}
	for _, x := range t.bottomSdfs {
		bb = bb.Extend(x.BoundingBox())
	}
	t.bb = sdf.Box3{Min: sdf.V3{X: bb.Min.X, Y: bb.Min.Y, Z: -t.height}.SubScalar(round), Max: sdf.V3{X: bb.Max.X, Y: bb.Max.Y, Z: t.height}.AddScalar(round)}
	return &t
}

func (t *treeLoftSDF3) Evaluate(p sdf.V3) float64 {
	// work out the mix value as a function of height
	k := sdf.Clamp((0.5*p.Z/t.height)+0.5, 0, 1)
	// mix the 2D SDFs
	pt := t.twist(p.Sub(t.bb.Center())).Add(sdf.V2{X: t.bb.Center().X, Y: t.bb.Center().Y})

	var a0 float64
	for i, sh := range t.topSdfs {
		ps := t.shiftPoint(pt, sh.BoundingBox(), k)
		if i == 0 {
			a0 = sh.Evaluate(ps)
		} else {
			a0 = t.min(a0, sh.Evaluate(ps))
		}
	}

	var a1 float64
	for i, sh := range t.bottomSdfs {
		ps := t.shiftPoint(pt, sh.BoundingBox(), 1-k)
		if i == 0 {
			a1 = sh.Evaluate(ps)
		} else {
			a1 = t.min(a1, sh.Evaluate(ps))
		}
	}

	a := beefyMix(a0, a1, k)

	b := sdf.Abs(p.Z) - t.height
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
			d = sdf.Max(a, b)
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
	c3 := t.bb.Center()
	d := sdf.V2{X: c3.X, Y: c3.Y}.Sub(bb.Center())
	return sdf.V2{X: p.X, Y: p.Y}.Sub(d.MulScalar(k))
}

func beefyMix(x, y, a float64) float64 {
	//min := sdf.Min(x, y)
	//return x + (a * (y - x))
	return sdf.Min(x, y)
}
