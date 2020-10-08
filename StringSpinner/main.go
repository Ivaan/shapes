package main

import (
	"flag"

	"github.com/deadsy/sdfx/sdf"
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
	//p := loadParams()
	h := newHyperbol(.5, 3)
	solid := sdf.Extrude3D(h, 2)
	sdf.RenderSTL(solid, 300, "solid.stl")

}

type hyperbol struct {
	minX, maxX float64
	bb         sdf.Box2
}

func newHyperbol(minX, maxX float64) sdf.SDF2 {
	bb := sdf.Box2{
		Min: sdf.V2{minX, 0},
		//Max: sdf.V2{maxX, 1 / minX},
		Max: sdf.V2{maxX, -minX/2.0 + 2.0},
	}
	return hyperbol{
		minX: minX,
		maxX: maxX,
		bb:   bb,
	}
}

func (hyp hyperbol) Evaluate(p sdf.V2) float64 {
	x := hyp.nearestXPointTo(p)
	var q sdf V2

	if x < hyp.minX || p.X < hyp.minX { 
		x = hyp.MinX
		if p.Y < hyp.maxY {
			q = sdf.V2(x, p.Y)
		}
	} else if q.X < hyp.maxX {
		adjustedX = x
	} else {
		x = hyp.maxX
		adjustedX = hyp.maxX
	}
	return sdf.V2{x, hyp.curveFunc(adjustedX)}
	//p = -2/5 (x - 2 (y + 1))
	q := hyp.nearestPointTo(p)
	if q.X < hyp.minX { 
		if q.Y < 0 {
			return q.Sub(hyp.bb.Min).Length()
		} else if q.Y < hyp.bb.Max.Y {
			return hyp.minX - q.X
		} else {
			return q.Sub(sdf.V2{hyp.minX, hyp.bb.Max.Y}).Length()
		}
	} else if q.X < hyp.maxX {
		p = hyp.nearestPointTo(q)
		var s float64
		if q.Y > -q.X/2.0+2.0 {
			s = 1
		} else {
			s = -1
		}
		return s * q.Sub(p).Length()
	} else {
		if q.Y < 0 {
			return q.Sub(sdf.V2{hyp.maxX, 0}).Length()
		} else if q.Y < hyp.bb.Max.Y {
			return q.X - hyp.maxX - q.X
		} else {
			return q.Sub(sdf.V2{hyp.maxX, -hyp.maxX/2 + 2}).Length()
		}
	}
}

func (hyp hyperbol) BoundingBox() sdf.Box2 {
	return hyp.bb
}

func (hyp hyperbol) nearestXPointTo(q sdfV2) {
	return -2.0 / 5.0 * (q.X - 2.0*(q.Y+1.0))
}

func (hyp hyperbol) curveFunc(x) {
	return -x/2.0 + 2.0
}

func makeImpellerDisk(p params) sdf.SDF3 {
	//return rotate_extrude( polygon({points: points}) );
	return nil
}

type impellerDiskProfile struct {
	bearingRadius float64
	thickness     float64
	radius        float64
	turbineRadius float64
	tollerance    float64

	impellerCompM float64
	impellerCompB float64
	turbineCompM  float64
	turbineCompB  float64

	bb sdf.Box2
}

func newImpellerDiskProfile(p params) sdf.SDF2 {
	imp := impellerDiskProfile{}
	imp.bearingRadius = p.bearingRaceOD / 2
	imp.thickness = p.impellerThickness / 2
	imp.radius = p.impellerDiameter / 2
	imp.turbineRadius = imp.radius + p.turbineRadius
	imp.tollerance = p.impellerToTurbineDiskTollerance
	compFactor := 1 / p.impellerCompressionFactor
	imp.impellerCompM = (compFactor - 1) / (imp.radius - imp.bearingRadius)
	imp.impellerCompB = 1 - imp.impellerCompM*imp.bearingRadius
	imp.turbineCompM = (1 - compFactor) / (imp.turbineRadius - imp.radius)
	imp.turbineCompB = compFactor - imp.turbineCompM*imp.radius
	imp.bb = sdf.Box2{sdf.V2{0, 0}, sdf.V2{imp.turbineRadius, imp.thickness + imp.tollerance}}
	return imp
}

//Evaluate implements sdf
func (imp impellerDiskProfile) Evaluate(p sdf.V2) float64 {
	//p = -1/2 sqrt((sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3)/(3 2^(1/3)) + (2^(1/3) (x y - 4))/(sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3) + y^2/4) - 1/2 sqrt(-(sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3)/(3 2^(1/3)) - (2^(1/3) (x y - 4))/(sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3) - (y^3 - 8 x)/(4 sqrt((sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3)/(3 2^(1/3)) + (2^(1/3) (x y - 4))/(sqrt((27 x^2 - 27 y^2)^2 - 4 (3 x y - 12)^3) + 27 x^2 - 27 y^2)^(1/3) + y^2/4)) + y^2/2) + y/4
	if p.X < 0 { //left of origin
		if p.Y < 0 {
			//return distance to (0,0)
		} else if p.Y < imp.thickness+imp.tollerance {
			//return p.Y (negative?)
		} else {
			//return distance to (0, imp.thickness + imp.tollerance)
		}
	} else if p.X < imp.bearingRadius { //left of where curve starts

	} else if p.X < imp.turbineRadius { // within curve

	} else { //past (right of) model

	}
	/*
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
	*/
	return 0
}

//BoundingBox implements sdf
func (imp impellerDiskProfile) BoundingBox() sdf.Box2 {
	return imp.bb
}
