package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	number_teeth := 19
	gearDiameter := 50.0
	module := gearDiameter / float64(number_teeth)
	pa := sdf.DtoR(20.0)
	h := 10.0
	gear_2d := sdf.InvoluteGear(
		number_teeth, // number_teeth
		module,       // gear_module
		pa,           // pressure_angle
		0.1,          // backlash
		0.1,          // clearance
		10,           // ring_width
		7,            // facets
	)
	gear3d := sdf.TwistExtrude3D(gear_2d, h, sdf.Tau/float64(number_teeth))
	gear3d2 := sdf.TwistExtrude3D(gear_2d, h, -sdf.Tau/float64(number_teeth))

	disk := sdf.Transform3D(
		sdf.Cylinder3D(h, gearDiameter/2, .5),
		sdf.Translate3d(sdf.V3{0, 0, h}),
	)
	pair := sdf.Union3D(gear3d, disk)
	other := sdf.Transform3D(
		sdf.Union3D(gear3d2, disk),
		sdf.Translate3d(sdf.V3{gearDiameter, 0, 0}),
	)
	sdf.RenderSTLSlow(sdf.Union3D(pair, other), 300, "gear.stl")
}
