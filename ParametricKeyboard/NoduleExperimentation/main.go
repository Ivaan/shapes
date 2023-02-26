package main

import (
	"fmt"

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
		switchFlatzoneWidth:              16,
		keycapWidth:                      18.5,
		keycapLength:                     18.6,
		keycapBottomHeightAbovePlateDown: 3,
		keycapHeight:                     13,
		keycapBottomHeightAbovePlateUp:   6.4,
		keycapClearanced:                 2.5,
		keycapRound:                      1.2,
		huggingCylinderRound:             0.6,
		laneWidth:                        19,
		insertLength:                     8.0 + 0.4, // per memory (and mostly a guess)
		insertDiameter:                   2.9 + 0.4, // 0.4 per experiment 62f4f30
		insertWallThickness:              2.0 + 0.3, // 0.3 per experiment 62f4f30
		screwThreadDiameter:              2.0 + 0.4, // 0.4 per experiment 62f4f30
		screwThreadLength:                12.0,
		screwHeadDiameter:                3.8 + 0.4, // 0.4 per experiment 62f4f30,
	}

	cols := []NoduleSource{
		Column{ //H
			offset:       sdf.V3{X: -19.1},
			splayAngle:   0,
			convexAngle:  0,
			numberOfKeys: 3,
			startAngle:   -20,
			startRadius:  60,
			endAngle:     75,
			endRadius:    85,
			keySpacing:   19.4,
			columnType:   LeftColumn,
		},
		Column{ //J
			offset:       sdf.V3{},
			splayAngle:   0,
			convexAngle:  0,
			numberOfKeys: 3,
			startAngle:   -20,
			startRadius:  60,
			endAngle:     75,
			endRadius:    85,
			keySpacing:   19.4,
			columnType:   MiddleColumn,
		},
		Column{ //K
			offset:       sdf.V3{X: 23},
			splayAngle:   5,
			convexAngle:  .5,
			numberOfKeys: 3,
			startAngle:   -20,
			startRadius:  65,
			endAngle:     75,
			endRadius:    95,
			keySpacing:   19.4,
			columnType:   MiddleColumn,
		},
		Column{ //L
			offset:       sdf.V3{X: 47},
			splayAngle:   10,
			convexAngle:  1,
			numberOfKeys: 3,
			startAngle:   -20,
			startRadius:  62.5,
			endAngle:     75,
			endRadius:    90,
			keySpacing:   19.4,
			columnType:   MiddleColumn,
		},
		Column{ //;
			offset:       sdf.V3{X: 77, Y: -4},
			splayAngle:   25,
			convexAngle:  3,
			numberOfKeys: 3,
			startAngle:   -20,
			startRadius:  55,
			endAngle:     75,
			endRadius:    70,
			keySpacing:   19.4,
			columnType:   MiddleColumn,
		},
		// Column{ //'
		// 	offset:       sdf.V3{X: 94.28675532291557, Y: -12.060946391727217, Z: -0.9996167642402273},
		// 	splayAngle:   25,
		// 	convexAngle:  3,
		// 	numberOfKeys: 4,
		// 	startAngle:   -20,
		// 	startRadius:  55,
		// 	endAngle:     75,
		// 	endRadius:    70,
		// 	keySpacing:   19.4,
		// 	columnType:   RightColumn,
		// },
		ConeRow{
			offsetToPoint:    sdf.V3{X: -20, Y: -132, Z: -24},
			centerLine:       sdf.V3{X: -45, Y: 92, Z: 5},
			firstKeyLocation: sdf.V3{X: -25, Y: -58, Z: -45},
			rowType:          OnlyRow,
			numberOfKeys:     2,
			keySpacing:       24,
		},
	}

	a := sdf.V3{X: 19.1}
	a = sdf.RotateZ(sdf.DtoR(-25)).Mul(sdf.RotateY(sdf.DtoR(3))).MulPosition(a)
	b := sdf.V3{X: 77, Y: -4}

	fmt.Println(b.Add(a))

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

	bubbleKeys := make(map[int64]KeyNodule)
	getBubbleKey := func(screwPossitionsBits int64) KeyNodule {
		k, ok := bubbleKeys[screwPossitionsBits]
		if !ok {
			k = knp.MakeBubbleKey(screwPossitionsBits)
			bubbleKeys[screwPossitionsBits] = k
		}
		return k
	}

	for i, p := range points {
		if p.noduleType == NoduleKey {
			topNodules[i] = getBubbleKey(p.screwPossitionsBits).Top.OrientAndMove(p.moveTo)
			bottomNodules[i] = getBubbleKey(p.screwPossitionsBits).Bottom.OrientAndMove(p.moveTo)
		} else if p.noduleType == NoduleDebug1 {
			topNodules[i] = MakeNoduleDebug1().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug1().OrientAndMove(p.moveTo)
		} else if p.noduleType == NoduleDebug2 {
			topNodules[i] = MakeNoduleDebug2().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug2().OrientAndMove(p.moveTo)

		} else if p.noduleType == NoduleDebug3 {
			topNodules[i] = MakeNoduleDebug3().OrientAndMove(p.moveTo)
			bottomNodules[i] = MakeNoduleDebug3().OrientAndMove(p.moveTo)
		}
	}
	top := NoduleCollection(topNodules).Combine()
	back := sdf.Difference3D(NoduleCollection(bottomNodules).Combine(), top)

	// s1, _ := sdf.Sphere3D(2)
	// s2, _ := sdf.Sphere3D(5)
	// c, _ := sdf.Cylinder3D(5, 5, 0)
	// r := cols[len(cols)-1].(ConeRow)

	// topDebug := make([]sdf.SDF3, 0)
	// topDebug = append(topDebug,
	// 	top,
	// 	sdf.Transform3D(s1, sdf.Translate3d(r.offsetToPoint)),
	// 	sdf.Transform3D(s2, sdf.Translate3d(r.offsetToPoint.Add(r.centerLine))),
	// 	sdf.Transform3D(c, sdf.Translate3d(r.firstKeyLocation)),
	// )

	// _, debugLocations := r.getKeyLocationsWithExtras()

	// for _, v := range debugLocations {
	// 	topDebug = append(topDebug, sdf.Transform3D(s1, sdf.Translate3d(v)))
	// }

	// topPlus := sdf.Union3D(topDebug...)
	_ = top
	_ = back
	//render.RenderSTLSlow(sdf.Intersect3D(top, back), 300, "overlap.stl")
	// render.RenderSTLSlow(top, 350, "top.stl")
	// render.RenderSTLSlow(back, 300, "back.stl")

	// render.RenderSTL(top, 350, "3x5plus2top.stl")
	// render.RenderSTL(back, 300, "3x5plus2back.stl")
	render.RenderSTL(top, 1100, "3x5plus2top.stl")
	render.RenderSTL(back, 1000, "3x5plus2back.stl")

	//testing RowCones
	// row := ConeRow{
	// 	offsetToPoint:    sdf.V3{X: -20, Y: -40, Z: -14},
	// 	centerLine:       sdf.V3{X: -5, Y: 15, Z: 5},
	// 	firstKeyLocation: sdf.V3{X: -20, Y: 15, Z: -25},
	// 	rowType:          OnlyRow,
	// 	numberOfKeys:     4,
	// 	keySpacing:       20,
	// }

	// render.RenderSTL(row.getKeyLocations(), 300, "testCone.stl")

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
