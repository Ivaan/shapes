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
	moveFunnelEnd := sdf.Translate3d(v3.Vec{Z: tubeLength/2.0 - funnelShorter/2.0})

	makeBendFunnel := func(length, r1, r2 float64) (sdf.SDF3, error) {
		funnel, err := sdf.Cone3D(length, r1, r2, 0) //using Inner Diameter as radius so it is doubled?
		if err != nil {
			return nil, err
		}
		funnelEnd, err := sdf.Sphere3D(r2)
		if err != nil {
			return nil, err
		}
		funnelEnd = sdf.Cut3D(funnelEnd, v3.Vec{}, v3.Vec{Z: 1})
		funnelEnd = sdf.Cut3D(funnelEnd, v3.Vec{}, v3.Vec{X: -1})
		funnel = sdf.Cut3D(funnel, v3.Vec{}, v3.Vec{X: -1})
		funnelEnd = sdf.Transform3D(funnelEnd, moveFunnelEnd)
		funnel = sdf.Union3D(funnel, funnelEnd)

		funnel = sdf.Transform3D(funnel, positionFunnel)
		funnel = bendIt(funnel)
		return funnel, nil
	}

	funnelHole, err := makeBendFunnel(tubeLength-funnelShorter, 0, funnelID/2.0)
	if err != nil {
		panic(err)
	}
	funnel, err := makeBendFunnel(tubeLength-funnelShorter, thickness, funnelID/2.0+thickness) //using Inner Diameter as radius so it is doubled?
	if err != nil {
		panic(err)
	}

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
