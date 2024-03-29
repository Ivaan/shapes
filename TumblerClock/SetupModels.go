package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//ClockSetup holds the details of the clock
type ClockSetup struct {
	Tumbler
	Bearing
	BearingHolder
	Shaft
	Spacer
	Transmission
	Gear
	Frame
	MotorMount
}

func (c *ClockSetup) computeSynthetics() ClockSetup {
	c.Tumbler = c.Tumbler.computeSynthetics()
	return *c
}

//Tumbler describe the phisical characteristics of the parts you see
//Width and spacing are configured the rest are computed
type Tumbler struct {
	FaceEdgeWidth  float64 `yaml:"faceEdgeWidth"`
	FaceEdgeHeight float64
	Spacing        float64 `yaml:"spacing"`
	Radius         float64
	CornerRound    float64
	ShortRadius    float64
}

func (t *Tumbler) computeSynthetics() Tumbler {
	t.FaceEdgeHeight = t.FaceEdgeWidth / 3
	t.Radius = t.FaceEdgeWidth / 2 / math.Cos(sdf.Tau/12)
	t.CornerRound = t.Radius * 0.05
	t.ShortRadius = math.Sqrt(t.Radius*t.Radius - (t.FaceEdgeWidth/2)*(t.FaceEdgeWidth/2))
	return *t
}

//Bearing describes the bearing used to hold tumblers
type Bearing struct {
	OD float64 `yaml:"OD"`
	//ID float64 `yaml:"ID"`
	Thickness float64 `yaml:"thickness"`
}

//BearingHolder describes the details of the part of the tumbler that holds the bearing
type BearingHolder struct {
	StopConstriction float64 `yaml:"stopConstriction"`
	Tolerance        float64 `yaml:"tolerance"`
	Thickness        float64 `yaml:"thickness"`
}

//Shaft describes the shaft that runs through the tumblers
type Shaft struct {
	OD float64 `yaml:"OD"`
}

//Spacer describes the spacer that holds tumblers appart
type Spacer struct {
	ShaftTollerance float64 `yaml:"saftTollerance"`
	//BearingTollerance := 0.1
	GapAngle  float64 `yaml:"gapAngle"`
	DiskWidth float64 `yaml:"diskWidth"`
}

//Transmission descirbes the pusher nibs and the track sizes
type Transmission struct {
	NibSize         float64 `yaml:"nibSize"`
	NibLength       float64 `yaml:"length"`
	TrackTollerance float64 `yaml:"trackTollerance"`
}

//Gear describes the details of the gears in the clock
type Gear struct {
	Thickness                float64 `yaml:"thickness"`
	backlash                 float64 `yaml:"backlash"`
	clearance                float64 `yaml:"clearance"`
	DrivenGearNumberOfTeeth  int     `yaml:"drivenGearNumberOfTeeth"`
	CouplerGearNumberOfTeeth int     `yaml:"couplerGearNumberOfTeeth"`
	ColonGearNumberOfTeeth   int     `yaml:"colonGearNumberOfTeeth"`
	MotorGearNumberOfTeeth   int     `yaml:"motorGearNumberOfTeeth"`
}

//Frame desscribes the details of the frame that holds the clock together
type Frame struct {
	Thickness                   float64 `yaml:"thickness"`
	WidthAsFractionOfShorRadius float64 `yaml:"widthAsFractionOfShorRadius"`
	ShaftHolderLength           float64 `yaml:"shaftHolderLength"`
}

type MotorMount struct {
	MotorAcross           float64 `yaml:"motorAcross"`
	MotorCornerToCorner   float64 `yaml:"motorCornerToCorner"`
	MotorShaftAllowance   float64 `yaml:"motorShaftAllowance"`
	MountThickness        float64 `yaml:"mountThickness"`
	MountDepth            float64 `yaml:"mountDepth"`
	ScrewDistanceFromEdge float64 `yaml:"screwDistanceFromEdge"`
	ScrewHoleDiameter     float64 `yaml:"screwHoleDiameter"`
}
