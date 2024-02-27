package main

// import (
// 	//"math"

// 	//"github.com/deadsy/sdfx/obj"
// 	"errors"
// 	"fmt"
// 	"math"

// 	"github.com/deadsy/sdfx/render"
// 	"github.com/deadsy/sdfx/sdf"
// 	v2 "github.com/deadsy/sdfx/vec/v2"
// 	v3 "github.com/deadsy/sdfx/vec/v3"
// )

// type circ struct {
// 	Loc     v2.Vec
// 	R       float64
// 	Concave bool
// 	CutFrom v2.Vec
// 	CutTo   v2.Vec
// }

// func main() {
// 	circs := []circ{ //circles that describe the 2D profile of the keyboard's right half base
// 		{Loc: v2.Vec{X: 28.28427124746189, Y: -25.4558441227157}, R: 6},                               //6
// 		{Loc: v2.Vec{X: 12, Y: -9.171572875253815}, R: 6},                                             //5
// 		{Loc: v2.Vec{X: 12.00000000000001, Y: 73.5999813335492}, R: 6},                                //4
// 		{Loc: v2.Vec{X: 69.2653918090254, Y: -122.097465432435}, R: 209.903937602407},                 //3
// 		{Loc: v2.Vec{X: 136.5307836180504, Y: 70.39197021156274}, R: 6},                               //2
// 		{Loc: v2.Vec{X: 136.5307836180513, Y: 1.627100563552649}, R: 6},                               //1
// 		{Loc: v2.Vec{X: 86.14330244975497, Y: -72.6729449002023}, R: 68.3000454637575, Concave: true}, //0
// 	}
// 	screwLocations := []v2.Vec{
// 		{X: 102.1072999123607, Y: 80.15479331598767},
// 		{X: 80.56007930185542, Y: -0.9709640961869095},
// 		{X: 80.14283201431901, Y: 57.62620932295196},
// 		{X: 101.3815303433159, Y: 40.36481970108936},
// 		{X: 28.28427124746186, Y: -27.45584412271571},
// 		{X: 14.95467069791763, Y: -15.78309782266383},
// 		{X: 42.97660077726053, Y: 13.41342961679601},
// 		{X: 45.3166570024517, Y: 82.46787066561872},
// 		{X: 138.1307699285115, Y: 71.59320822977331},
// 		{X: 137.9449971804252, Y: 0.2128870011787197},
// 		{X: 103.3410329784794, Y: -0.372899436447425},
// 	}

// 	fitTolerance := 0.2
// 	basePlateThickness := 3.0
// 	wallThickness := 5.0
// 	keyboardThickness := 32.0

// 	rotationAxisVector := v3.Vec{Y: 5, Z: 3}
// 	rotationAxisY := circs[2].Loc.Y / 2
// 	// rotationAngle := sdf.Tau / 10
// 	rotationAngle := 0.0

// 	screwHoleRadius := 2.2

// 	footRadius := 3.0

// 	footLength := 39.0 //distance between circle centers
// 	footThickness := 2.4
// 	feet := []sdf.SDF2{

// 		footBetween(circs[1].Loc, circs[2].Loc, footRadius+fitTolerance, footLength),
// 		footBetween(circs[2].Loc, circs[4].Loc, footRadius+fitTolerance, footLength),
// 		footBetween(circs[4].Loc, circs[5].Loc, footRadius+fitTolerance, footLength),
// 		footAt(screwLocations[4].Add(screwLocations[5]).DivScalar(2), circs[4].Loc.Sub(screwLocations[4].Add(screwLocations[5]).DivScalar(2)), footRadius+fitTolerance, footLength),
// 	}

// 	bottom2D, err := makeProfile(circs, fitTolerance)
// 	if err != nil {
// 		panic(err)
// 	}

// 	bottomOutling2D, err := makeProfile(circs, wallThickness)
// 	if err != nil {
// 		panic(err)
// 	}

// 	bottom3D := extrudeFromTo(bottom2D, basePlateThickness, keyboardThickness+basePlateThickness)
// 	bottomOutline3D := extrudeFromTo(bottomOutling2D, 0, keyboardThickness+basePlateThickness)
// 	plate := sdf.Difference3D(
// 		bottomOutline3D,
// 		bottom3D,
// 	)

// 	plate, err = sdf.Loft3D(bottomOutling2D, bottom2D, wallThickness, 0)
// 	if err != nil {
// 		panic(err)
// 	}

// 	screwHoles := make([]sdf.SDF3, len(screwLocations))
// 	for i, sl2D := range screwLocations {
// 		sl3D := v3.Vec{X: sl2D.X, Y: sl2D.Y, Z: basePlateThickness}
// 		sh, err := sphereAt(sl3D, screwHoleRadius)
// 		if err != nil {
// 			panic(err)
// 		}
// 		screwHoles[i] = sh
// 	}

// 	feet3D := make([]sdf.SDF3, len(feet))
// 	for i, f2D := range feet {
// 		feet3D[i] = extrudeFromThickness(f2D, basePlateThickness, -footThickness)
// 	}
// 	plate = sdf.Difference3D(plate, sdf.Union3D(screwHoles...))
// 	plate = sdf.Difference3D(plate, sdf.Union3D(feet3D...))
// 	platesApart := 20.0 //distance plates are from eachother
// 	plate = sdf.Transform3D(
// 		plate,
// 		sdf.Translate3d(v3.Vec{Y: rotationAxisY}).Mul(
// 			sdf.Rotate3d(rotationAxisVector, rotationAngle).Mul(
// 				sdf.Translate3d(v3.Vec{X: platesApart / 2, Y: -rotationAxisY}),
// 			),
// 		),
// 	)
// 	plateLeft := sdf.Transform3D(plate, sdf.MirrorYZ())

// 	plates := sdf.Union3D(
// 		plate,
// 		plateLeft,
// 		cylinderFromTo(v3.Vec{Y: rotationAxisY}, v3.Vec{Y: rotationAxisY}.Add(rotationAxisVector.MulScalar(10.0)), 3, .5),
// 	)
// 	_ = plates
// 	// render.ToSTL(plates, "loftTest.stl", render.NewMarchingCubesOctree(1500))
// 	//render.ToSTL(plates, "loftTest.stl", render.NewMarchingCubesUniform(1500))

// 	profileCache := Cache2DFunc(profileExtrude(circs, wallThickness, wallThickness))
// 	extrudeTest := ExtrudeBy2DFunction(profileCache.GetShapeAt, wallThickness, v3.Vec{}.AddScalar(wallThickness))
// 	render.ToSTL(extrudeTest, "profileExtrudeTest.stl", render.NewMarchingCubesUniform(1500))
// 	fmt.Println(profileCache)
// }

// func profileExtrude(circs []circ, height, expand float64) func(float64) sdf.SDF2 {
// 	return func(z float64) sdf.SDF2 {
// 		k := Clamp(1-(z+height/2)/height, 0, 1)
// 		// mix the 2D SDFs
// 		a := Mix(0, expand, k)
// 		p, err := makeProfile(circs, a)
// 		if err != nil {
// 			panic(err)
// 		}
// 		return p
// 	}
// }

// func makeProfile(circsIn []circ, expand float64) (sdf.SDF2, error) {
// 	circs := make([]circ, len(circsIn))
// 	for i := range circs {
// 		circs[i].Loc = circsIn[i].Loc
// 		circs[i].Concave = circsIn[i].Concave
// 		if circs[i].Concave {
// 			circs[i].R = circsIn[i].R - expand
// 		} else {
// 			circs[i].R = circsIn[i].R + expand
// 		}
// 	}
// 	circs2Ds := make([]sdf.SDF2, 0, len(circs))
// 	circs2DHoles := make([]sdf.SDF2, 0)
// 	polyPoints := make([]v2.Vec, len(circs)*2)
// 	for i := range circs {
// 		cFrom := &circs[i]
// 		cTo := &circs[(i+1)%len(circs)]
// 		va, vb, err := tangentPoints(*cFrom, *cTo)
// 		if err != nil {
// 			fmt.Println("errored on: ", i)
// 			panic(err)
// 		}
// 		polyPoints[i*2] = va
// 		polyPoints[i*2+1] = vb
// 		cFrom.CutFrom = va
// 		cTo.CutTo = vb
// 	}
// 	// fmt.Printf("%+v", circs)
// 	for _, c := range circs {
// 		if !c.Concave {
// 			circs2Ds = append(circs2Ds, sdf.Cut2D(circAt(c), c.CutFrom, c.CutTo.Sub(c.CutFrom)))
// 			// circs2Ds = append(circs2Ds, circAt(c))
// 		} else {
// 			circs2DHoles = append(circs2DHoles, circAt(c))
// 		}
// 		// lines2Ds[i] = lineFromTo(circs[i].Loc, circs[(i+1)%len(circs)].Loc)
// 	}
// 	p, err := sdf.Polygon2D(polyPoints)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bottom2D := sdf.Union2D(sdf.Union2D(circs2Ds...), p)
// 	return sdf.Difference2D(bottom2D, sdf.Union2D(circs2DHoles...)), nil
// }

// func tangentPoints(c1, c2 circ) (v2.Vec, v2.Vec, error) {
// 	crunch := func(p, cc circ, whichPerpendicular float64, useRightTangent, swapEnds bool) (v2.Vec, v2.Vec, error) {
// 		tr, t, err := circleTangents(cc, p.Loc)
// 		if err != nil {
// 			return v2.Vec{}, v2.Vec{}, err
// 		}
// 		if useRightTangent {
// 			tr, t = t, tr
// 		}
// 		d := sdf.Rotate2d(whichPerpendicular).MulPosition(t.Sub(p.Loc)).Normalize().MulScalar(p.R)
// 		va := p.Loc.Add(d)
// 		vb := t.Add(d)
// 		if swapEnds {
// 			va, vb = vb, va
// 		}
// 		return va, vb, nil

// 	}
// 	var p, cc circ
// 	var whichPerpendicular float64
// 	var useRightTangent, swapEnds bool
// 	//var innerTangent bool

// 	if !c1.Concave && !c2.Concave {
// 		if c1.R < c2.R {
// 			p = c1
// 			cc = circ{Loc: c2.Loc, R: c2.R - c1.R}
// 			whichPerpendicular = sdf.Tau / 4
// 			useRightTangent = false
// 			swapEnds = false

// 		} else {
// 			p = c2
// 			cc = circ{Loc: c1.Loc, R: c1.R - c2.R}
// 			whichPerpendicular = -sdf.Tau / 4
// 			useRightTangent = true
// 			swapEnds = true
// 		}
// 	} else if !c1.Concave && c2.Concave {
// 		p = c1
// 		cc = circ{Loc: c2.Loc, R: c2.R + c1.R}
// 		whichPerpendicular = sdf.Tau / 4
// 		useRightTangent = true
// 		swapEnds = false
// 	} else if c1.Concave && !c2.Concave {
// 		p = c1
// 		cc = circ{Loc: c2.Loc, R: c2.R + c1.R}
// 		whichPerpendicular = -sdf.Tau / 4
// 		useRightTangent = false
// 		swapEnds = false
// 	} else if c1.Concave && c2.Concave {
// 		if c1.R < c2.R {
// 			p = c1
// 			cc = circ{Loc: c2.Loc, R: c2.R - c1.R}
// 			whichPerpendicular = -sdf.Tau / 4
// 			useRightTangent = true
// 			swapEnds = false

// 		} else {
// 			p = c2
// 			cc = circ{Loc: c1.Loc, R: c1.R - c2.R}
// 			whichPerpendicular = sdf.Tau / 4
// 			useRightTangent = false
// 			swapEnds = true
// 		}
// 	}
// 	return crunch(p, cc, whichPerpendicular, useRightTangent, swapEnds)

// }

// func circAt(c circ) sdf.SDF2 {
// 	c2D, _ := sdf.Circle2D(c.R)
// 	return sdf.Transform2D(c2D, sdf.Translate2d(c.Loc))
// }

// func sphereAt(loc v3.Vec, r float64) (sdf.SDF3, error) {
// 	s, err := sdf.Sphere3D(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sdf.Transform3D(s, sdf.Translate3d(loc)), nil
// }

// func cylinderFromTo(bottom, top v3.Vec, radius, round float64) sdf.SDF3 {
// 	height := top.Sub(bottom).Length()
// 	c, err := sdf.Cylinder3D(height, radius, round)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return sdf.Transform3D(
// 		c,
// 		sdf.Translate3d(bottom).Mul(
// 			sdf.RotateToVector(v3.Vec{Z: 1}, top.Sub(bottom)).Mul(
// 				sdf.Translate3d(v3.Vec{Z: height / 2}),
// 			),
// 		),
// 	)
// }
// func extrudeFromTo(s2 sdf.SDF2, fromZ, toZ float64) sdf.SDF3 {
// 	if fromZ > toZ {
// 		fromZ, toZ = toZ, fromZ
// 	}
// 	return sdf.Transform3D(
// 		sdf.Extrude3D(s2, toZ-fromZ),
// 		sdf.Translate3d(v3.Vec{Z: (toZ-fromZ)/2 + fromZ}),
// 	)
// }
// func extrudeFromThickness(s2 sdf.SDF2, fromZ, thickness float64) sdf.SDF3 {
// 	return extrudeFromTo(s2, fromZ, fromZ+thickness)
// }
// func footAt(start, direction v2.Vec, r, length float64) sdf.SDF2 {
// 	circs := []circ{
// 		{Loc: start, R: r},
// 		{Loc: start.Add(direction.Normalize().MulScalar(length)), R: r},
// 	}
// 	foot, err := makeProfile(circs, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return foot
// }

// func footBetween(p1, p2 v2.Vec, r, length float64) sdf.SDF2 {
// 	center := p1.Add(p2.Sub(p1).DivScalar(2))
// 	q1 := center.Add(center.Sub(p1).Normalize().MulScalar(length / 2))
// 	q2 := center.Add(center.Sub(p2).Normalize().MulScalar(length / 2))
// 	circs := []circ{
// 		{Loc: q1, R: r},
// 		{Loc: q2, R: r},
// 	}
// 	foot, err := makeProfile(circs, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return foot
// }

// func lineFromTo(from, to v2.Vec) sdf.SDF2 {
// 	s := to.Sub(from)
// 	l := sdf.Line2D(s.Length(), 0.1)
// 	return sdf.Transform2D(l,
// 		sdf.Translate2d(from.Add(s.DivScalar(2))).Mul(
// 			sdf.Rotate2d(math.Atan2(s.Y, s.X)),
// 		),
// 	)
// }

// //circleTangents returns the two tangets of the given circle that intersect the given point
// //an error is returned if the point is within the circle
// //the case of the point being on the edge of the circle is also considered this error (perhaps erroneously)
// func circleTangents(c circ, p v2.Vec) (v2.Vec, v2.Vec, error) {
// 	p = p.Sub(c.Loc)

// 	pm := p.Length()

// 	// if p is inside the circle, there ain't no tangents.
// 	if pm <= c.R {
// 		// fmt.Println("pm <= c.R", pm, c.R, " equals ", pm == c.R)
// 		return v2.Vec{}, v2.Vec{}, errors.New("Can't find tangent: Point is inside the circle")
// 	}

// 	a := c.R * c.R / pm
// 	q := c.R * math.Sqrt((pm*pm)-(c.R*c.R)) / pm

// 	pN := p.DivScalar(pm)
// 	pNP := v2.Vec{X: -pN.Y, Y: pN.X}
// 	va := pN.MulScalar(a)

// 	tanPosA := va.Add(pNP.MulScalar(q))
// 	tanPosB := va.Sub(pNP.MulScalar(q))

// 	return tanPosA.Add(c.Loc), tanPosB.Add(c.Loc), nil
// }

// /* this algorythm borrowed from
// https://discussions.unity.com/t/finding-a-tangent-vector-from-a-given-point-and-circle/221943/2
// not clear on how it works but it's the kind of solution I was looking for and works perfectly
//   bool CircleTangents_2(Vector2 center, float r, Vector2 p, ref Vector2 tanPosA, ref Vector2 tanPosB) {
//     p -= center;

//     float P = p.magnitude;

//     // if p is inside the circle, there ain't no tangents.
//     if (P <= r) {
//       return false;
//     }

//     float a = r * r                                          / P;
//     float q = r * (float)System.Math.Sqrt((P * P) - (r * r)) / P;

//     Vector2 pN  = p / P;
//     Vector2 pNP = new Vector2(-pN.y, pN.x);
//     Vector2 va  = pN * a;

//     tanPosA = va + pNP * q;
//     tanPosB = va - pNP * q;

//     tanPosA += center;
//     tanPosB += center;

//     return true;
//   }*/
