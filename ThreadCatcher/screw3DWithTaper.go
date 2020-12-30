package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

// ScrewSDF3WithTaper is a 3d screw form with a taper.
type ScrewSDF3WithTaper struct {
	thread         sdf.SDF2 // 2D thread profile
	pitch          float64  // thread to thread distance
	lead           float64  // distance per turn (starts * pitch)
	length         float64  // total length of screw
	starts         int      // number of thread starts
	taperLength    float64  // the amount of thread that is tapered
	taperAtHeight  float64  //the height at which the taper begins
	maxTaperAmount float64  // the maximum difference between thread and tapered thread
	bb             sdf.Box3 // bounding box
}

// Screw3DWithTaper returns a screw SDF3.
func Screw3DWithTaper(
	thread sdf.SDF2, // 2D thread profile
	length float64, // length of screw
	pitch float64, // thread to thread distance
	starts int, // number of thread starts (< 0 for left hand threads)
	taperLength float64, // the amount of thread that is tapered
	maxTaperAmount float64, // the maximum difference between thread and tapered thread
) sdf.SDF3 {
	s := ScrewSDF3WithTaper{}
	s.thread = thread
	s.pitch = pitch
	s.length = length / 2
	s.lead = -pitch * float64(starts)
	s.taperLength = taperLength
	s.taperAtHeight = length/2 - taperLength
	s.maxTaperAmount = maxTaperAmount
	// Work out the bounding box.
	// The max-y axis of the sdf2 bounding box is the radius of the thread.
	bb := s.thread.BoundingBox()
	r := bb.Max.Y
	s.bb = sdf.Box3{sdf.V3{-r, -r, -s.length}, sdf.V3{r, r, s.length}}
	return &s
}

// Evaluate returns the minimum distance to a 3d screw form.
func (s *ScrewSDF3WithTaper) Evaluate(p sdf.V3) float64 {
	// map the 3d point back to the xy space of the profile
	p0 := sdf.V2{}
	// the distance from the 3d z-axis maps to the 2d y-axis
	p0.Y = math.Sqrt(p.X*p.X + p.Y*p.Y)
	if p.Z > s.taperAtHeight {
		p0.Y += (p.Z - s.taperAtHeight) * s.maxTaperAmount / s.taperLength
	}
	// the x/y angle and the z-height map to the 2d x-axis
	// ie: the position along thread pitch
	theta := math.Atan2(p.Y, p.X)
	z := p.Z + s.lead*theta/sdf.Tau
	p0.X = sdf.SawTooth(z, s.pitch)
	// get the thread profile distance
	d0 := s.thread.Evaluate(p0)
	// create a region for the screw length
	d1 := sdf.Abs(p.Z) - s.length
	// return the intersection
	return sdf.Max(d0, d1)
}

// BoundingBox returns the bounding box for a 3d screw form.
func (s *ScrewSDF3WithTaper) BoundingBox() sdf.Box3 {
	return s.bb
}
