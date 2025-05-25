package light

// PointLight is a Light that emits light in all directions from a single point.
// The Point property determines the location of the point light. Note,
// though, that it may also be moved through the transformation of the light.
// A point light is also called an "Omni" light in some 3D editors.
func PointLight() Light {
	return Of(PointKind)
}
