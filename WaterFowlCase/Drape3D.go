package main

import (
	// "math"

	//"github.com/deadsy/sdfx/obj"
	// "errors"
	// "fmt"
	// "math"
	// "sync"

	"sync"
	// "github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type DrapeSDF3 struct {
	shape     sdf.SDF3
	grainSize float64
	heights   map[v2.Vec]float64
	mutex     sync.RWMutex
	bb        sdf.Box3
}

func Drape3D(shape sdf.SDF3, grainSize float64) sdf.SDF3 {
	d := DrapeSDF3{}
	d.shape = shape
	d.grainSize = grainSize
	d.heights = make(map[v2.Vec]float64)
	d.mutex = sync.RWMutex{}
	d.bb = shape.BoundingBox().Enlarge(v3.Vec{X: 2 * grainSize, Y: 2 * grainSize, Z: 2 * grainSize})
	return &d
}

func (d *DrapeSDF3) BoundingBox() sdf.Box3 {
	return d.bb
}

func (d *DrapeSDF3) getHeight(col v2.Vec) float64 {
	d.mutex.RLock()
	f, ok := d.heights[col]
	d.mutex.RUnlock()
	if !ok {
		d.mutex.Lock()
		f, ok := d.heights[col]
		if !ok {
			zMax := d.shape.BoundingBox().Max.Z
			zMin := d.shape.BoundingBox().Min.Z
			f = zMin
			for z := zMax; z >= zMin; z -= d.grainSize {
				if d.shape.Evaluate(v3.Vec{X: col.X, Y: col.Y, Z: z}) < 0 {
					f = z
					break
				}
			}
			d.heights[col] = f
			d.mutex.Unlock()
			return f //this is the case where we had to compute heights because of a cache miss
		}
		d.mutex.Unlock()
		return f //this is the case where we had a cache miss but some other thread computed it while this thread waited for the lock
	}
	return f //cache hit
}

func (d *DrapeSDF3) Evaluate(p v3.Vec) float64 {
	h := d.getHeight(v2.Vec{X: p.X, Y: p.Y})
	zMin := d.shape.BoundingBox().Min.Z

	if p.Z > h || h == zMin {
		return d.shape.Evaluate(p)
	}
	if p.Z <= h && p.Z > zMin {
		return -1.0 //inside
	} else {
		return d.shape.BoundingBox().Min.Z - p.Z
	}
}

func TestDrape() sdf.SDF3 {

	b, err := sdf.Box3D(v3.Vec{X: 4, Y: 4, Z: 4}, .5)
	if err != nil {
		panic(err)
	}
	b = sdf.Transform3D(b, sdf.RotateX(sdf.Tau/6))
	c, err := sdf.Sphere3D(5.5)
	if err != nil {
		panic(err)
	}
	c = sdf.Transform3D(
		c,
		sdf.Translate3d(v3.Vec{X: 5, Y: 5, Z: 5}),
	)
	u := sdf.Union3D(b, c)

	return Drape3D(u, .3)
}
