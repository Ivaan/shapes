package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {

	outerDiameter := 45.0       //Outer diameter of cylinder
	holeDiameter := 22.0        //diameter for finger hole
	totalHeight := 32.0         //total finished height of cylinder
	innerSphearDiameter := 42.0 //diameter of inner ball
	sphearTollerance := 1.5     //gap between inner ball and socket
	round := 1.0                //roundness of cylinder corners
	//main cylinder
	//diameter of outerDiameter
	//height of innerSphearDiameter or totalHeight which ever is greater
	var height float64
	if innerSphearDiameter > totalHeight {
		height = innerSphearDiameter
	} else {
		height = totalHeight
	}
	mainCylinder := sdf.Cylinder3D(height, outerDiameter/2, round)

	//minus the socket
	//which is a sphear of innerSphearDiameter plus sphearTollerance
	socket := sdf.Sphere3D((innerSphearDiameter + sphearTollerance) / 2)

	//plus the inner ball
	//which is a sphear of innerSphear diameter
	innerBall := sdf.Sphere3D(innerSphearDiameter / 2)

	//minus the finger hole
	//diameter of holeDiameter
	//with height of main cylinder
	fingerHole := sdf.Cylinder3D(height, holeDiameter/2, round)

	//truncate to totalHeight
	//by intersecting to a total height cylinder
	truncatingCylinder := sdf.Cylinder3D(totalHeight, outerDiameter/2, round)

	//fiddler
	//do the math
	fiddler := sdf.Intersect3D(
		sdf.Union3D(
			sdf.Difference3D(mainCylinder, socket),
			sdf.Difference3D(innerBall, fingerHole),
		),
		truncatingCylinder,
	)
	sdf.RenderSTL(fiddler, 300, "fiddler.stl")
}

func cylinderFromZToZ(fromZ, toZ, radius float64, round float64) sdf.SDF3 {
	return sdf.Transform3D(
		sdf.Cylinder3D(toZ-fromZ, radius, round),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: fromZ + (toZ-fromZ) / 2}))
}
