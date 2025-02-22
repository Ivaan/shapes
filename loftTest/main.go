package main

import (
	//"math"

	//"github.com/deadsy/sdfx/obj"
	// "errors"
	// "fmt"
	// "math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

func main() {
	c1, err := sdf.Circle2D(2.0)
	// c1 = sdf.Cut2D(c1, v2.Vec{Y: -2}, v2.Vec{Y: 1})
	if err != nil {
		panic(err)
	}

	sq1 := sdf.Box2D(v2.Vec{X: 4, Y: 4}, 0.0)
	sq1 = sdf.Transform2D(sq1, sdf.Translate2d(v2.Vec{X: -2}))
	c1 = sdf.Difference2D(c1, sq1)
	top := sdf.Union2D(
		c1,
		sq1,
	)
	c2, err := sdf.Circle2D(4.0)
	// c2 = sdf.Cut2D(c2, v2.Vec{Y: -2}, v2.Vec{Y: 1})
	if err != nil {
		panic(err)
	}
	sq2 := sdf.Box2D(v2.Vec{X: 8, Y: 8}, 0.0)
	sq2 = sdf.Transform2D(sq2, sdf.Translate2d(v2.Vec{X: -4}))
	c2 = sdf.Difference2D(c2, sq2)
	bot := sdf.Union2D(
		c2,
		sq2,
	)

	coneish, err := sdf.Loft3D(bot, top, 4.0, 0.0)
	if err != nil {
		panic(err)
	}
	coneish2 := sdf.ScaleExtrude3D(top, 4.0, v2.Vec{X: 2, Y: 2})

	coneish3 := ExtrudeBy2DFunction(ExpandExtrude(top, 4, 4), 4, v3.Vec{}.AddScalar(4))

	render.ToSTL(coneish, "loftTest.stl", render.NewMarchingCubesUniform(150))
	render.ToSTL(coneish2, "scaleExtrudeTest.stl", render.NewMarchingCubesUniform(150))
	render.ToSTL(coneish3, "expandExtrude.stl", render.NewMarchingCubesUniform(150))
}
