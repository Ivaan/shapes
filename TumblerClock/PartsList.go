package main

import (
	"regexp"
	"strings"
)

type partsList struct {
	topFrame    bool
	bottomFrame bool
	colonGear   bool
	tumblers    []tumblerPart
}

type tumblerPart struct {
	hours   bool        //other wise minute
	tens    bool        //otherwise units or ones
	tumbler tumblerKind //include faces otherwise only gear
	gear    bool        //include gear otherwise only tumbler (one of tumbler or gear must be true)
}

func parsePartsString(partsString string) partsList {
	var parts partsList
	re := regexp.MustCompile(`^(h|m)(t|u)(a|b)?(g)?$`)

	bits := strings.Split(partsString, ",")
	for _, bit := range bits {
		bit = strings.Trim(bit, " ")
		switch bit {
		case "topFrame":
			parts.topFrame = true
		case "bottomFrame":
			parts.bottomFrame = true
		case "colonGear":
			parts.colonGear = true
		default:
			m := re.FindStringSubmatch(bit)
			if m != nil {
				var tp tumblerPart
				tp.hours = m[1] == "h"
				tp.tens = m[2] == "t"
				switch m[3] {
				case "a":
					tp.tumbler = aFace
				case "b":
					tp.tumbler = bFace
				default:
					tp.tumbler = none
				}
				tp.gear = m[4] == "g"
				parts.tumblers = append(parts.tumblers, tp)
			}
		}
	}
	return parts
}

type tumblerKind int

//Speed options
const ( //(top, middle, bottom tumblers in a digit are aFace, second and forth are bFace)
	none tumblerKind = iota
	aFace
	bFace
)

var tumblerKindEnglish = [...]string{
	"none",
	"aFace",
	"bFace",
}

func (t tumblerKind) String() string { return tumblerKindEnglish[t] }
