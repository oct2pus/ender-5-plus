package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

const (
	WIDTH      = 20.0
	FOOT_WIDTH = WIDTH * 3
)

func main() {
	render.ToSTL(feet(), "feet.stl", render.NewMarchingCubesUniform(600))
}

//3D

func feet() sdf.SDF3 {

	feet := sdf.Extrude3D(feet2D(), WIDTH)
	holes := sdf.Extrude3D(holes2d(), WIDTH)
	holes = sdf.Transform3D(holes, sdf.RotateX(sdf.DtoR(90)))
	return sdf.Difference3D(feet, holes)
}

//2D

func feet2D() sdf.SDF2 {
	start := 4.0
	bump := 6.0

	base := sdf.Box2D(v2.Vec{X: FOOT_WIDTH, Y: start}, 0)
	foot := sdf.Box2D(v2.Vec{X: FOOT_WIDTH / 3, Y: bump}, 0)
	rightAngle, err := triangle(-bump/2, bump)
	if err != nil {
		log.Printf("%v\n", err)
		return base
	}

	foot = sdf.Transform2D(foot, sdf.Translate2d(v2.Vec{X: -FOOT_WIDTH/2 - (-FOOT_WIDTH / 6), Y: bump/2 + start/2}))
	rightAngle = sdf.Transform2D(rightAngle, sdf.Translate2d(v2.Vec{X: -FOOT_WIDTH/6 - (-bump / 4), Y: bump/2 + start/2}))

	return sdf.Union2D(base, foot, rightAngle)
}

func holes2d() sdf.SDF2 {
	m4 := 4.0
	holes, err := sdf.Circle2D(m4 / 2)
	if err != nil {
		log.Printf("%v\n", err)
		return holes
	}
	holes = sdf.Union2D(
		holes,
		sdf.Transform2D(holes, sdf.Translate2d(v2.Vec{X: -FOOT_WIDTH/2 - (-FOOT_WIDTH / 6), Y: 0})),
		sdf.Transform2D(holes, sdf.Translate2d(v2.Vec{X: FOOT_WIDTH/2 - (FOOT_WIDTH / 6), Y: 0})),
	)
	return holes
}

// helper

// triangle returns a right angle triangle
func triangle(length, height float64) (sdf.SDF2, error) {
	vecs := make([]v2.Vec, 0)
	vecs = append(vecs, v2.Vec{X: 0, Y: 0})
	vecs = append(vecs, v2.Vec{X: length, Y: 0})
	vecs = append(vecs, v2.Vec{X: length, Y: height})
	poly, err := sdf.Polygon2D(vecs)
	if err != nil {
		return poly, err
	}
	poly = sdf.Center2D(poly)
	return poly, nil
}
