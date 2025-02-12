package main

import (
	"fmt"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"math"
	"os"
	// v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	tubeID := 22.0
	funnelID := 55.0
	tubeLipID := 25.0
	tubeLength := 120.0
	funnelShorter := 11.0
	thickness := 2.0
	bendRadius := 800.0
	// this is the amount the pipe dips due to the bend C - A where C^2 = A^2 + B^2
	bendDip := math.Sqrt(bendRadius*bendRadius+tubeLength*tubeLength) - bendRadius
	openingAngle := 100.0 / 360.0 * sdf.Tau
	textHeight := 15.0
	textPositionFromEnd := 20.0

	base, err := sdf.Box3D(v3.Vec{X: tubeLength, Y: tubeID/2 + thickness + bendDip, Z: tubeID + thickness*2}, thickness)
	if err != nil {
		panic(err)
	}
	f, err := sdf.LoadFont("Arial.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	text := sdf.NewText("№6")

	text2d, err := sdf.Text2D(f, text, textHeight)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	text3d := sdf.Extrude3D(text2d, thickness*1.25)
	if err != nil {
		fmt.Printf("can't generate text sdf3 %s\n", err)
		os.Exit(1)
	}
	base = sdf.Union3D(
		base,
		sdf.Transform3D(
			text3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: tubeID/2.0 + thickness}),
		),
		sdf.Transform3D(
			text3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: -tubeID/2.0 - thickness}).Mul(sdf.RotateY(1.0/2.0*sdf.Tau)),
		),
	)

	base = sdf.Transform3D(base, sdf.Translate3d(v3.Vec{Y: -tubeID/4.0 - thickness/2.0 + bendDip/2.0}))

	bendIt := func(in sdf.SDF3) sdf.SDF3 {
		out := sdf.Transform3D(in, sdf.RotateX(1.0/4.0*sdf.Tau))
		out = sdf.Transform3D(out, sdf.Translate3d(v3.Vec{X: bendRadius + tubeID/2.0, Y: -tubeLength / 2.0}))
		out = bend3d(out, bendRadius)
		out = sdf.Transform3D(out, sdf.Translate3d(v3.Vec{X: -bendRadius - tubeID/2.0 + bendDip, Y: tubeLength / 2.0}))
		out = sdf.Transform3D(out, sdf.RotateZ(1.0/4.0*sdf.Tau))
		return out
	}
	channel, err := sdf.Cylinder3D(tubeLength, tubeID/2.0+thickness, 0)
	if err != nil {
		panic(err)
	}

	channelOpeningKeepNormal := sdf.RotateZ(-openingAngle / 2.0).MulPosition(v3.Vec{Y: 1.0})
	channelOpening := sdf.Cut3D(channel, v3.Vec{}, channelOpeningKeepNormal)
	channelOpening = sdf.Cut3D(channelOpening, v3.Vec{}, sdf.MirrorXZ().MulPosition(channelOpeningKeepNormal))
	channelOpening = sdf.Intersect3D(channel, channelOpening)
	channelHole, err := sdf.Cylinder3D(tubeLength, tubeID/2.0, 0)
	if err != nil {
		panic(err)
	}

	positionFunnel := sdf.RotateY(math.Atan2(tubeID, tubeLength)).Mul(sdf.Translate3d(v3.Vec{X: tubeID / 2.0, Z: -funnelShorter / 2.0}))
	// positionFunnel := sdf.RotateY(0).Mul(sdf.Translate3d(v3.Vec{X: tubeID/2.0}))
	moveFunnelEnd := sdf.Translate3d(v3.Vec{Z: tubeLength/2.0 - funnelShorter/2.0})

	funnelHole, err := sdf.Cone3D(tubeLength-funnelShorter, 0, funnelID/2.0, 0)
	if err != nil {
		panic(err)
	}
	funnelEndHole, err := sdf.Sphere3D(funnelID / 2.0)
	if err != nil {
		panic(err)
	}
	funnelEndHole = sdf.Cut3D(funnelEndHole, v3.Vec{}, v3.Vec{Z: 1})
	funnelEndHole = sdf.Transform3D(funnelEndHole, moveFunnelEnd)
	funnelHole = sdf.Union3D(funnelHole, funnelEndHole)

	funnel, err := sdf.Cone3D(tubeLength-funnelShorter, thickness, funnelID/2.0+thickness, 0) //using Inner Diameter as radius so it is doubled?
	if err != nil {
		panic(err)
	}
	funnelEnd, err := sdf.Sphere3D(funnelID/2.0 + thickness)
	if err != nil {
		panic(err)
	}
	funnelEnd = sdf.Cut3D(funnelEnd, v3.Vec{}, v3.Vec{Z: 1})
	funnelEnd = sdf.Cut3D(funnelEnd, v3.Vec{}, v3.Vec{X: -1})
	funnel = sdf.Cut3D(funnel, v3.Vec{}, v3.Vec{X: -1})
	funnelEnd = sdf.Transform3D(funnelEnd, moveFunnelEnd)
	funnel = sdf.Union3D(funnel, funnelEnd)

	funnelHole = sdf.Transform3D(funnelHole, positionFunnel)
	funnel = sdf.Transform3D(funnel, positionFunnel)
	funnelHole = bendIt(funnelHole)
	funnel = bendIt(funnel)

	channelHole = sdf.Union3D(channelHole, channelOpening)
	channelHole = sdf.Transform3D(channelHole, sdf.Translate3d(v3.Vec{Z: -thickness})) //This moves the hole so the other end of the channel isn't cut away (leaving a wall)
	channelHole = bendIt(channelHole)
	channel = bendIt(channel)

	holes := sdf.Union3D(channelHole, funnelHole)
	pipe := sdf.Union3D(base, channel, funnel)
	pipe = sdf.Difference3D(pipe, holes)

	tubeLipHole, err := sdf.Cylinder3D(tubeLipID-tubeID, tubeLipID/2.0, 0)
	if err != nil {
		panic(err)
	}
	tubeLip, err := sdf.Cone3D(tubeLipID-tubeID, tubeLipID/2.0+thickness, tubeLipID/2.0+thickness+tubeLipID-tubeID, 0)
	if err != nil {
		panic(err)
	}
	// tubeLip = sdf.Cut3D(tubeLip, v3.Vec{Y: -thickness}, v3.Vec{Y: -1})
	tubeLip = sdf.Cut3D(tubeLip, v3.Vec{}, channelOpeningKeepNormal)
	tubeLip = sdf.Cut3D(tubeLip, v3.Vec{}, sdf.MirrorXZ().MulPosition(channelOpeningKeepNormal))

	tubeLip = sdf.Difference3D(tubeLip, tubeLipHole)
	tubeLip = sdf.Transform3D(
		tubeLip,
		sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 - (tubeLipID-tubeID)/2.0, Y: bendDip}).Mul(
			sdf.RotateY(1.0/4.0*sdf.Tau).Mul(
				sdf.RotateZ(-1.0/4.0*sdf.Tau),
			),
		),
	)

	pipe = sdf.Union3D(pipe, tubeLip)

	render.ToSTL(pipe, "pipe.stl", render.NewMarchingCubesUniform(600))

}

func urinal_main() {
	cupInnerWidth := 50.0 //describing the cub by its hole because this is the important bit
	cupInnerLength := 100.0
	cupRadius := cupInnerWidth / 2.0
	cupOtherRadius := cupInnerLength - cupRadius
	cupThickness := 2.0
	stripeWidth := 26.0
	stripeThickness := 2.5
	spoutOuterRadius := 11.0 //describing the outer spout because this needs to fit in the flask/vial opening
	spoutLength := 5.0
	spoutThickness := 1.0
	rotateAngle := 1.0 / 16.0 * sdf.Tau
	rotate := sdf.RotateY(rotateAngle)
	cutAngle := 1.0 / 16.0 * sdf.Tau
	cutRoundRadius := 30.0

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
	f, err := sdf.LoadFont("Arial.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	text := sdf.NewText("№6")

	text2d, err := sdf.Text2D(f, text, 13.0)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	text3d := sdf.Extrude3D(text2d, cupRadius+stripeThickness)
	if err != nil {
		fmt.Printf("can't generate text sdf3 %s\n", err)
		os.Exit(1)
	}
	text3d = sdf.Union3D(
		sdf.Transform3D(
			text3d,
			sdf.RotateZ(1.0/16.0*sdf.Tau).Mul(
				sdf.RotateX(1.0/4.0*sdf.Tau).Mul(
					sdf.RotateZ(1.0/4.0*sdf.Tau).Mul(
						sdf.Translate3d(v3.Vec{Z: (cupRadius + stripeThickness) / 2.0}),
					),
				),
			),
		),
		sdf.Transform3D(
			text3d,
			sdf.RotateZ(-1.0/16.0*sdf.Tau).Mul(
				sdf.RotateX(-1.0/4.0*sdf.Tau).Mul(
					sdf.RotateZ(1.0/4.0*sdf.Tau).Mul(
						sdf.Translate3d(v3.Vec{Z: (cupRadius + stripeThickness) / 2.0}),
					),
				),
			),
		),
	)
	stripe = sdf.Difference3D(stripe, text3d)
	stripe = sdf.Transform3D(
		stripe,
		sdf.Translate3d(v3.Vec{Z: -stripeWidth / 2}),
	)
	cup = sdf.Union3D(cup, stripe)
	cup = sdf.Difference3D(cup, cupHole)

	cutWidth := cupInnerWidth + 2*cupThickness + 2*stripeThickness
	cutLength := cupInnerLength + 2*cupThickness

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
	// funnel = sdf.Union3D(funnel, text3d)

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
