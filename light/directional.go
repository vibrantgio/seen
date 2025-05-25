package light

// DirectionalLight is a light that emits light in parallel lines,
// not eminating from any single point. For these lights, only the Normal
// property is used to determine the direction of the light. This may also
// be transformed.
func DirectionalLight() Light {
	return Of(DirectionalKind)
}
