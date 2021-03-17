package rendr
/* The rndKernel stores everything about a reconstruction kernel. The kernel
   is non-zero only within [-support/2,support/2], for integer "support"
   which may be odd or even (but always positive). The kernels are set up at
   compile-time in such a way that each kernel knows its own derivative; the
   derivative of rndKernel *k is k->deriv. */
type evalFunc func(xx float64) float64
type applyFunc func(ww []float64, xa float64)
type rndKernel struct {
	Name, Desc string // identifying and descriptive string
	Support int // # samples needed for convolution
	Eval evalFunc // evaluate the kernel once
	Apply applyFunc // evaluate the kernel support times
	Deriv *rndKernel // derivative of this kernel; will point back to itself when kernel is zero
}

// kernels defined in kernel.go
//var rndKernelZero *rndKernel = &zero1Kernel
//var rndKernelCtmr *rndKernel = &ctmr

const RND_KERNEL_SUPPORT_MAX int = 6

// kernel parse
//func rndKernelParse(_kstr string) *rndKernel;

type rndType int
//rndType
const (
	rndTypeUnknown rndType=iota          // (0) no type known
	rndTypeUChar            // (1) uchar (8-bit unsigned)
	rndTypeShort            // (2) short (16-bit signed)
	rndTypefloat64             // (2) float64 (floating point)
)

type rndSpace int
//rndSpace
const (
	rndSpaceUnknown rndSpace=iota
	rndSpaceRGB
	rndSpaceHSV
	rndSpaceAlpha
)

type rndProbe int
//rndProbe
const (
	rndProbeUnknown rndProbe=iota /* (0) */
	rndProbePosView   /* (1) 3-vector: view-space position of probe
   						(for debugging ray-casting geometry) */
	rndProbePosWorld  /* (2) 3-vector: world-space position of probe
						(for debugging ray-casting geometry) */
	rndProbePosIndex  /* (3) 3-vector: index-space position of probe
   						(for debugging ray-casting geometry) */
	rndProbeInside    /* (4) scalar: 1 if kernel support entirely inside data,
   						0 if outside (copied from rndConvo->inside) */
	rndProbeValue     /* (5) scalar: value measured by convolution */
	rndProbeGradVec   /* (6) 3-vector: gradient (in world-space) measured by
   						convolution */
	rndProbeGradMag   /* (7) scalar: gradient magnitude */
	// in both of the following, the "a" is opacity *after* opacity correction
	rndProbeRgba      /* (8) 4-vector: RGBA from the transfer function: the
						   univariate LUT set by rndCtxTxfSet(). Also, if
						   rndCtxLevoySet() is called, the opacity is also
						   multiplied by the Levoy opacity functions. The RGB
						   part of this is the Blinn-Phong "material color" */
	rndProbeRgbaLit   /* (9) 4-vector: rgba from transfer function (above),
	   with RGB shaded according to Blinn-Phong, and then
	   multiplied by the lerp'd depth cueing color */
)

type rndBlend int
//rndBlend
const (
	rndBlendUnknown rndBlend=iota/* (0) blending not known */
	rndBlendMax      /* (1) max of all values (in vectors, per-component) */
	rndBlendSum      /* (2) sum of all values (in vectors, per-component)
					   It is possible to do max and sum blending of
					   probes rndProbeRgba and rndProbeRgbaLit. */
	rndBlendOver     /* (3) "over" operator of RGB color and alpha. This
					   is the one blend that breaks the orthogonality of
					   blending and probing; rndBlendOver only makes sense
					   for rndProbeRgba and rndProbeRgbaLit */
)

/*
  struct rndCamera: A container for information about how a camera is viewing
  some part of world-space.  The look-from, look-at, and up vectors are in
  world-space coordinates (using an orthonormal basis).

  All the fields are named here the same as in FSV Section 5.3
*/
type rndCamera struct {
	fr,  // look-from point
	at,  // look-at point
	up [3]float64 // up
	nc, fc, 		/* near,far clip plane distances,
	   				relative to the look-at point */
	FOV float64   // vertical field-of-view, in degrees
	size [2]uint /* # horz,vert samples of image plane; this
	   				determines aspect ratio ("ar") */
	ortho int /* if non-zero, use orthographic instead of
			   perspective projection. This does not affect how
			   any variables below are calculated, but the
			   rndCamera is best place to store this info. */
	// "output" parameters set by rndCameraUpdate(), based on the above
	ar,        // ((float64)size[0])/size[1]
	d,           // distance between fr and at
	ncv, fcv,    /* (positive) distances, from eye to near and far
	   clipping planes */
	hght, wdth float64  // height and width of visible image plane
	u,        // right-ward basis vector of view space
	v,        // up-ward basis vector of view space
	n [3]float64        // back-ward (into eye) basis vector of view space
	VtoW [16]float64    // view-to-world transorm
}

