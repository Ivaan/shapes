package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func makeTopFrame(setup ClockSetup) sdf.SDF3 {
	frame, _ := sdf.Box3D(sdf.V3{
		X: setup.Tumbler.ShortRadius*setup.Frame.WidthAsFractionOfShorRadius*2 + setup.Tumbler.Radius*2,
		Y: setup.Tumbler.ShortRadius * setup.Frame.WidthAsFractionOfShorRadius * 2,
		Z: setup.Frame.Thickness,
	},
		setup.Frame.Thickness/3,
	)

	shaftHolder, _ := sdf.Cylinder3D(setup.Frame.ShaftHolderLength-setup.Frame.Thickness/2, setup.Shaft.OD*1.5, setup.Shaft.OD/2)
	shaftHole, _ := sdf.Cylinder3D(setup.Frame.ShaftHolderLength, setup.Shaft.OD/2, 0)
	holderLift := (setup.Frame.ShaftHolderLength - setup.Frame.Thickness/2) / 2
	holeLift := setup.Frame.ShaftHolderLength/2 - setup.Frame.Thickness/2
	frame = sdf.Union3D(
		frame,
		sdf.Transform3D(
			shaftHolder,
			sdf.Translate3d(sdf.V3{-setup.Tumbler.Radius, 0, holderLift}),
		),
		sdf.Transform3D(
			shaftHolder,
			sdf.Translate3d(sdf.V3{setup.Tumbler.Radius, 0, holderLift}),
		),
	)
	frame.(*sdf.UnionSDF3).SetMin(sdf.RoundMin(setup.Frame.Thickness / 2))
	//frame.(*sdf.UnionSDF3).SetMin(sdf.ExpMin(32))

	shaftHoles := sdf.Union3D(
		sdf.Transform3D(
			shaftHole,
			sdf.Translate3d(sdf.V3{-setup.Tumbler.Radius, 0, holeLift}),
		),
		sdf.Transform3D(
			shaftHole,
			sdf.Translate3d(sdf.V3{setup.Tumbler.Radius, 0, holeLift}),
		),
	)

	return sdf.Difference3D(
		frame,
		shaftHoles,
	)
}
