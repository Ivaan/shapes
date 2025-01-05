package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"math"
	// v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	cupInnerWidth := 40.0 //describing the cub by its hole because this is the important bit
	cupInnerLength := 70.0
	cupRadius := cupInnerWidth / 2.0
	cupOtherRadius := cupInnerLength - cupRadius
	cupThickness := 2.0
	stripeWidth := 20.0
	stripeThickness := 2.5
	spoutOuterRadius := 10.0 //describing the spout because this needs to fit in the flask/vial opening
	spoutLength := 5.0
	spoutThickness := 1.0
	rotateAngle := 1.0 / 16.0 * sdf.Tau
	rotate := sdf.RotateY(rotateAngle)

	cupHole, err := ReallyWarpedSphere3D(cupRadius, cupOtherRadius)
	if err != nil {
		panic(err)
	}
	cupHole = sdf.Transform3D(cupHole, rotate)
	cup, err := ReallyWarpedSphere3D(cupRadius+cupThickness, cupOtherRadius+cupThickness)
	if err != nil {
		panic(err)
	}
	cup = sdf.Transform3D(cup, rotate)
	// cup = sdf.Cut3D(cup, v3.Vec{}, v3.Vec{X: 1})

	spout, err := sdf.Cylinder3D(spoutLength, spoutOuterRadius, spoutThickness)
	if err != nil {
		panic(err)
	}
	// r^2 = x^2 + y^2
	// x^2 = r^2 - y^2
	// x = sqrt(r^2 - y^2)
	// x = sqrt(cupRadius^2 - (spoutRadius-spoutThickness)^2)*cupOtherRadius/cupkradius
	spoutLocation := (math.Sqrt(cupRadius*cupRadius-(spoutOuterRadius-spoutThickness)*(spoutOuterRadius-spoutThickness)) * cupOtherRadius / cupRadius)
	spout = sdf.Transform3D(
		spout,
		sdf.Translate3d(v3.Vec{Z: -spoutLocation - spoutLength/2}),
	)
	spout = sdf.Transform3D(spout, rotate)
	spoutHoleLength := math.Max(cupOtherRadius+cupThickness-spoutLocation, spoutLength)
	spoutHole, err := sdf.Cylinder3D(spoutHoleLength, spoutOuterRadius-spoutThickness, 0)
	if err != nil {
		panic(err)
	}
	spoutHole = sdf.Transform3D(
		spoutHole,
		sdf.Translate3d(v3.Vec{Z: -spoutLocation - spoutHoleLength/2}),
	)
	spoutHole = sdf.Transform3D(spoutHole, rotate)
	stripe, err := sdf.Cylinder3D(stripeWidth, cupRadius+stripeThickness, stripeThickness)
	stripe = sdf.Transform3D(
		stripe,
		sdf.Translate3d(v3.Vec{Z: -stripeWidth / 2}),
	)
	cup = sdf.Union3D(cup, stripe)
	cup = sdf.Difference3D(cup, cupHole)

	cutWidth := cupInnerWidth + 2*cupThickness + 2*stripeThickness
	cutLength := cupInnerLength + 2*cupThickness
	cutAngle := 4.0 / 32.0 * sdf.Tau
	cutRoundRadius := 25.0

	cutCylinder, err := sdf.Cylinder3D(cutWidth, cutRoundRadius, 0)
	if err != nil {
		panic(err)
	}
	cutBox, err := sdf.Box3D(v3.Vec{X: cutLength, Y: cutWidth, Z: cutWidth}, 0)
	if err != nil {
		panic(err)
	}
	cutBox0 := sdf.Transform3D(cutBox, sdf.Translate3d(v3.Vec{X: cutLength / 2.0, Y: -cutWidth/2.0 + cutRoundRadius}))
	cutBox1 := sdf.Transform3D(cutBox, sdf.Translate3d(v3.Vec{X: -cutLength / 2.0, Y: -cutWidth/2.0 + cutRoundRadius}))
	cutBox1 = sdf.Transform3D(cutBox1, sdf.RotateZ(cutAngle))

	cutOpening := sdf.Union3D(cutCylinder, cutBox0, cutBox1)
	render.ToSTL(cutOpening, "cutOpening.stl", render.NewMarchingCubesUniform(300))
	cutOpening = sdf.Transform3D(
		cutOpening,
		sdf.Translate3d(v3.Vec{X: -cutRoundRadius}).Mul(
			sdf.RotateZ(-1.0/4.0*sdf.Tau).Mul(
				sdf.RotateY(-1.0/4.0*sdf.Tau),
			),
		),
	)
	render.ToSTL(cutOpening, "cutOpeningMoved.stl", render.NewMarchingCubesUniform(300))

	funnel := sdf.Union3D(cup, spout)
	funnel = sdf.Difference3D(funnel, spoutHole)
	funnel = sdf.Difference3D(funnel, cutOpening)
	funnel = sdf.Cut3D(funnel, v3.Vec{X: cupRadius + cupThickness}, v3.Vec{X: -1})

	render.ToSTL(funnel, "funnel.stl", render.NewMarchingCubesUniform(300))
}

type ReallyWarpedSphereSDF struct {
	radius          float64
	negativeZRadius float64
	negZScale       float64
	bb              sdf.Box3
}

func ReallyWarpedSphere3D(radius, negativeZRadius float64) (sdf.SDF3, error) {
	if radius <= 0 {
		return nil, sdf.ErrMsg("radius <= 0")
	}
	s := ReallyWarpedSphereSDF{}
	s.radius = radius
	s.negZScale = radius / negativeZRadius
	min := v3.Vec{X: -radius, Y: -radius, Z: -negativeZRadius}
	max := v3.Vec{X: radius, Y: radius, Z: radius}
	s.bb = sdf.Box3{Min: min, Max: max}
	return &s, nil
}

// Evaluate returns the minimum distance to a sphere.
func (s *ReallyWarpedSphereSDF) Evaluate(p v3.Vec) float64 {
	if p.Z >= 0 {
		return p.Length() - s.radius
	} else {
		// return p.Length() - s.radius
		q := v3.Vec{X: p.X, Y: p.Y, Z: p.Z * s.negZScale}
		return q.Length() - s.radius
	}
}

// BoundingBox returns the bounding box for a sphere.
func (s *ReallyWarpedSphereSDF) BoundingBox() sdf.Box3 {
	return s.bb
}
