package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	knp := BubbleKeyNoduleProperties{
		sphereRadius:             20.0,
		sphereCut:                9,
		plateThickness:           4,
		sphereThicknes:           3,
		backCoverkcut:            11,
		switchHoleLength:         14,
		switchHoleWidth:          14,
		switchLatchWidth:         4,
		switchLatchGrabThickness: 1.5,
		switchFlatzoneLength:     16,
		switchFlatzoneWidth:      15,
		keycapWidth:              18.5,
		keycapLength:             18.6,
		keycapMinHeight:          3,
		keycapMaxHeight:          13,
		keycapClearanced:         2.5,
		keycapRound:              2,
		keycapOffset:             7.2,
		laneWidth:                19,
	}

	cols := []Column{
		{ //H
			offset:       sdf.V3{X: -19.1},
			splayAngle:   0,
			convexAngle:  0,
			numberOfKeys: 4,
			startAngle:   -20,
			startRadius:  60,
			endAngle:     75,
			endRadius:    85,
			keySpacing:   19.1,
		},
		{ //J
			offset:       sdf.V3{},
			splayAngle:   0,
			convexAngle:  0,
			numberOfKeys: 4,
			startAngle:   -20,
			startRadius:  60,
			endAngle:     75,
			endRadius:    85,
			keySpacing:   19.1,
		},
		{ //K
			offset:       sdf.V3{X: 21},
			splayAngle:   5,
			convexAngle:  0,
			numberOfKeys: 4,
			startAngle:   -20,
			startRadius:  65,
			endAngle:     75,
			endRadius:    95,
			keySpacing:   19.1,
		},
		{ //L
			offset:       sdf.V3{X: 42},
			splayAngle:   10,
			convexAngle:  0,
			numberOfKeys: 4,
			startAngle:   -20,
			startRadius:  62.5,
			endAngle:     75,
			endRadius:    90,
			keySpacing:   19.1,
		},
	}

	// points := col1.getKeyLocations()
	// points = append(points, col2.getKeyLocations()...)
	// nodes := make([]Nodule, len(points))

	// var err error
	// for i, p := range points {
	// 	nodes[i], err = knp.MakeBubbleKey(p)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// }
	colNodules := make([]Nodule, len(cols))
	for i, col := range cols {
		colNodules[i] = Nodule(col.getColumnNodule(knp.MakeBubbleKey))
	}
	nodesC := NoduleCollection(colNodules)
	//nodes = NoduleCollection([]Nodule{node1})

	top := sdf.Difference3D(sdf.Union3D(nodesC.GetTops()...), sdf.Union3D(nodesC.GetTopHoles()...))
	back := sdf.Difference3D(sdf.Union3D(nodesC.GetBacks()...), sdf.Union3D(nodesC.GetBackHoles()...))

	render.RenderSTLSlow(sdf.Intersect3D(top, back), 300, "overlap.stl")
	render.RenderSTLSlow(top, 350, "top.stl")
	render.RenderSTLSlow(back, 300, "back.stl")
	// render.RenderSTL(top, 350, "top.stl")
	// render.RenderSTL(back, 300, "back.stl")

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

/*
knp := FlatterKeyNoduleProperties{
	sphereRadius:             20.0,
	sphereCut:                7,
	plateThickness:           12,
	sphereThicknes:           3,
	backCoverLipCut:          2,
	switchHoleLength:         14,
	switchHoleWidth:          14,
	switchHoleDepth:          3.5,
	switchLatchWidth:         4,
	switchLatchGrabThickness: 1.5,
	switchFlatzoneLength:     16,
	switchFlatzoneWidth:      15,
	pcbLength:                17,
	pcbWidth:                 17,
	keycapWidth:              18.5,
	keycapHeight:             18.6,
	keycapRound:              2,
	keycapOffset:             7.2,
}
*/
