package solid

// Cube constructs an axis-aligned solid cuboid. Optional parameters are
// `Center` and `Size`, which default to `Center(0, 0, 0)` and
// `Size(2, 2, 2)`. The Size is specified a list of three numbers, one for
// each axis.
//
// Example code:
//
//	cube := solid.Cube(
//	  solid.Center(0, 0, 0),
//	  solid.Size(2, 2, 2))
func Cube(options ...Option) CSG {
	o := OptionsFrom(options)
	faces := [][]int{
		{0, 4, 6, 2},
		{1, 3, 7, 5},
		{0, 1, 5, 4},
		{2, 6, 7, 3},
		{0, 2, 3, 1},
		{4, 5, 7, 6}}
	normals := []Vector{
		{X: -1, Y: 0, Z: 0},
		{X: +1, Y: 0, Z: 0},
		{X: 0, Y: -1, Z: 0},
		{X: 0, Y: +1, Z: 0},
		{X: 0, Y: 0, Z: -1},
		{X: 0, Y: 0, Z: +1}}
	polygons := []Polygon{}
	c := o.Center
	r := o.Size.DividedBy(2)
	for i, corners := range faces {
		vertices := make([]Vertex, 0, len(corners))
		for _, corner := range corners {
			position := Vector{
				X: c.X + r.X*normals[0+corner&1/1].X,
				Y: c.Y + r.Y*normals[2+corner&2/2].Y,
				Z: c.Z + r.Z*normals[4+corner&4/4].Z,
			}
			vertices = append(vertices, Vertex{Pos: position, Normal: normals[i]})
		}
		polygons = append(polygons, PolygonFromVertices(vertices...))
	}
	return CSG(polygons)
}
