package main

import (
	//"math"

	//"github.com/deadsy/sdfx/obj"
	// "errors"
	// "fmt"
	"math"
	// "sync"

	// "github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// type ExtrudeFuncer interface{
// 	GetShapeAt(float64) sdf.SDF2
// }
// type ExtrudeFunc struct{
// 	shapeFunc func(float64) sdf.SDF2
// }
// func (e *ExtrudeFunc) GetShapeAt(z float64) sdf.SDF2{
// 	return e.shapeFunc(z)
// }
// type ExtrudeBy2DFunction3D struct {
// 	shapeFunc func(z float64) sdf.SDF2
// 	//probably something here to cache SDF2 per cache
// 	bb     sdf.Box3
// 	height float64
// }
type ExtrudeBy2DFunction3D struct {
	shapeFunc func(z float64) sdf.SDF2
	//probably something here to cache SDF2 per cache
	bb     sdf.Box3
	height float64
}

func ExtrudeBy2DFunction(shapeFunc func(float64) sdf.SDF2, height float64, boundingBoxAdjust v3.Vec) sdf.SDF3 {
	e := ExtrudeBy2DFunction3D{}
	e.shapeFunc = shapeFunc
	e.height = height / 2
	b2 := shapeFunc(0).BoundingBox()
	bb := sdf.Box3{
		Min: v3.Vec{X: b2.Min.X, Y: b2.Min.Y, Z: -e.height}.Sub(boundingBoxAdjust),
		Max: v3.Vec{X: b2.Max.X, Y: b2.Max.Y, Z: e.height}.Add(boundingBoxAdjust),
	}
	e.bb = bb
	return &e
}

func (e *ExtrudeBy2DFunction3D) BoundingBox() sdf.Box3 {
	return e.bb
}

func (e *ExtrudeBy2DFunction3D) Evaluate(p v3.Vec) float64 {
	// sdf for the projected 2d surface
	a := e.shapeFunc(p.Z).Evaluate(v2.Vec{X: p.X, Y: p.Y})
	// sdf for the extrusion region: z = [-height, height]
	b := math.Abs(p.Z) - e.height
	// return the intersection
	return math.Max(a, b)
}

func ExpandExtrude(sdf2 sdf.SDF2, height, expand float64) func(float64) sdf.SDF2 {
	return func(z float64) sdf.SDF2 {
		k := Clamp(1-(z+height/2)/height, 0, 1)
		// mix the 2D SDFs
		a := sdf.Mix(0, height, k)
		return sdf.Offset2D(sdf2, a)
	}
}

// Clamp x between a and b, assume a <= b
func Clamp(x, a, b float64) float64 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}
