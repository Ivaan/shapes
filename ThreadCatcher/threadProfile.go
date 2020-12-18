package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

// ISOThread returns the 2d profile for an ISO/UTS thread.
// https://en.wikipedia.org/wiki/ISO_metric_screw_thread
// https://en.wikipedia.org/wiki/Unified_Thread_Standard
func threadProfile(
	radius float64, // radius of thread
	pitch float64, // thread to thread distance
	angle float64,
	mode string, // internal/external thread
) sdf.SDF2 {

	theta := sdf.DtoR(angle)
	h := pitch / (2.0 * math.Tan(theta))
	rMajor := radius
	r0 := rMajor - (7.0/8.0)*h

	iso := sdf.NewPolygon()
	if mode == "external" {
		rRoot := (pitch / 8.0) / math.Cos(theta)
		xOfs := (1.0 / 16.0) * pitch
		iso.Add(pitch, 0)
		iso.Add(pitch, r0+h)
		iso.Add(pitch/2.0, r0).Smooth(rRoot, 5)
		iso.Add(xOfs, rMajor)
		iso.Add(-xOfs, rMajor)
		iso.Add(-pitch/2.0, r0).Smooth(rRoot, 5)
		iso.Add(-pitch, r0+h)
		iso.Add(-pitch, 0)
	} else if mode == "internal" {
		rMinor := r0 + (1.0/4.0)*h
		rCrest := (pitch / 16.0) / math.Cos(theta)
		xOfs := (1.0 / 8.0) * pitch
		iso.Add(pitch, 0)
		iso.Add(pitch, rMinor)
		iso.Add(pitch/2-xOfs, rMinor)
		iso.Add(0, r0+h).Smooth(rCrest, 5)
		iso.Add(-pitch/2+xOfs, rMinor)
		iso.Add(-pitch, rMinor)
		iso.Add(-pitch, 0)
	} else {
		panic("bad mode")
	}
	//iso.Render("iso.dxf")
	return sdf.Polygon2D(iso.Vertices())
}
