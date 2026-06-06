// Package mocap turns a parsed BVH [bvh.Hierarchy] into an animatable seen
// scene-graph skeleton.
//
// It mirrors the Mocap/MocapModel/MocapAnimator types from the original
// seen.js library: [New] walks the joint tree building a nested [seen.Group]
// per joint (so each joint inherits its parent's transform), draws a bone
// shape from every joint to each of its children, and precomputes the
// per-frame translation+rotation for every animated joint. [Model.Apply] then
// poses the skeleton at a given motion frame.
//
//	h, _ := bvh.Load("dance.bvh")
//	m := mocap.New(h, nil)
//	scene.Group.Add(m.Group)
//	m.Apply(frameIndex) // pose the skeleton, e.g. once per animation tick
package mocap

import (
	"math"
	"time"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/bvh"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
)

// endSiteID is the Id the parser assigns to a leaf "End Site" joint, whose
// offset is stored in [bvh.Joint.EndSite] rather than Offset.
const endSiteID = "END SITE"

// ShapeFactory builds the bone shape drawn from a joint's origin to endpoint,
// where endpoint is the offset of one of the joint's children expressed in the
// joint's local frame. Returning nil omits that bone (e.g. for a zero-length
// segment).
type ShapeFactory func(joint bvh.Joint, endpoint point.Point) seen.Object

// defaultFill is applied to bones built by DefaultShapeFactory so a skeleton
// created with a nil factory is visible without further setup.
var defaultFill = color.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0}

// DefaultShapeFactory draws a light-grey unit-radius pipe from the joint origin
// to endpoint, skipping zero-length segments.
func DefaultShapeFactory(_ bvh.Joint, endpoint point.Point) seen.Object {
	if endpoint.Length() < 1e-9 {
		return nil
	}
	pipe := shape.Pipe(point.Pt(0, 0, 0), endpoint)
	pipe.Faces().SetFill(defaultFill)
	return pipe
}

// Model is a skeleton built from a parsed BVH hierarchy together with the
// per-frame transforms needed to animate it. Add [Model.Group] to a scene and
// call [Model.Apply] to pose the skeleton at a frame.
type Model struct {
	// Group is the root of the skeleton's scene-graph subtree.
	Group *seen.Group
	// FrameTime is the sampling interval between motion frames.
	FrameTime time.Duration

	joints []animJoint
	frames [][]pose
}

// animJoint is a joint whose orientation/position is driven by motion channels.
type animJoint struct {
	node     *seen.Group
	offset   [3]float64
	channels []bvh.Channel
}

// pose is a precomputed transform for one joint node at one frame.
type pose struct {
	node       *seen.Group
	tx, ty, tz float64
	r          quaternion.Quat
}

// New builds an animatable skeleton from h, using factory to create the bone
// shapes. If factory is nil, [DefaultShapeFactory] is used.
func New(h bvh.Hierarchy, factory ShapeFactory) *Model {
	if factory == nil {
		factory = DefaultShapeFactory
	}
	m := &Model{Group: seen.NewGroup(), FrameTime: h.Motion.FrameTime}
	m.attach(m.Group, h.Root, factory)
	m.frames = make([][]pose, len(h.Motion.Frames))
	for i, f := range h.Motion.Frames {
		m.frames[i] = m.computePoses(f)
	}
	return m
}

// Frames returns the number of motion frames available.
func (m *Model) Frames() int { return len(m.frames) }

// Apply poses the skeleton at frame index i, wrapping modulo [Model.Frames].
func (m *Model) Apply(i int) {
	n := len(m.frames)
	if n == 0 {
		return
	}
	for _, p := range m.frames[((i%n)+n)%n] {
		p.node.SetTranslation(p.tx, p.ty, p.tz)
		p.node.SetRotation(p.r)
	}
}

// attach recursively builds the scene-graph subtree for joint j under parent,
// recording every channel-bearing joint in m.joints in BVH depth-first order
// (the order motion-frame values are laid out).
func (m *Model) attach(parent *seen.Group, j bvh.Joint, factory ShapeFactory) {
	node := seen.NewGroup()
	node.SetTranslation(j.Offset[0], j.Offset[1], j.Offset[2])
	parent.Add(node)

	if len(j.Channels) > 0 {
		m.joints = append(m.joints, animJoint{node: node, offset: j.Offset, channels: j.Channels})
	}

	if len(j.Joints) == 0 {
		return
	}

	// Bones and child-joint subtrees live in a child group so they inherit
	// this joint's per-frame rotation.
	bones := seen.NewGroup()
	node.Add(bones)
	for _, child := range j.Joints {
		end := child.Offset
		if child.Id == endSiteID {
			end = child.EndSite
		}
		if s := factory(j, point.Pt(end[0], end[1], end[2])); s != nil {
			bones.Add(s)
		}
		if child.Id != endSiteID {
			m.attach(bones, child, factory)
		}
	}
}

// computePoses converts one flat motion frame into a per-joint transform.
// Each joint's local matrix is Translate(offset+position) · R, where R is the
// product of the rotation channels in the order they are listed (the BVH
// convention).
func (m *Model) computePoses(frame bvh.Frame) []pose {
	poses := make([]pose, len(m.joints))
	fi := 0
	for k, aj := range m.joints {
		tx, ty, tz := 0.0, 0.0, 0.0
		r := quaternion.Identity
		for ci, ch := range aj.channels {
			v := frame[fi+ci]
			switch ch {
			case bvh.Xposition:
				tx += v
			case bvh.Yposition:
				ty += v
			case bvh.Zposition:
				tz += v
			case bvh.Xrotation:
				r = r.Mul(quaternion.RotX(radians(v)))
			case bvh.Yrotation:
				r = r.Mul(quaternion.RotY(radians(v)))
			case bvh.Zrotation:
				r = r.Mul(quaternion.RotZ(radians(v)))
			}
		}
		fi += len(aj.channels)
		poses[k] = pose{
			node: aj.node,
			tx:   aj.offset[0] + tx,
			ty:   aj.offset[1] + ty,
			tz:   aj.offset[2] + tz,
			r:    r,
		}
	}
	return poses
}

func radians(deg float64) float64 { return deg * math.Pi / 180.0 }
