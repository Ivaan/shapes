// wallaby camshaft

package main

import (
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	c := sdf.Cylinder3D(20, 10, 0)
	t := sdf.Translate3d((sdf.V3{X: 10, Y: 0, Z: 10}))
	negC := sdf.Transform3D(c, t)

	model := sdf.Difference3D(c, negC)

	sdf.RenderSTL(NewGrowSDF3(model, -3.0), 300, "twoCylinders.stl")
	//sdf.RenderSTL(c, 300, "twoCylinders.stl")

}

//GrowSDF3 is an experiment on single distance functions (a learning experience)
type GrowSDF3 struct {
	sdf    sdf.SDF3
	growBy float64
	bb     sdf.Box3
}

//GrowSDF3 creates a new GrowSDF3
func NewGrowSDF3(sdfIn sdf.SDF3, growBy float64) sdf.SDF3 {
	s := GrowSDF3{}
	s.sdf = sdfIn
	s.growBy = growBy
	bb := sdfIn.BoundingBox()
	s.bb = sdf.Box3{bb.Min.AddScalar(-growBy), bb.Max.AddScalar(growBy)}
	return &s
}

//Evaluate implements sdf
func (s *GrowSDF3) Evaluate(p sdf.V3) float64 {
	return s.sdf.Evaluate(p) - s.growBy
}

//BoundingBox implements sdf
func (s *GrowSDF3) BoundingBox() sdf.Box3 {
	return s.bb
}
