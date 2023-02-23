package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

type Column struct { // a column of keys for a single finger
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

type ConeRow struct { //a row of keys in the thumb cluster deffined as a path along a cone (the code described by the movement of the thumb where the end of the thumb closest to the wrist is the point of the cone)
	offsetToPoint sdf.V3
	centerLine    sdf.V3 //vector down the middle of this cone
	//distance float64 //from the point to the center of the keys (height from the cone perspective)
	//radius float64 //from the centerLine out to the keys
	//and angle somehow ...
	//or maybe we just define
	firstKeyLocation sdf.V3 //and compute radius and distance from that
	numberOfKeys     int
	keySpacing       float64
	rowType          RowType
}

type RowType int

const (
	//top most row
	TopRow RowType = iota
	//middle row
	MiddleRow
	//bottom most row
	BottomRow
	//only row (topmost and bottommost)
	OnlyRow
)

type NoduleType int

const (
	NoduleKey NoduleType = iota
	NoduleDebug1
	NoduleDebug2
	NoduleDebug3
)

type NoduleTypeAndPoint struct {
	moveTo              sdf.M44    //the transformation to move a nodule into position
	noduleType          NoduleType //the nodule type - reserved for later use (probably key vs other kind of key vs encoder, maybe mpu ...)
	screwPossitionsBits int64      //(the grid of which screws/holes need to be there)
}

type NoduleSource interface {
	getKeyLocations() []NoduleTypeAndPoint
}

//  1    2
// 2 0  4 1
//  3    8
// 06 02 03
// 04 00 01
// 12 08 09

func (col Column) getKeyLocations() []NoduleTypeAndPoint {
	points := spacedPointsOnAnArc(sdf.DtoR(col.startAngle), col.startRadius, sdf.DtoR(col.endAngle), col.endRadius, col.keySpacing, col.numberOfKeys)
	places := make([]NoduleTypeAndPoint, len(points))
	var firstType, middleType, lastType int64
	switch col.columnType {
	case LeftColumn:
		firstType = 6
		middleType = 4
		lastType = 12
	case MiddleColumn:
		firstType = 2
		middleType = 0
		lastType = 8
	case RightColumn:
		firstType = 3
		middleType = 1
		lastType = 9
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
			places[i].screwPossitionsBits = lastType
		} else if i == len(points)-1 {
			places[i].screwPossitionsBits = firstType
		} else {
			places[i].screwPossitionsBits = middleType
		}
		places[i].noduleType = NoduleKey
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
		//fmt.Println(arpd.point, a)
		previousArpd = arpd
	}
	return points
}

func (con ConeRow) getKeyLocations() []NoduleTypeAndPoint {
	/*
		offsetToPoint sdf.V3
		centerLine    sdf.V3 //vector down the middle of this cone
		//distance float64 //from the point to the center of the keys (height from the cone perspective)
		//radius float64 //from the centerLine out to the keys
		//and angle somehow ...
		//or maybe we just define
		firstKeyLocation sdf.V3 //and compute radius and distance from that
		rowType          RowType
		numberOfKeys     int
		keySpacing       float64

	*/

	toOrigin := sdf.Translate3d(con.offsetToPoint.Neg()) //translate offsetpoint to origin
	//fmt.Println("rotateToXZPlane:", sdf.RtoD(-math.Atan2(con.centerLine.Y, con.centerLine.X)))
	rotateToXZPlane := sdf.RotateZ(-math.Atan2(con.centerLine.Y, con.centerLine.X))     //rotate cone centerline about Z axis to XZ plane (i.e. the Y = zero plane)
	rotatedCenterLine := rotateToXZPlane.MulPosition(con.centerLine)                    //centerline after rotation to XZ plane
	rotateToZAxis := sdf.RotateY(-math.Atan2(rotatedCenterLine.X, rotatedCenterLine.Z)) //rotate cone centerline about Y axis to Z axis (i.e. the now X also = zero)

	movedFirstKeyLocation := rotateToZAxis.Mul(rotateToXZPlane).Mul(toOrigin).MulPosition(con.firstKeyLocation) //location of the first key relative to the cone center vector after being moved to the Z axis
	rotateFirstKeyToX := sdf.RotateZ(-math.Atan2(movedFirstKeyLocation.Y, movedFirstKeyLocation.X))             //rotate the moved First Key to the XZ plane
	firstKeyOnXZ := rotateFirstKeyToX.MulPosition(movedFirstKeyLocation)                                        //first key rotated to X axis relative to cone on Z axis

	rotateKeyToZ := sdf.RotateY(-sdf.Tau / 4).Mul(sdf.RotateZ(-sdf.Tau / 4)) //rotates a key facing up along Y axis so as to face it's bottom along X (to the inner surface of cone we're aligning things to)
	moveUpCone := sdf.Translate3d(sdf.V3{Z: firstKeyOnXZ.Length()})
	leanKeyOut := sdf.RotateY(math.Atan2(firstKeyOnXZ.X, firstKeyOnXZ.Z)) //rotation to lean a key on Z axis out to inner surface of cone at first key location

	rotateAnglePerKey := math.Atan2(con.keySpacing, firstKeyOnXZ.X) //if the keys are on the surface of a cone at a certain height, then the first key Z is that height, and the first key X is the distance from the center. The keys will describe a circle at radius X.

	putTheConeBack := toOrigin.Inverse().Mul(rotateToXZPlane.Inverse()).Mul(rotateToZAxis.Inverse().Mul(rotateFirstKeyToX.Inverse())) //reverse the transformation that brought the cone centerline to the Z axis

	places := make([]NoduleTypeAndPoint, con.numberOfKeys)
	// 06 02 03
	// 04 00 01
	// 12 08 09
	var firstType, middleType, lastType int64
	switch con.rowType {
	case TopRow:
		firstType = 6
		middleType = 2
		lastType = 3
	case MiddleRow:
		firstType = 4
		middleType = 0
		lastType = 1
	case BottomRow:
		firstType = 12
		middleType = 8
		lastType = 9
	case OnlyRow:
		firstType = 14
		middleType = 10
		lastType = 11
	}
	for i := 0; i < con.numberOfKeys; i++ {
		places[i] = NoduleTypeAndPoint{
			moveTo:              putTheConeBack.Mul(sdf.RotateZ(float64(i) * rotateAnglePerKey)).Mul(leanKeyOut).Mul(moveUpCone).Mul(rotateKeyToZ),
			screwPossitionsBits: 10, //TODO: compute screw possitions by rowType
			noduleType:          NoduleKey,
		}
		if i == 0 {
			places[i].screwPossitionsBits = lastType
		} else if i == len(places)-1 {
			places[i].screwPossitionsBits = firstType
		} else {
			places[i].screwPossitionsBits = middleType
		}
	}

	// points = append(points,
	// 	NoduleTypeAndPoint{
	// 		moveTo:     sdf.Translate3d(con.offsetToPoint),
	// 		noduleType: NoduleDebug1,
	// 	},
	// 	NoduleTypeAndPoint{
	// 		moveTo:     sdf.Translate3d(con.offsetToPoint.Add(con.centerLine)),
	// 		noduleType: NoduleDebug2,
	// 	},
	// 	NoduleTypeAndPoint{
	// 		moveTo:     sdf.Translate3d(con.firstKeyLocation),
	// 		noduleType: NoduleDebug2,
	// 	},
	// 	// NoduleTypeAndPoint{
	// 	// 	moveTo:     sdf.Identity3d(),
	// 	// 	noduleType: NoduleDebug1,
	// 	// },
	// 	// NoduleTypeAndPoint{
	// 	// 	moveTo:     sdf.Translate3d(movedFirstKeyLocation),
	// 	// 	noduleType: NoduleDebug2,
	// 	// },
	// 	// NoduleTypeAndPoint{
	// 	// 	moveTo:     sdf.Translate3d(firstKeyOnXZ),
	// 	// 	noduleType: NoduleDebug3,
	// 	// },
	// )
	return places

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
