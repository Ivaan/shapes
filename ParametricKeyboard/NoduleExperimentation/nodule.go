package main

import "github.com/deadsy/sdfx/sdf"

type Nodule interface {
	GetTops() []sdf.SDF3
	GetTopHoles() []sdf.SDF3
	GetBacks() []sdf.SDF3
	GetBackHoles() []sdf.SDF3
	GetHitBoxes() []sdf.SDF3
}

type NoduleCollection []Nodule

func (nc NoduleCollection) GetTops() []sdf.SDF3 {
	totalLength := 0
	for _, n := range nc {
		totalLength += len(n.GetTops())
	}
	tops := make([]sdf.SDF3, totalLength)
	var i int
	for _, n := range nc {
		i += copy(tops[i:], n.GetTops())
	}
	return tops
}

func (nc NoduleCollection) GetTopHoles() []sdf.SDF3 {
	totalLength := 0
	for _, n := range nc {
		totalLength += len(n.GetTopHoles())
	}
	holes := make([]sdf.SDF3, totalLength)
	var i int
	for _, n := range nc {
		i += copy(holes[i:], n.GetTopHoles())
	}
	return holes
}

func (nc NoduleCollection) GetBacks() []sdf.SDF3 {
	totalLength := 0
	for _, n := range nc {
		totalLength += len(n.GetBacks())
	}
	backs := make([]sdf.SDF3, totalLength)
	var i int
	for _, n := range nc {
		i += copy(backs[i:], n.GetBacks())
	}
	return backs
}

func (nc NoduleCollection) GetBackHoles() []sdf.SDF3 {
	totalLength := 0
	for _, n := range nc {
		totalLength += len(n.GetBackHoles())
	}
	holes := make([]sdf.SDF3, totalLength)
	var i int
	for _, n := range nc {
		i += copy(holes[i:], n.GetBackHoles())
	}
	return holes
}

func (nc NoduleCollection) GetHitBoxes() []sdf.SDF3 {
	totalLength := 0
	for _, n := range nc {
		totalLength += len(n.GetHitBoxes())
	}
	hitBoxes := make([]sdf.SDF3, totalLength)
	var i int
	for _, n := range nc {
		i += copy(hitBoxes[i:], n.GetHitBoxes())
	}
	return hitBoxes
}
