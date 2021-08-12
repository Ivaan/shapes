package main

import (
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

func main() {
	//today these measurment include tolerances
	screwHeadDiameter := 3.7 //3.34
	screwHeadHeight := 1.2   //1.14
	magnetDiameter := 6.5    //6.03
	magnetHeight := 1.7      //1.49
	cubeSize := 8.0

	fmt.Printf("Pause print at: %vmm\n", cubeSize-screwHeadHeight-magnetHeight)

	cube := sdf.Difference3D(
		sdf.Box3D(sdf.V3{X: cubeSize, Y: cubeSize, Z: cubeSize}, .5),
		sdf.Union3D(
			sdf.Transform3D(
				sdf.Cylinder3D(screwHeadHeight, screwHeadDiameter/2, 0),
				sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: cubeSize/2 - screwHeadHeight/2}),
			),
			sdf.Transform3D(
				sdf.Cylinder3D(magnetHeight, magnetDiameter/2, 0),
				sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: (cubeSize/2 - screwHeadHeight - magnetHeight/2)}),
			),
		),
	)
	sdf.RenderSTLSlow(cube, 400, "cube.stl")
}
