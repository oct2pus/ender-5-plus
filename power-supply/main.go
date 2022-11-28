package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

/* measurements:
adapter width: 58.0mm
adapter height: 25.0mm
adapter thickness: 5.0mm
*/

const (
	ADAPTER_WIDTH       = 58.0
	ADAPTER_HEIGHT      = 15.0
	ADAPTER_THICKNESS   = 5.0
	TOLERANCE           = 1.2
	BUFFER              = 10.0
	ROUND               = 6.0
	PROFILE_DISTANCE    = 10.0
	SCREW_HEAD_DIAMETER = 6.6 + TOLERANCE
	BUFFERED_WIDTH      = ADAPTER_WIDTH + BUFFER + TOLERANCE
	BUFFERED_HEIGHT     = ADAPTER_HEIGHT + BUFFER + TOLERANCE
)

func main() {
	render.ToDXF(holder2D(), "holder.dxf", render.NewMarchingSquaresUniform(600))
	render.ToSTL(holder(), "holder.stl", render.NewMarchingCubesUniform(800))
}

// 3D

func back() sdf.SDF3 {
	back := sdf.Extrude3D(back2D(), ADAPTER_THICKNESS*2)

	back = sdf.Difference3D(
		back,
		sdf.Transform3D(screw(), sdf.Translate3d(v3.Vec{X: PROFILE_DISTANCE, Y: 0, Z: 0})),
	)
	back = sdf.Difference3D(
		back,
		sdf.Transform3D(screw(), sdf.Translate3d(v3.Vec{X: -PROFILE_DISTANCE, Y: 0, Z: 0})),
	)
	return back
}

func holder() sdf.SDF3 {
	top := sdf.Extrude3D(front2D(), ADAPTER_THICKNESS)
	middle := sdf.Extrude3D(holder2D(), ADAPTER_THICKNESS)
	back := back()
	return sdf.Union3D(
		middle,
		sdf.Transform3D(top, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: ADAPTER_THICKNESS})),
		sdf.Transform3D(back, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -ADAPTER_THICKNESS * 1.5})),
	)
}

func screw() sdf.SDF3 {
	m4Diameter := 4.0
	headHeight := 3.8 + TOLERANCE

	screwLength := ADAPTER_THICKNESS * 2
	screwHead, err := sdf.Circle2D((SCREW_HEAD_DIAMETER) / 2)
	if err != nil {
		log.Fatal(err)
	}
	screwBody, err := sdf.Circle2D(m4Diameter / 2)
	if err != nil {
		log.Fatal(err)
	}

	head := sdf.Extrude3D(screwHead, headHeight)
	body := sdf.Extrude3D(screwBody, screwLength)

	return sdf.Union3D(
		body,
		sdf.Transform3D(head, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: screwLength/2 - headHeight/2})),
	)
}

// 2D

func back2D() sdf.SDF2 {
	return sdf.Box2D(v2.Vec{X: BUFFERED_WIDTH, Y: BUFFERED_HEIGHT}, ROUND)
}

func holder2D() sdf.SDF2 {
	base := back2D()
	adapter := sdf.Box2D(v2.Vec{X: ADAPTER_WIDTH + TOLERANCE, Y: BUFFERED_HEIGHT}, ROUND)
	adapter = sdf.Transform2D(adapter, sdf.Translate2d(v2.Vec{X: 0, Y: (BUFFERED_HEIGHT)/2 - ADAPTER_HEIGHT/2}))
	return sdf.Difference2D(base, adapter)
}

func front2D() sdf.SDF2 {
	psuWidth := 42.1
	screwHead, err := sdf.Circle2D(SCREW_HEAD_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}
	cutout := sdf.Box2D(v2.Vec{X: psuWidth + TOLERANCE, Y: BUFFERED_HEIGHT}, 0)
	front := sdf.Difference2D(back2D(), cutout)
	front = sdf.Difference2D(front, sdf.Transform2D(screwHead, sdf.Translate2d(v2.Vec{X: PROFILE_DISTANCE, Y: 0})))
	front = sdf.Difference2D(front, sdf.Transform2D(screwHead, sdf.Translate2d(v2.Vec{X: -PROFILE_DISTANCE, Y: 0})))

	return front
}
