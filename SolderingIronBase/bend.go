package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
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
	b.bb = sdf.NewBox3(sdf.V3{X: 0, Y: 0, Z: 0}, sdf.V3{X: d, Y: d, Z: shape.BoundingBox().Max.Z - shape.BoundingBox().Min.Z}) //this is far too pesimistic, we can take all corners of the cube and back them through the polar math to get the extends (probably?)

	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	// consider the vertices
	vs := shape.BoundingBox().Vertices()
	for _, v := range vs {
		polar := sdf.P2{R: v.X, Theta: math.Mod(v.Y/isollinearRadius, sdf.Pi)}
		mapped := polar.PolarToCartesian()
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
			maxX = math.Max(sdf.V2{X: maxX, Y: minY}.Length(), sdf.V2{X: maxX, Y: maxY}.Length())
		}
		if minX < 0 { //we bulge past on the negative side of X, same as above except for minX
			minX = math.Min(sdf.V2{X: minX, Y: minY}.Length(), sdf.V2{X: minX, Y: maxY}.Length())
		}
	}
	if minX < 0 && maxX > 0 {
		if maxY > 0 {
			maxY = math.Max(sdf.V2{X: minX, Y: maxY}.Length(), sdf.V2{X: maxX, Y: maxY}.Length())
		}
		if minY < 0 {
			minY = math.Min(sdf.V2{X: minX, Y: minY}.Length(), sdf.V2{X: maxX, Y: minY}.Length())
		}
	}

	b.bb = sdf.Box3{Min: sdf.V3{X: minX, Y: minY, Z: shape.BoundingBox().Min.Z}, Max: sdf.V3{X: maxX, Y: maxY, Z: shape.BoundingBox().Max.Z}}

	return &b
}

func (b *bendSDF3) Evaluate(p sdf.V3) float64 {
	c2d := sdf.V2{X: p.X, Y: p.Y}
	p2d := c2d.CartesianToPolar()
	return b.sdf.Evaluate(sdf.V3{X: p2d.R, Y: p2d.Theta * b.isollinearRadius, Z: p.Z})
}

func (b *bendSDF3) BoundingBox() sdf.Box3 {
	return b.bb
}
