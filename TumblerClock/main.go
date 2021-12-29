package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"gopkg.in/yaml.v3"
)

//-----------------------------------------------------------------------------

func main() {
	//setup := makeDefaultClockSetup()
	//bytes, _ := yaml.Marshal(setup)
	//ioutil.WriteFile("setup.yaml", bytes, 0644)

	filename := flag.String("SetupFile", "setup.yaml", "File name with the setup for the Tumbler Clock")
	partsString := flag.String("Parts", "all", "Parts list to print") //(h|m)(t|u)(a|b|ag|g)

	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	var setup ClockSetup
	err = yaml.Unmarshal(yamlFile, &setup)
	if err != nil {
		panic(err)
	}
	setup = setup.computeSynthetics()

	//fmt.Printf("%+v\n", setup) //uncomment to write setup conents

	partsList := parsePartsString(*partsString)
	fmt.Printf("%+v\n", partsList)

	if partsList.topFrame {
		//render.RenderSTL(makeTopFrame(setup), 200, "TopFrame.stl")
		render.RenderSTLSlow(makeTopFrame(setup), 400, "High-TopFrame.stl")
	}
	if partsList.bottomFrame {
		//todo: genrate bottom frame model and render to file
	}
	if partsList.motorMount {
		mount := makeMotorMount(setup)
		render.RenderSTLSlow(mount, 400, "MotorMount.stl")
	}
	if partsList.colonGear {
		//todo: genrate colon gear model and render to file
	}

	tb := makeTumblerBuilder(setup)

	for _, p := range partsList.tumblers {
		t := tb.makeTumbler(p)
		//sdf.RenderSTL(t, 200, p.fileName())
		render.RenderSTLSlow(t, 400, "High-"+p.fileName())
	}

}

func mainOld(setup ClockSetup) {
	facesA := [9]int{0, 0, 1, 1, 0, 1, 1, 1, 1} // A faces
	facesB := [9]int{0, 0, 1, 1, 0, 0, 1, 0, 1} // B faces
	tumblerOutsideA := makeTumblerOutside(setup.Tumbler, facesA)
	tumblerOutsideB := makeTumblerOutside(setup.Tumbler, facesB)

	insideHole := makeBearingHole(setup.Bearing, setup.BearingHolder, setup.Tumbler.FaceEdgeHeight)

	tracks, nibs := makePusherTracksAndNibs(setup.Transmission, setup.Tumbler, setup.Bearing.OD/2+setup.BearingHolder.Thickness)

	holes := sdf.Union3D(insideHole, tracks)
	tumblerA := sdf.Difference3D(tumblerOutsideA, holes)
	tumblerA = sdf.Union3D(tumblerA, nibs)

	tumblerB := sdf.Difference3D(tumblerOutsideB, holes)
	tumblerB = sdf.Union3D(tumblerB, nibs)

	//spacerDisk := makeSpacerDisk(shaftOD, spacerShaftTollerance, bearingID, spacerBearingTollerance, spacerBearingPenetrationDepth, tumblerSpacing, spacerGapAngle)
	spacerDisk := makeSimpleSpacerDisk(setup.Shaft, setup.Spacer, setup.Tumbler)
	// sdf.RenderSTLSlow(tumblerA, 400, "tumblerA.stl")
	// sdf.RenderSTLSlow(tumblerB, 400, "tumblerB.stl")
	//sdf.RenderSTLSlow(spacerDisk, 100, "spacerDisk.stl")
	render.RenderSTL(tumblerA, 200, "tumblerA.stl")
	render.RenderSTL(tumblerB, 200, "tumblerB.stl")
	render.RenderSTL(spacerDisk, 400, "spacerDiskWood.stl")

}

func makeShaftAndSpacerHole(shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness float64) sdf.SDF3 {
	innerChamfer := (tumblerShortRadius - tumblerMinimumWallThickness) - (shaftOD/2 + tumblerShaftTollerance)
	//innerChamfer := 10.0
	fmt.Printf("shaftOD: %v, tumblerShaftTollerance: %v, tumblerShaftBearingSurfaceLength: %v, spacerThickness: %v, spacerShaftTollerance: %v, spacerEdgeInFromTumblerEdge: %v, spacerTumblerTollerance: %v, tumblerFaceEdgeWidth: %v, tumblerFaceEdgeHeight: %v, tumblerShortRadius: %v, tumblerSpacing: %v, tumblerMinimumWallThickness: %v, innerChamfer: %v", shaftOD, tumblerShaftTollerance, tumblerShaftBearingSurfaceLength, spacerThickness, spacerShaftTollerance, spacerEdgeInFromTumblerEdge, spacerTumblerTollerance, tumblerFaceEdgeWidth, tumblerFaceEdgeHeight, tumblerShortRadius, tumblerSpacing, tumblerMinimumWallThickness, innerChamfer)

	poly, _ := sdf.Polygon2D([]sdf.V2{
		{0, tumblerFaceEdgeHeight / 2},
		{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeHeight / 2},
		{tumblerShortRadius - spacerEdgeInFromTumblerEdge, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance},
		{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance},
		{shaftOD/2 + tumblerShaftTollerance, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength},
		{shaftOD/2 + tumblerShaftTollerance + innerChamfer, tumblerFaceEdgeHeight/2 - spacerThickness - spacerTumblerTollerance - tumblerShaftBearingSurfaceLength - innerChamfer},
		{shaftOD/2 + tumblerShaftTollerance + innerChamfer, -tumblerFaceEdgeHeight/2 + tumblerShaftBearingSurfaceLength + innerChamfer},
		{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeHeight/2 + tumblerShaftBearingSurfaceLength},
		{shaftOD/2 + tumblerShaftTollerance, -tumblerFaceEdgeHeight / 2},
		{0, -tumblerFaceEdgeHeight / 2},
		{0, tumblerFaceEdgeHeight / 2},
	})

	rev, _ := sdf.Revolve3D(poly)
	return rev
}

func makeSpacerDisk(shaftOD, spacerShaftTollerance, bearingID, spacerBearingTollerance, spacerBearingPenetrationDepth, tumblerSpacing, spacerGapAngle float64) sdf.SDF3 {
	spacerHeight := spacerBearingPenetrationDepth*2 + tumblerSpacing
	diskProfile, _ := sdf.Polygon2D([]sdf.V2{
		{shaftOD/2 + spacerShaftTollerance, -spacerHeight / 2},
		{bearingID/2 - spacerBearingTollerance, -spacerHeight / 2},
		{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth},
		{bearingID/2 - spacerBearingTollerance + tumblerSpacing/2, -spacerHeight/2 + spacerBearingPenetrationDepth + tumblerSpacing/2},
		{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth + tumblerSpacing},
		{bearingID/2 - spacerBearingTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth*2 + tumblerSpacing},
		{shaftOD/2 + spacerShaftTollerance, -spacerHeight/2 + spacerBearingPenetrationDepth*2 + tumblerSpacing},
		{shaftOD/2 + spacerShaftTollerance, -spacerHeight / 2},
	})
	disk, _ := sdf.RevolveTheta3D(
		diskProfile,
		sdf.Tau-spacerGapAngle,
	)
	return disk
}

func makeSimpleSpacerDisk(shaft Shaft, spacer Spacer, tumbler Tumbler) sdf.SDF3 {
	diskProfile, _ := sdf.Polygon2D([]sdf.V2{
		{shaft.OD/2 + spacer.ShaftTollerance, 0},
		{shaft.OD/2 + spacer.ShaftTollerance + spacer.DiskWidth, 0},
		{shaft.OD/2 + spacer.ShaftTollerance + spacer.DiskWidth, tumbler.Spacing},
		{shaft.OD/2 + spacer.ShaftTollerance, tumbler.Spacing},
		{shaft.OD/2 + spacer.ShaftTollerance, 0},
	})
	disk, _ := sdf.RevolveTheta3D(
		diskProfile,
		sdf.Tau-spacer.GapAngle,
	)
	return disk
}
