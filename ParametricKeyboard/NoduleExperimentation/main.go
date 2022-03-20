package main

import (
	"fmt"
	"math"

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
		switchLatchGrabThickness: 4,
		switchFlatzoneWidth:      17,
		switchFlatzoneHeight:     17,
		keycapWidth:              15,
		keycapHeight:             15,
		keycapRound:              4,
		keycapOffset:             4,
	}
	points := spacedPointsOnAnArc(-sdf.Tau/13, 100, sdf.Tau/13, 180, 20, 6)
	nodes := make([]Nodule, len(points))
	/*
		node1, err := knp.MakeKey(sdf.V3{X: 0, Y: 0, Z: 0}, sdf.V3{X: 0, Y: 0, Z: 1})
		if err != nil {
			panic(err)
		}
		node2, err := knp.MakeKey(sdf.V3{X: 30, Y: 0, Z: 0}, sdf.V3{X: 0, Y: 0, Z: 1})
		if err != nil {
			panic(err)
		}
		node3, err := knp.MakeKey(sdf.V3{X: 30, Y: 20, Z: 0}, sdf.V3{X: -.3, Y: -.2, Z: 1})
		if err != nil {
			panic(err)
		}
	*/
	var err error
	for i, p := range points {
		nodes[i], err = knp.MakeKey(sdf.Translate3d(sdf.V3{X: 0, Y: p.location.Y, Z: -p.location.X}).Mul(sdf.RotateX(p.angle)))
		if err != nil {
			panic(err)
		}

	}

	nodesC := NoduleCollection(nodes)
	//nodes = NoduleCollection([]Nodule{node1})

	top := sdf.Difference3D(sdf.Union3D(nodesC.GetTops()...), sdf.Union3D(nodesC.GetHoles()...))
	back := sdf.Difference3D(sdf.Union3D(nodesC.GetBacks()...), sdf.Union3D(nodesC.GetHoles()...))

	render.RenderSTLSlow(top, 300, "top.stl")
	render.RenderSTLSlow(back, 300, "back.stl")

}

type locationWithAngle struct {
	location sdf.V2
	angle    float64
}

func pointsOnAnArc(startAngle, startRadius, endAngle, endRadius float64, count int) []locationWithAngle {
	deltaAngle := (endAngle - startAngle)
	deltaRadius := (endRadius - startRadius)
	angleIncrement := deltaAngle / float64(count)
	radiusIncrement := deltaRadius / float64(count)
	points := make([]locationWithAngle, count)

	for i := 0; i < count; i++ {
		angle := startAngle + angleIncrement*float64(i)
		radius := startRadius + radiusIncrement*float64(i)
		p := sdf.PolarToXY(radius, angle)
		run := deltaAngle * radius
		rise := deltaRadius
		angleAdjust := math.Atan2(rise, run)
		a := angle - angleAdjust
		points[i] = locationWithAngle{location: p, angle: a}
	}
	return points
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

			for j := 0; j < 10; j++ {
				var newAbsChangeAngle float64
				arpd = computePointAndistance(oldArpd.angle + changeAngle)
				fmt.Println(oldArpd, arpd, changeAngle)
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
