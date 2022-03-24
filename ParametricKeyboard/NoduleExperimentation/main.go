package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	knp := KeyNoduleProperties{
		sphereRadius:             20.0,
		sphereCut:                9,
		plateThickness:           4,
		sphereThicknes:           3,
		backCoverkcut:            11,
		switchHoleLength:         14,
		switchHoleWidth:          14,
		switchLatchWidth:         4, //todo:measure
		switchLatchGrabThickness: 1.5,
		switchFlatzoneWidth:      15,
		switchFlatzoneHeight:     16,
		keycapWidth:              18.5,
		keycapHeight:             18.6,
		keycapRound:              2,
		keycapOffset:             7.2,
	}

	col1 := Column{
		offset:       sdf.V3{},
		splayAngle:   0,
		convexAngle:  0,
		numberOfKeys: 4,
		startAngle:   -20,
		startRadius:  60,
		endAngle:     75,
		endRadius:    85,
		keySpacing:   19.1,
	}
	col2 := Column{
		offset:       sdf.V3{X: 19.1},
		splayAngle:   0,
		convexAngle:  0,
		numberOfKeys: 5,
		startAngle:   -20,
		startRadius:  70,
		endAngle:     75,
		endRadius:    104,
		keySpacing:   19.1,
	}

	points := col1.getKeyLocations()
	points = append(points, col2.getKeyLocations()...)
	nodes := make([]Nodule, len(points))

	var err error
	for i, p := range points {
		nodes[i], err = knp.MakeKey(p)
		if err != nil {
			panic(err)
		}

	}

	nodesC := NoduleCollection(nodes)
	//nodes = NoduleCollection([]Nodule{node1})

	top := sdf.Difference3D(sdf.Union3D(nodesC.GetTops()...), sdf.Union3D(nodesC.GetTopHoles()...))
	back := sdf.Difference3D(sdf.Union3D(nodesC.GetBacks()...), sdf.Union3D(nodesC.GetBackHoles()...))

	render.RenderSTLSlow(top, 300, "top.stl")
	render.RenderSTLSlow(back, 300, "back.stl")

}

/*
circomferance = 2 pi radius
circomferance = tau radius
arc = angle tau radius

crude slope over arc:
rise / run
run: delta angle * average radius
rise: detla radius
*/

/*
Distance calculation

angleMoved = (angle-startAngle)
radius = startRadius + radiusPerAngleIncrement*angleMoved
x2 = radius * COS(angle)
y2 = radius * SIN(angle)

distance = SQRT( (x2-x)^2 + (y2-y)^2 )

x2 = (startRadius + radiusPerAngleIncrement*(angle-startAngle)) * COS(angle)
y2 = (startRadius + radiusPerAngleIncrement*(angle-startAngle)) * SIN(angle)

distance = SQRT( ((startRadius + radiusPerAngleIncrement*(angle-startAngle)) * COS(angle)-x)^2 + ((startRadius + radiusPerAngleIncrement*(angle-startAngle)) * SIN(angle)-y)^2 )

*/
