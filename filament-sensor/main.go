package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

const (
	BODY_X          = 43.5
	BODY_Y          = 38.1
	BODY_Z          = 13.2
	ROUND           = 0.5
	SCREW_SPACING_X = 14.7
	SCREW_SPACING_Y = 27.8
)

func main() {
	render.ToDXF(bodyTop2D(), "top.dxf", render.NewMarchingSquaresUniform(600))
	render.ToSTL(body(), "body.stl", render.NewMarchingCubesUniform(600))
}

// 3D

func body() sdf.SDF3 {
	return sdf.Extrude3D(bodyTop2D(), 13.2)
}

// 2D
func bodyTop2D() sdf.SDF2 {
	body := sdf.Box2D(v2.Vec{X: BODY_X, Y: BODY_Y}, 0.5)
	// spacing 14.7 X 27.8 Y
	screwHole, _ := sdf.Circle2D(3.75 / 2)
	screwHoles := sdf.Union2D(
		sdf.Transform2D(screwHole, sdf.Translate2d(v2.Vec{X: SCREW_SPACING_X / 2, Y: SCREW_SPACING_Y / 2})),
		sdf.Transform2D(screwHole, sdf.Translate2d(v2.Vec{X: -SCREW_SPACING_X / 2, Y: SCREW_SPACING_Y / 2})),
		sdf.Transform2D(screwHole, sdf.Translate2d(v2.Vec{X: -SCREW_SPACING_X / 2, Y: -SCREW_SPACING_Y / 2})),
		sdf.Transform2D(screwHole, sdf.Translate2d(v2.Vec{X: SCREW_SPACING_X / 2, Y: -SCREW_SPACING_Y / 2})),
	)
	// this is for my sanity, far right pegs should be 25.0 from left side
	moveBy := 25.0
	screwHoles = sdf.Center2D(screwHoles)

	screwHoles = sdf.Transform2D(screwHoles, sdf.Translate2d(v2.Vec{X: -BODY_X / 2, Y: 0}))
	screwHoles = sdf.Transform2D(screwHoles, sdf.Translate2d(v2.Vec{X: moveBy / 2, Y: 0}))
	return sdf.Difference2D(body, screwHoles)
}

func bodyBottom2D() sdf.SDF2 {
	return nil
}
