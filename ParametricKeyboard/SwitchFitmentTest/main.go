package main

import (
	"fmt"
	"os"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/golang/freetype/truetype"
)

var _font *truetype.Font

func main() {
	//Grid of test holes for testing fitment of Cherry MX Blue switches

	minSpacingBetween := 5.0 //also edge border

	thicknessStart := 1.5
	thicknessIncrement := 0.5
	thicknessSteps := 3

	widthStart := 14.0
	widthIncrement := 0.2
	heightStart := 14.0
	heightIncrement := 0.2
	sizeSteps := 3

	maxThickness := thicknessStart + thicknessIncrement*float64(thicknessSteps)
	plateHeight := heightStart + heightIncrement*float64(sizeSteps) + minSpacingBetween*2

	f, err := sdf.LoadFont("cmr10.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}
	_font = f
	lableHeight := 3.0

	holes := make([]sdf.SDF3, sizeSteps)
	labels := make([]sdf.SDF3, sizeSteps)
	currentEdge := minSpacingBetween
	for sizeIndex := 0; sizeIndex < sizeSteps; sizeIndex++ {
		width := widthStart + widthIncrement*float64(sizeIndex)
		height := heightStart + heightIncrement*float64(sizeIndex)
		hole, _ := sdf.Box3D(sdf.V3{X: width, Y: height, Z: maxThickness}, 0)
		hole = sdf.Transform3D(hole, sdf.Translate3d(sdf.V3{X: currentEdge + width/2, Y: plateHeight / 2, Z: maxThickness / 2}))
		holes[sizeIndex] = hole
		label := sdf.Transform3D(lable(fmt.Sprint(width, height), lableHeight), sdf.Translate3d(sdf.V3{X: currentEdge + width/2}))
		labels[sizeIndex] = label
		currentEdge += width + minSpacingBetween
	}
	holes3D := sdf.Union3D(holes...)
	labels3D := sdf.Union3D(labels...)

	plates := make([]sdf.SDF3, thicknessSteps)
	for thicknessIndex := 0; thicknessIndex < thicknessSteps; thicknessIndex++ {
		thickness := thicknessStart + thicknessIncrement*float64(thicknessIndex)
		plate, _ := sdf.Box3D(sdf.V3{X: currentEdge, Y: plateHeight, Z: thickness}, 0)
		plate = sdf.Transform3D(plate, sdf.Translate3d(sdf.V3{X: currentEdge / 2, Y: plateHeight / 2, Z: thickness / 2}))
		plate = sdf.Difference3D(plate, holes3D)
		plate = sdf.Difference3D(plate, sdf.Transform3D(labels3D, sdf.Translate3d(sdf.V3{Z: thickness})))
		plate = sdf.Transform3D(plate, sdf.Translate3d(sdf.V3{Y: plateHeight * float64(thicknessIndex)}))
		plates[thicknessIndex] = plate
	}

	plate3D := sdf.Union3D(plates...)
	_ = plate3D
	render.RenderSTLSlow(plate3D, 800, "testPlate.stl")
	//render.RenderSTLSlow(holes3D, 800, "holes.stl")

}

func lable(label string, height float64) sdf.SDF3 {
	fmt.Println(label)
	labelText := sdf.NewText(label)
	label2d, err := sdf.TextSDF2(_font, labelText, height)
	if err != nil {
		fmt.Printf("error rendering text %s\n", err)
		os.Exit(1)
	}
	//label3d, _ := sdf.Loft3D(sdf.Offset2D(label2d, -1.0), label2d, 1.0, 0)
	label3d, _ := sdf.ExtrudeRounded3D(label2d, 2, 0.2)
	label3d = sdf.Transform3D(
		label3d,
		sdf.Translate3d(sdf.V3{Y: height / 2}),
	)
	return label3d
}
