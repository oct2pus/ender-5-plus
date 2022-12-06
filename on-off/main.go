package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
	ROUND           = 1.5
	BODY_X          = 40.0
	BODY_Y          = 40.0
	PROFILE         = 20.0
	OFFSET          = 5.0
	MOUNT_THICKNESS = 4.6
	WALL_HEIGHT     = 22.0
)

func main() {
	render.ToSTL(onoff(), "onoff.stl", render.NewMarchingCubesUniform(600))
}

// 3D
func onoff() sdf.SDF3 {
	bodyThickness := 3.0
	wallHeight := 22.0
	triangleThickness := OFFSET - ROUND*1.725

	body := sdf.Extrude3D(body(), bodyThickness)
	walls := sdf.Extrude3D(walls(), wallHeight)
	mount := sdf.Extrude3D(mount(), MOUNT_THICKNESS)
	triangle := sdf.Extrude3D(triangleSides(), triangleThickness)
	support := sdf.Extrude3D(triangleSupport(), BODY_Y)

	mount = sdf.Transform3D(mount, sdf.RotateY(sdf.DtoR(90)))
	triangle = sdf.Transform3D(triangle, sdf.RotateY(sdf.DtoR(270)))
	triangle = sdf.Transform3D(triangle, sdf.RotateZ(sdf.DtoR(90)))
	support = sdf.Transform3D(support, sdf.RotateX(sdf.DtoR(90)))

	triangles := sdf.Union3D(
		sdf.Transform3D(triangle, sdf.Translate3d(v3.Vec{X: 0, Y: BODY_Y/2 - triangleThickness/2, Z: wallHeight/2 + PROFILE/2})),
		sdf.Transform3D(triangle, sdf.Translate3d(v3.Vec{X: 0, Y: -BODY_Y/2 - (-triangleThickness / 2), Z: wallHeight/2 + PROFILE/2})),
	)

	body = sdf.Transform3D(body, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -wallHeight/2 - bodyThickness/2}))
	mount = sdf.Transform3D(mount, sdf.Translate3d(v3.Vec{X: BODY_X/2 - MOUNT_THICKNESS/2, Y: 0, Z: wallHeight/2 + PROFILE/2}))
	support = sdf.Transform3D(support, sdf.Translate3d(v3.Vec{X: BODY_X/2 - MOUNT_THICKNESS/2, Y: 0, Z: 0}))

	return sdf.Union3D(body, walls, mount, triangles, support)
}

// 2D

func body() sdf.SDF2 {
	holeDiameter := 20.0

	body := sdf.Box2D(v2.Vec{X: BODY_X, Y: BODY_Y}, 0)
	hole, err := sdf.Circle2D(holeDiameter / 2)
	if err != nil {
		log.Fatal(err)
	}

	return sdf.Difference2D(body, hole)
}

func mount() sdf.SDF2 {
	holeDiameter := 5.3

	mount := sdf.Box2D(v2.Vec{X: PROFILE, Y: BODY_Y}, 0)
	hole, err := sdf.Circle2D(holeDiameter / 2)
	if err != nil {
		log.Fatal(err)
	}

	holes := sdf.Union2D(
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: OFFSET * 1.5})),
		sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: 0, Y: -OFFSET * 1.5})),
	)

	return sdf.Difference2D(mount, holes)
}

func triangleSides() sdf.SDF2 {
	coords := []v2.Vec{
		{X: 0, Y: 0},
		{X: PROFILE, Y: 0},
		{X: 0, Y: BODY_X},
	}
	triangle, err := sdf.Polygon2D(coords)
	if err != nil {
		log.Fatal(err)
	}
	triangle = sdf.Center2D(triangle)

	return triangle
}

func triangleSupport() sdf.SDF2 {
	coords := []v2.Vec{
		{X: 0, Y: 0},
		{X: MOUNT_THICKNESS, Y: 0},
		{X: MOUNT_THICKNESS, Y: -WALL_HEIGHT},
	}
	triangle, err := sdf.Polygon2D(coords)
	if err != nil {
		log.Fatal(err)
	}
	triangle = sdf.Center2D(triangle)

	return triangle
}

func walls() sdf.SDF2 {
	body := sdf.Box2D(v2.Vec{X: BODY_X, Y: BODY_Y}, 0)
	inner := sdf.Box2D(v2.Vec{X: BODY_X - OFFSET, Y: BODY_Y - OFFSET}, ROUND)

	return sdf.Difference2D(body, inner)
}
