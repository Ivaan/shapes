package main

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

var _tolerance = 1.0e-10

func main() {
	//playing with the idea of a 2D shape extruded along a path (Probably just a 2D path for starts, my math isn't so strong)
	//path := NewPartialCirclePath2D(10.0, -sdf.Tau/7, sdf.Tau/7)
	path := NewPartialCirclePath2D(10.0, 0*sdf.Tau/24, 15.0, 6*sdf.Tau/24)
	shape, _ := sdf.Circle2D(3.0)
	othershape, _ := sdf.Circle2D(1.0)
	squareshape := sdf.Box2D(sdf.V2{X: 1.0, Y: 1.5}, 0.2)
	shape = sdf.Union2D(
		shape,
		sdf.Transform2D(othershape, sdf.Translate2d(sdf.V2{X: 3, Y: 3})),
		sdf.Transform2D(squareshape, sdf.Translate2d(sdf.V2{X: 3, Y: -2.3})),
		sdf.Transform2D(squareshape, sdf.Translate2d(sdf.V2{X: -3, Y: -2.3})),
	)
	box, _ := sdf.Box3D(sdf.V3{X: 10, Y: 10, Z: 10}, 0)
	box = sdf.Transform3D(box, sdf.Translate3d(sdf.V3{Z: 10}))
	tor := ExtrudeAlongPath3D(path, shape)
	tor = sdf.Union3D(box, tor)
	tor = sdf.Transform3D(tor, sdf.RotateZ(1*sdf.Tau/24))

	render.RenderSTLSlow(tor, 300, "pathTor.stl")

}

//Path2D is an interface that given a point in 3D space returns the closest point along the 2D curve and the normal vector at that point
type Path2D interface {
	Evaluate(p sdf.V2) (c sdf.V2, normal sdf.V2, inside bool)
	BoundingBox() sdf.Box2
}

//CirclePath2D is a Path2D of a circular path with a radius
type CirclePath2D struct {
	radius float64
	bb     sdf.Box2
}

//NewCirclePath2D makes a new CirclePath2D
func NewCirclePath2D(radius float64) Path2D {
	cp := CirclePath2D{}
	cp.radius = radius
	cp.bb = sdf.NewBox2(sdf.V2{}, sdf.V2{X: radius * 2, Y: radius * 2})
	return &cp
}

func (cp *CirclePath2D) Evaluate(p sdf.V2) (c sdf.V2, n sdf.V2, inside bool) {
	return p.Normalize().MulScalar(cp.radius), p.Normalize(), true
}

func (cp *CirclePath2D) BoundingBox() sdf.Box2 {
	return cp.bb
}

//PartialCirclePath2D is a Path2D of a circular path with a radius
type PartialCirclePath2D struct {
	startRadius float64
	startAngle  float64
	startPoint  sdf.V2
	startNormal sdf.V2
	endRadius   float64
	endAngle    float64
	endPoint    sdf.V2
	endNormal   sdf.V2
	bb          sdf.Box2
}

//NewCirclePath2D makes a new CirclePath2D
func NewPartialCirclePath2D(startRadius, startAngle, endRadius, endAngle float64) Path2D {
	cp := PartialCirclePath2D{}
	cp.startRadius = startRadius
	cp.startAngle = startAngle
	cp.startPoint = sdf.Rotate(startAngle).MulPosition(sdf.V2{X: startAngle})
	cp.startNormal = cp.startPoint.Normalize()
	cp.endRadius = endRadius
	cp.endAngle = endAngle
	cp.endPoint = sdf.Rotate(endAngle).MulPosition(sdf.V2{X: endRadius})
	cp.endNormal = cp.endPoint.Normalize()
	cp.bb = sdf.NewBox2(sdf.V2{}, sdf.V2{X: math.Max(startRadius, endRadius) * 2, Y: math.Max(startRadius, endRadius) * 2})
	return &cp
}

func (cp *PartialCirclePath2D) Evaluate(p sdf.V2) (c sdf.V2, n sdf.V2, inside bool) {

	angle := math.Atan2(p.Y, p.X)
	if angle >= cp.endAngle {
		return cp.endPoint, cp.endNormal, false
	} else if angle <= cp.startAngle {
		return cp.startPoint, cp.startNormal, false
	} else {
		deltaRadius := cp.endRadius - cp.startRadius
		deltaAngle := cp.endAngle - cp.startAngle
		pangle := (angle - cp.startAngle) / deltaAngle
		radius := pangle*deltaRadius + cp.startRadius
		return p.Normalize().MulScalar(radius), p.Normalize(), true
	}
}

func (cp *PartialCirclePath2D) BoundingBox() sdf.Box2 {
	return cp.bb
}

//ExtrudeAlongPathSDF3 is an sdf.SDF3 that takes a 2D shape and a Path2D and extrudes the shape along the path
type ExtrudeAlongPathSDF3 struct {
	path  Path2D
	shape sdf.SDF2
	bb    sdf.Box3
}

func ExtrudeAlongPath3D(path Path2D, shape sdf.SDF2) sdf.SDF3 {
	elp := ExtrudeAlongPathSDF3{}
	elp.path = path
	elp.shape = shape

	minZ := path.BoundingBox().Min.X
	maxZ := path.BoundingBox().Max.X
	minX := path.BoundingBox().Min.X - shape.BoundingBox().Max.Y
	maxX := path.BoundingBox().Max.X + shape.BoundingBox().Max.Y
	minY := path.BoundingBox().Min.Y - shape.BoundingBox().Max.Y
	maxY := path.BoundingBox().Max.Y + shape.BoundingBox().Max.Y
	elp.bb = sdf.Box3{Min: sdf.V3{X: minX, Y: minY, Z: minZ}, Max: sdf.V3{X: maxX, Y: maxY, Z: maxZ}}

	return &elp
}

func (elp *ExtrudeAlongPathSDF3) Evaluate(p sdf.V3) float64 {
	p2D := sdf.V2{X: p.X, Y: p.Y}
	nearestP, normalV, inside := elp.path.Evaluate(p2D)

	shiftedP2 := p2D.Sub(nearestP)

	angle := math.Atan2(normalV.Y, normalV.X)
	rotp2 := sdf.Rotate(-angle).MulPosition(shiftedP2)
	q2 := sdf.V2{X: p.Z, Y: rotp2.X}
	distanceAtProjectedSurface := elp.shape.Evaluate(q2)
	if inside {
		return distanceAtProjectedSurface
	} else {
		return math.Max(distanceAtProjectedSurface, math.Abs((rotp2.Y)))
		/*if distanceAtProjectedSurface < 0 {
			return math.Abs(rotp2.Y)
		} else {
			return sdf.V2{X: distanceAtProjectedSurface, Y: rotp2.Y}.Length()
		}*/
	}

}

func (elp *ExtrudeAlongPathSDF3) BoundingBox() sdf.Box3 {
	return elp.bb
}
