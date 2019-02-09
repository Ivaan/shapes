
package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

func main() {

	NumberOfLoops :=  7 //Number Of Loops"
	TotalHeight :=  1 //"Total Height"
	DiskDiameter :=  3 //"Diameter of 3D printed disk"
	Thickness :=  2 //Thickness of loop ring
	p := sdf.NewPolygon()
    c := sdf.Cylinder3D
	cutOffHeight := params.TotalHeight / 2;
    // var vcyl = cylinderFromZToZ;

    disk := makeFancyDisk(params.Thickness, cutOffHeight, params.DiskDiameter, params.NumberOfLoops);

    // return disk
    //     .subtract(vcyl(-cutOffHeight, cutOffHeight, 1, 6).translate([3,0,0]))
    //     .subtract(vcyl(-cutOffHeight, cutOffHeight, 1, 6).translate([-3,0,0]))
    //     .translate([0,0,cutOffHeight]);
}

func cylinderFromZToZ(fromZ, toZ, radius float64, round float64) Cylinder3D {
	return sdf.Transform3D(
        sdf.Cylinder3D(toZ - fromZ, radius, round),
        sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: fromZ}))
}

func repeateAround(solid sdf.SDF3, times int){
    // var copies = [];
    // for(var i = 0; i < times; ++i){
    //     copies.push(solid.rotateZ(i * 360 / times));
    // }
    // return copies;
}

// function makeFancyDisk(thickness, cutOffHeight, diskDiameter, numberOfLoops)
// {
//     var angle = atan((cutOffHeight - (thickness / 2)) / (diskDiameter / 2));
//     var tor = torus({ri:thickness, ro: diskDiameter / 4, fni:12}).rotateX(angle).translate([diskDiameter/4, 0, 0]);
//     return union(repeateAround(tor, numberOfLoops));
// }