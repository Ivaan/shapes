package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func makeMotorMount(setup ClockSetup) sdf.SDF3 {
	cornersBox, err := sdf.Box3D(sdf.V3{
		X: setup.MotorMount.MotorCornerToCorner + setup.MotorMount.MountThickness,
		Y: setup.MotorMount.MotorCornerToCorner + setup.MotorMount.MountThickness,
		Z: setup.MotorMount.MountDepth + setup.MotorMount.MountThickness,
	}, 0)
	if err != nil {
		panic(err)
	}
	cornersBox = sdf.Transform3D(
		cornersBox,
		sdf.RotateZ(sdf.Tau/8).Mul(
			sdf.Translate3d(sdf.V3{Z: (setup.MotorMount.MountThickness + setup.MotorMount.MountDepth) / 2}),
		),
	)

	edgesBox, _ := sdf.Box3D(sdf.V3{
		X: setup.MotorMount.MotorAcross + setup.MotorMount.MountThickness,
		Y: setup.MotorMount.MotorAcross + setup.MotorMount.MountThickness,
		Z: setup.MotorMount.MountDepth + setup.MotorMount.MountThickness}, 0)
	edgesBox = sdf.Transform3D(
		edgesBox,
		sdf.Translate3d(sdf.V3{Z: (setup.MotorMount.MountThickness + setup.MotorMount.MountDepth) / 2}),
	)

	innerCornersBox, _ := sdf.Box3D(sdf.V3{
		X: setup.MotorMount.MotorCornerToCorner,
		Y: setup.MotorMount.MotorCornerToCorner,
		Z: setup.MotorMount.MountDepth,
	}, 0)
	innerCornersBox = sdf.Transform3D(
		innerCornersBox,
		sdf.RotateZ(sdf.Tau/8).Mul(
			sdf.Translate3d(sdf.V3{Z: setup.MotorMount.MountDepth/2 + setup.MotorMount.MountThickness}),
		),
	)
	innerEdgesBox, _ := sdf.Box3D(sdf.V3{
		X: setup.MotorMount.MotorAcross,
		Y: setup.MotorMount.MotorAcross,
		Z: setup.MotorMount.MountDepth,
	}, 0)
	innerEdgesBox = sdf.Transform3D(
		innerEdgesBox,
		sdf.Translate3d(sdf.V3{Z: setup.MotorMount.MountDepth/2 + setup.MotorMount.MountThickness}),
	)

	shaftAllowance, _ := sdf.Cylinder3D(setup.MotorMount.MountThickness, setup.MotorMount.MotorShaftAllowance/2, 0)
	shaftAllowance = sdf.Transform3D(
		shaftAllowance,
		sdf.Translate3d(sdf.V3{
			Z: setup.MotorMount.MountThickness / 2}),
	)

	screwHole, _ := sdf.Cylinder3D(setup.MotorMount.MountThickness, setup.MotorMount.ScrewHoleDiameter/2, 0)
	screwHole = sdf.Transform3D(
		screwHole,
		sdf.Translate3d(sdf.V3{
			X: setup.MotorMount.MotorAcross/2 - setup.MotorMount.ScrewDistanceFromEdge,
			Y: setup.MotorMount.MotorAcross/2 - setup.MotorMount.ScrewDistanceFromEdge,
			Z: setup.MotorMount.MountThickness / 2}),
	)
	screwHoles := sdf.Union3D(
		screwHole,
		sdf.Transform3D(screwHole, sdf.RotateZ(sdf.Tau/4)),
		sdf.Transform3D(screwHole, sdf.RotateZ(sdf.Tau/2)),
		sdf.Transform3D(screwHole, sdf.RotateZ(sdf.Tau*3/4)),
	)

	return sdf.Difference3D(
		sdf.Intersect3D(cornersBox, edgesBox),
		sdf.Union3D(
			sdf.Intersect3D(innerCornersBox, innerEdgesBox),
			screwHoles,
			shaftAllowance,
		),
	)
}
