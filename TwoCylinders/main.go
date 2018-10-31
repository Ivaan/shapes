// wallaby camshaft

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	c := sdf.Cylinder3D(20, 10, 2)
	t := sdf.Translate3d((sdf.V3{X: 10, Y: 0, Z: 10}))
	negC := sdf.Transform3D(c, t)

	model := sdf.Difference3D(c, negC)

	sdf.RenderSTL(model, 800, "twoCylinders.stl")

}
