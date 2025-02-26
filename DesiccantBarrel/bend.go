package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/deadsy/sdfx/vec/conv"
	p2 "github.com/deadsy/sdfx/vec/p2"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type bendSDF3 struct {
	sdf              sdf.SDF3
	isollinearRadius float64
	bb               sdf.Box3
}

func bend3d(shape sdf.SDF3, isollinearRadius float64) sdf.SDF3 {
	b := bendSDF3{
		sdf:              shape,
		isollinearRadius: isollinearRadius,
	}
	d := 2 * math.Max(math.Abs(shape.BoundingBox().Max.X), math.Abs(shape.BoundingBox().Min.X))
	b.bb = sdf.NewBox3(v3.Vec{X: 0, Y: 0, Z: shape.BoundingBox().Min.Z}, v3.Vec{X: d, Y: d, Z: d}) //this is far too pesimistic, we can take all corners of the cube and back them through the polar math to get the extends (probably?)

	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	// consider the vertices
	vs := shape.BoundingBox().Vertices()
	for _, v := range vs {
		polar := p2.Vec{R: v.X, Theta: math.Mod(v.Y/isollinearRadius, sdf.Pi)}
		mapped := conv.P2ToV2(polar)
		if mapped.X < minX {
			minX = mapped.X
		}
		if mapped.X > maxX {
			maxX = mapped.X
		}
		if mapped.Y < minY {
			minY = mapped.Y
		}
		if mapped.Y > maxY {
			maxY = mapped.Y
		}
	}

	if minY < 0 && maxY > 0 { // if the projected shape straddles the X axis then the convex curve will bulge past min/max
		if maxX > 0 { //we bulge past on the positive side of X, new maxX is the max of the distance x,y is from origin
			maxX = math.Max(v2.Vec{X: maxX, Y: minY}.Length(), v2.Vec{X: maxX, Y: maxY}.Length())
		}
		if minX < 0 { //we bulge past on the negative side of X, same as above except for minX
			minX = math.Min(v2.Vec{X: minX, Y: minY}.Length(), v2.Vec{X: minX, Y: maxY}.Length())
		}
	}
	if minX < 0 && maxX > 0 {
		if maxY > 0 {
			maxY = math.Max(v2.Vec{X: minX, Y: maxY}.Length(), v2.Vec{X: maxX, Y: maxY}.Length())
		}
		if minY < 0 {
			minY = math.Min(v2.Vec{X: minX, Y: minY}.Length(), v2.Vec{X: maxX, Y: minY}.Length())
		}
	}

	b.bb = sdf.Box3{Min: v3.Vec{X: minX, Y: minY, Z: shape.BoundingBox().Min.Z}, Max: v3.Vec{X: maxX, Y: maxY, Z: shape.BoundingBox().Max.Z}}

	// fmt.Printf("Bend BB: %v\n", b.bb)
	return &b
}

func (b *bendSDF3) Evaluate(p v3.Vec) float64 {
	c2d := v2.Vec{X: p.X, Y: p.Y}
	p2d := conv.V2ToP2(c2d)
	return b.sdf.Evaluate(v3.Vec{X: p2d.R, Y: p2d.Theta * b.isollinearRadius, Z: p.Z})
}

// BoundingBox logic in this SDF is terrible atm - needs a total rethink!
func (b *bendSDF3) BoundingBox() sdf.Box3 {
	return b.bb
}

func (b *bendSDF3) SetBoundingBox(bb sdf.Box3) {
	b.bb = bb
}
