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

func makeImpellerDisk(p params){
	bearingRadius := p.bearingRaceOD / 2;
	thickness := p.impellerThickness / 2;
	radius := p.impellerDiameter/2;
	turbineRadius := radius + p.turbineRadius;
	tollerance := p.impellerToTurbineDiskTollerance;
	compFactor := 1 / p.impellerCompressionFactor;

	var points = [];
	points.push([0, -tollerance]);
	points.push([0, thickness]);
	var m = (compFactor-1)/(radius-bearingRadius);
	var b = 1 - m*bearingRadius;
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


	return rotate_extrude( polygon({points: points}) );
}
