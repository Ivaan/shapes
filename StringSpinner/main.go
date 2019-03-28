package main

import (
	"flag"
	"fmt"
)

type params struct {
	bearingRaceOD                   float64
	bearingRaceID                   float64
	bearingRaceThickness            float64
	bearingRaceHeight               float64
	impellerDiameter                float64
	impellerThickness               float64
	impellerBladeCount              int
	impellerBladeThickness          float64
	impellerToTurbineDiskTollerance float64
	impellerCompressionFactor       float64
	fanBladeCount                   int
	turbineVaneCount                int
	turbineRadius                   float64
}

func loadParams() params {
	p := params{}
	flag.Float64Var(&p.bearingRaceOD, "BearingRaceOD", 13., "Bearing Race Outside Diameter")
	flag.Float64Var(&p.bearingRaceID, "BearingRaceID", 5, "Bearing Race Inside Diameter")
	flag.Float64Var(&p.bearingRaceThickness, "BearingRaceThickness", 1, "Bearing Race Thickness")
	flag.Float64Var(&p.bearingRaceHeight, "BearingRaceHeight", 4, "Bearing Race Height")
	flag.Float64Var(&p.impellerDiameter, "ImpellerDiameter", 60, "Diameter of impeller")
	flag.Float64Var(&p.impellerThickness, "ImpellerThickness", 10, "Space between bearings (impeller disk thickness)")
	flag.IntVar(&p.impellerBladeCount, "ImpellerBladeCount", 11, "Numer of impeller blades")
	flag.Float64Var(&p.impellerBladeThickness, "ImpellerBladeThickness", 2, "Impeller blade thickness")
	flag.Float64Var(&p.impellerToTurbineDiskTollerance, "ImpellerToTurbineDiskTollerance", .2, "Tollerance between impeller and outer disk")
	flag.Float64Var(&p.impellerCompressionFactor, "ImpellerCompressionFactor", 2, "Impeller Compression factor, amount air compresses from atmo")
	flag.IntVar(&p.fanBladeCount, "FanBladeCount", 7, "Numer of fan blades")
	flag.IntVar(&p.turbineVaneCount, "TurbineVaneCount", 17, "Numer of turbine vanes")
	flag.Float64Var(&p.turbineRadius, "TurbineRadius", 12, "Turbine radius beyond impeller")
	flag.Parse()
	return p
}

func main() {
	p := loadParams()
	

}

func makeImpellerDisk(p params) sdf.SDF3{



	return rotate_extrude( polygon({points: points}) );
}

type impellerDiskProfile struct {
	bearingRadius float64
	thickness float64
	radius float64
	turbineRadius float64
	tollerance float64

	impellerCompM float64
	impellerCompB float64
	turbineCompM float64
	turbineCompB float64

	bb     sdf.Box2
}

func newImpellerDiskProfile(p params) sdf.SDF2{
	imp := impellerDisk{}
	imp.bearingRadius := p.bearingRaceOD / 2
	imp.thickness := p.impellerThickness / 2
	imp.radius := p.impellerDiameter/2
	imp.turbineRadius := radius + p.turbineRadius
	imp.tollerance := p.impellerToTurbineDiskTollerance
	compFactor := 1 / p.impellerCompressionFactor
	imp.impellerCompM := (compFactor-1)/(imp.radius-imp.bearingRadius)
	imp.impellerCompB := 1 - imp.impellerCompM*imp.bearingRadius
	imp.turbineCompM := (1-compFactor)/(imp.turbineRadius-imp.radius)
	imp.turbineCompB := compFactor - turbineCompM*imp.radius
	imp.bb = sdf.Box2{sdv.V2{0,0}, sdf.V2{imp.turbineRadius, imp.thickness + imp.tollerance}}
}

//Evaluate implements sdf
func (imp *impellerDiskProfile) Evaluate(p sdf.V2) float64 {

	if(p.X < bearingRadius)
	var points = [];
	points.push([0, -tollerance]);
	points.push([0, thickness]);
	var m = 
	var b = ;
	for(var x = bearingRadius; x < radius; x = x + 1){
		y = bearingRadius / x  * thickness * (m*x + b);
		points.push([x,y]);
	}
	m = (1-compFactor)/(turbineRadius-radius);
	b = compFactor - m*radius;
	for(x = radius; x < turbineRadius; x = x + 1){
		y = bearingRadius / x  * thickness * (m*x + b);
		points.push([x,y]);
	}
	
	points.push([turbineRadius+1, bearingRadius / x  * thickness * (m*turbineRadius + b)]);
	points.push([turbineRadius+1, -tollerance]);
	points.push([0, -tollerance]);
}

//BoundingBox implements sdf
func (imp *impellerDisk) BoundingBox() sdf.Box2 {
	return imp.bb
}
