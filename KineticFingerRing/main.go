package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	fingerDiameter := 22.0     //diameter of finger ring is to be worn on
	innerRingThickness := 1.0  // thickness of material against finger
	innerRingTolerance := 0.0  // gap between non moving inner rings
	middleRingThickness := 1.0 // thickness of material of bearing surface
	outerRingTolerance := 0.5  // gap between moving bit and middle ring
	outerRingThickenss := 3.0  // thickness of outer, knerled, ring
	keeperWallLip := 0.75

	overallWidth := 9.0        // total width of band
	keeperWallThickness := 1.0 // thickness of wall holding ring in place
	wallTolerance := 0.5       // gap between wall and knerled ring

	//side wall
	//used for both inner and middle ring
	//cyl at keeperWallThickness high and fingerDiameter/2 + ring thicknesses + tolerances + lip
	// with a whole cut out finger diamter
	wall := sdf.Transform3D(
		sdf.Difference3D(
			sdf.Cylinder3D(keeperWallThickness, fingerDiameter/2+innerRingThickness+innerRingTolerance+middleRingThickness+outerRingTolerance+keeperWallLip, 0.1), //disk
			sdf.Cylinder3D(keeperWallThickness, fingerDiameter/2, 0.1), //finger hole
		),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: keeperWallThickness / 2})) //move up to lay on xy plane

	//innerRing
	// wall plus ring at inside most edge against finger
	innerRing := sdf.Union3D(
		wall,
		sdf.Transform3D(
			sdf.Difference3D(
				sdf.Cylinder3D(overallWidth-keeperWallThickness, fingerDiameter/2+innerRingThickness, 0.1), //finger hole plus thinckness
				sdf.Cylinder3D(overallWidth-keeperWallThickness, fingerDiameter/2, 0.1),                    //minus finger hole
			),
			sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: (overallWidth - keeperWallThickness) / 2}),
		), //move up to lay on xy plane
	)

	//middleRing
	// wall plus ring at inside most edge between inner and what the turning ring runs along
	middleRing := sdf.Union3D(
		wall,
		sdf.Transform3D(
			sdf.Difference3D(
				sdf.Cylinder3D(overallWidth-keeperWallThickness, fingerDiameter/2+innerRingThickness+innerRingTolerance+middleRingThickness, 0.1), //innerring plus tolerance plus middle thickness
				sdf.Cylinder3D(overallWidth-keeperWallThickness, fingerDiameter/2+innerRingThickness+innerRingTolerance, 0.1),                     // minus innerring and tolerance
			),
			sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: (overallWidth - keeperWallThickness) / 2}),
		), //move up to lay on xy plane
	)

	//outerRing
	//kinetic bit that moves, knerled and all that
	//overall thickness minus wall and walls tolerance
	//diameter adds all the rings and tolerances
	r := fingerDiameter/2 + innerRingThickness + innerRingTolerance + middleRingThickness + outerRingTolerance + outerRingThickenss
	outerRing := sdf.Transform3D(
		sdf.Difference3D(
			sdf.KnurledHead3D(r, overallWidth-2*(keeperWallThickness+wallTolerance), r*0.25),
			sdf.Cylinder3D(overallWidth-2*(keeperWallThickness+wallTolerance), fingerDiameter/2+innerRingThickness+innerRingTolerance+middleRingThickness+outerRingTolerance, 0.1),
		),
		sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: (overallWidth - 2*(keeperWallThickness+wallTolerance)) / 2}),
	)

	sdf.RenderSTL(innerRing, 600, "innerRing.stl")
	sdf.RenderSTL(middleRing, 600, "middleRing.stl")
	sdf.RenderSTL(outerRing, 600, "outerRing.stl")

}
