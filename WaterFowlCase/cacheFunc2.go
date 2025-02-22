//-----------------------------------------------------------------------------
/*

2D Function Evaluation Cache


*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"sync"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// CacheSDF2 is an SDF2 cache.
type CacheSDF2Func struct {
	sdfFunc     func(float64) sdf.SDF2
	cache       map[float64]sdf.SDF2
	mutex       sync.RWMutex
	reads, hits uint
}

// Cache2D wraps the passed SDF2 with an evaluation cache.
func Cache2DFunc(sdfFunc func(float64) sdf.SDF2) *CacheSDF2Func {
	return &CacheSDF2Func{
		sdfFunc: sdfFunc,
		cache:   make(map[float64]sdf.SDF2),
		mutex:   sync.RWMutex{},
	}
}

func (s *CacheSDF2Func) String() string {
	r := float64(s.hits) / float64(s.reads)
	return fmt.Sprintf("reads %d hits %d (%.2f)", s.reads, s.hits, r)
}

// Evaluate returns the minimum distance to a cached 2d sdf.
func (s *CacheSDF2Func) GetShapeAt(z float64) sdf.SDF2 {
	s.reads++
	s.mutex.RLock()
	d, ok := s.cache[z]
	s.mutex.RUnlock()
	if ok {
		s.hits++
		return d
	}
	s.mutex.Lock()
	d, ok = s.cache[z]
	if ok {
		s.hits++
		s.mutex.Unlock()
		return d
	}
	d = s.sdfFunc(z)
	s.cache[z] = d
	s.mutex.Unlock()
	return d
}

//-----------------------------------------------------------------------------
