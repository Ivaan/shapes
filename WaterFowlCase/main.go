package main

import (
	//"math"

	//"github.com/deadsy/sdfx/obj"
	"errors"
	"fmt"
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	//v3 "github.com/deadsy/sdfx/vec/v3"
)

type circ struct {
	Loc     v2.Vec
	R       float64
	Concave bool
	CutFrom v2.Vec
	CutTo   v2.Vec
}

func main() {
	circs := []circ{
		{Loc: v2.Vec{X: 28.28427124746189, Y: -25.4558441227157}, R: 6},                               //6
		{Loc: v2.Vec{X: 12, Y: -9.171572875253815}, R: 6},                                             //5
		{Loc: v2.Vec{X: 12.00000000000001, Y: 73.5999813335492}, R: 6},                                //4
		{Loc: v2.Vec{X: 69.2653918090254, Y: -122.097465432435}, R: 209.903937602407},                 //3
		{Loc: v2.Vec{X: 136.5307836180504, Y: 70.39197021156274}, R: 6},                               //2
		{Loc: v2.Vec{X: 136.5307836180513, Y: 1.627100563552649}, R: 6},                               //1
		{Loc: v2.Vec{X: 86.14330244975497, Y: -72.6729449002023}, R: 68.3000454637575, Concave: true}, //0
	}

	c := circ{Loc: v2.Vec{X: 0, Y: 0}, R: 5}
	c2D := sdf.Cut2D(circAt(c), v2.Vec{X: 5}, v2.Vec{X: -5, Y: 2})
	c3D := sdf.Extrude3D(c2D, 1)
	render.ToSTL(c3D, "cutTest.stl", render.NewMarchingCubesOctree(300))

	circs2Ds := make([]sdf.SDF2, 0, len(circs))
	circs2DHoles := make([]sdf.SDF2, 0)
	polyPoints := make([]v2.Vec, len(circs)*2)
	for i := range circs {
		cFrom := &circs[i]
		cTo := &circs[(i+1)%len(circs)]
		va, vb, err := tangentPoints(*cFrom, *cTo)
		if err != nil {
			fmt.Println("errored on: ", i)
			panic(err)
		}
		polyPoints[i*2] = va
		polyPoints[i*2+1] = vb
		cFrom.CutFrom = va
		cTo.CutTo = vb
	}
	fmt.Printf("%+v", circs)
	for _, c := range circs {
		if !c.Concave {
			circs2Ds = append(circs2Ds, sdf.Cut2D(circAt(c), c.CutFrom, c.CutTo.Sub(c.CutFrom)))
			// circs2Ds = append(circs2Ds, circAt(c))
		} else {
			circs2DHoles = append(circs2DHoles, circAt(c))
		}
		// lines2Ds[i] = lineFromTo(circs[i].Loc, circs[(i+1)%len(circs)].Loc)
	}
	p, err := sdf.Polygon2D(polyPoints)
	if err != nil {
		panic(err)
	}
	fmt.Println("len(circs2Ds)", len(circs2Ds))
	fmt.Println("len(circs2DHoles)", len(circs2DHoles))
	bottom2D := sdf.Union2D(sdf.Union2D(circs2Ds...), p)
	bottom2D = sdf.Difference2D(bottom2D, sdf.Union2D(circs2DHoles...))
	bottom3D := sdf.Extrude3D(bottom2D, 5)

	render.ToSTL(bottom3D, "bottom.stl", render.NewMarchingCubesUniform(3000))
	// c1 := circ{Loc: v2.Vec{X: 10, Y: 0}, R: 4}
	// c2 := circ{Loc: v2.Vec{X: 0, Y: -10}, R: 6}
	// c3 := circ{Loc: v2.Vec{X: -10, Y: 0}, R: 4, Concave: true}
	// c4 := circ{Loc: v2.Vec{X: 0, Y: 6}, R: 6, Concave: false}

	// va, va2, err := tangentPoints(c1, c2)
	// if err != nil {
	// 	panic(err)
	// }
	// vb, vb2, err := tangentPoints(c2, c3)
	// if err != nil {
	// 	panic(err)
	// }
	// vc, vc2, err := tangentPoints(c3, c4)
	// if err != nil {
	// 	panic(err)
	// }
	// vd, vd2, err := tangentPoints(c4, c1)
	// if err != nil {
	// 	panic(err)
	// }
	// circs2D := sdf.Union2D(circAt(c1), circAt(c2), circAt(c3), circAt(c4), lineFromTo(va2, va), lineFromTo(vb2, vb), lineFromTo(vc2, vc), lineFromTo(vd2, vd))
	// c3D := sdf.Extrude3D(circs2D, 0.1)
	// //p := sdf.Polygon2D()
	// render.ToSTL(c3D, "circ.stl", render.NewMarchingCubesUniform(500))
}

func tangentPoints(c1, c2 circ) (v2.Vec, v2.Vec, error) {
	crunch := func(p, cc circ, whichPerpendicular float64, useRightTangent, swapEnds bool) (v2.Vec, v2.Vec, error) {
		tr, t, err := circleTangents(cc, p.Loc)
		if err != nil {
			return v2.Vec{}, v2.Vec{}, err
		}
		if useRightTangent {
			tr, t = t, tr
		}
		d := sdf.Rotate2d(whichPerpendicular).MulPosition(t.Sub(p.Loc)).Normalize().MulScalar(p.R)
		va := p.Loc.Add(d)
		vb := t.Add(d)
		if swapEnds {
			va, vb = vb, va
		}
		return va, vb, nil

	}
	var p, cc circ
	var whichPerpendicular float64
	var useRightTangent, swapEnds bool
	//var innerTangent bool

	if !c1.Concave && !c2.Concave {
		if c1.R < c2.R {
			p = c1
			cc = circ{Loc: c2.Loc, R: c2.R - c1.R}
			whichPerpendicular = sdf.Tau / 4
			useRightTangent = false
			swapEnds = false

		} else {
			p = c2
			cc = circ{Loc: c1.Loc, R: c1.R - c2.R}
			whichPerpendicular = -sdf.Tau / 4
			useRightTangent = true
			swapEnds = true
		}
	} else if !c1.Concave && c2.Concave {
		p = c1
		cc = circ{Loc: c2.Loc, R: c2.R + c1.R}
		whichPerpendicular = sdf.Tau / 4
		useRightTangent = true
		swapEnds = false
	} else if c1.Concave && !c2.Concave {
		p = c1
		cc = circ{Loc: c2.Loc, R: c2.R + c1.R}
		whichPerpendicular = -sdf.Tau / 4
		useRightTangent = false
		swapEnds = false
	} else if c1.Concave && c2.Concave {
		if c1.R < c2.R {
			p = c1
			cc = circ{Loc: c2.Loc, R: c2.R - c1.R}
			whichPerpendicular = -sdf.Tau / 4
			useRightTangent = true
			swapEnds = false

		} else {
			p = c2
			cc = circ{Loc: c1.Loc, R: c1.R - c2.R}
			whichPerpendicular = sdf.Tau / 4
			useRightTangent = false
			swapEnds = true
		}
	}
	return crunch(p, cc, whichPerpendicular, useRightTangent, swapEnds)

}

func circAt(c circ) sdf.SDF2 {
	c2D, _ := sdf.Circle2D(c.R)
	return sdf.Transform2D(c2D, sdf.Translate2d(c.Loc))
}

func lineFromTo(from, to v2.Vec) sdf.SDF2 {
	s := to.Sub(from)
	l := sdf.Line2D(s.Length(), 0.1)
	return sdf.Transform2D(l,
		sdf.Translate2d(from.Add(s.DivScalar(2))).Mul(
			sdf.Rotate2d(math.Atan2(s.Y, s.X)),
		),
	)
}

func circleTangents(c circ, p v2.Vec) (v2.Vec, v2.Vec, error) {
	p = p.Sub(c.Loc)

	pm := p.Length()

	// if p is inside the circle, there ain't no tangents.
	if pm <= c.R {
		fmt.Println("pm <= c.R", pm, c.R, " equals ", pm == c.R)
		return v2.Vec{}, v2.Vec{}, errors.New("Can't find tangent: Point is inside the circle")
	}

	a := c.R * c.R / pm
	q := c.R * math.Sqrt((pm*pm)-(c.R*c.R)) / pm

	pN := p.DivScalar(pm)
	pNP := v2.Vec{X: -pN.Y, Y: pN.X}
	va := pN.MulScalar(a)

	tanPosA := va.Add(pNP.MulScalar(q))
	tanPosB := va.Sub(pNP.MulScalar(q))

	return tanPosA.Add(c.Loc), tanPosB.Add(c.Loc), nil
}

/* this algorythm borrowed from
https://discussions.unity.com/t/finding-a-tangent-vector-from-a-given-point-and-circle/221943/2
not clear on how it works but it's the kind of solution I was looking for and works perfectly
  bool CircleTangents_2(Vector2 center, float r, Vector2 p, ref Vector2 tanPosA, ref Vector2 tanPosB) {
    p -= center;

    float P = p.magnitude;

    // if p is inside the circle, there ain't no tangents.
    if (P <= r) {
      return false;
    }

    float a = r * r                                          / P;
    float q = r * (float)System.Math.Sqrt((P * P) - (r * r)) / P;

    Vector2 pN  = p / P;
    Vector2 pNP = new Vector2(-pN.y, pN.x);
    Vector2 va  = pN * a;

    tanPosA = va + pNP * q;
    tanPosB = va - pNP * q;

    tanPosA += center;
    tanPosB += center;

    return true;
  }*/
