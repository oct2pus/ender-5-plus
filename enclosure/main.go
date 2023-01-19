package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/render/dc"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

const (
	X               = 345.0
	Y               = 230.0
	SIDE_THICKNESS  = 8.0
	HOLE_DIAMETER   = 5.94
	PANEL_THICKNESS = 8.0
)

func main() {
	//	render.ToDXF(panel2D(100, 200), "test.dxf", render.NewMarchingSquaresQuadtree(600))
	render.ToSTL(sdf.Extrude3D(panel2D(100, 80), PANEL_THICKNESS), "test_wall1.stl", dc.NewDualContouringDefault(600))
	render.ToSTL(sdf.Extrude3D(panel2D(150, 50), PANEL_THICKNESS), "test_wall2.stl", dc.NewDualContouringDefault(600))

}

// 2D
func panel2D(l, w float64) sdf.SDF2 {
	plane := sdf.Box2D(v2.Vec{X: l, Y: w}, 0)
	cutout := sdf.Box2D(v2.Vec{X: l - SIDE_THICKNESS, Y: w - SIDE_THICKNESS}, 0)
	panel := sdf.Difference2D(plane, cutout)

	lattice := lattice(cutout.BoundingBox().Size(), l/20, w/10)

	// corners
	corner := make([]sdf.SDF2, 0, 4)
	corner = append(corner, cornerTriangle(0))   // bottom left
	corner = append(corner, cornerTriangle(90))  // bottom right
	corner = append(corner, cornerTriangle(180)) // top right
	corner = append(corner, cornerTriangle(270)) // top left
	corner[0] = sdf.Transform2D(corner[0], sdf.Translate2d(v2.Vec{X: -cutout.BoundingBox().Size().X/2 - (-corner[0].BoundingBox().Size().X / 2), Y: -cutout.BoundingBox().Size().Y/2 - (-corner[0].BoundingBox().Size().Y / 2)}))
	corner[1] = sdf.Transform2D(corner[1], sdf.Translate2d(v2.Vec{X: cutout.BoundingBox().Size().X/2 - (corner[1].BoundingBox().Size().X / 2), Y: -cutout.BoundingBox().Size().Y/2 - (-corner[1].BoundingBox().Size().Y / 2)}))
	corner[2] = sdf.Transform2D(corner[2], sdf.Translate2d(v2.Vec{X: cutout.BoundingBox().Size().X/2 - (corner[2].BoundingBox().Size().X / 2), Y: cutout.BoundingBox().Size().Y/2 - (corner[2].BoundingBox().Size().Y / 2)}))
	corner[3] = sdf.Transform2D(corner[3], sdf.Translate2d(v2.Vec{X: -cutout.BoundingBox().Size().X/2 - (-corner[3].BoundingBox().Size().X / 2), Y: cutout.BoundingBox().Size().Y/2 - (corner[3].BoundingBox().Size().Y / 2)}))
	corners := sdf.Union2D(corner...)

	// holes
	mount, err := sdf.Circle2D(HOLE_DIAMETER / 2)
	if err != nil {
		log.Fatalln(err)
	}

	hole := make([]sdf.SDF2, 0, 4)
	hole = append(hole, mount) // bottom left
	hole = append(hole, mount) // bottom right
	hole = append(hole, mount) // top right
	hole = append(hole, mount) // top left
	hole[0] = sdf.Transform2D(hole[0], sdf.Translate2d(v2.Vec{X: -cutout.BoundingBox().Size().X/2 - (-HOLE_DIAMETER / 4), Y: -cutout.BoundingBox().Size().Y/2 - (-HOLE_DIAMETER / 4)}))
	hole[1] = sdf.Transform2D(hole[1], sdf.Translate2d(v2.Vec{X: cutout.BoundingBox().Size().X/2 - (HOLE_DIAMETER / 4), Y: -cutout.BoundingBox().Size().Y/2 - (-HOLE_DIAMETER / 4)}))
	hole[2] = sdf.Transform2D(hole[2], sdf.Translate2d(v2.Vec{X: cutout.BoundingBox().Size().X/2 - (HOLE_DIAMETER / 4), Y: cutout.BoundingBox().Size().Y/2 - (HOLE_DIAMETER / 4)}))
	hole[3] = sdf.Transform2D(hole[3], sdf.Translate2d(v2.Vec{X: -cutout.BoundingBox().Size().X/2 - (-HOLE_DIAMETER / 4), Y: cutout.BoundingBox().Size().Y/2 - (HOLE_DIAMETER / 4)}))
	holes := sdf.Union2D(hole...)

	// combine parts
	panel = sdf.Union2D(panel, lattice, corners)
	panel = sdf.Difference2D(panel, holes)

	return panel
}

// cornerTriangle returns a triangle for the corner, rotation is the rotation of the triangle
func cornerTriangle(rotation float64) sdf.SDF2 {
	xyLengthHeight := SIDE_THICKNESS + 4
	dimensions := []v2.Vec{
		{X: 0, Y: 0},
		{X: xyLengthHeight, Y: 0},
		{X: 0, Y: xyLengthHeight},
	}
	triangle, err := sdf.Polygon2D(dimensions)
	if err != nil {
		log.Fatalln(err)
	}
	triangle = sdf.Center2D(triangle)
	return sdf.Transform2D(triangle, sdf.Rotate2d(sdf.DtoR(rotation)))
}

func lattice(vec v2.Vec, width, points float64) sdf.SDF2 {
	xDis, yDis := vec.X/points, vec.Y/points
	log.Printf("\nxDis, yDis: %v,%v\nvec.X, vec.Y: %v,%v\n", xDis, yDis, vec.X, vec.Y)
	lines := make([]sdf.SDF2, 0)

	//bottom to right
	for i := 0.0; i < points; i++ {
		lineVec := []v2.Vec{
			{X: xDis * i, Y: 0},
			{X: (xDis * i) + width, Y: 0},
			{X: vec.X, Y: vec.Y - (yDis * i) - width},
			{X: vec.X, Y: vec.Y - (yDis * i)},
		}
		line, err := sdf.Polygon2D(lineVec)
		if err != nil {
			log.Fatalln(err)
		}
		lines = append(lines, line)
	}
	lattice := sdf.Center2D(sdf.Union2D(lines...))
	lattice = sdf.Union2D(
		lattice,
		sdf.Transform2D(lattice, sdf.MirrorX()),
	)
	lattice = sdf.Union2D(
		lattice,
		sdf.Transform2D(lattice, sdf.Rotate2d(sdf.DtoR(180))),
	)

	return lattice
}
