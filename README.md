# seen

A 3-D scene graph and software renderer in pure Go, for
[Vibrant Gio](https://github.com/vibrantgio), a design system for native desktop
applications on macOS, Windows and Linux, written in pure Go on
[Gio](https://gioui.org). Ported from Gabe Lerner's CoffeeScript library
[seen.js](https://github.com/themadcreator/seen).

Gio draws two dimensions. If you want a rotating solid, a motion-captured
skeleton, or a slowly deforming surface behind a native window, the choices are
to hand-roll the projection with `clip.Path` or to bring a GPU stack into a
process that did not need one. seen is the third choice: it runs the whole
pipeline — transform, project, light, shade, depth-order — on the CPU in Go, and
then emits the result as ordinary Gio `clip` and `paint` operations into the
`op.Ops` it was handed. The output is a component. It composites with the rest of
the tree exactly like a button does, honours the layout constraints it is given,
and its library packages add no dependency beyond Gio and the
`golang.org/x/image` Gio already requires.

The root module is the strongest statement of that: `go.mod` has no `require`
block at all and `go.sum` is empty. Scene graph, matrices, quaternions,
shading, BSP trees, Wavefront `.obj` and BVH motion capture are the standard
library and nothing else. Gio enters only in the nested `context/gio` module.

Being a software renderer is the trade. There is no GPU, no shader compilation
and no driver to fail; there is also a polygon budget, and the depth ordering is
a painter's algorithm rather than a z-buffer, so what you get is a
polygon-accurate render of a modest scene rather than a fast render of a large
one.

## Where it sits

Outside the tier table — a support library. The
[organization page](https://github.com/vibrantgio) has the full tier table.

```sh
go get github.com/vibrantgio/seen
go get github.com/vibrantgio/seen/context/gio
```

Two modules. `github.com/vibrantgio/seen` at the root, on Go 1.25.1 with no
dependencies; `github.com/vibrantgio/seen/context/gio` nested, whose library packages
add gioui.org v0.10.2; its examples also require `vibrantgio/noise`,
`zserge/lorca` and `golang.org/x/exp`. Nested-module tags carry the directory as a prefix:
`context/gio/vX.Y.Z`.

## Packages

| Package | |
| --- | --- |
| `seen` | The scene graph. `Scene` holds a `Group`, a `camera.Camera`, a `viewport.Viewport` and a `shader.Shader`; `NewDefaultScene` adds three-point lighting. `FitCenter` and `FitOrigin` set up the view for a given pixel rectangle. `Group` nests, holds lights, and holds `Node`s. |
| `shape` | Geometry: `Cube`, `Sphere(subdivisions)`, `Icosahedron`, `Rectangle`, `Patch`, `Pipe`, `Extrude`, `Text`. |
| `face` | The faces a shape is built from: `Face` and `Faces`, `Facet` and `Facets` (the point indices of a face), and `Coordinates`. |
| `solid` | Constructive solid geometry — `Union`, `Subtract`, `Intersect` over BSP trees, with `Cube`, `Sphere` and `Cylinder` primitives. An adaptation of [csg](https://github.com/vibrantgio/csg), not a dependency on it: the algorithm is the same, the types are seen's, so a solid from `NewSolid` is a `seen.Object` that drops straight into a `Group`. |
| `obj`, `bvh`, `mocap` | Loading. `obj.Parse` reads Wavefront meshes; `bvh` parses Biovision motion capture; `mocap` turns a BVH hierarchy into a seen skeleton. |
| `point`, `matrix`, `quaternion`, `transform`, `affine`, `float` | The maths. `Point` is three `float64`s; `Matrix` is 4×4; `quaternion` has both plain and dual quaternions; `transform.Transform` is the translate/rotate/scale block every node embeds. |
| `camera`, `projection`, `viewport` | The view pipeline. `projection` has `Perspective`, `Ortho` and `Frustum`; `viewport` is purely the NDC-to-pixel mapping. |
| `light`, `intensity`, `shader`, `color` | Shading. `light` has ambient, directional and point lights; `shader` has the shading functions `AmbientShade`, `FlatShade`, `PhongDiffuseShade` and `PhongDiffuseAndSpecularShade`, each also exported as a `Shader` value (`DiffusePhong` for `PhongDiffuseShade`, `DiffuseAndSpecularPhong` and its short form `Phong` for `PhongDiffuseAndSpecularShade`), with `Default` the diffuse-and-specular one; `color` adds HSL, named colours, and the random and drifting `Source`s that give an untextured mesh its per-face colour. |
| `layer` | Painter's-algorithm ordering, three ways. `nsort` is the Newell–Newell–Sancha view-dependent sort with cycle cutting, and the one the examples use. `zsort` sorts on barycentre z and never splits. `bsort` builds a BSP tree and splits; `layer/bsort/bsp` is that tree: `Plane`, `NewTree` and `Compare`. `backdrop` is a solid rounded rectangle behind them. |
| `canvas` | The backend contract: `Canvas`, `PathPainter`, `RectPainter`, `CirclePainter`, `TextPainter`. Implement these and seen renders to your target. |
| `context` | `Context` has five methods — `SetLayers`, `Render`, `Animate`, `Drag`, `Zoom`. `context/svg` writes SVG; `context/gio` supplies `Widget`. |
| `animation`, `drag`, `zoom` | Interaction. `drag` has inertia; `zoom` maps scroll deltas to a scale factor. |

## Usage

Every scene is the same four steps: build geometry, add it to a scene, wrap the
scene in a layer, hand the layer to a context. A randomly coloured sphere you
can spin with the pointer:

```go
shape := shape.Sphere(2)
shape.SetScale(size, size, size)
shape.Faces().SetColorFrom(color.NewDriftingSourceWith(color.Opacity(0.9)))

scene := seen.NewDefaultScene()
scene.ShowBackfaces = true
scene.Group.Add(shape)

context.SetLayers(nsort.NewLayerForScene(scene))

context.Animate().OnBefore(func(t, dt time.Duration) {
	dtms := float64(dt.Milliseconds())
	shape.RotX(dtms * 1e-4).RotY(0.7 * dtms * 1e-4)
}).Start()

context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
	shape.RotX(e.Dy / 150).RotY(e.Dx / 150)
	context.Render()
})

return gio.Widget(context, func(w, h unit.Dp) {
	scene.FitCenter(0, 0, float64(w), float64(h))
})
```

`gio.Widget` is the whole Gio integration. It returns a `layout.Widget`; the
callback it takes runs once per frame with the component's own size in `unit.Dp`,
which is where `FitCenter` belongs — the scene refits itself to whatever the
layout gave it. `context.Render()` is `window.Invalidate()`, nothing more:
seen never draws outside a frame.

`gio.Widget` adds a scale transform to the ops it renders into and leaves it
there. A scene drawn behind other content is therefore isolated from the op
list it renders into:

```go
view := seengio.Widget(ctx, func(w, h unit.Dp) {
	scene.FitCenter(0, 0, float64(w), float64(h), cameraDist)
})

background := func(gtx layout.Context) layout.Dimensions {
	defer op.Offset(image.Point{}).Push(gtx.Ops).Pop()
	return view(gtx)
}
```

The identity `op.Offset` push and pop discards the transform the seen component
leaves behind, so the content drawn above the background is never moved by it.

## Examples

`context/gio/example` has eleven Go programs: `combinedsolid`, `giftbox`,
`helloworld`, `mocap`, `multipleangles`, `noisysphere`, `noisywavepatch`,
`poem`, `rectangle`, `solids` and `text`. Six of them keep the seen.js
CoffeeScript demo they port beside the Go. Seven more directories — `2048`,
`audioequalizer`, `canvasvssvg`, `depthoffield`, `gallery`, `masksandeffects`
and `nbody` — hold the CoffeeScript demo and a screenshot and no Go. The
directories also hold a minified `seen.min.js` (in `gallery`) and a compiled
`main.wasm` (in `combinedsolid`).

## For coding assistants

Read the org guide before you write code against this module: the plan
root's [`AGENTS.md`](https://github.com/vibrantgio/organisation/blob/master/AGENTS.md).
The build and test commands are reached from that file.

## Status

The gaps in what ships.

- **The Gio rendering path has no tests.** No package in `context/gio` has a
  test file, including `context/gio/canvas`, which is the code that turns a
  scene into `clip` and `paint` operations. In the root module `animation`,
  `canvas`, `drag`, `face`, `intensity`, `layer`, `layer/backdrop`,
  `layer/zsort`, `light`, `obj`, `shader`, `transform`, `viewport` and `zoom`
  have no tests either.
- **Panics instead of errors.** `canvas.FontFaces` in `context/gio` panics when
  one of its embedded fonts fails to parse, and it runs when the package is
  initialised. `Points.MulB` and `Points.Clip` in `point` panic on a
  slice-length mismatch.
- **The doc of `Scene.Regenerate` points at a method that is not there.** It
  says "To flush the cache, call FlushCache()", but `Scene` has no such method.
  The only `FlushCache` in the repository is on `zsort.Layer` — a different
  type in a different package. `nsort` and `bsort` cache too, and neither
  offers any way to flush.
- **`zsort` is approximate.** It never cuts a polygon, has no tests and no
  package comment, and the order-checking harness does not run it.
- **There are no golden-image tests, and nothing to regenerate.**
  `layer/internal/ordercheck` renders three scenes — a cross, an occlusion
  cycle and a coplanar control — through the SVG context, parses the emitted
  `<path>` elements back (document order is paint order), and checks the
  last-painted polygon at sample points against a ray/plane intersection
  computed in the test. It is analytic, so there is no reference artifact and
  no `-update` flag; `go test ./...` is the whole story.
  `mocap/render_test.go` is a smoke test, not a pixel comparison.
- **The package documentation is close to absent.** `doc.go` says only
  `// Package seen`, so pkg.go.dev's front page for the module says nothing.
  Most root packages have no package comment, and neither of `context/gio`'s
  two library packages has one.
- **The LICENSE is Apache-2.0 with the template placeholder left in** —
  `Copyright [yyyy] [name of copyright owner]`, unfilled — and there is no
  NOTICE file. Neither seen.js nor csg.js is credited there, although `solid`
  is an adaptation of csg.js, which is MIT. seen.js is credited in this README
  and in comments of the `mocap` package and the `mocap` and
  `multipleangles` examples; csg.js is credited nowhere in the repository.
- **The animation clock is a package-level global.** `seen.Scheduler` and the
  interval table in `interval.go` are `var`s in the root package, so two scenes
  in one process share one scheduler. Nothing in the API lets you give a scene
  its own.
