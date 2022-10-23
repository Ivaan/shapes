package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	knp := BubbleKeyNoduleProperties{
		sphereRadius:                     20.0,
		plateTopAtRadius:                 9,
		plateThickness:                   4,
		sphereThicknes:                   3,
		backCoverCutAtRadius:             0,
		switchHoleLength:                 14,
		switchHoleWidth:                  14,
		switchLatchWidth:                 4,
		switchLatchGrabThickness:         1.5,
		switchFlatzoneLength:             16,
		switchFlatzoneWidth:              15,
		keycapWidth:                      18.5,
		keycapLength:                     18.6,
		keycapBottomHeightAbovePlateDown: 3,
		keycapHeight:                     13,
		keycapBottomHeightAbovePlateUp:   6.4,
		keycapClearanced:                 2.5,
		keycapRound:                      1.2,
		huggingCylinderRound:             0.6,
		laneWidth:                        19,
		insertLength:                     8.0,
		insertDiameter:                   2.9,
		insertWallThickness:              2.0,
		screwThreadDiameter:              2.0,
		screwThreadLength:                12.0,
		screwHeadDiameter:                3.8,
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
			columnType:   LeftColumn,
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
			columnType:   MiddleColumn,
		},
		{ //K
			offset:       sdf.V3{X: 23},
			splayAngle:   5,
			convexAngle:  0,
			numberOfKeys: 4,
			startAngle:   -20,
			startRadius:  65,
			endAngle:     75,
			endRadius:    95,
			keySpacing:   19.1,
			columnType:   RightColumn,
		},
		// { //L
		// 	offset:       sdf.V3{X: 42},
		// 	splayAngle:   10,
		// 	convexAngle:  0,
		// 	numberOfKeys: 4,
		// 	startAngle:   -20,
		// 	startRadius:  62.5,
		// 	endAngle:     75,
		// 	endRadius:    90,
		// 	keySpacing:   19.1,
		// },
	}

	// //dual key
	// cols = []Column{
	// 	{ //H
	// 		offset:       sdf.V3{X: -19.1},
	// 		splayAngle:   0,
	// 		convexAngle:  0,
	// 		numberOfKeys: 2,
	// 		startAngle:   0,
	// 		startRadius:  60,
	// 		endAngle:     95,
	// 		endRadius:    85,
	// 		keySpacing:   19.1,
	// 	},
	// }

	// //single key
	// cols = []Column{
	// 	{
	// 		offset:       sdf.V3{},
	// 		splayAngle:   0,
	// 		convexAngle:  0,
	// 		numberOfKeys: 1,
	// 		startAngle:   0,
	// 		startRadius:  60,
	// 		endAngle:     95,
	// 		endRadius:    85,
	// 		keySpacing:   19.1,
	// 	},
	// }

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

	points := make([]NoduleTypeAndPoint, 0)
	for _, col := range cols {
		points = append(points, col.getKeyLocations()...)
	}

	topNodules := make([]Nodule, len(points))
	bottomNodules := make([]Nodule, len(points))

	// 0 1 2    1
	// 3 4 5   2 0
	// 6 7 8    3
	bubbleKeys := make([]KeyNodule, 9)
	bubbleKeys[0] = knp.MakeBubbleKey([]int{1, 2})
	bubbleKeys[1] = knp.MakeBubbleKey([]int{1})
	bubbleKeys[2] = knp.MakeBubbleKey([]int{0, 1})
	bubbleKeys[3] = knp.MakeBubbleKey([]int{2})
	bubbleKeys[4] = knp.MakeBubbleKey([]int{})
	bubbleKeys[5] = knp.MakeBubbleKey([]int{0})
	bubbleKeys[6] = knp.MakeBubbleKey([]int{2, 3})
	bubbleKeys[7] = knp.MakeBubbleKey([]int{3})
	bubbleKeys[8] = knp.MakeBubbleKey([]int{0, 3})

	for i, p := range points {
		topNodules[i] = bubbleKeys[p.noduleType].Top.OrientAndMove(p.moveTo)
		bottomNodules[i] = bubbleKeys[p.noduleType].Bottom.OrientAndMove(p.moveTo)
	}
	top := NoduleCollection(topNodules).Combine()
	back := NoduleCollection(bottomNodules).Combine()

	//render.RenderSTLSlow(sdf.Intersect3D(top, back), 300, "overlap.stl")
	// render.RenderSTLSlow(top, 350, "top.stl")
	// render.RenderSTLSlow(back, 300, "back.stl")
	render.RenderSTL(top, 350, "top.stl")
	render.RenderSTL(back, 300, "back.stl")

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
