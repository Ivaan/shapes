package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func MakeNoduleDebug1() Nodule {
	s, _ := sdf.Sphere3D(2)
	return MakeNodule(
		[]sdf.SDF3{},
		[]sdf.SDF3{s},
	)
}

func MakeNoduleDebug2() Nodule {
	c, _ := sdf.Cylinder3D(2, 4, 0)
	return MakeNodule(
		[]sdf.SDF3{},
		[]sdf.SDF3{c},
	)
}

func MakeNoduleDebug3() Nodule {
	c, _ := sdf.Cylinder3D(2, 4, 0)
	s, _ := sdf.Sphere3D(2)
	s = sdf.Transform3D(s, sdf.Translate3d(sdf.V3{Z: 2}))
	return MakeNodule(
		[]sdf.SDF3{s},
		[]sdf.SDF3{c},
	)
}
