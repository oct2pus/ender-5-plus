package main

import (
	"log"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
	BASE_X             = 40.0
	BASE_Y             = 60.0
	BASE_POST_OFFSET_Y = BASE_Y - POST_OFFSET
	PROFILE_DISTANCE   = 10.0
	PILLAR_HEIGHT      = BOARD_Y - 5
	TRIANGLE_HEIGHT    = BASE_X - POST_OFFSET
)

// 3D

func pillar() sdf.SDF3 {
	length := BASE_POST_OFFSET_Y - INSERT_DIAMETER/2 - POST_OFFSET
	offset := 3.0
	longInsertLength := 6.35

	bottomTri, err := triangle(length, TRIANGLE_HEIGHT, offset)
	if err != nil {
		log.Fatal(err)
	}
	topTri, err := triangle(length/2, TRIANGLE_HEIGHT/2, offset/4)
	if err != nil {
		log.Fatal(err)
	}
	mount := sdf.Extrude3D(mountingHole(), longInsertLength)

	topTri = sdf.Transform2D(topTri, sdf.Translate2d(v2.Vec{X: 0, Y: -(TRIANGLE_HEIGHT / 2) - -(TRIANGLE_HEIGHT / 2 / 2)}))

	pillar, err := sdf.Loft3D(bottomTri, topTri, PILLAR_HEIGHT, 0)
	if err != nil {
		log.Fatal(err)
	}

	mount = sdf.Transform3D(mount, sdf.RotateX(sdf.DtoR(90)))
	mount = sdf.Transform3D(mount, sdf.Translate3d(v3.Vec{X: 0, Y: -TRIANGLE_HEIGHT/2 - (-longInsertLength / 2), Z: -PILLAR_HEIGHT/2 + POST_Y + PROFILE_DISTANCE/2 - INSERT_HEIGHT}))

	return sdf.Difference3D(pillar, mount)
}

func support() sdf.SDF3 {
	base := sdf.Extrude3D(base2D(), INSERT_HEIGHT)
	pillar := pillar()

	pillar = sdf.Transform3D(pillar, sdf.RotateZ(sdf.DtoR(90)))
	pillar = sdf.Transform3D(pillar, sdf.Translate3d(v3.Vec{X: BASE_X/2 - TRIANGLE_HEIGHT/2, Y: 0, Z: 0}))

	base = sdf.Union3D(
		base,
		sdf.Transform3D(pillar, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: INSERT_HEIGHT/2 + PILLAR_HEIGHT/2})),
	)
	return base
}

// 2D

func mountingHole() sdf.SDF2 {
	hole, err := sdf.Circle2D(INSERT_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}

	return hole
}

func base2D() sdf.SDF2 {

	base := sdf.Box2D(v2.Vec{X: BASE_X, Y: BASE_Y}, ROUND)
	mount, err := sdf.Circle2D(INSERT_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}

	mounts := sdf.Union2D(
		sdf.Transform2D(mount, sdf.Translate2d(v2.Vec{X: PROFILE_DISTANCE, Y: (BASE_Y - POST_OFFSET) / 2})),
		sdf.Transform2D(mount, sdf.Translate2d(v2.Vec{X: PROFILE_DISTANCE, Y: -(BASE_Y - POST_OFFSET) / 2})),
		sdf.Transform2D(mount, sdf.Translate2d(v2.Vec{X: -PROFILE_DISTANCE, Y: -(BASE_Y - POST_OFFSET) / 2})),
		sdf.Transform2D(mount, sdf.Translate2d(v2.Vec{X: -PROFILE_DISTANCE, Y: (BASE_Y - POST_OFFSET) / 2})),
	)

	return sdf.Difference2D(base, mounts)
}

func triangle(length, height, offset float64) (sdf.SDF2, error) {
	x, y := length-offset, height-offset
	dimens := []v2.Vec{ // left to top to right
		{X: 0, Y: -offset},
		{X: 0, Y: offset},
		{X: x/2 - offset, Y: y},
		{X: x/2 + offset, Y: y},
		{X: x, Y: offset},
		{X: x, Y: -offset},
	}
	poly, err := sdf.Polygon2D(dimens)
	return sdf.Center2D(poly), err
}
