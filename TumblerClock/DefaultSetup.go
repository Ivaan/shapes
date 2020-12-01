package main

import "github.com/deadsy/sdfx/sdf"

func makeDefaultClockSetup() ClockSetup {
	return (&ClockSetup{
		Tumbler: Tumbler{
			FaceEdgeWidth: 75.0,
			Spacing:       2.0,
		},
		Bearing: Bearing{
			OD: 22.0,
			//ID := 8.0
			Thickness: 7.0,
		},
		BearingHolder: BearingHolder{
			StopConstriction: 1.0, //horizontle and vertical chamfer distance
			Tolerance:        0.1,
			Thickness:        4.0,
		},
		Shaft: Shaft{
			OD: 8.0,
		},
		Spacer: Spacer{
			ShaftTollerance: -0.2,
			//BearingTollerance := 0.1
			GapAngle:  6.0 / 360.0 * sdf.Tau,
			DiskWidth: 2.5,
			//BearingPenetrationDepth := 0.0
		},
		Transmission: Transmission{
			NibSize:         3.0,
			NibLength:       8.5,
			TrackTollerance: 1.5,
		},
	}).computeSynthetics()
}
