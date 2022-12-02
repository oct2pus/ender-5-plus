package main

import (
	"log"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
// move consts only related to stuff used in frame.go here
)

// 3D

func frame() sdf.SDF3 {
	mounts := sdf.Extrude3D(mounts2D(), INSERT_HEIGHT)
	frame := sdf.Extrude3D(frame2D(), INSERT_HEIGHT/2)
	lattice := sdf.Extrude3D(lattice2D(), INSERT_HEIGHT/2)
	frame = sdf.Transform3D(frame, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -INSERT_HEIGHT / 2}))
	lattice = sdf.Transform3D(lattice, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -INSERT_HEIGHT / 2}))

	return sdf.Union3D(mounts, frame, lattice)
}

// 2D

func lattice2D() sdf.SDF2 {
	base := sdf.Box2D(v2.Vec{X: POST_X, Y: FRAME_WIDTH / 2}, 0)

	lines := sdf.Union2D(
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: -POST_Y/2 + POST_Y/8})),
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: -POST_Y/2 + POST_Y/4})),
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: -POST_Y/2 + POST_Y/2.65})),
		base,
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y / 8})),
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y / 4})),
		sdf.Transform2D(base, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y / 2.65})),
	)
	diagLines := sdf.Union2D(
		lines,
		sdf.Transform2D(lines, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y / 2})),
		sdf.Transform2D(lines, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y})),
	)
	diagLines = sdf.Center2D(diagLines)
	diagLines = sdf.Elongate2D(diagLines, v2.Vec{X: 20, Y: 0})
	lines = sdf.Union2D(
		sdf.Transform2D(diagLines, sdf.Rotate2d(sdf.DtoR(30))),
		sdf.Transform2D(diagLines, sdf.Rotate2d(sdf.DtoR(330))),
	)
	lines = sdf.Intersect2D(lines, sdf.Box2D(v2.Vec{X: POST_X, Y: POST_Y}, 0))

	return lines
}

func frame2D() sdf.SDF2 {
	return sdf.Union2D(mounts2D(), walls2D())
}

func mounts2D() sdf.SDF2 {
	post := sdf.Box2D(v2.Vec{X: MOUNT_SIZE, Y: MOUNT_SIZE}, ROUND)
	insert, err := sdf.Circle2D(INSERT_DIAMETER / 2)
	if err != nil {
		log.Fatal(err)
	}
	post = sdf.Difference2D(post, insert)
	return sdf.Union2D(
		sdf.Transform2D(post, sdf.Translate2d(v2.Vec{X: POST_X / 2, Y: POST_Y / 2})),
		sdf.Transform2D(post, sdf.Translate2d(v2.Vec{X: -POST_X / 2, Y: POST_Y / 2})),
		sdf.Transform2D(post, sdf.Translate2d(v2.Vec{X: -POST_X / 2, Y: -POST_Y / 2})),
		sdf.Transform2D(post, sdf.Translate2d(v2.Vec{X: POST_X / 2, Y: -POST_Y / 2})),
	)
}

func walls2D() sdf.SDF2 {
	wallX, wallY := sdf.Box2D(v2.Vec{X: POST_X - MOUNT_SIZE, Y: FRAME_WIDTH}, 0), sdf.Box2D(v2.Vec{X: FRAME_WIDTH, Y: POST_Y - MOUNT_SIZE}, 0)
	return sdf.Union2D(
		sdf.Transform2D(wallX, sdf.Translate2d(v2.Vec{X: 0, Y: POST_Y / 2})),
		sdf.Transform2D(wallX, sdf.Translate2d(v2.Vec{X: 0, Y: -POST_Y / 2})),
		sdf.Transform2D(wallY, sdf.Translate2d(v2.Vec{X: POST_X / 2, Y: 0})),
		sdf.Transform2D(wallY, sdf.Translate2d(v2.Vec{X: -POST_X / 2, Y: 0})),
	)
}
