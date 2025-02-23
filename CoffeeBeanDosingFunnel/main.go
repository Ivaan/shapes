package main

import (
	"fmt"
	"math"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"

	// v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	tubeID := 22.0
	funnelID := 55.0
	vialRound := 3.0
	tubeHolderStack := []Dimensions{
		{Height: 1.0, Radius: 25.0 / 2.0}, //lip
		{Height: 9.0, Radius: 27.3 / 2.0}, //threads
		// {Height: 1.0, Radius: 25.0 / 2.0}, //after threads, before body of vial
		{Height: 10.0, Radius: 33.0 / 2.0, Round: vialRound},
	}
	// tubeLipID := 25.0
	tubeLength := 120.0
	funnelShorter := 11.0
	thickness := 2.0
	bendRadius := 800.0
	// this is the amount the pipe dips due to the bend C - A where C^2 = A^2 + B^2
	bendDip := math.Sqrt(bendRadius*bendRadius+tubeLength*tubeLength) - bendRadius
	openingAngle := 100.0 / 360.0 * sdf.Tau
	thumbSpotDepth := thickness * 2
	thumbSpotDepthSphereDistance := 30.0

	doText := true
	// no6Text := false //true is No6 text, false is NAZ fab text
	textHeight := 13.0
	cutTextHeight := 15.0
	textPositionFromEnd := 20.0
	textLegWidth := 35.0
	backLegCenter := tubeLength - textLegWidth/2.0 - funnelShorter

	// base, err := sdf.Box3D(v3.Vec{X: tubeLength, Y: tubeID/2 + thickness + bendDip, Z: tubeID + thickness*2}, thickness)
	leg, err := sdf.Box3D(v3.Vec{X: textLegWidth, Y: tubeID/2 + thickness + bendDip, Z: tubeID + thickness*2}, thickness)
	base := sdf.Union3D(
		sdf.Transform3D(leg, sdf.Translate3d(v3.Vec{X: -tubeLength/2 + textPositionFromEnd})),
		sdf.Transform3D(leg, sdf.Translate3d(v3.Vec{X: -tubeLength/2 + backLegCenter})),
	)
	if err != nil {
		panic(err)
	}
	f, err := sdf.LoadFont("Arial.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	// text := sdf.NewText("№6")
	// cutText := sdf.NewText(" ")
	text := sdf.NewText("fab")
	cutText := sdf.NewText("NAZ")

	text2d, err := sdf.Text2D(f, text, textHeight)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	cutText2d, err := sdf.Text2D(f, cutText, cutTextHeight)
	if err != nil {
		fmt.Printf("can't generate X sdf2 %s\n", err)
		os.Exit(1)
	}

	text3d := sdf.Extrude3D(text2d, thickness)
	if err != nil {
		fmt.Printf("can't generate text sdf3 %s\n", err)
		os.Exit(1)
	}
	text3d = sdf.Union3D(
		sdf.Transform3D(
			text3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: tubeID/2.0 + thickness}),
		),
		sdf.Transform3D(
			text3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: -tubeID/2.0 - thickness}).Mul(sdf.RotateY(1.0/2.0*sdf.Tau)),
		),
	)
	cutText3d := sdf.Extrude3D(cutText2d, thickness)
	if err != nil {
		fmt.Printf("can't generate text sdf3 %s\n", err)
		os.Exit(1)
	}
	cutText3d = sdf.Union3D(
		sdf.Transform3D(
			cutText3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: tubeID/2.0 + thickness}),
		),
		sdf.Transform3D(
			cutText3d,
			sdf.Translate3d(v3.Vec{X: -tubeLength/2.0 + textPositionFromEnd, Z: -tubeID/2.0 - thickness}).Mul(sdf.RotateY(1.0/2.0*sdf.Tau)),
		),
	)
	cutText3d = sdf.Difference3D(cutText3d, text3d)
	if doText {
		base = sdf.Union3D(
			base,
			text3d,
		)
		base = sdf.Difference3D(base, cutText3d)
	}

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

	channelOpeningKeepNormal := sdf.RotateZ(-openingAngle / 2.0).MulPosition(v3.Vec{Y: -1.0})
	channelCut1 := sdf.Cut3D(channel, v3.Vec{}, channelOpeningKeepNormal)
	channelCut2 := sdf.Cut3D(channel, v3.Vec{}, sdf.MirrorXZ().MulPosition(channelOpeningKeepNormal))
	channel = sdf.Union3D(channelCut1, channelCut2)
	// channelOpening = sdf.Intersect3D(channel, channelOpening)
	channelHole, err := sdf.Cylinder3D(tubeLength, tubeID/2.0, 0)
	if err != nil {
		panic(err)
	}
	channelHole = sdf.Transform3D(channelHole, sdf.Translate3d(v3.Vec{Z: -thickness})) //This moves the hole so the other end of the channel isn't cut away (leaving a wall)
	channelHole = bendIt(channelHole)
	channel = bendIt(channel)

	positionFunnel := sdf.RotateY(math.Atan2(tubeID, tubeLength)).Mul(sdf.Translate3d(v3.Vec{X: tubeID / 2.0, Z: -funnelShorter / 2.0}))
	moveFunnelEnd := sdf.Translate3d(v3.Vec{Z: tubeLength/2.0 - funnelShorter/2.0})

	makeBendFunnel := func(length, r1, r2 float64, cut bool) (sdf.SDF3, error) {
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
		if cut {
			funnel = sdf.Cut3D(funnel, v3.Vec{}, v3.Vec{X: -1})
		}
		funnelEnd = sdf.Transform3D(funnelEnd, moveFunnelEnd)
		funnel = sdf.Union3D(funnel, funnelEnd)

		funnel = sdf.Transform3D(funnel, positionFunnel)
		funnel = bendIt(funnel)
		return funnel, nil
	}

	funnelHole, err := makeBendFunnel(tubeLength-funnelShorter, 0, funnelID/2.0, false)
	if err != nil {
		panic(err)
	}
	funnel, err := makeBendFunnel(tubeLength-funnelShorter, thickness, funnelID/2.0+thickness, true) //using Inner Diameter as radius so it is doubled?
	if err != nil {
		panic(err)
	}

	vialHolderHole, err := stackedCylinders(0, tubeHolderStack...)
	if err != nil {
		panic(err)
	}

	d := getMaxRadiusAndSumHeight(tubeHolderStack...)
	d.Radius += thickness
	d.Height += -vialRound
	d.Round = thickness
	vialHolder, err := stackedCylinders(0, d)
	if err != nil {
		panic(err)
	}

	thumbSpotSphere, err := sdf.Sphere3D(thumbSpotDepthSphereDistance)
	if err != nil {
		panic(err)
	}

	vialHolder = sdf.Difference3D(vialHolder, sdf.Transform3D(thumbSpotSphere, sdf.Translate3d(v3.Vec{Y: thumbSpotDepthSphereDistance + d.Radius - thumbSpotDepth})))
	vialHolder = sdf.Difference3D(vialHolder, sdf.Transform3D(thumbSpotSphere, sdf.Translate3d(v3.Vec{Y: -thumbSpotDepthSphereDistance - d.Radius + thumbSpotDepth})))

	vialHolderCut1 := sdf.Cut3D(vialHolder, v3.Vec{}, channelOpeningKeepNormal)
	vialHolderCut2 := sdf.Cut3D(vialHolder, v3.Vec{}, sdf.MirrorXZ().MulPosition(channelOpeningKeepNormal))
	vialHolder = sdf.Union3D(vialHolderCut1, vialHolderCut2)
	vialHolderHole = sdf.Transform3D(
		vialHolderHole,
		sdf.Translate3d(v3.Vec{X: -tubeLength / 2.0, Y: bendDip}).Mul(
			sdf.RotateY(-1.0/4.0*sdf.Tau).Mul(
				sdf.RotateZ(1.0/4.0*sdf.Tau),
			),
		),
	)
	vialHolder = sdf.Transform3D(
		vialHolder,
		sdf.Translate3d(v3.Vec{X: -tubeLength / 2.0, Y: bendDip}).Mul(
			sdf.RotateY(-1.0/4.0*sdf.Tau).Mul(
				sdf.RotateZ(1.0/4.0*sdf.Tau),
			),
		),
	)

	holes := sdf.Union3D(channelHole, funnelHole, vialHolderHole)
	pipe := sdf.Union3D(channel, funnel)
	pipe = sdf.Union3D(pipe, base, vialHolder)
	pipe.(*sdf.UnionSDF3).SetMin(sdf.RoundMin(thickness))
	pipe = sdf.Difference3D(pipe, holes)

	render.ToSTL(pipe, "pipe.stl", render.NewMarchingCubesUniform(600))

}

type Dimensions struct {
	Height float64
	Radius float64
	Round  float64
}

func stackedCylinders(startZ float64, dimensionses ...Dimensions) (sdf.SDF3, error) {
	cylinders := make([]sdf.SDF3, len(dimensionses))
	targetZ := startZ
	for i, d := range dimensionses {
		c, err := sdf.Cylinder3D(d.Height, d.Radius, d.Round)
		if err != nil {
			return nil, err
		}
		c = sdf.Transform3D(c, sdf.Translate3d(v3.Vec{Z: d.Height/2.0 + targetZ}))
		targetZ += d.Height
		cylinders[i] = c
	}
	return sdf.Union3D(cylinders...), nil
}

func getMaxRadiusAndSumHeight(dimensionses ...Dimensions) Dimensions {
	maxd := Dimensions{Height: 0, Radius: dimensionses[0].Radius}
	for _, d := range dimensionses {
		if maxd.Radius < d.Radius {
			maxd.Radius = d.Radius
		}
		maxd.Height += d.Height
	}
	return maxd
}
