package main

import (
	//"math"

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
	shape sdf.SDF3
	cache map[v2.Vec]columnCache
}

type columnCache struct {
	mutex sync.RWMutex
}

func Drape3D(shape sdf.SDF3) sdf.SDF3 {
	d := DrapeSDF3{}
	d.shape = shape
	d.cache = make(map[v2.Vec]columnCache)
	d.mutex = sync.RWMutex{}
	return &d
}

func (d *DrapeSDF3) BoundingBox() sdf.Box3 {
	return d.shape.BoundingBox()
}

func (d *DrapeSDF3) Evaluate(p v3.Vec) float64 {
	p2 := v2.Vec{X: p.X, Y: p.Y}
	d.mutex.RLock()
	c, ok := d.cache[p2]
	d.mutex.RUnlock()
	e := d.shape.Evaluate(p)
	if !ok || e < c {
		d.mutex.Lock()
		d.cache[p2] = e
		d.mutex.Unlock()
		return e
	} else {
		return c
	}
}

func TestDrape() sdf.SDF3 {

	b, err := sdf.Box3D(v3.Vec{X: 4, Y: 4, Z: 4}, .5)
	if err != nil {
		panic(err)
	}
	c, err := sdf.Sphere3D(1.5)
	if err != nil {
		panic(err)
	}
	c = sdf.Transform3D(
		c,
		sdf.Translate3d(v3.Vec{X: 5, Y: 5, Z: 5}),
	)
	u := sdf.Union3D(b, c)

	return Drape3D(u)
}
