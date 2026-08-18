# AGENTS.md — seen

A 3-D scene graph and software renderer, ported from seen.js: `Scene` and
`Group` over `shape`, `solid` and `obj` geometry; `camera`, `projection`
and `viewport` for the view pipeline; `light` and `shader` for shading;
`layer` for painter's-algorithm ordering — `zsort`, `nsort`'s Newell sort
and `bsort`'s BSP splitting; `mocap` for BVH motion capture; `drag` and
`zoom` for interaction. Rendering goes through a `canvas` context:
`context/gio` draws into Gio ops and supplies a `Widget`, `context/svg`
writes SVG.

**Layer.** Outside ADR-001's tier table: a support library, which the rule
binds in one direction only — every tier may import it, and it may import
nothing in the table itself. Its root module imports nothing else in the
organization. Its nested `seen/context/gio` module adds `noise` — that edge
is the nested module's and not the root's. That direction is measured
rather than typed — `scripts/check-layers.sh --edges` reports the graph and
`scripts/sync-agents.sh` renders these sentences from it — so correcting
them here changes nothing. The other direction is measured too and
deliberately not written down: the gate checks the graph both ways, but a
public API's consumers are unknowable, so this file says what its module
needs and never who needs it.

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
refactor of the camera and viewport pipeline and every task in it is checked
off — `mdplan next PLAN.md` prints DONE, which is how to re-check that. Its
header describes sibling checkouts directly under the workspace root with no
repository at that root — which was wrong for the whole span when the clones
hid under `.github`'s gitignored `.repos/` (B2.1 through G0E.1), and this note
said so. G0E.1 flattened the tree back to exactly that shape, so the header's
paths resolve again; only its claim that `go.work` is missing stays wrong —
`clone-all.sh` generates one at the workspace root. The organization's plan
lives in the `.github` sibling. Do not follow this one unless the task
you were handed is one of its tasks.

**`solid/` is an adaptation of `vibrantgio/csg`, not a dependency on it.** The
BSP algorithm is the same; the types are not. seen's copy is rewritten onto
`point.Point`, `face.Faces` and `transform.Transform` so that a `*Solid` is a
`seen.Object` that drops straight into a scene, and this repository requires
nothing to get it. A fix in one is not a fix in the other.
