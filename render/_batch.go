package render

/*
func (pb *build) batch() {
	trac.Log("batch")

	var batched []*buildPass
	var last *buildPass

	for i, p := range pb.passes {
		if i == 0 {
			last = p
			continue
		}
		if pb.mergeable(last, p) {
			pb.merge(last, p)
		} else {
			batched = append(batched, last)
			last = p
		}
	}
	batched = append(batched, last)
	trac.Log("  merged passes %d to %d", len(pb.passes), len(batched))
	pb.passes = batched
}

func (pb *build) mergeable(a, b *buildPass) bool {
	if a.prog.id != b.prog.id {
		trac.Log("program IDs differ")
		return false
	}
	if len(a.bindings) != len(b.bindings) {
		trac.Log("bindings differ", len(a.bindings), len(b.bindings))
		return false
	}
	if a.instanceCount > 0 || b.instanceCount > 0 {
		// TODO probably can be supported, but has been buggy
		trac.Log("instanced")
		return false
	}
	if len(a.uniforms) != len(b.uniforms) {
		trac.Log("uniforms differ")
		return false
	}
	if a.layer != b.layer {
		trac.Log("layers differ")
		return false
	}
	for i := range b.bindings {
		if a.bindings[i].attr != b.bindings[i].attr {
			return false
		}
		if a.bindings[i].divisor != b.bindings[i].divisor {
			return false
		}
	}
	for k, v := range a.uniforms {
		c, ok := b.uniforms[k]
		if !ok {
			return false
		}
		if c != v {
			return false
		}
	}
	return true
}

func (pb *build) merge(a, b *buildPass) {

	// merge faces
	vc := uint32(a.vertexCount)
	for i := 0; i < b.faceCount; i++ {
		pb.faces[b.faceOffset+i] += vc
	}
	a.faceCount += b.faceCount
	a.vertexCount += b.vertexCount

	for i := range a.bindings {
		for _, bv := range b.bindings[i].values {
			a.bindings[i].values = append(a.bindings[i].values, bv)
		}
	}
}
*/
