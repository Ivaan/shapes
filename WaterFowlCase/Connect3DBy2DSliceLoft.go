package main

import (
	//"math"

	//"github.com/deadsy/sdfx/obj"
	// "errors"
	// "fmt"
	"math"
	// "sync"

	// "github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	// v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Connect3DBy2DSliceLoft3D struct {
	SdfA            sdf.SDF3
	PointInAToSlice v3.Vec
	SdfB            sdf.SDF3
	PointInBToSlice v3.Vec
	bb              sdf.Box3
	min             sdf.MinFunc
}

func Connect3DBy2DSliceLoft(sdfa, sdfb sdf.SDF3, pointInAToSlice, pointInBToSlice v3.Vec) sdf.SDF3 {
	s := Connect3DBy2DSliceLoft3D{
		SdfA:            sdfa,
		PointInAToSlice: pointInAToSlice,
		SdfB:            sdfb,
		PointInBToSlice: pointInBToSlice,
		bb:              sdfa.BoundingBox().Extend(sdfb.BoundingBox()),
		min:             math.Min,
	}
	return &s
}

func (s *Connect3DBy2DSliceLoft3D) BoundingBox() sdf.Box3 {
	return s.bb
}

func (s *Connect3DBy2DSliceLoft3D) Evaluate(p v3.Vec) float64 {
	// p1 := sdf.Slice2D(s.SdfA, s.PointInAToSlice, s.PointInBToSlice.Sub(s.PointInAToSlice))
	// p2 := sdf.Slice2D(s.SdfB,
	ap := p.Sub(s.PointInAToSlice)                 //Vector from A to P
	ab := s.PointInBToSlice.Sub(s.PointInAToSlice) //Vector from A to B

	magnitudeAB := ab.Length2()           //Magnitude of AB vector (it's length squared)
	abapProduct := ap.Dot(ab)             //The DOT product of a_to_p and a_to_b
	distance := abapProduct / magnitudeAB //The normalized "distance" from a to your closest point

	u := s.min(s.SdfA.Evaluate(p), s.SdfB.Evaluate(p))
	np := s.PointInAToSlice.Add(ab.MulScalar(distance)) // nearest point to p on line ab
	dv := p.Sub(np)                                     //delta vector, from nearest point to p
	apnt := s.PointInAToSlice.Add(dv)
	bpnt := s.PointInBToSlice.Add(dv)
	if distance < 0 { //Check if p projection is over vectorAB
		pe := p.Sub(apnt).Length()
		return s.min(u, pe)
	} else if distance > 1 {
		pe := p.Sub(bpnt).Length()
		return s.min(u, pe)
	} else {
		ad := s.SdfA.Evaluate(apnt) //compute loft
		bd := s.SdfB.Evaluate(bpnt)
		m := sdf.Mix(ad, bd, distance)
		return s.min(u, m)
	}
}

// SetMin sets the minimum function to control blending.
func (s *Connect3DBy2DSliceLoft3D) SetMin(min sdf.MinFunc) {
	s.min = min
}

func ShowConnect() sdf.SDF3 {
	b, err := sdf.Box3D(v3.Vec{X: 4, Y: 4, Z: 4}, .5)
	if err != nil {
		panic(err)
	}
	c, err := sdf.Sphere3D(1.5)
	if err != nil {
		panic(err)
	}
	c = sdf.Transform3D(
		c,
		sdf.Translate3d(v3.Vec{X: 5, Y: 5, Z: 5}),
	)
	connect := Connect3DBy2DSliceLoft(b, c, v3.Vec{X: 1, Y: 1, Z: 1}, v3.Vec{X: 4.5, Y: 4.5, Z: 4.5})

	// connect.(*Connect3DBy2DSliceLoft3D).SetMin(sdf.PolyMin(.1))
	return connect
}

func ShowDebug() sdf.SDF3 {
	o := make([]sdf.SDF3, 0)
	a := v3.Vec{X: 3, Y: 0, Z: 0}
	b := v3.Vec{X: 10, Y: 5, Z: 4}
	p := v3.Vec{X: 0, Y: 0, Z: 0}

	ap := p.Sub(a) //Vector from A to P
	ab := b.Sub(a) //Vector from A to B

	magnitudeAB := ab.Length2()           //Magnitude of AB vector (it's length squared)
	abapProduct := ap.Dot(ab)             //The DOT product of a_to_p and a_to_b
	distance := abapProduct / magnitudeAB //The normalized "distance" from a to your closest point

	// u := s.min(s.SdfA.Evaluate(p), s.SdfB.Evaluate(p))
	// if distance < 0 { //Check if p projection is over vectorAB
	// 	return u
	// } else if distance > 1 {
	// 	return u
	// } else {
	np := a.Add(ab.MulScalar(distance)) // nearest point to p on line ab
	dv := p.Sub(np)                     //delta vector, from nearest point to p
	apnt := a.Add(dv)
	bpnt := b.Add(dv)
	// ad := s.SdfA.Evaluate(ap) //compute loft
	// bd := s.SdfA.Evaluate(bp)
	// a := sdf.Mix(ad, bd, distance)
	// return s.min(u, a)
	// }

	o = append(o, debug(a, 1, .1))
	o = append(o, debug(b, 2, .1))
	o = append(o, debug(p, 3, .1))
	o = append(o, debug(np, 4, .1))
	o = append(o, debug(apnt, 5, .1))
	o = append(o, debug(bpnt, 6, .1))
	return sdf.Union3D(o...)
}

func debug(p v3.Vec, id int, scale float64) sdf.SDF3 {
	c, err := sdf.Sphere3D(scale)
	if err != nil {
		panic(err)
	}
	if id < 1 {
		panic("id of 1 or more supported")
	} else if id == 1 {
		return sdf.Transform3D(c, sdf.Translate3d(p))
	} else {
		o := make([]sdf.SDF3, id)
		for i := 0; i < id; i++ {
			o[i] = sdf.Transform3D(
				c,
				sdf.RotateZ(float64(i)*sdf.Tau/float64(id)).Mul(
					sdf.Translate3d(v3.Vec{X: scale / 2}),
				),
			)
		}
		return sdf.Transform3D(
			sdf.Union3D(o...),
			sdf.Translate3d(p),
		)
	}
}
