package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func main() {

	NumberOfLoops := 7   //Number Of Loops"
	TotalHeight := 12.0  //"Total Height"
	DiskDiameter := 30.0 //"Diameter of 3D printed disk"
	Thickness := 3.0     //Thickness of loop ring
	cutOffHeight := TotalHeight / 2
	// var vcyl = cylinderFromZToZ;

	disk := sdf.Difference3D(
		makeFancyDisk(Thickness, cutOffHeight, DiskDiameter, NumberOfLoops),
		sdf.MultiCylinder3D(TotalHeight*4, 1, sdf.V2Set([]sdf.V2{sdf.V2{3, 0}, sdf.V2{-3, 0}})))

	sdf.RenderSTL(disk, 300, "ButtonCorse.stl")
	// return disk
	//     .subtract(vcyl(-cutOffHeight, cutOffHeight, 1, 6).translate([3,0,0]))
	//     .subtract(vcyl(-cutOffHeight, cutOffHeight, 1, 6).translate([-3,0,0]))
	//     .translate([0,0,cutOffHeight]);
}

func cylinderFromZToZ(fromZ, toZ, radius float64, round float64) sdf.SDF3 {
	return sdf.Transform3D(
		sdf.Cylinder3D(toZ-fromZ, radius, round),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: fromZ}))
}

func repeateAround(solid sdf.SDF3, times int) []sdf.SDF3 {
	var copies = make([]sdf.SDF3, times)
	for i := 0; i < times; i++ {
		copies[i] = sdf.Transform3D(
			solid,
			sdf.RotateZ(float64(i)*sdf.TAU/float64(times)))
	}

	return copies
}

// function makeFancyDisk(thickness, cutOffHeight, diskDiameter, numberOfLoops)
// {
//     var angle = atan((cutOffHeight - (thickness / 2)) / (diskDiameter / 2));
//     var tor = torus({ri:thickness, ro: diskDiameter / 4, fni:12}).rotateX(angle).translate([diskDiameter/4, 0, 0]);
//     return union(repeateAround(tor, numberOfLoops));
// }

func makeFancyDisk(thickness, cutOffHeight, diskDiameter float64, numberOfLoops int) sdf.SDF3 {
	angle := math.Atan((cutOffHeight - (thickness / 2)) / (diskDiameter / 2))
	tor := sdf.Transform3D(
		torus(thickness, diskDiameter/4),
		sdf.RotateX(angle).Mul(sdf.Translate3d(sdf.V3{X: diskDiameter / 4, Y: 0, Z: 0})))
	return sdf.Union3D(repeateAround(tor, numberOfLoops)...)
}

func torus(r1, r2 float64) sdf.SDF3 {
	return sdf.Revolve3D(
		sdf.Transform2D(
			sdf.Circle2D(r1),
			sdf.Translate2d(sdf.V2{X: r2, Y: 0})))

}
