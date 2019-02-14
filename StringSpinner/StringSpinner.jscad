function getParameterDefinitions() {
	return [

        { name: "BearingRaceOD", type: "float", "rangeMin": 1, "rangeMax": 50, initial: 13., caption: "Bearing Race Outside Diameter" },
        { name: "BearingRaceID", type: "float", "rangeMin": 1, "rangeMax": 50, initial: 5, caption: "Bearing Race Inside Diameter" },
        { name: "BearingRaceThickness", type: "float", "rangeMin": 1, "rangeMax": 50, initial: 1, caption: "Bearing Race Thickness" },
        { name: "BearingRaceHeight", type: "float", "rangeMin": 1, "rangeMax": 50, initial: 4, caption: "Bearing Race Height" },
		{ name: "ImpellerDiameter", type: "float", "rangeMin": 0, "rangeMax": 100, initial: 60, caption: "Diameter of impeller" },
		{ name: "ImpellerThickness", type: "float", "rangeMin": 1, "rangeMax": 50, initial: 10, caption: "Space between bearings (impeller disk thickness)" },
		{ name: "ImpellerBladeCount", type: "int", "rangeMin": 1, "rangeMax": 50, initial: 11, caption: "Numer of impeller blades" },
		{ name: "ImpellerBladeThickness", type: "float", "rangeMin": 1, "rangeMax": 15, initial: 2, caption: "Impeller blade thickness" },
		{ name: "ImpellerToTurbineDiskTollerance", type: "float", "rangeMin": 0, "rangeMax": 5, initial: .2, caption: "Tollerance between impeller and outer disk"},
		{ name: "ImpellerCompressionFactor", type: "float", "rangeMin": 0, "rangeMax": 5, initial: 2, caption: "Impeller Compression factor, amount air compresses from atmo"},
		{ name: "FanBladeCount", type: "int", "rangeMin": 1, "rangeMax": 50, initial: 7, caption: "Numer of fan blades" },
		{ name: "TurbineVaneCount", type: "int", "rangeMin": 1, "rangeMax": 50, initial: 17, caption: "Numer of turbine vanes" },
		{ name: "TurbineRadius", type: "float", "rangeMin": 0, "rangeMax": 5, initial: 12, caption: "Turbine radius beyond impeller"},
	];
}

var cyl = cylinderFromXToX;
var vcyl = cylinderFromZToZ;
var vcone = cylinderFromZRToZR;

function main(params){
	var impellerDisk = makeImpellerDisk(params);
	var impellerBlades = makeImpellerBlades(params);
	var impellerHub = makeImpellerHub(params);
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	var bearingThickness = params.BearingRaceHeight;
	var impellerThickness = params.ImpellerThickness / 2;
	var turbineThickness = impellerThickness + tollerance + bearingThickness;

	impellerBlades = intersection(impellerDisk, union(impellerBlades));

	var impellerDriveShaft = makeDriveShaft(params, impellerHub);
	var bearing = makeBearing(params).translate([0, 0, params.ImpellerThickness / 2 + tollerance]);
	var impeller = union(impellerHub, impellerBlades).subtract(impellerDriveShaft).setColor(1, .5, 0);

	var turbine = makeTurbine(params, impellerDisk);

	//var bearingIR = params.BearingRaceID / 2;
	//var bearingThickness = params.BearingRaceHeight;
	//return [union( impellerDriveShaft, mirror([0,0,-1], impellerDriveShaft)).translate([0,0, thickness + tollerance + bearingThickness + bearingIR])];	


	//var thickness = params.ImpellerThickness / 2;
	//var tollerance = params.ImpellerToTurbineDiskTollerance;
	//return union(impeller, mirror([0,0,-1], impeller)).translate([0,0, thickness + tollerance])

	return union([impellerDriveShaft, impeller, bearing, turbine, mirror([0,0,-1], [impellerDriveShaft, impeller, bearing, turbine])]).translate([0,0, turbineThickness]);
}

function makeImpellerDisk(params){
	var bearingRadius = params.BearingRaceOD / 2;
	var thickness = params.ImpellerThickness / 2;
	var radius = params.ImpellerDiameter/2;
	var turbineRadius = radius + params.TurbineRadius;
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	var compFactor = 1 / params.ImpellerCompressionFactor;

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

function makeImpellerBlades(params){
	var bearingRadius = params.BearingRaceOD / 2;
	var thickness = params.ImpellerThickness / 2;
	var radius = params.ImpellerDiameter / 2;
	var bladeThickness = params.ImpellerBladeThickness;
	var bladeCount = params.ImpellerBladeCount;

	return repeateAround(
		linear_extrude({ height: bladeThickness }, square([radius, thickness])).translate([0, 0, -1]).rotateX(90),
		bladeCount);
}

function makeImpellerHub(params){
	var bearingRadius = params.BearingRaceOD / 2;
	var shaftRadius = params.BearingRaceID / 2;
	var thickness = params.ImpellerThickness / 2;
	var raceThickness = params.BearingRaceThickness;
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	
	return union(
		vcyl(0, thickness, bearingRadius, 72),
		vcyl(thickness, thickness + tollerance, shaftRadius + raceThickness, 72)
	);
	
}

function makeDriveShaft(params, hub){
	var thickness = params.ImpellerThickness / 2;
	var bearingIR = params.BearingRaceID / 2;
	var squareDriveIR = Math.sqrt(bearingIR * bearingIR * 2);
	var bearingThickness = params.BearingRaceHeight;
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	var stringHoleHeight = thickness + bearingThickness + bearingIR / 2;

	var squareDrive = intersection(vcyl(0, thickness + tollerance, squareDriveIR, 4), hub);
	var bearingSpot = vcyl(thickness + tollerance, thickness + tollerance + bearingThickness, bearingIR, 72);
	var pastBearing = vcyl(thickness + tollerance + bearingThickness, thickness + tollerance + bearingThickness + bearingIR, bearingIR * 0.9, 32);
	var stringHole = translate([0, 0, stringHoleHeight],cyl(-bearingIR, bearingIR, bearingIR * 0.3, 12));

	return union([squareDrive, bearingSpot, pastBearing]).subtract(stringHole).setColor(.5, .5, 0);
}

function makeTurbine(params, impellerDisk){
	var radius = params.ImpellerDiameter/2;
	var bearingRadius = params.BearingRaceOD / 2;
	var turbineRadius = radius + params.TurbineRadius;
	var impellerThickness = params.ImpellerThickness / 2;
	var bearingThickness = params.BearingRaceHeight;
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	var bearingMaterialThickness = params.BearingRaceThickness;


	var rawCyl = vcyl(0, impellerThickness + tollerance + bearingThickness, turbineRadius, 72);

	var bearingBlank = makeBearingBlank(params).translate([0, 0, impellerThickness / 2 + tollerance]);
	var fanHoleOuterR = bearingRadius + impellerThickness + bearingMaterialThickness;
	var fanHoleInnerR = bearingRadius + bearingMaterialThickness;
	
	var fanHole = vcyl(0, bearingThickness + 10, fanHoleOuterR, 72)
		.subtract(vcyl(0, bearingThickness + 10, fanHoleInnerR, 72))
		.translate([0, 0, impellerThickness / 2 + tollerance]);

	var fan = makeFan(params, fanHole, fanHoleOuterR, impellerThickness + tollerance + bearingThickness);

	var vanes = makeVanes(params);
	
	return union([
			rawCyl.subtract([impellerDisk.translate([0, 0, tollerance]), bearingBlank, fanHole]),
			fan,
			vanes])
			.setColor(0, 0.1, 0.8, 0.7)


}

function makeFan(params, fanHole, fanHoleOuterR, fanHoleTop){
	var fanBladeCount = params.FanBladeCount;
	var blade = intersection(fanHole, makeFanblade(params, fanHoleOuterR).translate([0, 0, fanHoleTop]));
	var blades = repeateAround(blade, fanBladeCount);
	return blades;
}

//incomplete makeFanblade - params not even solid nevermind used
function makeFanblade(params, fanHoleOuterR){
	var arc = [];
    for(var x = 0; x < 5; x += .1){
        y = x * x/25;
        arc.push(circle(.2).translate([x, y, 0]));    
    }
    var hull = chain_hull(arc);
    var extrude = hull.extrude({offset: [0, 0, fanHoleOuterR], twistangle: 45, twiststeps: 16})
    return extrude
        .rotateX(-90);
}

function makeVanes(params){
	var turbineVaneCount = params.TurbineVaneCount;
	var vane = makeVane(params);
	var vanes = repeateAround(vane, turbineVaneCount);
	return vanes;
}

function makeVane(params){
	var impellerThickness = params.ImpellerThickness / 2;
	var impellerRadius = params.ImpellerDiameter / 2;
	var bearingThickness = params.BearingRaceHeight;
	var tollerance = params.ImpellerToTurbineDiskTollerance;
	var turbineThickness = impellerThickness + tollerance + bearingThickness;
	var bladeThickness = params.ImpellerBladeThickness;
	var bladeCount = params.ImpellerBladeCount;
	var turbineVaneCount = params.TurbineVaneCount;
	var vaneThicknessRadius = bladeThickness * bladeCount / turbineVaneCount / 2;
	var vaneParabolicOffset = 3;
	var vaneParabolicScale = 10;
	var turbineRadius = params.TurbineRadius;
	var xOffset = vaneParabolicOffset * vaneParabolicOffset / vaneParabolicScale
	var arc = [];
	for(var y = vaneThicknessRadius; y < turbineRadius - vaneThicknessRadius; y++){
		x = (y+vaneParabolicOffset)*(y+vaneParabolicOffset)/vaneParabolicScale - xOffset;

        arc.push(circle(vaneThicknessRadius).translate([0, y + impellerRadius, 0]).rotateZ(-x));    
	}
	var hull = chain_hull(arc);
	var extrude = hull.extrude({offset: [0, 0, turbineThickness]});
	return extrude;
}

function cylinderFromXToX(fromX, toX, radius, resolution){
	return CSG.cylinder({start: [fromX, 0, 0], end: [toX, 0, 0], radius: radius, resolution: resolution});
}

function cylinderFromZToZ(fromZ, toZ, radius, resolution){
	return CSG.cylinder({start: [0, 0, fromZ], end: [0, 0, toZ], radius: radius, resolution: resolution});
}

function cylinderFromZRToZR(fromZ, toZ, fromRadius, toRadius, resolution){
	return CSG.cylinder({start: [0, 0, fromZ], end: [0, 0, toZ], radiusStart: fromRadius, radiusEnd: toRadius, resolution: resolution});
}

function repeateAround(solid, times){
	var copies = [];
	for(var i = 0; i < times; ++i){
		copies.push(solid.rotateZ(i * 360 / times));
	}
	return copies;
}

function makeBearingBlank(params){
    var or = params.BearingRaceOD / 2;
    var ir = params.BearingRaceID / 2;
    var thickness = params.BearingRaceHeight;

	return vcyl(0, thickness+10, or, 72).setColor(0,0,0);
	
}

function makeBearing(params){
    var or = params.BearingRaceOD / 2;
    var ir = params.BearingRaceID / 2;
    var thickness = params.BearingRaceHeight;
    var raceThickness = params.BearingRaceThickness;
    var ballCenter = (or + ir) / 2;
    var ballRadius = (or - ir - raceThickness) / 2;
	var bottom = [0, 0, thickness / -2];
	var top = [0, 0, thickness / 2];

	var outer = CSG.cylinder({start: bottom, end: top, radius: or, resolution: 72})
		.subtract(CSG.cylinder({start: bottom, end: top, radius: or - raceThickness, resolution: 72}));
	var inner = CSG.cylinder({start: bottom, end: top, radius: ir + raceThickness, resolution: 72})
		.subtract(CSG.cylinder({start: bottom, end: top, radius: ir, resolution: 72}));

	var races = union(outer, inner).subtract(torus({ri:ballRadius, ro: ballCenter, fni:64}));

	var ball = sphere({r: ballRadius, center: true, fn: 20, type: 'geodesic'}).translate([ballCenter, 0, 0]);

	var bits = repeateAround(ball, 7);
	bits.push(races)

	return union(bits).translate([0, 0, thickness/2]).setColor(0.5, 0.5, 0.5);
}

