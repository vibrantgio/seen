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
`op.Ops` it was handed. The output is a widget. It composites with the rest of
the tree exactly like a button does, honours the layout constraints it is given,
and adds no dependency beyond Gio itself.

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

Outside ADR-001's tier table — a support library the design system consumes and
never depends on. The [organization page](https://github.com/vibrantgio) has the
full tier table.

No tier of the design system imports seen. Its two consumers are
`workbench/launcher`, whose animated background is a seen scene, and
[svg](https://github.com/vibrantgio/svg)'s `driver/seen`, which turns SVG paths
into seen geometry.

```sh
go get github.com/vibrantgio/seen
go get github.com/vibrantgio/seen/context/gio
```

Two modules. `github.com/vibrantgio/seen` at the root, on Go 1.25.1 with no
dependencies; `github.com/vibrantgio/seen/context/gio` nested, which adds
gioui.org v0.10.2 like the rest of the organization. Nested-module tags carry
the directory as a prefix: `context/gio/v0.0.7`, not `v0.0.7`.

## Packages

Thirty-three packages in the root module, fourteen more in `context/gio` — of
which twelve are example commands.

| Package | |
| --- | --- |
| `seen` | The scene graph. `Scene` holds a `Group`, a `camera.Camera`, a `viewport.Viewport` and a `shader.Shader`; `NewDefaultScene` adds three-point lighting. `FitCenter` and `FitOrigin` set up the view for a given pixel rectangle. `Group` nests, carries lights, and holds `Node`s. |
| `shape` | Geometry: `Cube`, `Sphere(subdivisions)`, `Icosahedron`, `Rectangle`, `Patch`, `Pipe`, `Extrude`, `Text`. |
| `solid` | Constructive solid geometry — `Union`, `Subtract`, `Intersect` over BSP trees, with `Cube`, `Sphere` and `Cylinder` primitives. An adaptation of [csg](https://github.com/vibrantgio/csg), not a dependency on it: the algorithm is the same, the types are seen's, so a `*Solid` is a `seen.Object` that drops straight into a `Group`. |
| `obj`, `bvh`, `mocap` | Loading. `obj.Parse` reads Wavefront meshes; `bvh` parses Biovision motion capture; `mocap` turns a BVH hierarchy into a seen skeleton. |
| `point`, `matrix`, `quaternion`, `transform`, `affine`, `float` | The maths. `Point` is three `float64`s; `Matrix` is 4×4; `quaternion` has both plain and dual quaternions; `transform.Transform` is the translate/rotate/scale block every node embeds. |
| `camera`, `projection`, `viewport` | The view pipeline. `projection` has `Perspective`, `Ortho` and `Frustum`; `viewport` is purely the NDC-to-pixel mapping. |
| `light`, `intensity`, `shader`, `color` | Shading. `light` has ambient, directional and point lights; `shader` has `Ambient`, `Flat`, `Phong` and the two diffuse variants; `color` adds HSL, named colours, and the random and drifting `Source`s that give an untextured mesh its per-face colour. |
| `layer` | Painter's-algorithm ordering, three ways. `nsort` is the Newell–Newell–Sancha view-dependent sort with cycle cutting — the one everything actually uses. `zsort` sorts on barycentre z and never splits. `bsort` builds a BSP tree and splits. `backdrop` is a solid rounded rectangle behind them. |
| `canvas` | The backend contract: `Canvas`, `PathPainter`, `RectPainter`, `CirclePainter`, `TextPainter`. Implement these and seen renders to your target. |
| `context` | `Context` — five methods: `SetLayers`, `Render`, `Animate`, `Drag`, `Zoom`. `context/svg` writes an SVG document. `context/gio` (the nested module) draws into `op.Ops` and supplies `Widget`. |
| `animation`, `drag`, `zoom` | Interaction. `drag` has inertia; `zoom` maps scroll deltas to a scale factor. |

## Usage

Every scene is the same four steps: build geometry, add it to a scene, wrap the
scene in a layer, hand the layer to a context. From
`context/gio/example/helloworld/main.go:54` — a randomly coloured sphere you can
spin with the pointer and scale with the wheel:

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
callback it takes runs once per frame with the widget's own size in `unit.Dp`,
which is where `FitCenter` belongs — the scene refits itself to whatever the
layout gave it. `context.Render()` is `window.Invalidate()`, nothing more:
seen never draws outside a frame.

The production consumer is `workbench/launcher`, which puts a seen scene behind
a whole application. Its `field.go:138` is the same shape with one addition
worth copying — the scene is isolated from the op list it renders into:

```go
view := seengio.Widget(f.ctx, func(w, h unit.Dp) {
	f.scene.FitCenter(0, 0, float64(w), float64(h), cameraDist)
})

f.view = func(gtx layout.Context) layout.Dimensions {
	defer op.Offset(image.Point{}).Push(gtx.Ops).Pop()
	return view(gtx)
}
```

The identity `op.Offset` push and pop discards any transform the seen widget
leaves behind, so the UI drawn above the background can never be disturbed by
it. Without it, a scene used as a backdrop can move the layers on top of it.

## For coding assistants

Read the canonical guide before writing code against this module — the module
inventory with current tags, the application skeleton, MVU and rx semantics,
typography, and the pitfalls that are not guessable:

<https://raw.githubusercontent.com/vibrantgio/workbench/master/llms.txt>

[`AGENTS.md`](./AGENTS.md) in this repository has the build and test commands.

## Status

Honest about what does not work yet. Every count below is measured.

- **`PLAN.md` at the root is not the organization's plan, and it is finished.**
  It covers one refactor of the camera and viewport pipeline, all twenty-two of
  its boxes are checked, and its header describes a workspace that does not
  exist: sibling checkouts directly under `~/code/w/vibrantgio` wired together
  by a `go.work`. There is no such `go.work` and the clones live one directory
  deeper. Do not follow it, and do not trust a build that it says was verified —
  a workspace would have masked exactly the version pins that break the
  consumers below. The organization's plan lives in `vibrantgio/.github`.
- **`svg/driver/seen` does not build**, and neither does `workbench/launcher`.
  Both fail identically on `verifying github.com/vibrantgio/seen/context/gio@v0.0.7:
  checksum mismatch`. This is a stale `go.sum` pin on the *consumer* side, not a
  defect in seen — seen's own two modules build and test clean — but it means
  the module's only two real consumers are currently red, and nothing exercises
  the Gio render path end to end.
- **The Gio rendering path has no tests at all.** All fourteen packages in
  `context/gio` report `[no test files]`, including `context/gio/canvas`, which
  is the code that turns a scene into `clip` and `paint` operations. In the root
  module sixteen of thirty-three packages are untested — among them the whole of
  `shader` (all four shading models), all of `light`, all of `face`, all of
  `transform`, and `layer/zsort`. Seventeen packages do have tests and all pass.
- **`context/gio/canvas/text.go:50` panics inside a render path** when a font
  fails to parse. A library that is called once per frame should return an
  error, not take the process down. `point/points.go` panics twice more on a
  slice-length mismatch.
- **`Scene.Regenerate`'s doc points at a method that is not there.** `scene.go:41`
  says "To flush the cache, call FlushCache()", but `Scene` has no such method.
  The only `FlushCache` in the repository is on `layer/zsort.Layer` — a
  different type in a different package, and the one sorter nothing uses.
  `nsort` and `bsort` cache too, and neither offers any way to flush.
- **Two of the three sorters are effectively dead.** Every one of the eleven Go
  examples and the launcher use `nsort`. `zsort` and `bsort` are reached only
  from seen's own tests and from `svg/driver/seen`, which does not build. `bsort`
  survives because the order-checking harness pins it; `zsort` has no tests, no
  package comment, and its own harness notes say it "never cuts, is approximate
  by design and is not run against this harness".
- **Other packages with no consumer anywhere in the organization:** `solid`,
  `mocap`, `bvh`, `affine`, `canvas` and `layer/backdrop`. `obj` has exactly one,
  the `multipleangles` example.
- **There are no golden-image tests, and nothing to regenerate.** What exists
  instead is better and worth knowing about: `layer/internal/ordercheck` renders
  three scenes — a cross, an occlusion cycle and a coplanar control — through the
  SVG context, parses the emitted `<path>` elements back (document order is paint
  order), and checks the last-painted polygon at sample points against a
  ray/plane intersection computed in the test. It is analytic, so there is no
  reference artifact and no `-update` flag; `go test ./...` is the whole story.
  `mocap/render_test.go` is a smoke test that counts paths, not a pixel
  comparison.
- **The package documentation is close to absent.** `doc.go` is two lines long
  and its entire content is `// Package seen` — so pkg.go.dev's front page for
  the module says nothing at all. Eight of thirty-three root packages carry a
  package comment; the other twenty-five have none, and neither of `context/gio`'s
  two library packages has one.
- **The LICENSE is Apache-2.0 with the template placeholder left in** —
  `Copyright [yyyy] [name of copyright owner]`, unfilled — and there is no
  NOTICE file. Neither seen.js nor csg.js is credited there. That second one
  matters: `solid/` is a derivative of csg.js, which is MIT, and the attribution
  did not travel with the code. The provenance is recorded only in `AGENTS.md`
  and in scattered example comments.
- **Seven of the eighteen example directories were never ported.** `2048`,
  `audioequalizer`, `canvasvssvg`, `depthoffield`, `gallery`, `masksandeffects`
  and `nbody` contain a `main.coffee` and a screenshot and no Go at all. Nine
  more keep the original `.coffee` alongside the Go port. The repository also
  carries a checked-in minified `seen.min.js` build and a checked-in
  `main.wasm`.
- **The animation clock is a package-level global.** `seen.Scheduler` and the
  interval table in `interval.go` are `var`s in the root package, so two scenes
  in one process share one scheduler. Nothing in the API lets you give a scene
  its own.
- **`go vet` is not clean.** The root module reports ten unkeyed struct literals
  (`affine/affine.go` and `shape/shape_test.go`); `context/gio` reports
  unreachable code in five of the example mains. `solid/vector.go` also keeps
  sixty lines of csg's original `Vector` type in a block comment, beneath the
  `type Vector = point.Point` alias that replaced it.
