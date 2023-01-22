package main

import (
	"log"
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
	X                     = 345.0
	Y                     = 230.0
	SIDE_THICKNESS        = 8.0
	HOLE_DIAMETER         = 5.94
	PANEL_THICKNESS       = 8.0
	TRIANGLE_LENGTHHEIGHT = SIDE_THICKNESS + 4.0
)

func main() {
	//	render.ToSTL(sdf.Extrude3D(panel2D(150, 60), PANEL_THICKNESS), "test_lr.stl", dc.NewDualContouringDefault(600))
	//	render.ToSTL(sdf.Extrude3D(panel2D(100, 60), PANEL_THICKNESS), "test_fb.stl", dc.NewDualContouringDefault(600))
	//	render.ToSTL(sdf.Extrude3D(panel2D(150, 100), PANEL_THICKNESS), "test_tb.stl", dc.NewDualContouringDefault(600))
	render.ToSTL(Corner(), "corner.stl", render.NewMarchingCubesUniform(400))
}

// 3D

func Corner() sdf.SDF3 {
	thickness := PANEL_THICKNESS / 3
	corner2D := corner2D()
	//	cornerBottomLength := math.Sqrt(math.Pow(SIDE_THICKNESS, 2) + math.Pow(SIDE_THICKNESS, 2))
	edgeTriangles2D := make([]sdf.SDF2, 3)
	edgeTriangles2D[0] = cornerTriangle(thickness, 90)
	edgeTriangles2D[1] = cornerTriangle(thickness, 180)
	edgeTriangles2D[2] = cornerTriangle(thickness, 180)
	corner := sdf.Extrude3D(corner2D, thickness)
	edgeTriangles := make([]sdf.SDF3, 3)
	edgeTriangles[0] = sdf.Extrude3D(edgeTriangles2D[0], TRIANGLE_LENGTHHEIGHT)
	edgeTriangles[1] = sdf.Extrude3D(edgeTriangles2D[1], TRIANGLE_LENGTHHEIGHT)
	edgeTriangles[2] = sdf.Extrude3D(edgeTriangles2D[2], TRIANGLE_LENGTHHEIGHT)

	corner = sdf.Cut3D(corner, v3.Vec{X: 0, Y: 6, Z: 0}, v3.Vec{X: -1, Y: -1, Z: 1})

	//cornerTop := sdf.Extrude3D(cornerTop2D, thickness)
	corner = sdf.Union3D(
		corner,
		sdf.Transform3D(sdf.Transform3D(
			sdf.Transform3D(corner, sdf.RotateX(sdf.DtoR(90))),
			sdf.RotateY(sdf.DtoR(90))),
			sdf.Translate3d(v3.Vec{
				X: 0,
				Y: -corner2D.BoundingBox().Size().Y/2 - thickness/2,
				Z: -corner2D.BoundingBox().Size().Y/2 - thickness/2})),
		sdf.Transform3D(sdf.Transform3D(sdf.Transform3D(corner, sdf.MirrorXY()),
			sdf.RotateY(sdf.DtoR(90))),
			sdf.Translate3d(v3.Vec{
				X: -corner2D.BoundingBox().Size().X/2 - thickness/2,
				Y: 0,
				Z: -corner2D.BoundingBox().Size().Y/2 - thickness/2})),
		sdf.Transform3D(sdf.Transform3D(edgeTriangles[0], sdf.RotateX(sdf.DtoR(90))),
			sdf.Translate3d(v3.Vec{
				X: -corner2D.BoundingBox().Size().X/2 - edgeTriangles2D[0].BoundingBox().Size().X/2,
				Y: corner2D.BoundingBox().Size().Y/2 - TRIANGLE_LENGTHHEIGHT/2,
				Z: 0,
			})),
		sdf.Transform3D(sdf.Transform3D(edgeTriangles[1], sdf.RotateY(sdf.DtoR(90))),
			sdf.Translate3d(v3.Vec{
				X: corner2D.BoundingBox().Size().X/2 - TRIANGLE_LENGTHHEIGHT/2,
				Y: -corner2D.BoundingBox().Size().Y/2 - (thickness / 2),
				Z: 0,
			})),
		sdf.Transform3D(edgeTriangles[2], sdf.Translate3d(v3.Vec{
			X: -corner2D.BoundingBox().Size().X/2 - (thickness / 2),
			Y: -corner2D.BoundingBox().Size().Y/2 - (thickness / 2),
			Z: -corner2D.BoundingBox().Size().X/2 - (TRIANGLE_LENGTHHEIGHT / 2.25)})),
	)

	return corner
}

// 2D

func panel2D(l, w float64) sdf.SDF2 {
	plane := sdf.Box2D(v2.Vec{X: l, Y: w}, 0)
	cutout := sdf.Box2D(v2.Vec{X: l - SIDE_THICKNESS, Y: w - SIDE_THICKNESS}, 0)
	panel := sdf.Difference2D(plane, cutout)

	lattice := lattice(cutout.BoundingBox().Size(), math.Ceil(l/20), math.Ceil(w/20))

	// corners
	corner := make([]sdf.SDF2, 0, 4)
	corner = append(corner, cornerTriangle(TRIANGLE_LENGTHHEIGHT, 0))   // bottom left
	corner = append(corner, cornerTriangle(TRIANGLE_LENGTHHEIGHT, 90))  // bottom right
	corner = append(corner, cornerTriangle(TRIANGLE_LENGTHHEIGHT, 180)) // top right
	corner = append(corner, cornerTriangle(TRIANGLE_LENGTHHEIGHT, 270)) // top left
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

func corner2D() sdf.SDF2 {
	side := sdf.Box2D(v2.Vec{X: SIDE_THICKNESS, Y: TRIANGLE_LENGTHHEIGHT}, 0)
	triangleTR := cornerTriangle(TRIANGLE_LENGTHHEIGHT, 0)
	triangleBL := cornerTriangle(SIDE_THICKNESS, 180)
	side = sdf.Union2D(
		side,
		sdf.Transform2D(triangleTR, sdf.Translate2d(v2.Vec{X: SIDE_THICKNESS/2 + triangleTR.BoundingBox().Size().X/2, Y: 0.0})),
		sdf.Transform2D(sdf.Transform2D(side, sdf.Rotate2d(sdf.DtoR(90))), sdf.Translate2d(v2.Vec{X: SIDE_THICKNESS/2 + (TRIANGLE_LENGTHHEIGHT / 2), Y: -TRIANGLE_LENGTHHEIGHT/2 - (SIDE_THICKNESS / 2)})),
		sdf.Transform2D(triangleBL, sdf.Translate2d(v2.Vec{X: 0, Y: -TRIANGLE_LENGTHHEIGHT/2 - SIDE_THICKNESS/2})),
	)
	hole, err := sdf.Circle2D(HOLE_DIAMETER / 2)
	if err != nil {
		log.Fatalln(err)
	}
	hole = sdf.Transform2D(hole, sdf.Translate2d(v2.Vec{X: (SIDE_THICKNESS / 2) + (HOLE_DIAMETER / 4), Y: -(TRIANGLE_LENGTHHEIGHT / 2) + (HOLE_DIAMETER / 4)}))
	side = sdf.Difference2D(side, hole)

	return sdf.Center2D(side)
}

// cornerTriangle returns a triangle for the corner, rotation is the rotation of the triangle
func cornerTriangle(lengthheight, rotation float64) sdf.SDF2 {
	dimensions := []v2.Vec{
		{X: 0, Y: 0},
		{X: lengthheight, Y: 0},
		{X: 0, Y: lengthheight},
	}
	triangle, err := sdf.Polygon2D(dimensions)
	if err != nil {
		log.Fatalln(err)
	}
	triangle = sdf.Center2D(triangle)
	return sdf.Transform2D(triangle, sdf.Rotate2d(sdf.DtoR(rotation)))
}

func baseTriangle(lengthheight, rotation float64) sdf.SDF2 {
	dimensions := []v2.Vec{
		{X: -lengthheight / 2, Y: 0},
		{X: lengthheight / 2, Y: 0},
		{X: 0, Y: lengthheight},
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
