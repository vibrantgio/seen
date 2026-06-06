// Package bvh parses Biovision Hierarchical (BVH) motion capture files.
//
// A BVH file has two sections: a HIERARCHY describing a skeleton as a tree of
// [Joint]s (each with an offset from its parent and a set of animation
// [Channel]s), and a MOTION section holding the per-frame channel values. The
// parser is generated from grammar.peg with github.com/mna/pigeon; regenerate
// it with the ./p script after editing the grammar.
//
// The generated [Parse], [ParseFile] and [ParseReader] entry points return the
// result as an any holding a [Hierarchy]. Use [Load] for a typed convenience
// wrapper:
//
//	h, err := bvh.Load("testdata/01_06.bvh")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(h.Root.Id, len(h.Motion.Frames), h.Motion.FrameTime)
package bvh
