package main

import (
	"fmt"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

type Column struct {
	offset       sdf.V3
	splayAngle   float64
	convexAngle  float64
	numberOfKeys int
	startAngle   float64
	startRadius  float64
	endAngle     float64
	endRadius    float64
	keySpacing   float64
	columnType   ColumnType
}

type ColumnType int

const (
	// left most column
	LeftColumn ColumnType = iota
	// middle column
	MiddleColumn
	// right most column
	RightColumn
)

type NoduleTypeAndPoint struct {
	moveTo     sdf.M44 //the transformation to move a nodule into position
	noduleType int     //the nodule type (the grid of which screws/holes need to be there)
}

// func (col Column) getColumnNodule(makeBubbleKey func(sdf.M44) KeyNodule) ColumnNodule {
// 	points := col.getKeyLocations()
// 	nodes := make([]KeyNodule, len(points))

// 	for i, p := range points {
// 		nodes[i] = makeBubbleKey(p)
// 	}
// 	return ColumnNodule{keys: nodes}
// }

// 0 1 2    0
// 3 4 5   3 1
// 6 7 8    2
func (col Column) getKeyLocations() []NoduleTypeAndPoint {
	points := spacedPointsOnAnArc(sdf.DtoR(col.startAngle), col.startRadius, sdf.DtoR(col.endAngle), col.endRadius, col.keySpacing, col.numberOfKeys)
	places := make([]NoduleTypeAndPoint, len(points))
	var firstType, middleType, lastType int
	switch col.columnType {
	case LeftColumn:
		firstType = 0
		middleType = 3
		lastType = 6
	case MiddleColumn:
		firstType = 1
		middleType = 4
		lastType = 7
	case RightColumn:
		firstType = 2
		middleType = 5
		lastType = 8
	}
	for i, p := range points {
		places[i].moveTo = sdf.Translate3d(col.offset).Mul( //offset column per knucle possition
			sdf.RotateZ(sdf.DtoR(-col.splayAngle)), //rotate column for finger splay
		).Mul(
			sdf.RotateY(sdf.DtoR(col.convexAngle)), //rotate column for row convex curve
		).Mul(
			sdf.Translate3d(sdf.V3{X: 0, Y: p.location.Y, Z: -p.location.X}), //position key on sweap of column
		).Mul(
			sdf.RotateX(p.angle), //rotate key into column sweap angle
		)
		if i == 0 {
			places[i].noduleType = lastType
		} else if i == len(points)-1 {
			places[i].noduleType = firstType
		} else {
			places[i].noduleType = middleType
		}
	}
	return places
}

func spacedPointsOnAnArc(startAngle, startRadius, endAngle, endRadius, distanceTarget float64, count int) []locationWithAngle {
	deltaAngle := (endAngle - startAngle)
	deltaRadius := (endRadius - startRadius)

	radiusPerAngleIncrement := deltaRadius / deltaAngle
	points := make([]locationWithAngle, count)

	previousArpd := angleRadiusPointDistance{}

	computePointAndistance := func(ang float64) angleRadiusPointDistance {
		angleMoved := ang - startAngle
		r := startRadius + radiusPerAngleIncrement*angleMoved
		p := sdf.P2{R: r, Theta: ang}.PolarToCartesian()
		return angleRadiusPointDistance{angle: ang, radius: r, point: p, distance: previousArpd.point.Sub(p).Length() - distanceTarget}
	}
	for i := 0; i < count; i++ {
		var arpd angleRadiusPointDistance

		if i == 0 {
			arpd = computePointAndistance(startAngle)
		} else {
			oldArpd := computePointAndistance(previousArpd.angle)
			changeAngle := math.Atan2(distanceTarget, oldArpd.radius)

			for j := 0; j < 100; j++ {
				var newAbsChangeAngle float64
				arpd = computePointAndistance(oldArpd.angle + changeAngle)
				//fmt.Println(oldArpd, arpd, changeAngle)
				if arpd.distance > 0 && oldArpd.distance < 0 || arpd.distance < 0 && oldArpd.distance > 0 {
					//crossed 0
					//half change
					newAbsChangeAngle = math.Abs(changeAngle / 2.0)
				} else if math.Abs(arpd.distance) < math.Abs(oldArpd.distance)/2.0 {
					//same side and less than half way closer
					//double change
					newAbsChangeAngle = math.Abs(changeAngle) * 2
				} else {
					//otherwise
					//keep change
					newAbsChangeAngle = math.Abs(changeAngle)
				}

				if math.Abs(arpd.distance) < math.Abs(oldArpd.distance) {
					//closer
					//use new
					oldArpd = arpd
				} else {
					//further
					//revert
					arpd = oldArpd
				}

				if arpd.distance > 0 {
					changeAngle = -newAbsChangeAngle
				} else {
					changeAngle = newAbsChangeAngle
				}

			}
		}

		run := deltaAngle * arpd.radius
		rise := deltaRadius
		angleAdjust := math.Atan2(rise, run)
		a := arpd.angle - angleAdjust
		points[i] = locationWithAngle{location: arpd.point, angle: a}
		fmt.Println(arpd.point, a)
		previousArpd = arpd
	}
	return points
}

type angleRadiusPointDistance struct {
	angle    float64
	radius   float64
	point    sdf.V2
	distance float64
}

type locationWithAngle struct {
	location sdf.V2
	angle    float64
}

type ColumnNodule struct {
	keys []KeyNodule
}

// func (cn ColumnNodule) GetTops() []sdf.SDF3 {
// 	return []sdf.SDF3{
// 		sdf.Difference3D(
// 			sdf.Union3D(AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.tops })...),
// 			sdf.Union3D(AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.topColumnHoles })...),
// 		),
// 	}
// }
// func (cn ColumnNodule) GetTopHoles() []sdf.SDF3 {
// 	return AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.topHoles })
// }
// func (cn ColumnNodule) GetBacks() []sdf.SDF3 {
// 	return AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.backs })
// }
// func (cn ColumnNodule) GetBackHoles() []sdf.SDF3 {
// 	return AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.backHoles })
// }
// func (cn ColumnNodule) GetHitBoxes() []sdf.SDF3 {
// 	return AccumulateSDF3FromKeyNodule(cn.keys, func(kn KeyNodule) []sdf.SDF3 { return kn.GetHitBoxes() })
// }

func AccumulateSDF3FromKeyNodule(kns []KeyNodule, getSDF3s func(KeyNodule) []sdf.SDF3) []sdf.SDF3 {
	totalLength := 0
	for _, kn := range kns {
		totalLength += len(getSDF3s(kn))
	}
	sdfs := make([]sdf.SDF3, totalLength)
	var i int
	for _, kn := range kns {
		i += copy(sdfs[i:], getSDF3s(kn))
	}
	return sdfs
}
