package main

import (
	"fmt"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func main() {
	jackThreadsId := 7.56
	jackThreadsOd := 13.0
	jackThreadsLength := 4.76
	jackConnectionsId := 9.8
	jackConnectionsOd := 13.0
	jackConnectionsLength := 33.0 - jackThreadsLength

	ledHolderId := 14.95
	ledHolderOd := 18.95
	ledHolderLength := 6.6 + 2.1 + 1.1
	ledHolderShelfLength := 2.0

	ledRetainerThickness := 1.0
	ledRetainerLength := 6.0
	ledRetainerLipLength := 2.0

	lensHolderId := 20.25
	lensHolderOd := 24.25
	lensHolderLength := 24.0
	lensHolderLipDepth := 8.0
	lensHolderLipThickness := 2.5
	lensHolderThreadsLength := lensHolderLength - lensHolderLipDepth

	lensRetainerThickness := 1.0
	lensRetainerLength := 6.0

	threadPitch := 3.0
	threadRadius := ledHolderOd / 2
	threadTolerance := 0.16

	tpi45, err := threadProfile(threadRadius+threadTolerance, threadPitch, 45, "internal")
	if err != nil {
		panic(err)
	}

	tpe45, err := threadProfile(threadRadius-threadTolerance, threadPitch, 45, "external")
	if err != nil {
		panic(err)
	}

	jackThread, _ := sdf.Cylinder3D(jackThreadsLength, jackThreadsOd/2, 0)
	jackThreadHole, _ := sdf.Cylinder3D(jackThreadsLength, jackThreadsId/2, 0)
	jackConnections, _ := sdf.Cylinder3D(jackConnectionsLength, jackConnectionsOd/2, 0)
	jackConnectionsHole, _ := sdf.Cylinder3D(jackConnectionsLength, jackConnectionsId/2, 0)

	ledHolderShelf, _ := sdf.Cylinder3D(ledHolderShelfLength, ledHolderOd/2, 0)
	ledHolderShelfHole, _ := sdf.Cylinder3D(ledHolderShelfLength, jackConnectionsId/2, 0)

	ledHolder, _ := sdf.Cylinder3D(ledHolderLength+ledRetainerLength, threadRadius, 0)
	ledHolderHole, _ := sdf.Cylinder3D(ledHolderLength+ledRetainerLength, ledHolderId/2, 0)

	ledRetainer, _ := sdf.Cylinder3D(ledRetainerLength, ledHolderId/2, 0)
	ledRetainerHole, _ := sdf.Cylinder3D(ledRetainerLength, ledHolderId/2-ledRetainerThickness, 0)
	ledRetainerLip, _ := sdf.Cylinder3D(ledRetainerLipLength, ledHolderOd/2, 0)
	ledRetainerLipHole, _ := sdf.Cylinder3D(ledRetainerLipLength, ledHolderId/2-ledRetainerThickness, 0)

	ledStack := stackAndUnion(ledHolder, ledHolderShelf, jackConnections, jackThread)
	ledStackHole := stackAndUnion(ledHolderHole, ledHolderShelfHole, jackConnectionsHole, jackThreadHole)
	innerThreads, _ := sdf.Screw3D(
		tpe45, // 2D thread profile
		ledStack.BoundingBox().Max.Z-ledStack.BoundingBox().Min.Z, // length of screw
		threadPitch, // thread to thread distance
		1,           // number of thread starts (< 0 for left hand threads)
	)

	ledRetainerStack := stackAndUnion(ledRetainerLip, ledRetainer)
	ledRetainerStackHole := stackAndUnion(ledRetainerLipHole, ledRetainerHole)

	ledHolderFinal := sdf.Intersect3D(stackAndUnion(innerThreads), sdf.Difference3D(ledStack, ledStackHole))
	ledRetainerFinal := sdf.Difference3D(ledRetainerStack, ledRetainerStackHole)

	f, err := sdf.LoadFont("data-latin.ttf")
	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}
	logoText := sdf.NewText("NazFab")
	logo2d, err := sdf.TextSDF2(f, logoText, lensHolderLength/3)
	if err != nil {
		fmt.Printf("error rendering text %s\n", err)
		os.Exit(1)
	}
	logo3d, _ := sdf.ExtrudeRounded3D(logo2d, 10, 0.2) //hard coded 10, could be related to knerl somehow, never did work out what makes the knerl points tall vs not ...
	logo3d = sdf.Transform3D(
		logo3d,
		sdf.Translate3d(sdf.V3{X: lensHolderOd/2 + 10/2}).Mul(
			sdf.RotateX(-sdf.Tau/4).Mul(
				sdf.RotateY(sdf.Tau/4),
			),
		),
	)
	logo3d = bend3d(logo3d, lensHolderOd/2)
	lensHolderKnerl, _ := obj.KnurledHead3D(lensHolderOd/2, lensHolderLength, lensHolderOd/8)
	lensHolderKnerl = sdf.Difference3D(lensHolderKnerl, logo3d)
	lensHolderHole, _ := sdf.Cylinder3D(lensHolderLipDepth-lensHolderLipThickness, lensHolderId/2, 0)
	lensHolderLipHole, _ := sdf.Cone3D(lensHolderLipThickness, lensHolderId/2, lensHolderId/2-lensHolderLipThickness, 0)
	lensHolderThreadHoles, _ := sdf.Screw3D(
		tpi45,                   // 2D thread profile
		lensHolderThreadsLength, // length of screw
		threadPitch,             // thread to thread distance
		1,                       // number of thread starts (< 0 for left hand threads)
	)

	lensRetainer, _ := sdf.Cylinder3D(lensRetainerLength, lensHolderId/2, 0)
	lensRetainerHole, _ := sdf.Cylinder3D(lensRetainerLength-lensHolderLipThickness, lensHolderId/2-lensRetainerThickness, 0)
	lensRetainerLipHole, _ := sdf.Cone3D(lensHolderLipThickness, lensHolderId/2-lensHolderLipThickness, lensHolderId/2, 0)

	lensHoleStack := stackAndUnion(lensHolderHole, lensHolderLipHole, lensHolderThreadHoles)
	lensHolderFinal := sdf.Difference3D(stackAndUnion(lensHolderKnerl), lensHoleStack)

	lensRetainerHoleStack := stackAndUnion(lensRetainerHole, lensRetainerLipHole)
	lensRetainerFinal := sdf.Difference3D(stackAndUnion(lensRetainer), lensRetainerHoleStack)

	render.RenderSTLSlow(slice(ledHolderFinal), 400, "ledHolder.stl")
	render.RenderSTLSlow(slice(ledRetainerFinal), 400, "ledRetainer.stl")
	render.RenderSTLSlow(slice(lensHolderFinal), 400, "lensHolder.stl")
	render.RenderSTLSlow(slice(lensRetainerFinal), 400, "lensRetainer.stl")

}

func stackAndUnion(cylinders ...sdf.SDF3) sdf.SDF3 {
	translatedCylinders := make([]sdf.SDF3, len(cylinders))
	currentHeight := 0.0
	for i, c := range cylinders {
		translatedCylinders[i] = sdf.Transform3D(
			c,
			sdf.Translate3d(sdf.V3{Z: currentHeight - c.BoundingBox().Min.Z}),
		)
		currentHeight += c.BoundingBox().Max.Z - c.BoundingBox().Min.Z
	}
	return sdf.Union3D(translatedCylinders...)
}

func slice(shape sdf.SDF3) sdf.SDF3 {
	return sdf.Cut3D(shape, sdf.V3{}, sdf.V3{X: 1})
}
