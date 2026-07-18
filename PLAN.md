# seen: conventional camera/viewport (the Fit refactor)

Restructure seen's view pipeline into a conventional model — Camera owns the
full pose (eye position + orientation + normalization), Viewport is purely the
NDC→pixel screen mapping — while isolating the legacy seen.js "fitting"
behaviour (fill-the-view, page-space identity, distance-follows-height) into
two Scene setup methods, `FitCenter` and `FitOrigin`, that reproduce today's
rendering **exactly**. The refactor is a pure regrouping of the same matrix
product; identical output is proven by the existing test suites, not eyeballed.

Working directory: `/Users/rene/code/w/vibrantgio/seen` (git repo, branch
`master`; commit directly on master, never branch-first). The workspace root
`/Users/rene/code/w/vibrantgio` is NOT a repo; its `go.work` wires the local
checkouts (`seen`, `seen/context/gio`, `workbench/*`, …) together, so
cross-module builds use local code regardless of go.mod version pins.

Rules for every task:
- Commit per task, task heading in the commit message, body ending with
  `Co-Authored-By: Claude Fable 5 <noreply@anthropic.com>`. Stage only the
  files the task touched — never `ref/`, never unrelated working-tree changes.
- `gofmt` and `go vet` clean on files you touch (pre-existing vet warnings in
  `context/gio/example/*` about unreachable code are NOT yours to fix).
- Never leave build binaries in the repo: `go build ./...` for checking, and
  any runnable binary goes to the session scratchpad via `-o`.
- The bar for every phase: `go test ./...` green in the affected module, and
  the golden-producing tests byte-identical (see REF-verification).

## Out of scope

- `LookAt(eye, target, up)` and any new camera features — this plan only
  restructures; LookAt becomes a trivial follow-up once Camera owns the pose.
- Removing the bsort or zsort layer packages (bsort stays pinned by the
  ordercheck artifact harness; no example or app uses either anymore).
- Pushing, tagging (v0.0.5), or bumping consumer go.mod pins — local go.work
  covers all in-workspace consumers; release mechanics are the user's call.
- context/svg API changes beyond the mechanical call-site sweep.
- Any change to `projection` package matrices, shading, or geometry.

## Decisions

### ADR-001: Goals are context-bound, not time-bound

SMART's "time-bound" becomes "context-bound": a task is right-sized when an agent
can finish it inside one bounded context window. Target ~100K tokens for a task's
full working set — its `mdplan next` packet plus the files and output it needs.
If a task will not fit, split it until each one does.

### ADR-002: Fitting becomes Scene.FitCenter/FitOrigin; hard break, no shims

The old `viewport.Center`/`viewport.Origin` cannot survive as wrappers: they
return only a `Viewport`, but the fitting behaviour also positions the eye,
which now lives on the Camera. They are REPLACED by Scene methods that set
both. All consumers are in this workspace (go.work), no released tag ships the
new API yet, so a hard break with a mechanical sweep is cheaper and clearer
than a transition shim. `viewport.Default` disappears; the zero-config default
moves into `camera.Default` (Eye `(0,0,1)`, Norm identity) plus a
`viewport.Default` replacement computed by the same formula as `FitOrigin(0,0,1,1)`.

### ADR-003: Camera.Transform keeps world-transform (view-matrix) semantics

Today `Camera.Matrix()` is applied to world points directly — it IS a view
matrix, not a pose to invert. Existing users depend on that: mocap dollies with
`scene.Camera.SetTranslation(0, 0, h/2)`, the svg tests translate and rotate
the camera, turntable drags rotate about the world origin. The refactor does
NOT flip this to pose-and-invert; it only ADDS `Eye` (translation applied after
the transform) and `Norm` (the fit scale). Result: every existing camera
manipulation keeps its exact meaning, and rotation still orbits the origin.

### ADR-004: Identity is proven by regrouping plus byte-identical tests

The new pipeline multiplies the SAME factors in the same pairwise order; only
the outer association changes (`(P·(S·T))·M` becomes `P·((S·T)·M)`), so
numerical drift is ≤ ~1e-15 relative. The SVG pipeline rounds coordinates to 3
decimals, so golden files must come out byte-identical. Acceptance rule: if any
golden differs, inspect the diff — a last-decimal flicker on isolated
coordinates may be accepted by regenerating ONLY after confirming it is float
association noise (sub-0.001 px); anything larger means the refactor has a bug.
Do not regenerate goldens wholesale.

## Reference

### REF-current-pipeline

Matrix conventions (matrix package): builder chaining right-multiplies —
`A.Translate(t)` = `A · T(t)`, so the RIGHTMOST factor applies to the point
first. `m.Mul(n)` = `m · n` (n first). Transforms compose T·R·S (scale first).

Render pipeline, identical in `layer/{bsort,nsort,zsort}/layer.go` RenderOn:

```go
projection := scene.Camera.Projection.Mul(scene.Viewport.Prescale).Mul(scene.Camera.Matrix())
viewport   := scene.Viewport.Postscale
// per face: world points → Clip(projection, -2, …) (reject z_clip ≤ −2, divide
// by w) → MulB(viewport) → screen points     (point/points.go)
```

`viewport/viewport.go` today (`Viewport struct{ Prescale, Postscale matrix.Matrix }`):

```go
// Center(ox, oy, w, h, dist...): W,H,D = w,h,h — or dist,dist,dist when given
Prescale:  matrix.Scale(1/W, 1/H, 1/D).Translate(-ox, -oy, -D)
Postscale: matrix.Translate(ox+w/2, oy+h/2, D).Scale(W, -H, D)
// Origin: same Prescale; Postscale: matrix.Translate(ox, oy, D).Scale(W, -H, D)
// Default = Origin(0, 0, 1, 1)
```

`camera/camera.go`: `Camera{transform.Transform; Projection matrix.Matrix}`,
`Default = CameraWithProjection(projection.DefaultPerspective)`. The camera
transform is applied to WORLD points (see ADR-003).

Eye recovery, duplicated in bsort and nsort layer.go (zsort has none):

```go
view := scene.Viewport.Prescale.Mul(scene.Camera.Matrix())
if inv, ok := view.Invert(); ok { eye = inv · origin } else
    { eye = point.Pt(0, 0, -1.0/projection[2][2]) }  // legacy fallback
```

Key insight the refactor encodes: `Prescale` is HALF CAMERA — its translate
places the eye at world `(ox, oy, D)`, its scale is projection normalization.
Only `Postscale` is a true viewport. The offsets appear in both halves so that
world `(ox, oy, 0)` maps to the region's anchor — seen.js "page space".

### REF-target-model

```go
// camera package
type Camera struct {
    transform.Transform            // world transform — semantics unchanged (ADR-003)
    Projection matrix.Matrix       // perspective — unchanged
    Eye        point.Point         // eye position in world space
    Norm       matrix.Matrix       // view normalization (fit scale); matrix.Identity when unused
}

// View returns the world→view matrix: Norm · Translate(−Eye) · Transform.Matrix().
// Factor order matches the old Prescale·CameraM exactly.
func (c Camera) View() matrix.Matrix

// EyeInWorld returns the eye's world position: the preimage of the view-space
// origin under View(). Centralizes the recovery both sorting layers duplicate
// today, including the legacy -1/proj[2][2] fallback for degenerate views
// (fallback needs the composed projection: c.Projection.Mul(c.View())).
func (c Camera) EyeInWorld() point.Point

var Default = … // Transform: transform.Default, Projection: DefaultPerspective,
                // Eye: point.Pt(0, 0, 1), Norm: matrix.Identity
                // (equivalent to the old camera.Default + viewport.Default pair)

// viewport package
type Viewport struct{ Screen matrix.Matrix }  // NDC→pixels; nothing else
var Default = …                                // Screen of old Origin(0,0,1,1) Postscale

// seen package (scene.go)
// FitCenter configures Camera and Viewport with the legacy seen.js fitting
// behaviour — identical to the old scene.Viewport = viewport.Center(…):
// eye above world (x, y) at distance D, world (x, y, 0) at the region centre,
// scale following the region unless locked with dist.
func (s *Scene) FitCenter(x, y, w, h float64, dist ...float64)
func (s *Scene) FitOrigin(x, y, w, h float64, dist ...float64)
```

Fit formulas (W,H,D = w,h,h or dist,dist,dist — same rule as today):

| sets | FitCenter | FitOrigin |
|---|---|---|
| `Camera.Eye` | `(x, y, D)` | `(x, y, D)` |
| `Camera.Norm` | `Scale(1/W, 1/H, 1/D)` | same |
| `Viewport.Screen` | `Translate(x+w/2, y+h/2, D).Scale(W, −H, D)` | `Translate(x, y, D).Scale(W, −H, D)` |
| untouched | `Camera.Transform`, `Camera.Projection` | same |

Layers become:

```go
projection := cam.Projection.Mul(cam.View())
screenMap  := scene.Viewport.Screen
eye        := cam.EyeInWorld()          // bsort + nsort; delete their local recovery
```

Equivalence check (write it as a unit test): for both Fit variants and several
(x, y, w, h[, dist]) tuples including the launcher's `(0,0,1100,760,2200)`,
`cam.Projection.Mul(cam.View())` equals the old
`Projection·Prescale·CameraM` within 1e-12 per element, and `Screen` equals the
old `Postscale` exactly.

### REF-call-sites

Sweep pattern: `scene.Viewport = viewport.Center(args…)` becomes
`scene.FitCenter(args…)` (drop the `viewport` import if unused). All sites:

Module `seen` (Phase P1):
- `scene.go:49,60` — NewScene/NewDefaultScene defaults → `camera.Default` +
  new `viewport.Default` (no behaviour change).
- `context/svg/context_test.go:259,311,405` — Center(0,0,width,height); note
  this file also uses `Camera.SetTranslation/SetRotation` which stay untouched.
- `mocap/render_test.go:34` — Center(0,0,500,500).
- `layer/internal/ordercheck/ordercheck.go:69,86,199` — the harness MIRRORS the
  pipeline: rewrite `projectPoint` and `eyePoint` against the new model
  (build a Camera+Viewport via the same Fit formulas, compose
  `Projection·View`, use `EyeInWorld`); scene setup at :199 → `FitCenter`.
- `layer/{bsort,nsort}/layer.go` — pipeline + eye recovery (REF-target-model).
- `layer/zsort/layer.go` — pipeline only.
- `matrix/matrix_test.go:74` — comment mentions viewport.Center; update wording.

Module `seen/context/gio` (Phase P2) — all resize callbacks:
- helloworld:87, combinedsolid:114, noisywavepatch:85, mocap:115 (keep its
  `Camera.SetTranslation(0,0,h/2)` dolly — unchanged), rectangle:85,
  giftbox:150, text:103, noisysphere:51, solids:123, poem:74 — all
  `Center(0,0,w,h)` → `scene.FitCenter(0, 0, w, h)`.
- multipleangles:68–71,94–95 — REWRITE, not sweep (task G2.1.2): today each
  mini scene calls `Center(ox,oy,vw,vh)` plus a compensating
  `Camera.SetTranslation(ox, oy, 0)` to drag the eye back over the origin.
  New model: `scene.FitCenter(ox, oy, vw, vh)` then set `Camera.Eye =
  point.Pt(0, 0, D)` directly (D = vh) — the compensation and its long
  comment DIE; per-view `Camera.RotX/RotY` stay as-is. Same pixels.

Module `workbench` (Phase P3, separate repo at
`/Users/rene/code/w/vibrantgio/workbench`, branch master):
- `launcher/field.go:140` — `viewport.Center(0,0,w,h,cameraDist)` →
  `f.scene.FitCenter(0, 0, float64(w), float64(h), cameraDist)`.

### REF-verification

Automated (must be green; goldens byte-identical per ADR-004):
- `cd /Users/rene/code/w/vibrantgio/seen && go test ./...` — includes
  context/svg golden tests, mocap render test, layer/bsort+bsp tests,
  layer/nsort tests (incl. the ordercheck artifact harness for BOTH bsort and
  nsort — cross/cycle/coplanar sample-point correctness), solid volume tests,
  matrix tests. `git status` must show no modified `testdata/`.
- `cd seen/context/gio && go build ./... && go test ./...`.
- `cd /Users/rene/code/w/vibrantgio/workbench/launcher && go build ./...`.
- Grep gate when the sweep is done: `grep -rn "Prescale\|Postscale"
  --include="*.go" seen | grep -v ref/` must return nothing.

Visual (GUI apps on this Mac — hard-won recipe, follow it exactly):
- NEVER launch a Gio app from a sandboxed shell — it panics with
  `runtime/cgo: misuse of an invalid Handle`. Use Bash with
  `dangerouslyDisableSandbox: true`.
- The SAME panic fires if the display is asleep: run `caffeinate -u -t 2` first.
- Build with `-o <scratchpad>/name-bin`, run with the Bash tool's
  `run_in_background: true` (a nohup'd process dies with the tool call).
- Screenshot with `screencapture -x <scratchpad>/shot.png`, crop with
  `sips -c <h> <w> --cropOffset <y> <x>`, then READ the png. `pkill` the
  binary afterwards.
- What to check: giftbox (box + ribbon render, drag rotates), solids (three
  CSG solids left→right red union / green subtract / blue intersect),
  multipleangles (bunny in 5 viewports — main large top, four minis beneath:
  default, top view, underside, side; drag on main rotates ALL views),
  launcher (workbench root: triangle field covers window incl. top edge,
  cards float on it).

## Phase P1: Core model in the seen module

The seen module converts in one phase because Scene's Camera/Viewport types
change shape — the module must build and test green at every commit within
this phase's sequence. context/gio will NOT build until Phase P2; that interim
breakage across modules is accepted (single repo, sequential commits).

![[#REF-current-pipeline]]
![[#REF-target-model]]

### G1.1: New camera/viewport model with proven-equivalent Fit setup

- **Specific:** Camera carries Eye+Norm with View()/EyeInWorld(); Viewport is a
  single Screen matrix; Scene gains FitCenter/FitOrigin reproducing the old
  Center/Origin behaviour exactly.
- **Measurable:** the equivalence unit test (REF-target-model) passes at 1e-12;
  `go build ./...` in the seen module compiles everything except packages that
  still read Prescale/Postscale (they convert in G1.2).
- **Achievable:** pure additive matrix bookkeeping plus two small structs.
- **Relevant:** unlocks real multi-camera work (LookAt later) without
  sacrificing the zero-config fitting that every existing scene relies on.
- **Context-bound:** three files plus one test file — see ADR-001.

![[#ADR-002]]
![[#ADR-003]]

#### G1.1.1: Add Eye and Norm to Camera with View and EyeInWorld

Extend `camera/camera.go` per REF-target-model. Keep Transform/Projection
semantics untouched (ADR-003). `View()` must compose `Norm.Mul(matrix.Translate(
-Eye.X, -Eye.Y, -Eye.Z)).Mul(c.Matrix())` — exactly that factor order.
`EyeInWorld()` ports the recovery block from `layer/bsort/layer.go` (see
REF-current-pipeline) including the degenerate-view fallback, against
`c.View()` and `c.Projection.Mul(c.View())`. Update `camera.Default` (Eye
`(0,0,1)`, Norm `matrix.Identity`).

- [x] Add the `Eye`/`Norm` fields, `View()`, `EyeInWorld()`, and the new
      `Default`; document each with the rationale from ADR-003.
- [x] Unit tests in `camera/`: View() factor order against a hand-built
      matrix; EyeInWorld() for identity transform (returns Eye), for a
      translated camera (mocap dolly case), and the fallback path for a
      zero-scale Norm.

#### G1.1.2: Reduce Viewport to the screen mapping and add Scene.Fit methods

Rewrite `viewport/viewport.go`: `Viewport struct{ Screen matrix.Matrix }`,
new `Default`, delete Center/Origin (ADR-002). Add `FitCenter`/`FitOrigin`
methods on Scene in `scene.go` per the REF-target-model formula table, and
update NewScene/NewDefaultScene defaults. Port the doc comments from the old
Center/Origin (the dist-locking explanation is load-bearing — the launcher
depends on it) onto the Fit methods.

- [x] Rewrite the viewport package; move the fitting formulas into
      `Scene.FitCenter` / `Scene.FitOrigin`.
- [x] Equivalence test (in `seen` root package or `viewport`): for FitCenter
      and FitOrigin over several tuples incl. `(0,0,1100,760,2200)` and a
      non-square `(10,20,600,300)`, new `Projection·View` equals old
      `Projection·Prescale·CameraM` (old formulas hardcoded in the test) to
      1e-12, and `Screen` equals old `Postscale` exactly.

### G1.2: Layers and seen-module call sites on the new model

- **Specific:** all three layers compose `Projection·View` and map through
  `Viewport.Screen`; bsort/nsort use `EyeInWorld()`; every seen-module call
  site (REF-call-sites, module seen) uses Fit; ordercheck mirrors the new math.
- **Measurable:** `go test ./...` green in seen with byte-identical goldens
  (`git status` clean of testdata churn); the Prescale/Postscale grep gate
  passes for the seen module.
- **Achievable:** mechanical edits with the harness as the only thinking part.
- **Relevant:** deletes the duplicated eye recovery — the code that motivated
  the matrix.Invert epsilon fix — and makes the layers read conventionally.
- **Context-bound:** six small files plus the harness — see ADR-001.

![[#ADR-004]]
![[#REF-call-sites]]

#### G1.2.1: Convert the three layers

In `layer/bsort/layer.go`, `layer/nsort/layer.go`, `layer/zsort/layer.go`:
projection composition per REF-target-model; replace both local eye-recovery
blocks with `scene.Camera.EyeInWorld()`; keep everything else (tree caching,
shader, ordering) untouched. bsort's world-only rebuild trigger reads the
camera matrices via the shader — verify a camera Eye/Norm change does NOT
count as a world change (it must not: world planes are pre-view).

- [x] Convert bsort, nsort, zsort RenderOn to `Projection·View` +
      `Viewport.Screen`.
- [x] Replace the two eye-recovery blocks with EyeInWorld and delete the
      local fallbacks.
- [x] `go test ./layer/...` green (ordercheck still fails here if not yet
      converted — that is G1.2.2's box; run the non-harness tests).

#### G1.2.2: Rewrite the ordercheck harness math and sweep seen call sites

`layer/internal/ordercheck/ordercheck.go` mirrors the old pipeline in
`projectPoint`/`eyePoint`/`renderScene` — rebuild them on the new model
(construct the Camera/Viewport via the same formulas Scene.FitCenter uses,
compose `Projection·View`, `EyeInWorld` for the analytic rays). Then the
mechanical sweep: scene.go defaults, context/svg/context_test.go,
mocap/render_test.go, matrix_test comment (all listed in REF-call-sites).

- [ ] Rewrite ordercheck's projection/eye helpers on the new model.
- [ ] Sweep the remaining seen-module call sites to FitCenter.
- [ ] Full `go test ./...` green, goldens byte-identical (ADR-004), grep gate
      clean for the seen module; commit closes Phase P1.

## Phase P2: context/gio module sweep

The examples module. The user may have uncommitted work here — check
`git status` FIRST and never stage or revert files you did not change for
this plan. All ten mechanical sites and the one rewrite are listed in
REF-call-sites (module seen/context/gio).

![[#REF-call-sites]]
![[#REF-verification]]

### G2.1: Examples build and render identically on the new model

- **Specific:** ten resize callbacks swept to FitCenter; multipleangles
  rewritten to express its cameras directly (no compensation).
- **Measurable:** `go build ./...` in context/gio; giftbox, solids and
  multipleangles verified on screen per REF-verification.
- **Achievable:** one-line swaps plus one small rewrite with a known-identical
  target image.
- **Relevant:** proves the new model against every rendering path the examples
  exercise, and produces the first genuinely-conventional multi-camera code.
- **Context-bound:** eleven files, one of them thoughtful — see ADR-001.

#### G2.1.1: Sweep the ten mechanical example call sites

- [ ] Replace `scene.Viewport = viewport.Center(0,0,w,h)` with
      `scene.FitCenter(0, 0, w, h)` in helloworld, combinedsolid,
      noisywavepatch, mocap, rectangle, giftbox, text, noisysphere, solids,
      poem; drop unused viewport imports.
- [ ] `go build ./...` green in context/gio.

#### G2.1.2: Rewrite multipleangles on real cameras

Per REF-call-sites: FitCenter per region, then override `Camera.Eye` to
`(0, 0, D)` so every camera looks at the world origin; delete the
`SetTranslation(ox, oy, 0)` compensation and its explanatory comment (the
model now expresses the intent directly); keep the per-view rotations and the
main view's `SetScale(2,2,2)`.

- [ ] Rewrite the five-view setup; the compensation hack and comment are gone.
- [ ] Visual check per REF-verification: bunny in all five views, drag on the
      main view rotates all five in lockstep.

#### G2.1.3: Visual verification of giftbox and solids

- [ ] Run giftbox and solids per the REF-verification recipe and confirm the
      described images; screenshots seen, binaries killed.
- [ ] Commit closes Phase P2.

## Phase P3: workbench launcher and wrap-up

Separate repo: `/Users/rene/code/w/vibrantgio/workbench` (branch master, same
commit rules). The launcher's field locks its scale with the `dist` parameter
— the Fit methods must have preserved that (they did if G1.1.2 ported the doc
semantics; the equivalence test covers the `(0,0,1100,760,2200)` tuple).

![[#REF-verification]]

### G3.1: Launcher on the new model, workspace consistent

- **Specific:** launcher swept; whole workspace builds; no Prescale/Postscale
  or viewport.Center references remain anywhere.
- **Measurable:** launcher builds and its field renders full-window on screen;
  final grep gates return nothing.
- **Achievable:** one line plus verification.
- **Relevant:** the launcher is the flagship consumer and the fixed-dist
  stress case for the fitting semantics.
- **Context-bound:** one file plus checks — see ADR-001.

#### G3.1.1: Sweep launcher field.go and verify

- [ ] `launcher/field.go:140` → `f.scene.FitCenter(0, 0, float64(w),
      float64(h), cameraDist)`; drop the viewport import.
- [ ] Build, run, screenshot per REF-verification: triangle field covers the
      window including the top edge at the default size AND after resizing
      the window larger (the regrow path exercises Fit on every resize).
- [ ] Commit in the workbench repo.

#### G3.1.2: Final consistency gates

- [ ] Workspace-wide: `grep -rn "viewport.Center\|viewport.Origin\|Prescale\|
      Postscale" --include="*.go" /Users/rene/code/w/vibrantgio | grep -v
      ref/` returns nothing.
- [ ] Full test suites one last time: seen `go test ./...`, context/gio
      `go build ./... && go test ./...`, launcher `go build ./...`.
- [ ] Final commit in seen: note in the body that consumers outside go.work
      need a seen tag (v0.0.5) plus context/gio and workbench go.mod bumps
      before this is consumable remotely.
