// Package affine contains fake projections with affine transforms
//
// It is not possible to exactly render text in a scene with a perspective
// projection because Canvas and SVG support only affine transformations. So,
// in order to fake it, we create an affine transform that approximates the
// linear effects of a perspective projection on an unrendered planar face
// that represents the text's shape. We can use this transform directly in the
// text painter to warp the text.
//
// This fake projection will produce unrealistic results with large strings of
// text that are not broken into their own shapes.
package affine
