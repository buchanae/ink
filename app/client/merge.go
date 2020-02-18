package client

import "github.com/buchanae/ink/trac"

func mergeBatches(batches []*batch) []*batch {

	var merged []*batch
	var last *batch

	before := len(batches)
	defer func() {
		after := len(merged)
		trac.Log("merged %d to %b", before, after)
	}()

	for i, b := range batches {
		if i == 0 {
			last = b
			continue
		}
		if mergeable(last, b) {
			merge(last, b)
		} else {
			merged = append(merged, last)
			last = b
		}
	}
	merged = append(merged, last)
	return merged
}

func merge(a, b *batch) {
	a.meshes = append(a.meshes, b.meshes...)
	for k, v := range a.attrs {
		bv := b.attrs[k]
		v.vals = append(v.vals, bv.vals...)
	}
}

func mergeable(a, b *batch) bool {
	if a.pass.Shader != b.pass.Shader {
		trac.Log("shaders differ")
		return false
	}
	if a.pass.Layer != b.pass.Layer {
		trac.Log("layers differ")
		return false
	}
	if a.pass.Instances > 0 || b.pass.Instances > 0 {
		// TODO probably can be supported, but has been buggy
		trac.Log("instanced")
		return false
	}
	if len(a.pass.Uniforms) != len(b.pass.Uniforms) {
		trac.Log("uniform counts differ: %d != %d",
			len(a.pass.Attrs), len(b.pass.Attrs))
		return false
	}
	for k, v := range a.pass.Uniforms {
		c, ok := b.pass.Uniforms[k]
		if !ok {
			return false
		}
		if c != v {
			return false
		}
	}
	if len(a.attrs) != len(b.attrs) {
		trac.Log("attr counts differ: %d != %d",
			len(a.attrs), len(b.attrs))
		return false
	}
	for key, av := range a.attrs {
		bv, ok := b.attrs[key]
		if !ok {
			trac.Log("attrs differ: %s", key)
			return false
		}
		if av.divisor != bv.divisor {
			trac.Log("divisors differ: %s", key)
			return false
		}
	}
	return true
}
