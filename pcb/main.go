package main

import (
	"github.com/deadsy/sdfx/render"
)

const (
	// put constants that are used by multiple files here
	POST_OFFSET     = 10.0
	BOARD_X         = 160
	BOARD_Y         = 100
	POST_X          = BOARD_X - POST_OFFSET
	POST_Y          = BOARD_Y - POST_OFFSET
	INSERT_HEIGHT   = 3.81
	INSERT_DIAMETER = 5.23
	ROUND           = 0.5
)

func main() {
	//frame
	render.ToDXF(frame2D(), "frame.dxf", render.NewMarchingSquaresUniform(1200))
	render.ToSTL(support(), "support.stl", render.NewMarchingCubesUniform(600))
	render.ToSTL(frame(), "frame.stl", render.NewMarchingCubesUniform(2000))

}
