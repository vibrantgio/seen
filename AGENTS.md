# AGENTS.md — seen

A 3-D scene graph and software renderer, ported from seen.js: `Scene` and
`Group` over `shape`, `solid` and `obj` geometry; `camera`, `projection`
and `viewport` for the view pipeline; `light` and `shader` for shading;
`layer` for painter's-algorithm ordering — `zsort`, `nsort`'s Newell sort
and `bsort`'s BSP splitting; `mocap` for BVH motion capture; `drag` and
`zoom` for interaction. Rendering goes through a `canvas` context:
`context/gio` draws into Gio ops and supplies a `Widget`, `context/svg`
writes SVG.

**Layer.** Outside ADR-001's tier table: a support library the design
system consumes and never depends on. No tier of the design system imports
seen; its consumers are `workbench/launcher`, whose animated background is
a seen scene, and `svg/driver/seen`, which turns SVG paths into seen
geometry. The root module depends on nothing whatsoever, inside the
organization or out; `context/gio` adds Gio, and two of its example
programs reach for the noise support library.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Modules.** `github.com/vibrantgio/seen` at the repository root, and one
nested module: `context/gio/` (`github.com/vibrantgio/seen/context/gio`).
Nested-module tags carry the directory as a prefix — `context/gio/v0.0.7`,
not `v0.0.7`.

**Build and test.** From the repository root, and again inside each nested
module directory — `./...` does not cross a module boundary:

    go build ./... && go test ./...

**`PLAN.md` here is seen's own plan, not the organization's.** It covers one
refactor of the camera and viewport pipeline, and its header describes a
workspace that no longer exists — sibling checkouts directly under
`~/code/w/vibrantgio`, wired together by a `go.work` that is not there. The
organization's plan lives in `vibrantgio/.github`. Do not follow this one
unless the task you were handed is one of its tasks.

**`solid/` is an adaptation of `vibrantgio/csg`, not a dependency on it.** The
BSP algorithm is the same; the types are not. seen's copy is rewritten onto
`point.Point`, `face.Faces` and `transform.Transform` so that a `*Solid` is a
`seen.Object` that drops straight into a scene, and this repository requires
nothing to get it. A fix in one is not a fix in the other.
