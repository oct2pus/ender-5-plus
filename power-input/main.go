package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
	PROFILE_LENGTH = 20.0
	SOCKET_LENGTH  = 44.2 // as a note, this is mounted sideways, so length = height and height = length
	SOCKET_HEIGHT  = 33.0
	HEIGHT         = SOCKET_LENGTH + PADDING
	LENGTH         = SOCKET_HEIGHT + PADDING
	PADDING        = 10.0
	HOLE_DIAMETER  = 5.3
	HOLE_SPACING   = 36.0
	TOLERANCE      = 0.8
	ROUND          = 1.5
)

func main() {
	render.ToSTL(model(), "model.stl", render.NewMarchingCubesUniform(600))
}

// 3D

func model() sdf.SDF3 {
	thickness := 4.6
	cutoutThickness := 6.3 + 17.0

	correctionCube, err := sdf.Box3D(v3.Vec{X: 3, Y: HEIGHT, Z: thickness}, 0) // i literally dont know what to call this
	if err != nil {
		log.Fatal(err)
	}

	cutout := sdf.Extrude3D(cutout2D(), cutoutThickness)
	holder := sdf.Extrude3D(holder2D(), thickness)
	plug := sdf.Extrude3D(plug2D(), thickness)

	return sdf.Union3D(
		sdf.Transform3D(cutout, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -thickness/2 - cutoutThickness/2})),
		plug,
		sdf.Transform3D(holder, sdf.Translate3d(v3.Vec{X: -LENGTH/2 - PROFILE_LENGTH/2, Y: 0, Z: 0})),
		sdf.Transform3D(correctionCube, sdf.Translate3d(v3.Vec{X: -LENGTH / 2, Y: 0, Z: 0})),
	)

}

// 2D

func cutout2D() sdf.SDF2 {
	cutout := sdf.Box2D(v2.Vec{X: LENGTH - PADDING/2, Y: HEIGHT - PADDING/2}, ROUND)
	body := sdf.Box2D(v2.Vec{X: LENGTH - PADDING/2, Y: HEIGHT}, ROUND)

	body = sdf.Transform2D(body, sdf.Translate2d(v2.Vec{X: PADDING / 4, Y: 0}))
	return sdf.Difference2D(body, cutout)
}

func holder2D() sdf.SDF2 {
	holder := sdf.Box2D(v2.Vec{X: PROFILE_LENGTH, Y: HEIGHT}, 0)
	hole, err := sdf.Circle2D(HOLE_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}
	holes := sdf.Union2D(
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: -HOLE_SPACING / 2})),
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: HOLE_SPACING / 2})),
	)
	return sdf.Difference2D(holder, holes)
}

func plug2D() sdf.SDF2 {
	socketX, socketY := 31.2, 27.2
	plug := sdf.Box2D(v2.Vec{X: LENGTH, Y: HEIGHT}, ROUND)
	hole, err := sdf.Circle2D(HOLE_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}
	holes := sdf.Union2D(
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: -HOLE_SPACING / 2})),
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: HOLE_SPACING / 2})),
	)
	socket := sdf.Box2D(v2.Vec{X: socketX + TOLERANCE, Y: socketY + TOLERANCE}, ROUND)

	plug = sdf.Difference2D(plug, holes)
	plug = sdf.Difference2D(plug, socket)

	return plug
}
