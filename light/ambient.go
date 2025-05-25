package light

// AmbientLight is a light that emits a constant amount of light
// everywhere at once. Transformation of the light has no effect.
func AmbientLight() Light {
	return Of(AmbientKind)
}
